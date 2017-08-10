package binlog

import (
	gomysql "github.com/siddontang/go-mysql/mysql"
	"testing"
)

func TestParseMysqlGTIDSet(t *testing.T) {
	type args struct {
		gtidset string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{"t1", args{"de278ad0-2106-11e4-9f8e-6edd0ca20947:1:2:3:4:6:5"}, false},
		{"t2", args{"de278ad0-2106-11e4-9f8e-6edd0ca20947:1-5,de278ad0-2106-11e4-9f8e-6edd0ca20947:6-7:10-20"}, false},
		{"t3", args{"de278ad0-2106-11e4-9f8e-6edd0ca20947:6-7:1-4:10-20"}, false},
		{"t4", args{"de278ad0-2106-11e4-9f8e-6edd0ca20947:14972:14977:14983:14984:14989:14992:14995:14996:15000:15001:15002:15007:15009:15010:15011:15012:15014:15015:15017:15018:15019:15020:15021:15022:15024:15026:15028:15030:15032:15033:15035:15036:15041:15042:15043:15046:15047:15048:15049:15050:15051:15055:15060:15061:15062:15064:15066:15070:15071:15072:15073:15074:15075:15077:15078:15079:15081:15082:15084:15085:15086:15087:15089:15093:15094:15095:15097:15100:15102:15104:15106:15107:15109:15112:15114:15116:15117:15119:15120:15121:15122:15124:15125:15127:15128:15129:15130:15132:15136:15138:15139:15141:15142:15144:15145:15146:15148:15149:15150:15153:15155:15157:15159:15160:15161:15166:15168:15169:15172:15174:15175:15176:15177:15179:15181:15185:15186:15187:15189:15191:15192:15194:15195:15197:15199:15200:15201:15203:15205:15206:15208:15209:15211:15215:15216:15218:15219:15220:15221:15222:15225:15226:15227:15228:15229:15230:15232:15233:15235:15239:15240:15241:15242:15243:15244:15246:15248:15249:15250:15251:15252:15254:15255:15258:15259:15260:15265:15268:15271:15273:15274:15275:15276:15277:15279:15282:15285:15286:15288:15290:15292:15295:15296:15297:15298:15301:15303:15304:15306:15307:15308:15309:15310:15311:15312:15313:15314:15316:15317:15318:15319:15321:15324:15328:15329:15330:15331:15335:15337:15339:15341:15342:15344:15345:15348:15349:15350:15351:15352:15355:15361:15362:15363:15364:15369:15371:15374:15375:15376:15377:15379:15380:15383:15384:15385:15386:15389:15391:15392:15394:15396:15397:15399:15400:15401:15404:15405:15407:15409:15410:15411:15417:15419:15420:15421:15424:15425:15426:15428:15430:15431:15433:15434:15436:15439:15441:15442:15443:15444:15445:15447:15449:15451:15452:15453:15456:15461:15463:15466:15467:15470:15472:15473:15476:15480:15482:15483:15485:15489:15490:15492:15493:15494:15497:15501:15502:15503:15504:15506:15507:15509:15510:15511:15517:15518:15521:15522:15523:15524:15526:15529:15534:15535:15537:15538:15539:15540:15541:15542:15545:15546:15549:15552:15556:15557:15560:15562:15563:15564:15568:15570:15575:15577:15578:15579:15583:15589:15590:15592:15593:15595:15597:15599:15600:15601:15602:15604:15606:15607:15608:15609:15611:15614:15615:15616:15620:15621:15622:15626:15627:15628:15629:15630:15633:15634:15636:15637:15639:15640:15642:15645:15646:15647:15648:15649:15651:15652:15654:15655:15656:15657:15658:15659:15662:15663:15664:15667:15669:15671:15673:15674:15678:15680:15681:15682:15683:15692:15696:15700:15701:15702:15705:15706:15708:15711:15713:15714:15715:15716:15718:15720:15723:15724:15726:15727:15729:15730:15731:15735:15736:15737:15739:15740:15742:15744:15748:15749:15750:15752:15753:15754:15757:15758:15759:15760:15761:15762:15766:15767:15768:15771:15773:15777:15779:15781:15783:15785:15787:15788:15789:15793:15795:15797:15798:15799:15802:15803:15804:15806:15809:15810:15811:15812:15814:15815:15816:15819:15821:15822:15824:15826:15827:15828:15831:15832:15835:15836:15838:15840:15841:15842:15845:15847:15849:15850:15854:15856:15857:15858:15860:15861:15863:15864:15866:15868:15870:15871:15872:15876:15877:15879:15880:15881:15882:15883:15884:15885:15886:15890:15892:15893:15895:15896:15897:15899:15901:15903:15905:15910:15913:15914:15916:15917:15919:15920:15921:15922:15924:15925:15926:15927:15929:15931:15933:15935:15936:15939:15943:15945:15947:15949:15950:15951:15952:15953:15957:15960:15961:15962:15963:15964:15965:15967:15968:15969:15971:15973:15974:15976:15977:15982:15983:15984:15985:15986:15987:15988:15989:15990:15996:15998:15999:16000:16003:16004:16005:16006:16011:16013:16016:16020:16022:16028:16029:16031:16037:16039:16040:16041:16043:16044:16046:16047:16048:16050:16052:16054:16056:16059:16062:16064:16067:16068:16069:16070:16072:16075:16079:16080:16081:16082:16083:16084:16086:16087:16088:16089:16091:16093:16096:16097:16102:16103:16104:16105:16108:16109:16110:16111:16113:16115:16116:16118:16125:16127:16128:16129:16131:16132:16134:16135:16136:16137:16142:16143:16145:16146:16147:16148:16149:16150:16153:16155:16156:16157:16158:16159:16160:16163:16166:16167:16168:16169:16175:16177:16178:16179:16180:16181:16182:16184:16186:16187:16191:16193:16195:16196:16198:16199:16201:16203:16205:16206:16207:16209:16212:16213:16214:16215:16217:16218:16221:16222:16223:16226:16228:16229:16230:16232:16233:16234:16236:16238:16239:16242:16244:16248:16250:16253:16254:16256:16260:16261:16266:16269:16272:16273:16275:16276:16277:16279:16280:16281:16284:16285:16287:16291:16293:16294:16296:16298:16300:16303:16304:16305:16306:16310:16312:16313:16317:16319:16320:16323:16326:16328:16329:16331:16333:16336:16337:16339:16341:16344:16345:16347:16348:16352:16355:16359:16362:16364:16365:16366:16368:16372:16373:16375:16379:16380:16382:16383:16385:16387:16389:16390:16391:16392:16394:16396:16398:16399:16400:16402:16406:16408:16409:16411:16414:16416:16417:16422:16423:16424:16425:16427:16431:16432:16433:16436:16439:16440:16441:16443:16444:16446:16451:16452:16453:16455:16458:16460:16463:16466:16467:16469:16475:16476:16478:16479:16482:16483:16487:16488:16490:16494:16495:16498:16500:16502:16506:16508:16510:16511:16514:16517:16518:16521:16524:16525:16526:16528:16529:16530:16533:16535:16536:16537:16538:16540:16541:16545:16547:16549:16551:16553:16555:16556:16557:16560:16562:16563:16564:16565:16567:16572:16574:16576:16577:16579:16580:16584:16585:16587:16589:16590:16591:16593:16598:16599:16600:16602:16603:16604:16605:16606:16607:16609:16611:16613:16615:16616:16618:16620:16621:16623:16626:16627:16629:16630:16632:16633:16634:16635:16641:16642:16646:16650:16651:16653:16655:16656:16657:16658:16659:16663:16665:16667:16668:16669:16672:16674:16678:16680:16681:16683:16684:16685:16688:16690:16691:16693:16695:16697:16700:16703:16705:16707:16708:16709:16714:16715:16718:16721:16722:16723:16724:16726:16728:16729:16730:16735:16739:16740:16741:16742:16747:16748:16749:16751:16753:16755:16756:16758:16760:16761:16762:16764:16765:16766:16770:16771:16773:16774:16776:16781:16783:16784:16786:16789:16790:16792:16794:16795:16796:16798:16802:16805:16807:16808:16809:16813:16814:16815:16816:16818:16819:16821:16822:16827:16829:16831:16832:16833:16834:16835:16837:16838:16844:16846:16847:16850:16852:16853:16854:16855:16859:16861:16863:16866:16868:16869:16870:16872:16873:16876:16877:16878:16879:16882:16885:16888:16890:16891:16892:16893:16894:16895:16896:16899:16900:16903:16904:16906:16907:16908:16910:16913:16915:16916:16918:16920:16921:16923:16926:16928:16929:16932:16933:16934:16936:16938:16941:16942:16946:16947:16948:16951:16952:16953:16956:16957:16958:16959:16960:16961:16963:16965:16966:16967:16969:16971:16972:16974:16975:16977:16979:16980:16982:16984:16987:16989:16990:16991:16992:16994:16995:16996:16998:17001:17003:17004:17007:17010:17011:17013:17014:17015:17017:17019:17021:17023:17027:17029:17031:17034:17036:17038:17040:17041:17042:17045:17048:17049:17050:17051:17053:17055:17057:17059:17060:17062:17064:17068:17069:17071:17072:17075:17076:17078:17081:17082:17084:17085:17087:17091:17092:17094:17096:17097:17099:17101:17103:17104:17106:17107:17108:17110:17115:17116:17118:17120:17121:17124:17127:17129:17131:17132:17134:17135:17138:17142:17143:17145:17147:17148:17150:17152:17157:17158:17159:17161:17163:17165:17167:17171:17175:17180:17181:17183:17186:17187:17192:17193:17195:17197:17198:17199:17205:17206:17207:17209:17213:17215:17216:17218:17220:17224:17225:17227:17229:17231:17234:17236:17237:17240:17242:17243:17245:17246:17248:17249:17255:17257:17259:17262:17266:17267:17268:17270:17272:17274:17278:17279:17281:17282:17285:17287:17289:17290:17292:17295:17296:17297:17300:17301:17303:17304:17306:17307:17308:17309:17310:17311:17313:17315:17317:17321:17322:17326:17329:17330:17331:17332:17336:17338:17339:17341:17343:17346:17348:17349:17350:17352:17354:17357:17359:17360:17361:17362:17363:17367:17369:17370:17371:17373:17375:17377:17378:17379:17384:17385:17386:17387:17390:17393:17395:17396:17401:17402:17404:17406:17407:17408:17409:17413:17415:17417:17418:17421:17423:17426:17427:17430:17431:17432:17435:17438:17440:17441:17442:17446:17447:17450:17451:17453:17454:17455:17458:17460:17463:17465:17468:17471:17472:17474:17475:17476:17477:17478:17480:17481:17483:17484:17485:17486:17490:17492:17493:17503:17506:17511:17514:17515:17517:17519:17521:17523:17525:17527:17528:17530:17531:17532:17534:17535:17537:17539:17542:17543:17546:17547:17549:17551:17552:17553:17555:17558:17560:17561:17563:17565:17566:17567:17568:17569:17570:17572:17573:17574:17576:17580:17581:17582:17589:17591:17593:17594:17595:17597:17598:17600:17603:17606:17611:17613:17614:17616:17617:17618:17621:17623:17625:17627:17630:17631:17633:17638:17639:17640:17641:17644:17645:17648:17650:17652:17653:17656:17659:17661:17665:17666:17669:17670:17671:17674:17678:17681:17682:17684:17686:17688:17689:17692:17694:17695:17696:17698:17700:17701:17704:17705:17706:17710:17711:17714:17715:17718:17719:17720:17721:17724:17726:17728:17729:17733:17735:17738:17740:17741:17743:17748:17749:17750:17751:17752:17755:17758:17759:17760:17761:17762:17764:17765:17766:17768:17770:17771:17772:17774:17775:17777:17778:17779:17781:17783:17786:17787:17789:17790:17793:17795:17797:17802:17803:17804:17805:17809:17811:17812:17817:17818:17819:17822:17826:17828:17829:17834:17835:17836:17837:17838:17839:17840:17842:17844:17848:17850:17851:17852:17854:17856:17857:17858:17859:17861:17863:17868:17870:17871:17872:17873:17875:17876:17878:17881:17883:17884:17885:17887:17888:17889:17890:17891:17892:17893:17895:17897:17900:17901:17903:17904:17906:17910:17912:17914:17915:17917:17919:17921:17928:17929:17932:17934:17935:17940:17942:17943:17944:17948:17949:17950:17952:17954:17956:17959:17960:17961:17963:17964:17965:17967:17968:17970:17971:17973:17974:17977:17979:17982:17985:17990:17991:17992:17993:17998:18000:18001:18002:18003:18006:18008:18009:18010:18012:18015:18017:18018:18020:18022:18024:18026:18027:18031:18034:18036:18041:18042:18043:18045:18046:18047:18048:18049:18051:18055:18058:18060:18061:18062:18065:18067:18068:18071:18072:18074:18075:18077:18079:18080:18084:18087:18088:18090:18091:18092:18095:18097:18098:18099:18100:18101:18103:18107:18110:18112:18113:18118:18119:18121:18122:18123:18124:18126:18127:18133:18136:18143:18146:18147:18151:18152:18158:18166:18172:18174:18175:18177:18179:18180:18181:18182:18183:18184:18189:18193:18194:18198:18203:18206:18207:18209:18211:18212:18217:18218:18221:18225:18230:18235:18238:18240:18241:18245:18246:18249:18250:18251:18253:18256:18258:18259:18264:18268:18278:18279:18280:18281:18283:18284:18288:18298:18300:18302:18303:18307:18310:18311:18312:18314:18315:18317:18318:18322:18326:18327:18330:18332:18334:18335:18339:18341:18344:18345:18346:18348:18350:18354:18361:18362:18363:18373:18374:18375:18377:18378:18383:18388:18393:18394:18398:18400:18402:18403:18404:18406:18408:18410:18414:18415:18417:18422:18423:18424:18426:18431:18433:18438:18439:18440:18443:18457:18458:18467:18468:18469:18470:18474:18476:18478:18480:18481:18482:18484:18485:18487:18490:18491:18492:18494:18496:18498:18499:18502:18503:18505:18508:18511:18513:18515:18516:18518:18523:18524:18526:18527:18533:18534:18537:18539:18541:18542:18543:18545:18549:18551:18552:18553:18555:18556:18557:18558:18560:18562:18563:18564:18565:18566:18568:18571:18572:18575:18579:18580:18581:18582:18583:18585:18587:18588:18592:18593:18596:18597:18600:18604:18606:18607:18609:18611:18612:18616:18618:18619:18626:18630:18631:18632:18635:18636:18638:18639:18640:18641:18643:18645:18646:18651:18652:18656:18658:18659:18662:18666:18669:18671:18674:18676:18679:18681:18683:18684:18687:18688:18689:18691:18693:18694:18698:18700:18702:18705:18707:18709:18710:18712:18713:18715:18717:18719:18720:18721:18722:18724:18725:18730:18731:18735:18738:18740:18744:18747:18748:18751:18755:18757:18758:18760:18764:18765:18767:18769:18771:18773:18775:18776:18777:18782:18785:18788:18789:18790:18794:18795:18797:18798:18799:18806:18808:18809:18811:18813:18815:18816:18817:18818:18819:18820:18821:18822:18825:18828:18829:18830:18834:18837:18839:18841:18844:18845:18846:18849:18853:18855:18857:18859:18861:18862:18864:18865:18868:18869:18870:18872:18876:18878:18880:18881:18885:18886:18887:18888:18889:18890:18891:18894:18899:18900:18901:18902:18904:18905:18906:18907:18908:18909:18914:18918:18921:18923:18924:18927:18929:18930:18932:18936:18940:18948:18949:18950:18951:18958:18963:18965:18968:18973:18977:18985:18988:18990:18991:18996:18997:18999:19002:19005:19009:19010:19013:19019:19021:19022:19028:19033:19034:19036:19038:19041:19043:19045:19048:19049:19051:19052:19054:19057:19058:19060:19061:19063:19064:19066:19068:19069:19072:19074:19075:19076:19078:19080:19082:19084:19087:19090:19091:19092:19096:19099:19101:19102:19104:19105:19106:19108:19110:19113:19115:19117:19118:19119:19120:19122:19126:19128:19136:19137:19138:19139:19144:19145:19150:19152:19153:19157:19158:19160:19162:19166:19168:19170:19172:19174:19176:19179:19181:19183:19185:19186:19190:19193:19195:19196:19198:19199:19201:19204:19205:19206:19208:19211:19217:19218:19220:19224:19225:19226:19229:19231:19233:19235:19237:19238:19240:19244:19246:19248:19250:19252:19253:19255:19257:19259:19260:19262:19267:19272:19276:19279:19280:19282:19284:19285:19290:19291:19293:19295:19296:19299:19302:19304:19306:19309:19311:19313:19314:19316:19317:19319:19320:19324:19325:19326:19328:19330:19331:19332:19334:19339:19342:19344:19346:19347:19349:19351:19353:19357:19361:19362:19363:19364:19366:19368:19370:19371:19372:19378:19381:19383:19385:19386:19387:19391:19393:19396:19398:19399:19400:19402:19404:19408:19410:19413:19414:19415:19416:19417:19419:19421:19424:19428:19430:19431:19433:19434:19439:19442:19443:19448:19449:19450:19456:19457:19459:19461:19462:19465:19466:19467:19475:19476:19480:19483:19490:19492:19493:19494:19496:19499:19500:19501:19506:19508:19509:19510:19514:19517:19520:19521:19522:19526:19527:19528:19529:19532:19534:19535:19536:19538:19540:19542:19543:19545:19547:19551:19552:19553:19554:19556:19558:19559:19560:19563:19566:19568:19569:19570:19572:19574:19575:19578:19579:19584:19591:19594:19596:19605:19607:19608:19610:19611:19612:19618:19622:19630:19632:19638:19640:19649:19655:19657:19660:19661:19669:19674:19676:19677:19678:19679:19680:19681:19684:19686:19692:19693:19698:19703:19705:19711:19713:19716:19717:19720:19722:19725:19726:19728:19730:19732:19736:19737:19742:19743:19745:19748:19749:19750:19752:19755:19758:19759:19761:19762:19765:19770:19772:19773:19778:19779:19782:19784:19786:19787:19790:19792:19795:19797:19804:19806:19807:19809:19811:19816:19819:19822:19823:19824:19826:19828:19834:19835:19837:19838:19840:19846:19850:19851:19853:19855:19857:19858:19860:19862:19863:19866:19867:19868:19871:19875:19878:19879:19883:19888:19889:19890:19892:19896:19898:19899:19900:19902:19904:19907:19909:19912:19913:19915:19916:19918:19920:19922:19924:19927:19928:19929:19933:19936:19937:19938:19947:19948:19949:19953:19956:19958:19959:19962:19967:19968:19971:19972:19975:19978:19979:19980:19982:19984:19986:19987:19991:19993:19994:19995:19998:20002:20005:20007:20008:20012:20013:20014:20015:20018:20019:20022:20023:20025:20026:20032:20033:20034:20035:20037:20038:20040:20041:20043:20044:20046:20048:20052:20053:20055:20057:20058:20060:20061:20062:20063:20064:20067:20069:20070:20071:20072:20076:20078:20083:20087:20088:20091:20093:20094:20098:20101:20103:20105:20107:20108:20109:20112:20117:20119:20120:20121:20124:20127:20128:20130:20133:20136:20137:20139:20143:20147:20149:20150:20153:20156:20157:20159:20161:20169:20170:20171:20174:20177:20179:20181:20184:20185:20186:20187:20190:20195:20197:20199:20200:20202:20207:20209:20210:20211:20212:20214:20215:20217:20221:20222:20223:20225:20227:20229:20232:20233:20236:20238:20239:20241:20243:20245:20246:20247:20249:20250:20254:20256:20258:20259:20261:20263:20264:20265:20267:20269:20270:20272:20277:20280:20281:20282:20283:20285:20286:20287:20290:20291:20292:20293:20294:20295:20296:20297:20298:20299:20303:20304:20305:20306:20309:20311:20315:20318:20322:20323:20325:20328:20330:20334:20335:20337:20338:20341:20342:20344:20345:20348:20349:20350:20351:20352:20354:20356:20357:20358:20363:20364:20366:20369:20370:20371:20375:20378:20382:20385:20389:20391:20393:20398:20399:20400:20401:20404:20405:20406:20407:20411:20412:20418:20421:20422:20423:20426:20427:20429:20430:20431:20432:20434:20437:20440:20443:20444:20446:20451:20453:20454:20459:20460:20461:20464:20465:20466:20469:20471:20472:20474:20480:20482:20483:20487:20490:20492:20493:20496:20498:20499:20500:20504:20506:20508:20511:20513:20515:20519:20522:20523:20525:20527:20528:20529:20531:20532:20534:20535:20537:20538:20540:20542:20543:20544:20546:20548:20549:20558:20563:20566:20568:20570:20573:20574:20582:20585:20586:20587:20588:20589:20590:20596:20600:20601:20602:20605:20607:20608:20609:20610:20611:20612:20614:20615:20617:20623:20624:20625:20629:20632:20635:20636:20638:20641:20645:20646:20647:20648:20650:20651:20652:20654:20659:20660:20661:20664:20665:20666:20669:20672:20675:20678:20680:20683:20685:20687:20688:20691:20693:20695:20698:20700:20701:20702:20706:20707:20712:20715:20718:20719:20720:20721:20724:20727:20729:20730:20731:20732:20734:20736:20737:20739:20740:20742:20745:20747:20749:20754:20755:20757:20758:20761:20765:20767:20769:20770:20771:20772:20775:20781:20783:20784:20786:20787:20789:20791:20792:20794:20799:20802:20804:20805:20808:20810:20812:20814:20815:20817:20818:20819:20823:20824:20825:20826:20827:20828:20830:20834:20835:20836:20838:20839:20854:20855:20857:20861:20862:20870:20877:20881:20886:20888:20890:20891:20893:20895:20897:20900:20901:20906:20907:20911:20914:20916:20924:20925:20926:20927:20929:20930:20932:20933:20934:20936:20938:20941:20942:20948:20952:20953:20958:20961:20962:20965:20969:20974:20977:20978:20979:20981:20983:20990:20991:20993:20994:20995:20998:21000:21004:21008:21010:21011:21012:21014:21016:21017:21019:21022:21023:21025:21026:21028:21030:21033:21034:21035:21038:21040:21041:21047:21049:21053:21055:21056:21058:21060:21063:21065:21067:21069:21070:21076:21081:21091:21093:21094:21095:21097:21099:21103:21104:21106:21116:21120:21124:21125:21126:21128:21131:21133:21135:21138:21142:21146:21147:21150:21151:21155:21156:21157:21161:21164:21165:21166:21167:21173:21174:21175:21178:21179:21184:21188:21191:21192:21196:21197:21199:21200:21203:21206:21207:21210:21212:21213:21214:21216:21219:21222:21228:21231:21232:21234:21237:21240:21241:21243:21244:21249:21251:21252:21253:21257:21258:21260:21262:21264:21265:21267:21268:21272:21273:21274:21276:21278:21280:21281:21282:21283:21288:21290:21291:21292:21293:21296:21298:21299:21303:21312:21315:21317:21318:21320:21323:21324:21325:21326:21327:21330:21335:21337:21341:21342:21346:21347:21351:21353:21354:21355:21356:21357:21360:21361:21362:21365:21366:21368:21369:21371:21373:21376:21377:21378:21380:21384:21386:21387:21390:21392:21393:21394:21398:21401:21404:21407:21409:21413:21416:21417:21419:21420:21423:21424:21428:21429:21432:21435:21437:21439:21442:21445:21447:21448:21449:21455:21457:21459:21461:21467:21476:21478:21479:21481:21488:21489:21494:21496:21497:21499:21500:21502:21504:21505:21507:21512:21514:21523:21524:21527:21531:21533:21534:21537:21551:21553:21557"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotGtidset, err := gomysql.ParseMysqlGTIDSet(tt.args.gtidset)
			if err != nil {
				t.Errorf("ParseMysqlGTIDSet error = %v", err)
				return
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseMysqlGTIDSet error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			println(gotGtidset.String())
		})
	}
}
