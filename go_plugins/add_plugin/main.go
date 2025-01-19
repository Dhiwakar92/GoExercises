package main

import (
	"pluginsdk"

	"github.com/hashicorp/go-plugin"
)

type Addition struct {
}

func (add *Addition) DoSomething(inputs map[string]interface{}) interface{} {
	num1, _ := inputs["num1"].(int)
	num2, _ := inputs["num2"].(int)
	return num1 + num2
}

func main() {
	//fmt.Println("serving the add plugin")
	addImpl := &Addition{}
	addPlugin := &pluginsdk.MyPlugin{
		Impl: addImpl,
	}
	pluginSet := make(map[string]plugin.Plugin)
	pluginSet["add"] = addPlugin
	pluginCfg := plugin.ServeConfig{
		Plugins: pluginSet,
		HandshakeConfig: plugin.HandshakeConfig{
			MagicCookieKey:   "TEST",
			MagicCookieValue: "VALUE",
			ProtocolVersion:  1,
		},
	}
	plugin.Serve(&pluginCfg)
}
