package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	. "github.com/D7024E-Distributed-Systems/Kademlia/src/kademlia"
	"github.com/gorilla/mux"
)

type data struct {
	Value string
	TTL   string
}

func GetRoute(ip string, port int, kademlia *Kademlia) {
	address := fmt.Sprintf("%s:%d", ip, port)

	fmt.Println("ENTERING REST ON ADDRESS", address)
	r := mux.NewRouter()
	r.HandleFunc("/objects/{hash}", func(w http.ResponseWriter, r *http.Request) {
		str := mux.Vars(r)
		newKademliaID := ToKademliaID((str["hash"]))
		message := kademlia.LookupData(newKademliaID)
		fmt.Println("MESSAGE IS", message)
		if message == nil {
			fmt.Println("NO SUCH VALUE FOUND ON ID", newKademliaID.String())
			w.WriteHeader(400)
			w.Write([]byte("NO SUCH VALUE FOUND"))
			return
		}
		fmt.Println("FOUND", string(message))
		w.WriteHeader(200)
		w.Write([]byte("FOUND "))
		w.Write(message)
	}).Methods("GET")
	r.HandleFunc("/object", func(w http.ResponseWriter, r *http.Request) {
		var data data
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		ttl, err := time.ParseDuration(data.TTL)
		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte("FAILED TO PARSE TTL"))
			return
		}
		id, deadAt := kademlia.Store([]byte(data.Value), ttl)
		fmt.Println("STORED", id)
		w.WriteHeader(201)
		w.Write([]byte("Created data " + data.Value + " which will die at " + deadAt.String() + "\n"))
		w.Write([]byte("Location: " + address + "/objects/" + id.String()))

	}).Methods("POST")
	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(address, r))
}
