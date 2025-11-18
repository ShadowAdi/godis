package handlers

import (
	"fmt"
	"godis/helper"
	"godis/internals/resp"
	"strconv"
	"sync"
	"time"
)

var SETs = map[string]string{}
var SETsMS = sync.RWMutex{}
var HSETs = map[string]map[string]string{}
var HSETsMS = sync.RWMutex{}
var TTLs = map[string]int64{}
var TTLsMS = sync.RWMutex{}

func HSET(args []resp.Value) resp.Value {
	if len(args) != 3 {
		return resp.Value{Typ: "error", Str: "ERR wrong number of arguments for 'hset' command"}
	}

	hash := args[0].Bulk
	key := args[1].Bulk
	value := args[2].Bulk

	HSETsMS.Lock()
	defer HSETsMS.Unlock()

	// If hash doesn't exist, create it
	if _, ok := HSETs[hash]; !ok {
		HSETs[hash] = map[string]string{}
	}

	// Check if field already exists
	if _, exists := HSETs[hash][key]; exists {
		return resp.Value{Typ: "error", Str: "ERR field already exists"}
	}

	// Insert new field
	HSETs[hash][key] = value

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

func DELETE(args []resp.Value) resp.Value {
	if len(args) != 1 {
		return resp.Value{Typ: "error", Str: "ERR wrong number of arguments for 'del' command"}
	}

	key := args[0].Bulk

	SETsMS.Lock()
	defer SETsMS.Unlock()

	_, ok := SETs[key]
	if !ok {
		return resp.Value{Typ: "integer", Num: 0} // Redis returns 0 when key does not exist
	}

	delete(SETs, key)

	return resp.Value{Typ: "integer", Num: 1} // Redis returns 1 when deletion succeeds
}

func WILDCARD(args []resp.Value) resp.Value {
	if len(args) != 0 {
		return resp.Value{
			Typ: "error",
			Str: "ERR rong number of arguments for WILDCARD command",
		}
	}

	var all []resp.Value

	SETsMS.RLock()
	for key, value := range SETs {
		entry := fmt.Sprintf("SET %s %s", key, value)
		all = append(all, resp.Value{Typ: "string", Bulk: entry})
	}
	SETsMS.Unlock()

	HSETsMS.Lock()
	for hash, fields := range HSETs {
		for field, val := range fields {
			entry := fmt.Sprintf("HSET %s %s %s", hash, field, val)
			all = append(all, resp.Value{Typ: "string", Bulk: entry})
		}
	}
	HSETsMS.Unlock()

	return resp.Value{
		Typ:   "array",
		Array: all,
	}

}

func EXPIRE(args []resp.Value) resp.Value {
	if len(args) != 2 {
		return resp.Value{
			Typ: "error",
			Str: "ERR rong number of arguments for EXPIRE command",
		}
	}
	key := args[0].Bulk
	seconds, _ := strconv.Atoi(args[1].Bulk)

	SETsMS.Lock()
	_, ok := SETs[key]
	SETsMS.Unlock()

	if !ok {
		return resp.Value{Typ: "integer", Num: 0}
	}

	TTLsMS.Lock()
	TTLs[key] = time.Now().Unix() + int64(seconds)
	TTLsMS.Unlock()
	return resp.Value{Typ: "integer", Num: 1}

}

var Handlers = map[string]func([]resp.Value) resp.Value{
	"PING":     ping,
	"SET":      SET,
	"GET":      GET,
	"HSET":     HSET,
	"HGET":     hget,
	"HGETALL":  hGetAll,
	"EXISTS":   IsExists,
	"DELETE":   DELETE,
	"WILDCARD": WILDCARD,
}
