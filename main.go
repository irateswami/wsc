package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/irateswami/wsc/websocket"
)

var (
	upgrader = websocket.Upgrader{}
)

type Cache struct {
	Mut  *sync.RWMutex
	Data map[string]*[]byte
}

func (c *Cache) insert(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("connection upgrade error: ", err)
	}
	defer conn.Close()

	for {

		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("message read error: ", err)
			continue
		}

		_ = messageType // switch on message type

		var tempStruct map[string]any
		if err := json.Unmarshal(message, &tempStruct); err != nil {
			log.Println("message unmarshall error: ", err)
			continue
		}

		newKey := tempStruct["key"].(string)
		if len(newKey) == 0 {
			log.Println("key is empty")
			continue
		}

		messageBuff := bytes.NewBuffer(message)
		w := gzip.NewWriter(messageBuff)

		var newBuffer []byte
		w.Write(newBuffer)

		c.Mut.Lock()
		c.Data[newKey] = &newBuffer
		c.Mut.Unlock()
	}
}

func main() {
	fmt.Println("hi")
}
