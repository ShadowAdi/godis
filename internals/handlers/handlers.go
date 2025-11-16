package handlers

import "godis/internals/resp"

func ping(args []resp.Value) resp.Value {
	return resp.Value{
		Typ: "string",
		Str: "PONG",
	}
}

var Handlers = map[string]func([]resp.Value) resp.Value{
	"PING": ping,
}
