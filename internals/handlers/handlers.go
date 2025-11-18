package handlers

import (
	"fmt"
	"godis/helper"
	"godis/internals/resp"
	"sync"
)

var SETs = map[string]string{}
var SETsMS = sync.RWMutex{}
var HSETs = map[string]map[string]string{}
var HSETsMS = sync.RWMutex{}

func hset(args []resp.Value) resp.Value {
	if len(args) != 3 {
		helper.LogError("HSET: wrong number of arguments")
		return resp.Value{Typ: "error", Str: "ERR wrong number of arguments for 'hset' command"}
	}

	hash := args[0].Bulk
	key := args[1].Bulk
	value := args[2].Bulk

	HSETsMS.Lock()
	if _, ok := HSETs[hash]; !ok {
		HSETs[hash] = map[string]string{}
		helper.LogInfo(fmt.Sprintf("Created new hash: %s", hash))
	}
	HSETs[hash][key] = value
	HSETsMS.Unlock()
	helper.LogInfo(fmt.Sprintf("HSET: hash=%s, key=%s", hash, key))
	return resp.Value{Typ: "string", Str: "OK"}
}

func hget(args []resp.Value) resp.Value {
	if len(args) != 2 {
		helper.LogError("HGET: wrong number of arguments")
		return resp.Value{Typ: "error", Str: "ERR wrong number of arguments for 'hget' command"}
	}

	hash := args[0].Bulk
	key := args[1].Bulk

	HSETsMS.RLock()
	value, ok := HSETs[hash][key]
	HSETsMS.RUnlock()

	if !ok {
		helper.LogInfo(fmt.Sprintf("HGET: key not found - hash=%s, key=%s", hash, key))
		return resp.Value{Typ: "null"}
	}

	helper.LogInfo(fmt.Sprintf("HGET: hash=%s, key=%s", hash, key))
	return resp.Value{
		Typ:  "bulk",
		Bulk: value,
	}

}

func hGetAll(args []resp.Value) resp.Value {
	if len(args) != 1 {
		helper.LogError("HGETALL: wrong number of arguments")
		return resp.Value{Typ: "error", Str: "ERR wrong number of arguments for 'hgetall' command"}
	}

	hash := args[0].Bulk

	HSETsMS.RLock()
	values, ok := HSETs[hash]
	HSETsMS.RUnlock()

	if !ok {
		helper.LogInfo(fmt.Sprintf("HGETALL: hash not found - hash=%s", hash))
		return resp.Value{Typ: "null"}
	}

	arr := []resp.Value{}
	for k, v := range values {
		arr = append(arr, resp.Value{Typ: "bulk", Bulk: k})
		arr = append(arr, resp.Value{Typ: "bulk", Bulk: v})
	}

	helper.LogInfo(fmt.Sprintf("HGETALL: hash=%s, fields=%d", hash, len(values)))
	return resp.Value{
		Typ:   "array",
		Array: arr,
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

	// Check if key already exists
	SETsMS.RLock()
	_, exists := SETs[key]
	SETsMS.RUnlock()

	if exists {
		return resp.Value{Typ: "error", Str: "ERR key already exists"}
	}

	// Insert new key
	SETsMS.Lock()
	SETs[key] = value
	SETsMS.Unlock()

	return resp.Value{Typ: "string", Str: "OK"}
}

func GET(args []resp.Value) resp.Value {
	if len(args) != 1 {
		helper.LogError("GET: wrong number of arguments")
		return resp.Value{Typ: "error", Str: "ERR wrong number of arguments for 'get' command"}
	}
	key := args[0].Bulk
	SETsMS.RLock()
	value, ok := SETs[key]
	SETsMS.RUnlock()
	if !ok {
		helper.LogInfo(fmt.Sprintf("GET: key not found - key=%s", key))
		return resp.Value{Typ: "null"}
	}
	helper.LogInfo(fmt.Sprintf("GET: key=%s", key))
	return resp.Value{Typ: "bulk", Bulk: value}
}

func IsExists(args []resp.Value) resp.Value {
	if len(args) != 1 {
		helper.LogError("EXISTS: wrong number of arguments")
		return resp.Value{Typ: "error", Str: "ERR wrong number of arguments for 'exists' command"}
	}
	key := args[0].Bulk
	SETsMS.RLock()
	_, ok := SETs[key]
	SETsMS.RUnlock()
	if !ok {
		helper.LogInfo(fmt.Sprintf("EXISTS: key not found - key=%s", key))
		return resp.Value{Typ: "integer", Num: 0}
	}
	helper.LogInfo(fmt.Sprintf("EXISTS: key=%s", key))
	return resp.Value{Typ: "integer", Num: 1}
}

func DELETE(args []resp.Value) bool {
	if len(args) != 1 {
		helper.LogError("EXISTS: wrong number of arguments")
		return false
	}
	key := args[0].Bulk
	SETsMS.RLock()
	_, ok := SETs[key]
	if !ok {
		helper.LogInfo(fmt.Sprintf("GET: key not found - key=%s", key))
		return false
	}
	helper.LogInfo(fmt.Sprintf("GET: key=%s", key))
	return true
}

var Handlers = map[string]func([]resp.Value) resp.Value{
	"PING":    ping,
	"SET":     SET,
	"GET":     GET,
	"HSET":    hset,
	"HGET":    hget,
	"HGETALL": hGetAll,
	"EXISTS":  IsExists,
}
