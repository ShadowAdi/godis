package handlers

import (
	"godis/internals/resp"
	"sync"
)

var SETs = map[string]string{}
var SETsMS = sync.RWMutex{}
var HSETs = map[string]map[string]string{}
var HSETsMS = sync.RWMutex{}

func hset(args []resp.Value) resp.Value {
	if len(args) != 3 {
		return resp.Value{Typ: "error", Str: "ERR wrong number of arguments for 'hset' command"}
	}

	hash := args[0].Bulk
	key := args[1].Bulk
	value := args[2].Bulk

	HSETsMS.Lock()
	if _, ok := HSETs[hash]; !ok {
		HSETs[hash] = map[string]string{}
	}
	HSETs[hash][key] = value
	HSETsMS.Unlock()
	return resp.Value{Typ: "string", Str: "OK"}
}

func gset(args []resp.Value) resp.Value {
	if len(args) != 2 {
		return resp.Value{Typ: "error", Str: "ERR wrong number of arguments for 'hset' command"}
	}

	hash := args[0].Bulk
	key := args[1].Bulk

	HSETsMS.RLock()
	value, ok := HSETs[hash][key]
	HSETsMS.RUnlock()

	if !ok {
		return resp.Value{Typ: "null"}
	}

	return resp.Value{
		Typ:  "bulk",
		Bulk: value,
	}

}

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

func GET(args []resp.Value) resp.Value {
	if len(args) != 1 {
		return resp.Value{Typ: "error", Str: "ERR wrong number of arguments for 'set' command"}
	}
	key := args[0].Bulk
	SETsMS.RLock()
	value, ok := SETs[key]
	SETsMS.RUnlock()
	if !ok {
		return resp.Value{Typ: "null"}
	}
	return resp.Value{Typ: "bulk", Bulk: value}
}

var Handlers = map[string]func([]resp.Value) resp.Value{
	"PING": ping,
	"SET":  SET,
	"GET":  GET,
}
