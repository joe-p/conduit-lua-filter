package lua_processor

import (
	"context"
	_ "embed" // used to embed config

	log "github.com/sirupsen/logrus"
	lua "github.com/yuin/gopher-lua"
	luar "layeh.com/gopher-luar"

	"github.com/algorand/conduit/conduit/data"
	"github.com/algorand/conduit/conduit/plugins"
	"github.com/algorand/conduit/conduit/plugins/processors"
	sdk "github.com/algorand/go-algorand-sdk/v2/types"
)

type Config struct {
	// <code>omit-group-transactions</code> configures the filter processor to return the matched transaction without its grouped transactions.
	IncludeGroupTransactions bool `yaml:"omit-group-transactions"`
}

// PluginName to use when configuring.
const PluginName = "lua_processor"

// package-wide init function
func init() {
	processors.Register(PluginName, processors.ProcessorConstructorFunc(func() processors.Processor {
		return &LuaProcessor{}
	}))
}

type LuaProcessor struct {
	logger   *log.Logger
	ctx      context.Context
	luaState *lua.LState
}

// Metadata returns metadata
func (a *LuaProcessor) Metadata() plugins.Metadata {
	return plugins.Metadata{
		Name:         PluginName,
		Description:  "Filter transactions out via a lua script.",
		Deprecated:   false,
		SampleConfig: "",
	}
}

// Config returns the config
func (a *LuaProcessor) Config() string {
	return ""
}

// Init initializes the filter processor
func (a *LuaProcessor) Init(ctx context.Context, _ data.InitProvider, _ plugins.PluginConfig, logger *log.Logger) error {
	a.logger = logger
	a.ctx = ctx

	a.luaState = lua.NewState()
	defer a.luaState.Close()
	if err := a.luaState.DoFile("./static.lua"); err != nil {
		panic(err)
	}
	if err := a.luaState.DoFile("./filter.lua"); err != nil {
		panic(err)
	}

	return nil

}

// Close a no-op for this processor
func (a *LuaProcessor) Close() error {
	return nil
}

// Process processes the input data
func (a *LuaProcessor) Process(input data.BlockData) (data.BlockData, error) {
	var err error
	payset := input.Payset
	var filteredPayset []sdk.SignedTxnInBlock

	for i := 0; i < len(payset); i++ {
		err = a.luaState.CallByParam(
			lua.P{
				Fn:      a.luaState.GetGlobal("processTxn"),
				NRet:    1,
				Protect: true,
			},
			luar.New(a.luaState, payset[i]),
		)

		if err != nil {
			return input, err
		}

		ret := a.luaState.Get(-1)
		a.luaState.Pop(1)

		if ret.String() == "true" {
			a.logger.Infof("accepting txn %d", i)
			filteredPayset = append(filteredPayset, payset[i])
		} else {
			a.logger.Infof("rejecting txn %d", i)
		}
	}

	input.Payset = filteredPayset
	return input, err
}
