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

// func YourHandler(w http.ResponseWriter, r *http.Request, kademlia *Kademlia) {
// 	res := mux.Vars(r)
// 	restFul.kademlia.LookupData(res["hash"])
// 	fmt.Println(res["hash"])
// }

type data struct {
	Value string
	ttl   time.Duration
}

func GetRoute(ip string, port int, kademlia *Kademlia) {
	address := fmt.Sprintf("%s:%d", ip, port)

	fmt.Println("ENTERING REST ON ADDRESS", address)
	r := mux.NewRouter()
	r.HandleFunc("/objects/{hash}", func(w http.ResponseWriter, r *http.Request) {
		str := mux.Vars(r)
		newKademliaID := ToKademliaID((str["hash"]))
		message := kademlia.LookupData(newKademliaID)
		if message == nil {
			w.WriteHeader(400)
			w.Write([]byte("NO SUCH VALUE FOUND"))

		} else {
			fmt.Println("FOUND", string(message))
			w.WriteHeader(200)
			w.Write([]byte("FOUND "))
			w.Write(message)
		}
	}).Methods("GET")
	r.HandleFunc("/object", func(w http.ResponseWriter, r *http.Request) {
		var data data
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		id := kademlia.Store([]byte(data.Value), data.ttl)
		fmt.Println("STORED", id)
		w.WriteHeader(201)
		w.Write([]byte("CREATED "))
		w.Write([]byte(id.String()))

	}).Methods("POST")
	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(address, r))
}
