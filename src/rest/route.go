package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	. "github.com/D7024E-Distributed-Systems/Kademlia/src/kademlia"
	"github.com/gorilla/mux"
)

type data struct {
	Value string
	TTL   string
}

func SetGetRoute(w http.ResponseWriter, r *http.Request, kademlia *Kademlia) {
	str := mux.Vars(r)
	newKademliaID := ToKademliaID((str["hash"]))
	message, _ := kademlia.GetValue(newKademliaID)
	fmt.Println("MESSAGE IS", message)
	if message == nil {
		fmt.Println("NO SUCH VALUE FOUND ON ID", newKademliaID.String())
		w.WriteHeader(400)
		w.Write([]byte("NO SUCH VALUE FOUND"))
		return
	}
	fmt.Println("FOUND", string(*message))
	w.WriteHeader(200)
	w.Write([]byte("FOUND "))
	w.Write([]byte(*message))
}

func SetPostRoute(w http.ResponseWriter, r *http.Request, kademlia *Kademlia, address string) {
	var data data
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("COULD NOT DECODE JSON"))
		return
	}

	ttl, err := time.ParseDuration(data.TTL)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("FAILED TO PARSE TTL"))
		return
	}
	contacts := kademlia.StoreValue([]byte(data.Value), ttl)
	id := NewKademliaID(data.Value)
	placement := " node"
	if len(contacts) > 1 {
		placement = " nodes"
	}
	fmt.Println("STORED", id.String())
	w.WriteHeader(201)
	w.Write([]byte("Created data " + data.Value + "\n"))
	w.Write([]byte("The data was stored on " + strconv.Itoa(len(contacts)) + placement + "\n"))
	w.Write([]byte("Location: " + address + "/objects/" + id.String()))
}

func GetRoute(ip string, port int, kademlia *Kademlia) {
	address := fmt.Sprintf("%s:%d", ip, port)

	fmt.Println("ENTERING REST ON ADDRESS", address)
	r := mux.NewRouter()
	r.HandleFunc("/objects/{hash}", func(w http.ResponseWriter, r *http.Request) {
		SetGetRoute(w, r, kademlia)
	}).Methods("GET")
	r.HandleFunc("/object", func(w http.ResponseWriter, r *http.Request) {
		SetPostRoute(w, r, kademlia, address)
	}).Methods("POST")
	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(address, r))
}
