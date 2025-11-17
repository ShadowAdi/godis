package aof

import (
	"bufio"
	"fmt"
	"godis/helper"
	"godis/internals/resp"
	"io"
	"os"
	"sync"
	"time"
)

type Aof struct {
	file *os.File
	rd   *bufio.Reader
	mu   *sync.Mutex
}

func NewAoF(path string) (*Aof, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		helper.LogError(fmt.Sprintf("Failed to open AOF file: %v", err))
		return nil, err
	}
	helper.LogInfo(fmt.Sprintf("AOF file opened: %s", path))

	aof := &Aof{
		file: f,
		rd:   bufio.NewReader(f),
		mu:   &sync.Mutex{},
	}

	go func() {
		for {
			aof.mu.Lock()
			err := aof.file.Sync()
			aof.mu.Unlock()
			if err != nil {
				helper.LogError(fmt.Sprintf("AOF sync failed: %v", err))
			}
			time.Sleep(time.Second)
		}
	}()
	return aof, nil
}

func (aof *Aof) Close() error {
	aof.mu.Lock()
	defer aof.mu.Unlock()
	err := aof.file.Close()
	if err != nil {
		helper.LogError(fmt.Sprintf("Failed to close AOF file: %v", err))
		return err
	}
	helper.LogInfo("AOF file closed successfully")
	return nil
}

func (aof *Aof) Write(value resp.Value) error {
	aof.mu.Lock()
	defer aof.mu.Unlock()
	bytes := value.Marshal()
	_, err := aof.file.Write(bytes)
	if err != nil {
		helper.LogError(fmt.Sprintf("Failed to write to AOF: %v", err))
		return err
	}

	return nil
}

func (aof *Aof) Read(callback func(value resp.Value)) error {
	aof.mu.Lock()
	defer aof.mu.Unlock()
	resp := resp.NewResp(aof.file)
	commandCount := 0
	for {
		value, err := resp.Read()
		if err == nil {
			callback(value)
			commandCount++
		}

		if err == io.EOF {
			helper.LogInfo(fmt.Sprintf("AOF read complete: %d commands processed", commandCount))
			break
		}

		if err != nil {
			helper.LogError(fmt.Sprintf("Error reading AOF: %v", err))
			return err
		}
	}
	return nil
}
