package mysql

import (
	"bytes"
	"database/sql"
	"encoding/gob"
	"fmt"
	"regexp"
	"strings"
	"time"

	stan "github.com/nats-io/go-nats-streaming"
	"github.com/ngaut/log"
	gomysql "github.com/siddontang/go-mysql/mysql"

	uconf "udup/config"
	ubinlog "udup/plugins/mysql/binlog"
	usql "udup/plugins/mysql/sql"
)

const (
	EventsChannelBufferSize       = 1
	ReconnectStreamerSleepSeconds = 5
)

type Extractor struct {
	cfg                      *uconf.DriverConfig
	initialBinlogCoordinates *ubinlog.BinlogCoordinates
	tables                   map[string]*usql.Table
	db                       *sql.DB
	tb                       *ubinlog.TxBuilder
	eventsChannel            chan *ubinlog.BinlogEvent
	currentFde               string
	currentSqlB64            *bytes.Buffer
	reMap                    map[string]*regexp.Regexp
	bp                       *ubinlog.BinlogParser

	stanConn stan.Conn
}

func NewExtractor(cfg *uconf.DriverConfig) *Extractor {
	return &Extractor{
		cfg:           cfg,
		tables:        make(map[string]*usql.Table),
		eventsChannel: make(chan *ubinlog.BinlogEvent, EventsChannelBufferSize),
		reMap:         make(map[string]*regexp.Regexp),
	}
}

func (e *Extractor) InitiateExtractor() error {
	log.Infof("Extract binlog events from the datasource :%v", e.cfg.ConnCfg)
	time.Sleep(10 * time.Second)

	if err := e.initDBConnections(); err != nil {
		return err
	}

	if err := e.initNatsPubClient(); err != nil {
		return err
	}

	if err := e.initiateTxBuilder(); err != nil {
		return err
	}
	log.Infof("Beginning streaming")
	if err := e.streamEvents(); err != nil {
		return err
	}

	return nil
}

// initiateStreaming begins treaming of binary log events and registers listeners for such events
func (e *Extractor) initiateTxBuilder() error {
	e.tb = ubinlog.NewTxBuilder(e.cfg)
	go e.tb.Run()
	return nil
}

func (e *Extractor) initDBConnections() (err error) {
	if e.db, err = usql.CreateDB(e.cfg.ConnCfg); err != nil {
		return err
	}
	if err = e.mysqlGTIDMode(); err != nil {
		return err
	}

	// expecting 'foreign_key_checks = OFF'
	// NB you should use `set global foreign_key_checks = 0'.
	//
	if err = e.checkForeignKey(); err != nil {
		return err
	}

	if err := e.readCurrentBinlogCoordinates(); err != nil {
		return err
	}

	if err = e.initBinlogParser(e.initialBinlogCoordinates); err != nil {
		return err
	}

	return nil
}

func (e *Extractor) mysqlGTIDMode() error {
	query := `SELECT @@gtid_mode`
	var gtidMode string
	if err := e.db.QueryRow(query).Scan(&gtidMode); err != nil {
		return err
	}
	if gtidMode != "ON" {
		return fmt.Errorf("must have GTID enabled: %+v", gtidMode)
	}
	return nil
}

func (e *Extractor) checkForeignKey() error {
	query := `show variables like 'foreign_key_checks'`
	err := usql.QueryRowsMap(e.db, query, func(m usql.RowMap) error {
		if m["Value"].String == "ON" {
			return fmt.Errorf("foreign_key_checks == ON detected. need to be off")
		}

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

// readCurrentBinlogCoordinates reads master status from hooked server
func (e *Extractor) readCurrentBinlogCoordinates() error {
	query := `show /* udup readCurrentBinlogCoordinates */ master status`
	foundMasterStatus := false
	if e.cfg.Gtid != "" {
		gtidSet, err := gomysql.ParseMysqlGTIDSet(e.cfg.Gtid)
		if err != nil {
			return err
		}
		e.initialBinlogCoordinates = &ubinlog.BinlogCoordinates{
			GtidSet: gtidSet,
		}
		foundMasterStatus = true

		return nil
	}
	err := usql.QueryRowsMap(e.db, query, func(m usql.RowMap) error {
		// the mysql GTID set likes this "de278ad0-2106-11e4-9f8e-6edd0ca20947:1-2"
		/*executedGtidSet := strings.Split(rowMap["Executed_Gtid_Set"].String, ":")
		gtidSet, err := gomysql.ParseMysqlGTIDSet(fmt.Sprintf("%s:1", executedGtidSet[0]))*/
		gtidSet, err := gomysql.ParseMysqlGTIDSet(m["Executed_Gtid_Set"].String)
		if err != nil {
			return err
		}
		e.initialBinlogCoordinates = &ubinlog.BinlogCoordinates{
			GtidSet: gtidSet,
		}
		foundMasterStatus = true

		return nil
	})
	if err != nil {
		return err
	}
	if !foundMasterStatus {
		return fmt.Errorf("Got no results from SHOW MASTER STATUS. Bailing out")
	}
	log.Debugf("Streamer binlog coordinates: %+v", *e.initialBinlogCoordinates)
	return nil
}

// initBinlogParser creates and connects the reader: we hook up to a MySQL server as a replica
func (e *Extractor) initBinlogParser(binlogCoordinates *ubinlog.BinlogCoordinates) error {
	binlogParser, err := ubinlog.NewBinlogParser(e.cfg)
	if err != nil {
		return err
	}
	if err := binlogParser.ConnectBinlogStreamer(*binlogCoordinates); err != nil {
		return err
	}
	e.bp = binlogParser
	return nil
}

func (e *Extractor) initNatsPubClient() (err error) {
	sc, err := stan.Connect("test-cluster", "pub1", stan.NatsURL(fmt.Sprintf("nats://%s", e.cfg.NatsAddr)))
	if err != nil {
		log.Fatalf("Can't connect: %v.\nMake sure a NATS Streaming Server is running at: %s", err, fmt.Sprintf("nats://%s", e.cfg.NatsAddr))
	}
	e.stanConn = sc
	return nil
}

// Encode
func Encode(v interface{}) ([]byte, error) {
	b := new(bytes.Buffer)
	enc := gob.NewEncoder(b)
	if err := enc.Encode(v); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func (e *Extractor) GetCurrentBinlogCoordinates() *ubinlog.BinlogCoordinates {
	return e.bp.GetCurrentBinlogCoordinates()
}

func (e *Extractor) GetReconnectBinlogCoordinates() *ubinlog.BinlogCoordinates {
	return &ubinlog.BinlogCoordinates{LogFile: e.GetCurrentBinlogCoordinates().LogFile, LogPos: 4}
}

func (e *Extractor) stopFlag() bool {
	return e.cfg.Disabled
}

func (e *Extractor) streamEvents() error {
	go func() {
		for tx := range e.tb.TxChan {
			msg, err := Encode(tx)
			if err != nil {
				e.cfg.ErrCh <- err
			}
			if err := e.stanConn.Publish("subject", msg); err != nil {
				e.cfg.ErrCh <- err
			}
		}
	}()
	// The next should block and execute forever, unless there's a serious error
	var successiveFailures int64
	var lastAppliedRowsEventHint ubinlog.BinlogCoordinates
	for {
		if err := e.bp.StreamEvents(e.stopFlag, e.tb.EvtChan); err != nil {
			time.Sleep(ReconnectStreamerSleepSeconds * time.Second)
			if e.stopFlag() {
				return nil
			}
			log.Infof("StreamEvents encountered unexpected error: %+v", err)

			// See if there's retry overflow
			if e.bp.LastAppliedRowsEventHint.Equals(&lastAppliedRowsEventHint) {
				successiveFailures += 1
			} else {
				successiveFailures = 0
			}
			if successiveFailures > e.cfg.MaxRetries {
				return fmt.Errorf("%d successive failures in streamer reconnect at coordinates %+v", successiveFailures, e.GetReconnectBinlogCoordinates())
			}

			// Reposition at same binlog file.
			lastAppliedRowsEventHint = e.bp.LastAppliedRowsEventHint
			log.Infof("Reconnecting... Will resume at %+v", lastAppliedRowsEventHint)
			if err := e.initBinlogParser(e.GetReconnectBinlogCoordinates()); err != nil {
				return err
			}
			e.bp.LastAppliedRowsEventHint = lastAppliedRowsEventHint
		}
	}
}

func (e *Extractor) clearTables() {
	e.tables = make(map[string]*usql.Table)
}

func (e *Extractor) getTableFromDB(db *sql.DB, schema string, name string) (*usql.Table, error) {
	table := &usql.Table{}
	table.Schema = schema
	table.Name = name

	err := e.getTableColumns(db, table)
	if err != nil {
		return nil, err
	}

	err = e.getTableIndex(db, table)
	if err != nil {
		return nil, err
	}

	if len(table.Columns) == 0 {
		return nil, fmt.Errorf("invalid table %s.%s", schema, name)
	}

	return table, nil
}

func (e *Extractor) getTableColumns(db *sql.DB, table *usql.Table) error {
	if table.Schema == "" || table.Name == "" {
		return fmt.Errorf("schema/table is empty")
	}

	query := fmt.Sprintf("show columns from %s.%s", table.Schema, table.Name)
	rows, err := usql.QuerySQL(db, query)
	if err != nil {
		return err
	}
	defer rows.Close()

	rowColumns, err := rows.Columns()
	if err != nil {
		return err
	}

	idx := 0
	for rows.Next() {
		datas := make([]sql.RawBytes, len(rowColumns))
		values := make([]interface{}, len(rowColumns))

		for i := range values {
			values[i] = &datas[i]
		}

		err = rows.Scan(values...)
		if err != nil {
			return err
		}

		column := &usql.Column{}
		column.Idx = idx
		column.Name = string(datas[0])

		// Check whether column has unsigned flag.
		if strings.Contains(strings.ToLower(string(datas[1])), "unsigned") {
			column.Unsigned = true
		}

		table.Columns = append(table.Columns, column)
		idx++
	}

	if rows.Err() != nil {
		return rows.Err()
	}

	return nil
}

func (e *Extractor) getTableIndex(db *sql.DB, table *usql.Table) error {
	if table.Schema == "" || table.Name == "" {
		return fmt.Errorf("schema/table is empty")
	}

	query := fmt.Sprintf("show index from %s.%s", table.Schema, table.Name)
	rows, err := usql.QuerySQL(db, query)
	if err != nil {
		return err
	}
	defer rows.Close()

	rowColumns, err := rows.Columns()
	if err != nil {
		return err
	}

	var keyName string
	var columns []string
	for rows.Next() {
		datas := make([]sql.RawBytes, len(rowColumns))
		values := make([]interface{}, len(rowColumns))

		for i := range values {
			values[i] = &datas[i]
		}

		err = rows.Scan(values...)
		if err != nil {
			return err
		}

		nonUnique := string(datas[1])
		if nonUnique == "0" {
			if keyName == "" {
				keyName = string(datas[2])
			} else {
				if keyName != string(datas[2]) {
					break
				}
			}

			columns = append(columns, string(datas[4]))
		}
	}

	if rows.Err() != nil {
		return rows.Err()
	}

	table.IndexColumns = usql.FindColumns(table.Columns, columns)
	return nil
}

func (e *Extractor) getTable(schema string, table string) (*usql.Table, error) {
	key := fmt.Sprintf("%s.%s", schema, table)

	value, ok := e.tables[key]
	if ok {
		return value, nil
	}

	t, err := e.getTableFromDB(e.db, schema, table)
	if err != nil {
		return nil, err
	}

	e.tables[key] = t
	return t, nil
}

func (e *Extractor) matchDB(patternDBS []string, a string) bool {
	for _, b := range patternDBS {
		if e.matchString(b, a) {
			return true
		}
	}
	return false
}

func (e *Extractor) matchString(pattern string, t string) bool {
	if re, ok := e.reMap[pattern]; ok {
		return re.MatchString(t)
	}
	return pattern == t
}

func (e *Extractor) matchTable(patternTBS []uconf.TableName, tb uconf.TableName) bool {
	for _, ptb := range patternTBS {
		retb, oktb := e.reMap[ptb.Name]
		redb, okdb := e.reMap[ptb.Schema]

		if oktb && okdb {
			if redb.MatchString(tb.Schema) && retb.MatchString(tb.Name) {
				return true
			}
		}
		if oktb {
			if retb.MatchString(tb.Name) && tb.Schema == ptb.Schema {
				return true
			}
		}
		if okdb {
			if redb.MatchString(tb.Schema) && tb.Name == ptb.Name {
				return true
			}
		}

		//create database or drop database
		if tb.Name == "" {
			if tb.Schema == ptb.Schema {
				return true
			}
		}

		if ptb == tb {
			return true
		}
	}

	return false
}

func (e *Extractor) skipRowEvent(schema string, table string) bool {
	if e.cfg.ReplicateDoTable != nil || e.cfg.ReplicateDoDb != nil {
		table = strings.ToLower(table)
		//if table in tartget Table, do this event
		for _, d := range e.cfg.ReplicateDoTable {
			if e.matchString(d.Schema, schema) && e.matchString(d.Name, table) {
				return false
			}
		}

		//if schema in target DB, do this event
		if e.matchDB(e.cfg.ReplicateDoDb, schema) && len(e.cfg.ReplicateDoDb) > 0 {
			return false
		}

		return true
	}
	return false
}

func (e *Extractor) skipQueryEvent(sql string, schema string) bool {
	sql = strings.ToUpper(sql)

	if strings.HasPrefix(sql, "GRANT REPLICATION SLAVE ON") {
		return true
	}

	if strings.HasPrefix(sql, "BEGIN") {
		return true
	}

	if strings.HasPrefix(sql, "COMMIT") {
		return true
	}

	if strings.HasPrefix(sql, "FLUSH PRIVILEGES") {
		return true
	}

	return false
}

func (e *Extractor) skipQueryDDL(sql string, schema string) bool {
	tb, err := usql.ParserDDLTableName(sql)
	if err != nil {
		log.Warnf("[get table failure]:%s %s", sql, err)
	}

	if err == nil && (e.cfg.ReplicateDoTable != nil || e.cfg.ReplicateDoDb != nil) {
		//if table in target Table, do this sql
		if tb.Schema == "" {
			tb.Schema = schema
		}
		if e.matchTable(e.cfg.ReplicateDoTable, tb) {
			return false
		}

		// if  schema in target DB, do this sql
		if e.matchDB(e.cfg.ReplicateDoDb, tb.Schema) {
			return false
		}
		return true
	}
	return false
}

func (e *Extractor) Shutdown() error {
	if e.stopFlag() {
		return nil
	}
	e.bp.Close()
	e.stanConn.Close()
	close(e.eventsChannel)
	err := usql.CloseDBs(e.db)
	if err != nil {
		return err
	}
	e.cfg.Disabled = true
	log.Infof("Closed streamer connection.")
	return nil
}
