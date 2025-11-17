package resp

import (
	"bufio"
	"fmt"
	"godis/helper"
	"io"
	"strconv"
)

const (
	STRING  = '+'
	ERROR   = '-'
	INTEGER = ':'
	BULK    = '$'
	ARRAY   = '*'
)

type Value struct {
	Typ   string
	Str   string
	Num   int
	Bulk  string
	Array []Value
}
type Resp struct {
	reader *bufio.Reader
}

type Writer struct {
	writer io.Writer
}

func NewWriter(wt io.Writer) *Writer {
	return &Writer{
		writer: wt,
	}
}

func (w *Writer) Write(v Value) error {
	var bytes = v.Marshal()
	_, err := w.writer.Write(bytes)
	if err != nil {
		helper.LogError(fmt.Sprintf("Failed to write response: %v", err))
		return err
	}
	return nil
}

func NewResp(rd io.Reader) *Resp {
	return &Resp{reader: bufio.NewReader(rd)}
}

func (r *Resp) readLine() (line []byte, n int, err error) {
	for {
		b, err := r.reader.ReadByte()
		if err != nil {
			if err != io.EOF {
				helper.LogError(fmt.Sprintf("Error reading byte: %v", err))
			}
			return nil, 0, err
		}
		n += 1
		line = append(line, b)
		if len(line) >= 2 && line[len(line)-2] == '\r' {
			break
		}
	}
	return line[:len(line)-2], n, nil
}

func (r *Resp) readInteger() (x int, n int, err error) {
	line, n, err := r.readLine()
	if err != nil {
		return 0, 0, err
	}
	i64, err := strconv.ParseInt(string(line), 10, 64)
	if err != nil {
		helper.LogError(fmt.Sprintf("Failed to parse integer: %v", err))
		return 0, n, err
	}
	return int(i64), n, nil
}

func (r *Resp) Read() (Value, error) {
	_type, err := r.reader.ReadByte()

	if err != nil {
		if err != io.EOF {
			helper.LogError(fmt.Sprintf("Error reading type byte: %v", err))
		}
		return Value{}, err
	}

	switch _type {
	case ARRAY:
		return r.readArray()
	case BULK:
		return r.readBulk()
	default:
		helper.LogError(fmt.Sprintf("Unknown RESP type: %v", string(_type)))
		return Value{}, fmt.Errorf("unknown RESP type: %v", string(_type))
	}
}

func (r *Resp) readArray() (Value, error) {
	v := Value{}
	v.Typ = "array"

	// read length of array
	len, _, err := r.readInteger()
	if err != nil {
		helper.LogError(fmt.Sprintf("Failed to read array length: %v", err))
		return v, err
	}

	// foreach line, parse and read the value
	v.Array = make([]Value, 0)
	for i := 0; i < len; i++ {
		val, err := r.Read()
		if err != nil {
			helper.LogError(fmt.Sprintf("Failed to read array element %d: %v", i, err))
			return v, err
		}

		// append parsed value to array
		v.Array = append(v.Array, val)
	}

	return v, nil
}

func (r *Resp) readBulk() (Value, error) {
	v := Value{}

	v.Typ = "bulk"

	len, _, err := r.readInteger()
	if err != nil {
		helper.LogError(fmt.Sprintf("Failed to read bulk string length: %v", err))
		return v, err
	}

	bulk := make([]byte, len)

	_, err = r.reader.Read(bulk)
	if err != nil {
		helper.LogError(fmt.Sprintf("Failed to read bulk string data: %v", err))
		return v, err
	}

	v.Bulk = string(bulk)

	// Read the trailing CRLF
	_, _, err = r.readLine()
	if err != nil {
		helper.LogError(fmt.Sprintf("Failed to read trailing CRLF: %v", err))
		return v, err
	}

	return v, nil
}

func (v Value) marshalString() []byte {
	var bytes []byte
	bytes = append(bytes, STRING)
	bytes = append(bytes, []byte(v.Str)...)
	bytes = append(bytes, '\r', '\n')
	return bytes
}

func (v Value) marshalBulk() []byte {
	var bytes []byte
	bytes = append(bytes, BULK)
	bytes = append(bytes, []byte(strconv.Itoa(len(v.Bulk)))...)
	bytes = append(bytes, '\r', '\n')
	bytes = append(bytes, []byte(v.Bulk)...)
	bytes = append(bytes, '\r', '\n')
	return bytes
}

func (v Value) marshalArray() []byte {
	len := len(v.Array)
	var bytes []byte
	bytes = append(bytes, ARRAY)
	bytes = append(bytes, strconv.Itoa(len)...)
	bytes = append(bytes, '\r', '\n')

	for i := 0; i < len; i++ {
		bytes = append(bytes, v.Array[i].Marshal()...)
	}

	return bytes
}

func (v Value) marshallError() []byte {
	var bytes []byte
	bytes = append(bytes, ERROR)
	bytes = append(bytes, []byte(v.Str)...)
	bytes = append(bytes, '\r', '\n')

	return bytes
}

func (v Value) marshallNull() []byte {
	return []byte("$-1\r\n")
}

func (v Value) Marshal() []byte {
	switch v.Typ {
	case "array":
		return v.marshalArray()
	case "bulk":
		return v.marshalBulk()
	case "string":
		return v.marshalString()
	case "null":
		return v.marshallNull()
	case "error":
		return v.marshallError()
	default:
		return []byte{}
	}
}
