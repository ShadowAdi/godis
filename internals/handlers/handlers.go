package handlers

import (
	"godis/internals/resp"
	"sync"
)

var SETs = map[string]string{}
var SETsMS = sync.RWMutex{}

func ping(args []resp.Value) resp.Value {
	return resp.Value{
		Typ: "string",
		Str: "PONG",
	}
}

func SET(args []resp.Value) resp.Value {
	if len(args) != 2 {
		return resp.Value{Typ: "error", Str: "ERR wrong number of arguments for 'set' command"}
	}
	key := args[0].Bulk
	value := args[1].Bulk
	SETsMS.Lock()
	SETs[key] = value
	SETsMS.Unlock()
	return resp.Value{
		Typ: "string",
		Str: "OK",
	}
}

var Handlers = map[string]func([]resp.Value) resp.Value{
	"PING": ping,
}
