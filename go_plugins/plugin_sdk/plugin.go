package pluginsdk

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
)

type MyPluginImplementation interface {
	DoSomething(inputs map[string]interface{}) interface{}
}

// CLIENT
type MyPluginClient struct {
	c *rpc.Client
}

func (mpc *MyPluginClient) DoSomething(inputs map[string]interface{}) interface{} {
	var resp interface{}
	err := mpc.c.Call("Plugin.DoSomething", inputs, &resp)
	if err != nil {
		panic(err)
	}
	return resp
}

// SERVER
type MyPluginServer struct {
	impl MyPluginImplementation
}

func (mps *MyPluginServer) DoSomething(inputs map[string]interface{}, resp *interface{}) error {
	*resp = mps.impl.DoSomething(inputs)
	return nil
}

// ACTUAL PLUGIN INTERFACE IMPLEMENTATION
type MyPlugin struct {
	// Inputs map[string]interface{}
	Impl MyPluginImplementation
}

func (mp *MyPlugin) Server(mux *plugin.MuxBroker) (interface{}, error) {
	return &MyPluginServer{mp.Impl}, nil
}

func (mp *MyPlugin) Client(mux *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &MyPluginClient{c: c}, nil
}
