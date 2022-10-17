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

/*
Defines what to do when we receive a get to retrieve the value in the kademlia network
  - http http.ResponseWriter, to write back the value retrieved from the network
  - http *http.Request, to read the hash value from the request
  - kademlia *Kademlia, to search for the value in the kademlia network
*/
func SetGetRoute(w http.ResponseWriter, r *http.Request, kademlia *Kademlia) {
	str := mux.Vars(r)
	newKademliaID := ToKademliaID((str["hash"]))
	message, _ := kademlia.GetValue(newKademliaID)
	if message == nil {
		w.WriteHeader(400)
		w.Write([]byte("NO SUCH VALUE FOUND"))
		return
	}
	w.WriteHeader(200)
	w.Write([]byte("FOUND "))
	w.Write([]byte(*message))
}

/*
Defines what to do when we receive a post request to store data in the kademlia network
  - http http.ResponseWriter, to write back the hash of the data
  - http *http.Request, to get the data from the request
  - kademlia, *Kademlia to store the value sent to this node
  - address string, to be able to send back our address and port number
*/
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
	contacts, _ := kademlia.StoreValue([]byte(data.Value), ttl)
	id := NewKademliaID(data.Value)
	placement := " node"
	if len(contacts) > 1 {
		placement = " nodes"
	}
	w.WriteHeader(201)
	w.Write([]byte("Created data " + data.Value + "\n"))
	w.Write([]byte("The data was stored on " + strconv.Itoa(len(contacts)) + placement + "\n"))
	w.Write([]byte("Location: " + address + "/objects/" + id.String()))
}

func GetRoute(ip string, port int, kademlia *Kademlia) {
	address := fmt.Sprintf("%s:%d", ip, port)

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
