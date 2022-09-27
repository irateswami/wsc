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

type Database struct {
	Mut  *sync.RWMutex
	Data map[string][]byte
}

func (db *Database) insert(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("connection upgrade error: ", err)
	}
	defer conn.Close()

	for {
		messageType, message, err := conn.ReadMessage()
		if messageType == websocket.CloseMessage {
			break
		}
		if err != nil {
			log.Println("message read error: ", err)
			break
		}

		var tempStruct map[string]any
		if err := json.Unmarshal(message, &tempStruct); err != nil {
			log.Println("message unmarshall error: ", err)
			break
		}

		//log.Printf("%+v\n", tempStruct)
		newKey, ok := tempStruct["id"].(string)
		if !ok {
			log.Println("id conversion failed")
			break
		}

		var compressedBytes bytes.Buffer
		w := gzip.NewWriter(&compressedBytes)
		w.Write(message)

		db.Mut.Lock()
		if _, found := db.Data[newKey]; !found {
			db.Data[newKey] = compressedBytes.Bytes()
			log.Printf("new id: %s, len: %d", newKey, len(db.Data[newKey]))
		}
		db.Mut.Unlock()
	}
}

func main() {
	fmt.Println("hi")

	var mut sync.RWMutex
	data := make(map[string][]byte)

	db := Database{
		Data: data,
		Mut:  &mut,
	}

	http.HandleFunc("/insert", db.insert)
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
