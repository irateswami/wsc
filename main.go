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
	Data map[string][]byte
}

func (c *Cache) findOne(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("connection upgrade error: ", err)
	}
	defer conn.Close()

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("message read error: ", err)
			break
		}

		_ = messageType // switch on message type

		var tempStruct map[string]any
		if err := json.Unmarshal(message, &tempStruct); err != nil {
			log.Println("message unmarshall error: ", err)
			break
		}

		newKey, ok := tempStruct["id"].(string)
		if !ok {
			log.Println("id doesn't behave like a string type")
			break
		}
		if len(newKey) == 0 {
			log.Println("key is empty")
			break
		}

		messageBuff := bytes.NewBuffer(message)
		w := gzip.NewWriter(messageBuff)

		var newBuffer []byte
		w.Write(newBuffer)

		c.Mut.Lock()
		c.Data[newKey] = newBuffer
		c.Mut.Unlock()
	}
}

func (c *Cache) returnAll(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("connection upgrade error: ", err)
	}
	defer conn.Close()

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("message read error: ", err)
			break
		}

		_ = messageType // switch on message type

		var tempStruct map[string]any
		if err := json.Unmarshal(message, &tempStruct); err != nil {
			log.Println("message unmarshall error: ", err)
			break
		}

		newKey := tempStruct["id"].(string)
		if len(newKey) == 0 {
			log.Println("key is empty")
			continue
		}

		messageBuff := bytes.NewBuffer(message)
		w := gzip.NewWriter(messageBuff)

		var newBuffer []byte
		w.Write(newBuffer)

		c.Mut.Lock()
		c.Data[newKey] = newBuffer
		c.Mut.Unlock()
	}
}

func (c *Cache) update(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("connection upgrade error: ", err)
	}
	defer conn.Close()

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("message read error: ", err)
			break
		}

		_ = messageType // switch on message type

		var tempStruct map[string]any
		if err := json.Unmarshal(message, &tempStruct); err != nil {
			log.Println("message unmarshall error: ", err)
			break
		}

		newKey := tempStruct["id"].(string)
		if len(newKey) == 0 {
			log.Println("key is empty")
			continue
		}

		messageBuff := bytes.NewBuffer(message)
		w := gzip.NewWriter(messageBuff)

		var newBuffer []byte
		w.Write(newBuffer)

		c.Mut.Lock()
		c.Data[newKey] = newBuffer
		c.Mut.Unlock()
	}
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
			break
		}

		_ = messageType // switch on message type

		var tempStruct map[string]any
		if err := json.Unmarshal(message, &tempStruct); err != nil {
			log.Println("message unmarshall error: ", err)
			break
		}

		fmt.Printf("%+v\n", tempStruct)

		newKey := tempStruct["id"].(string)
		if len(newKey) == 0 {
			log.Println("key is empty")
			continue
		}

		messageBuff := bytes.NewBuffer(message)
		w := gzip.NewWriter(messageBuff)

		var newBuffer []byte
		w.Write(newBuffer)

		c.Mut.Lock()
		if _, found := c.Data[newKey]; !found {
			c.Data[newKey] = newBuffer
		}
		c.Mut.Unlock()
	}
}

func main() {
	fmt.Println("hi")

	var mut sync.RWMutex
	data := make(map[string][]byte)

	cache := Cache{
		Data: data,
		Mut:  &mut,
	}

	http.HandleFunc("/insert", cache.insert)
	http.HandleFunc("/update", cache.update)
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
