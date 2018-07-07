package fabricgosdkclientcore_test

import (
	"fmt"
	"testing"
	"time"

	hlfsdkutil "github.com/suddutt1/fabricgosdkclientcore"
)

func Test_FabricSDKClient_Init(t *testing.T) {

	clientsMap := initializeClients(t)
	defer cleanup(clientsMap)

}
func Test_ChannelCreation(t *testing.T) {
	clientsMap := initializeClients(t)
	defer cleanup(clientsMap)
	channelID := "settlementchannel"
	if !clientsMap["retail"].SaveChannelInOrderer(channelID, "/home/suddutt1/projects/producttracer/network/"+channelID+".tx", nil) {
		t.Logf("Save channel could not completed successfully ")
		t.FailNow()
	}
	if !clientsMap["retail"].JoinChannel(channelID, nil) {
		t.Logf("Join  channel could not completed successfully for retail")
		t.FailNow()
	}
	if !clientsMap["dist"].JoinChannel(channelID, nil) {
		t.Logf("Join  channel could not completed successfully for dist")
		t.FailNow()
	}
	if !clientsMap["manuf"].JoinChannel(channelID, nil) {
		t.Logf("Join  channel could not completed successfully for manuf")
		t.FailNow()
	}

}
func Test_Install_InitiateChainCode(t *testing.T) {
	clientsMap := initializeClients(t)
	defer cleanup(clientsMap)
	//First install chain code
	ccPath := "github.com/suddutt1/basechaincode"
	goPath := "/home/suddutt1/go"
	ccID := fmt.Sprintf("Basic_%d", time.Now().UnixNano())
	ccPolicy := "And ('ManufacturerMSP.member','DistributerMSP.member','RetailerMSP.member')"
	installInstantiate(clientsMap, "settlementchannel", ccPath, goPath, ccID, ccPolicy, t)
}
func Test_InvokeTrxn_Query(t *testing.T) {
	clientsMap := initializeClients(t)
	defer cleanup(clientsMap)
	ccPath := "github.com/suddutt1/basechaincode"
	goPath := "/home/suddutt1/go"
	ccID := fmt.Sprintf("Basic_%d", time.Now().UnixNano())
	ccPolicy := "And ('ManufacturerMSP.member','DistributerMSP.member','RetailerMSP.member')"
	channelName := "settlementchannel"
	installInstantiate(clientsMap, channelName, ccPath, goPath, ccID, ccPolicy, t)
	userID := "User1"
	ccFn := "save"
	key := fmt.Sprintf("KEY_%d", time.Now().Nanosecond())
	value := fmt.Sprintf("VALUE%d", time.Now().Nanosecond())

	invokeArgs := [][]byte{[]byte(key), []byte(value)}
	peers := []string{"peer0.manuf.net", "peer0.distributer.net", "peer0.retailer.com"}
	rsltBytes, isSuccess, err := clientsMap["manuf"].InvokeTrxn(channelName, userID, ccID, ccFn, invokeArgs, peers, nil)
	if !isSuccess || err != nil {
		t.Logf("Error in Invoke Trxn %v", err)
		t.FailNow()
	}
	t.Logf("Result Invoke Trxn %s", string(rsltBytes))
	//Need to query and verify
	queryArgs := [][]byte{[]byte(key)}
	ccFn = "retrieve"
	queryRsltBytes, isSuccess, err := clientsMap["dist"].Query(channelName, userID, ccID, ccFn, queryArgs, peers, nil)
	if !isSuccess || err != nil {
		t.Logf("Error in Query Trxn %v", err)
		t.FailNow()
	}
	t.Logf("Result Query Trxn %s", string(queryRsltBytes))
	if value != string(queryRsltBytes) {
		t.FailNow()
	}

}
func Test_InvokeTrxn_Query_Loop(t *testing.T) {
	clientsMap := initializeClients(t)
	defer cleanup(clientsMap)
	ccPath := "github.com/suddutt1/basechaincode"
	goPath := "/home/suddutt1/go"
	ccID := fmt.Sprintf("Basic_%d", time.Now().UnixNano())
	ccPolicy := "And ('ManufacturerMSP.member','DistributerMSP.member','RetailerMSP.member')"
	channelName := "settlementchannel"
	installInstantiate(clientsMap, channelName, ccPath, goPath, ccID, ccPolicy, t)
	for {
		userID := "User1"
		ccFn := "save"
		key := fmt.Sprintf("KEY_%d", time.Now().Nanosecond())
		value := fmt.Sprintf("VALUE%d", time.Now().Nanosecond())

		invokeArgs := [][]byte{[]byte(key), []byte(value)}
		peers := []string{"peer0.manuf.net", "peer0.distributer.net", "peer0.retailer.com"}
		rsltBytes, isSuccess, err := clientsMap["manuf"].InvokeTrxn(channelName, userID, ccID, ccFn, invokeArgs, peers, nil)
		if !isSuccess || err != nil {
			t.Logf("Error in Invoke Trxn %v", err)
			t.FailNow()
		}
		t.Logf("Result Invoke Trxn %s", string(rsltBytes))
		time.Sleep(15 * time.Second)
		//Need to query and verify
		queryArgs := [][]byte{[]byte(key)}
		ccFn = "retrieve"
		queryRsltBytes, isSuccess, err := clientsMap["dist"].Query(channelName, userID, ccID, ccFn, queryArgs, peers, nil)
		if !isSuccess || err != nil {
			t.Logf("Error in Query Trxn %v", err)
			t.FailNow()
		}
		t.Logf("Result Query Trxn %s", string(queryRsltBytes))
		if value != string(queryRsltBytes) {
			t.FailNow()
		}

	}

}
func installInstantiate(clientsMap map[string]*hlfsdkutil.FabricSDKClient, channelName, ccPath, goPath, ccID, ccPolicy string, t *testing.T) {
	initArgs := [][]byte{[]byte("init")}
	ccVersion := "1.0"
	isInstallSuccess := clientsMap["retail"].InstallChainCode(ccID, ccVersion, goPath, ccPath, nil)
	if !isInstallSuccess {
		t.Logf("Error in CC installation for  retail")
		t.FailNow()
	}
	isInstallSuccess = clientsMap["dist"].InstallChainCode(ccID, ccVersion, goPath, ccPath, nil)
	if !isInstallSuccess {
		t.Logf("Error in CC installation for  dist")
		t.FailNow()
	}
	isInstallSuccess = clientsMap["manuf"].InstallChainCode(ccID, ccVersion, goPath, ccPath, nil)
	if !isInstallSuccess {
		t.Logf("Error in CC installation for  manuf")
		t.FailNow()
	}
	//Now instantiate at any node
	isInstantiationSuccess, err := clientsMap["manuf"].InstantiateCC(channelName, ccID, ccPath, ccVersion, initArgs, ccPolicy, nil)
	if !isInstantiationSuccess || err != nil {
		t.Logf("Error in CC instantiation for  manuf")
		t.FailNow()
	}
}
func initializeClients(t *testing.T) map[string]*hlfsdkutil.FabricSDKClient {
	fabricNetworkClientManuf := new(hlfsdkutil.FabricSDKClient)
	rslt := fabricNetworkClientManuf.Init("./manuf-client-config.yaml")
	if !rslt {
		t.Logf("Error in sdk initialization manufacturer")
		t.FailNow()
	}
	fabricNetworkClientRetail := new(hlfsdkutil.FabricSDKClient)
	if !fabricNetworkClientRetail.Init("./retailer-client-config.yaml") {
		t.Logf("Error in sdk initialization retailer")
		t.FailNow()
	}
	fabricNetworkClientDist := new(hlfsdkutil.FabricSDKClient)
	if !fabricNetworkClientDist.Init("./dist-client-config.yaml") {
		t.Logf("Error in sdk initialization distributer")
		t.FailNow()
	}
	clientsMap := make(map[string]*hlfsdkutil.FabricSDKClient)
	clientsMap["retail"] = fabricNetworkClientRetail
	clientsMap["dist"] = fabricNetworkClientDist
	clientsMap["manuf"] = fabricNetworkClientManuf
	return clientsMap
}

func cleanup(clientMap map[string]*hlfsdkutil.FabricSDKClient) {
	clientMap["retail"].Shutdown()
	clientMap["dist"].Shutdown()
	clientMap["manuf"].Shutdown()
	fmt.Println("Cleanup completed")

}
