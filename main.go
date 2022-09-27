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
	Mut     *sync.RWMutex
	Data    map[string][]byte
	ErrChan chan error
}

type MyError struct {
	data json.RawMessage
}

type ThisError struct{}

func (m *MyError) Error() string {
	return "boom"
}

func (m *ThisError) Error() string {
	return "boom"
}

func errorMonitor(e chan error) {
	for {

		// consume an error, switch on the type
		switch err := <-e; err.(type) {
		case *MyError:

			// convert the type and consume it's data
			fmt.Println(string(err.(*MyError).data))
			fmt.Println("my error")

		case *ThisError:
			fmt.Println("this error")
		default:
			fmt.Println("undefined error")
		}
	}
}

func (db *Database) remove(w http.ResponseWriter, r *http.Request) {
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

		db.Mut.Lock()
		if obj, found := db.Data[newKey]; found {
			delete(db.Data, newKey)
			fmt.Printf("deleted: %+v\n", obj)
		}
		db.Mut.Unlock()
	}
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
			newErr := MyError{}
			newErr.data = message
			db.ErrChan <- &newErr
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

	ec := make(chan error)

	go errorMonitor(ec)

	db := Database{
		Data:    data,
		Mut:     &mut,
		ErrChan: ec,
	}

	http.HandleFunc("/insert", db.insert)
	http.HandleFunc("/remove", db.remove)
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
