package agent

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/rpc"

	"github.com/docker/libkv/store"
	"github.com/hashicorp/serf/serf"
	"github.com/ngaut/log"

	uconf "udup/config"
	"udup/plugins"
)

var (
	ErrForDeletedJob = errors.New("Received err for a deleted job.")
)

type RPCServer struct {
	agent *Agent
}

func (rpcs *RPCServer) GetJob(jobName string, job *Job) error {
	j, err := rpcs.agent.store.GetJob(jobName)
	if err != nil {
		return err
	}

	// Copy the data structure
	job = j

	return nil
}

func (rpcs *RPCServer) StartJob(j Job, reply *serf.NodeResponse) error {
	// Load the job from the store
	job, err := rpcs.agent.store.GetJob(j.Name)
	if err != nil {
		if err == store.ErrKeyNotFound {
			log.Warning(ErrForDeletedJob)
			return ErrForDeletedJob
		}
		return err
	}
	// Lock the job while editing
	if err = job.Lock(); err != nil {
		log.Fatal(err)
	}
	// Get the defined output types for the job, and call them
	for k, v := range job.Processors {
		log.Infof("Processing execution with plugin:%v", k)
		switch v.Driver {
		case plugins.MysqlDriverAttr:
			{
				driver, err := createDriver(v)
				if err != nil {
					return fmt.Errorf("failed to create driver of job '%s': %v",
						job.Name, err)
				}

				errCh := make(chan error)
				v.ErrCh = errCh
				go job.listenOnPanicAbort(v)

				gtidCh := make(chan string, 100)
				v.GtidCh = gtidCh
				go job.listenOnGtid(v)
				/*if k == plugins.ProcessorTypeExtract {
					log.Infof("k:%v",k)
					if job.ParentJob !="" {
						jp, err := rpcs.agent.store.GetJob(job.ParentJob)
						if err != nil {
							return fmt.Errorf("failed to get parent job '%s': %v",
								job.Name, err)
						}
						v.Gtid = jp.Processors["apply"].Gtid
					}else{
						v.Gtid = j.Processors["apply"].Gtid
					}
				}*/
				// Start the job
				err = driver.Start(k, v)
				if err != nil {
					return err
				}
				rpcs.agent.processorPlugins[jobDriver{name: job.Name, tp: k}] = driver
			}
		default:
			{
				return fmt.Errorf("Unknown job type : %+v", v.Driver)
			}
		}
	}

	job.Disabled = false
	for _, p := range job.Processors {
		p.Disabled = false
	}

	if err := rpcs.agent.store.UpsertJob(job); err != nil {
		log.Fatal(err)
	}

	// Release the lock
	if err = job.Unlock(); err != nil {
		log.Fatal(err)
	}

	reply.From = rpcs.agent.config.NodeName
	reply.Payload = []byte("saved")

	// Jobs that have dependent jobs are a bit more expensive because we need to call the Status() method for every execution.
	// Check first if there's dependent jobs and then check for the job status to begin executiong dependent jobs on success.
	if len(job.DependentJobs) > 0 && job.Status() == Success {
		for _, djn := range job.DependentJobs {
			dj, err := rpcs.agent.store.GetJob(djn)
			if err != nil {
				return err
			}
			dj.Run()
		}
	}

	return nil
}

func (rpcs *RPCServer) StopJob(j Job, reply *serf.NodeResponse) error {
	// Load the job from the store
	job, err := rpcs.agent.store.GetJob(j.Name)
	if err != nil {
		if err == store.ErrKeyNotFound {
			log.Warning(ErrForDeletedJob)
			return ErrForDeletedJob
		}
		return err
	}
	// Lock the job while editing
	if err = job.Lock(); err != nil {
		log.Fatal(err)
	}

	for k, v := range rpcs.agent.processorPlugins {
		if k.name == j.Name {
			if err = v.Stop(k.tp); err != nil {
				return err
			}
		}
	}
	job.Disabled = true
	for _, p := range job.Processors {
		p.Disabled = true
	}
	if err := rpcs.agent.store.UpsertJob(job); err != nil {
		log.Fatal(err)
	}

	// Release the lock
	if err = job.Unlock(); err != nil {
		log.Fatal(err)
	}

	reply.From = rpcs.agent.config.NodeName
	reply.Payload = []byte("saved")

	// Jobs that have dependent jobs are a bit more expensive because we need to call the Status() method for every execution.
	// Check first if there's dependent jobs and then check for the job status to begin executiong dependent jobs on success.
	if len(job.DependentJobs) > 0 && job.Status() == Success {
		for _, djn := range job.DependentJobs {
			dj, err := rpcs.agent.store.GetJob(djn)
			if err != nil {
				return err
			}
			dj.Run()
		}
	}

	return nil
}

// createDriver makes a driver for the task
func createDriver(cfg *uconf.DriverConfig) (plugins.Driver, error) {
	driverCtx := plugins.NewDriverContext(cfg.Driver, cfg)
	driver, err := plugins.DiscoverPlugins(cfg.Driver, driverCtx)
	if err != nil {
		return nil, fmt.Errorf("Failed to create driver '%s': %v",
			cfg.Driver, err)
	}
	return driver, err
}

var workaroundRPCHTTPMux = 0

func listenRPC(a *Agent) {
	r := &RPCServer{
		agent: a,
	}

	log.Infof("Registering RPC server: %v", a.getRPCAddr())

	rpc.Register(r)

	oldMux := http.DefaultServeMux
	if workaroundRPCHTTPMux > 0 {
		mux := http.NewServeMux()
		http.DefaultServeMux = mux
	}
	workaroundRPCHTTPMux = workaroundRPCHTTPMux + 1

	rpc.HandleHTTP()

	http.DefaultServeMux = oldMux

	l, e := net.Listen("tcp", a.getRPCAddr())
	if e != nil {
		log.Fatal(e)
	}
	go http.Serve(l, nil)
}

type RPCClient struct {
	//Addres of the server to call
	ServerAddr string
}

func (rpcc *RPCClient) callStartJob(job *Job) error {
	client, err := rpc.DialHTTP("tcp", rpcc.ServerAddr)
	if err != nil {
		log.Errorf("Error dialing: %v,server_addr:%v.", err, rpcc.ServerAddr)
		return err
	}
	defer client.Close()

	// Synchronous call
	var reply serf.NodeResponse
	err = client.Call("RPCServer.StartJob", job, &reply)
	if err != nil {
		log.Errorf("Error calling StartJob: %v", err)
		return err
	}

	return nil
}

func (rpcc *RPCClient) callStopJob(job *Job) error {
	client, err := rpc.DialHTTP("tcp", rpcc.ServerAddr)
	if err != nil {
		log.Errorf("Error dialing: %v,server_addr:%v.", err, rpcc.ServerAddr)
		return err
	}
	defer client.Close()

	// Synchronous call
	var reply serf.NodeResponse
	err = client.Call("RPCServer.StopJob", job, &reply)
	if err != nil {
		log.Errorf("Error calling Stop: %v", err)
		return err
	}

	return nil
}

func (rpcc *RPCClient) GetJob(jobName string) (*Job, error) {
	client, err := rpc.DialHTTP("tcp", rpcc.ServerAddr)
	if err != nil {
		log.Errorf("Error dialing: %v,server_addr:%v.", err, rpcc.ServerAddr)
		return nil, err
	}
	defer client.Close()

	// Synchronous call
	var job Job
	err = client.Call("RPCServer.GetJob", jobName, &job)
	if err != nil {
		log.Errorf("Error calling GetJob: %v", err)
		return nil, err
	}

	return &job, nil
}
