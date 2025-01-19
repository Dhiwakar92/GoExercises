package main

import (
	"fmt"
	"os"
	"os/exec"
	"pluginsdk"

	"github.com/hashicorp/go-plugin"
)

func main() {
	pluginSet := make(map[string]plugin.Plugin)
	pluginSet["add"] = &pluginsdk.MyPlugin{}
	pluginCfg := plugin.ClientConfig{
		Plugins: pluginSet,
		Cmd:     exec.Command("./add"),
		HandshakeConfig: plugin.HandshakeConfig{
			MagicCookieKey:   "TEST",
			MagicCookieValue: "VALUE",
			ProtocolVersion:  1,
		},
		Managed: true,
	}
	pluginClient := plugin.NewClient(&pluginCfg)
	addr, err := pluginClient.Start()
	if err != nil {
		fmt.Println("error starting plugin client (NOT TO BE CONFUSED WITH RPC CLIENT !) : ", err)
		os.Exit(1)
	}
	fmt.Println("RPC address to connect is :", addr.String())
	rpcClient, err := pluginClient.Client()
	if err != nil {
		fmt.Println("error fetching the RPC client : ", err)
		os.Exit(1)
	}
	myAddPlugin, err := rpcClient.Dispense("add")
	if err != nil {
		fmt.Println("error dispensing my add plugin : ", err)
		os.Exit(1)
	}

	addPlugin := myAddPlugin.(pluginsdk.MyPluginImplementation)
	num1, num2 := 20, 22
	fmt.Printf("%v+%v=%v\n", num1, num2, addPlugin.DoSomething(map[string]interface{}{"num1": num1, "num2": num2}))
}
