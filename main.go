package main

import (
	"bufio"
	"fabric-go-sdk/sdkInit"
	"fmt"
	"os"
	"strings"
	"time"
)

const (
	cc_name    = "simplecc"
	cc_version = "1.0.0"
)

var App sdkInit.Application

func main() {
	// init orgs information

	orgs := []*sdkInit.OrgInfo{
		{
			OrgAdminUser:  "Admin",
			OrgName:       "Org1",
			OrgMspId:      "Org1MSP",
			OrgUser:       "User1",
			OrgPeerNum:    2,
			OrgAnchorFile: "/home/usr/fabric-go-sdk/fixtures/channel-artifacts/Org1MSPanchors.tx",
		},
	}

	// init sdk env info
	info := sdkInit.SdkEnvInfo{
		ChannelID:        "mychannel",
		ChannelConfig:    "/home/usr/fabric-go-sdk/fixtures/channel-artifacts/channel.tx",
		Orgs:             orgs,
		OrdererAdminUser: "Admin",
		OrdererOrgName:   "OrdererOrg",
		OrdererEndpoint:  "orderer.example.com",
		ChaincodeID:      cc_name,
		ChaincodePath:    "/home/usr/fabric-go-sdk/chaincode/",
		ChaincodeVersion: cc_version,
	}

	// sdk setup
	sdk, err := sdkInit.Setup("config.yaml", &info)
	if err != nil {
		fmt.Println(">> SDK setup error:", err)
		os.Exit(-1)
	}

	// create channel and join
	if err := sdkInit.CreateAndJoinChannel(&info); err != nil {
		fmt.Println(">> Create channel and join error:", err)
		os.Exit(-1)
	}

	// create chaincode lifecycle
	if err := sdkInit.CreateCCLifecycle(&info, 1, false, sdk); err != nil {
		fmt.Printf(">> create chaincode lifecycle error: %v\n", err)
		os.Exit(-1)
	}

	// invoke chaincode set status
	fmt.Println(">> 通过链码外部服务设置链码状态......")

	if err := info.InitService(info.ChaincodeID, info.ChannelID, info.Orgs[0], sdk); err != nil {

		fmt.Println("InitService successful")
		os.Exit(-1)
	}

	App = sdkInit.Application{
		SdkEnvInfo: &info,
	}
	fmt.Println(">> 设置链码状态完成")

	defer info.EvClient.Unregister(sdkInit.BlockListener(info.EvClient))
	defer info.EvClient.Unregister(sdkInit.ChainCodeEventListener(info.EvClient, info.ChaincodeID))

	fmt.Println("==========command format==========")
	fmt.Println("input data: set [key] [value]")
	fmt.Println("query data: get [key]")
	fmt.Println("exit: exit")
	fmt.Println("Please input command: ")
	for {
		inputReader := bufio.NewReader(os.Stdin)
		input, _, err := inputReader.ReadLine()
		if err != nil {
			fmt.Println("input command error:", err)
		}
		cmd := strings.Split(string(input), " ")
		switch cmd[0] {
		case "set":
			ret, err := App.Set(cmd)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("<--- 添加信息　--->：", ret)
		case "get":
			response, err := App.Get(cmd)
			fmt.Println(cmd)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("<--- 查询信息　--->：", response)
		case "exit":
			break
		}
	}
	fmt.Println("==========program end==========")
	time.Sleep(time.Second * 5)

}
