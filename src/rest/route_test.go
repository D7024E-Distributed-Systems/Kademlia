package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	. "github.com/D7024E-Distributed-Systems/Kademlia/src/kademlia"
	"github.com/gorilla/mux"
)

func TestPostAndGetRoute(t *testing.T) {
	nodeID := NewRandomKademliaID()
	contact := NewContact(nodeID, "localhost:8080")
	network := NewNetwork(&contact)
	kademlia := NewKademliaStruct(network)

	go GetRoute("localhost", 8080, kademlia)
	time.Sleep(1 * time.Second)

	postBody, _ := json.Marshal(map[string]string{
		"value": "file",
		"ttl":   "10m",
	})

	responseBody := bytes.NewBuffer(postBody)

	req := httptest.NewRequest(http.MethodPost, "/object", responseBody)
	w := httptest.NewRecorder()
	SetPostRoute(w, req, kademlia, "localhost:8080")
	res := w.Result()
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)

	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	fmt.Println(string(data))
	if string(data)[:17] != "Created data file" {
		t.Fail()
	}
	// id := "3937316334313964643630393333313334336465"
	// url1 := path.Join("/objects/", id)

	req2, _ := http.NewRequest(http.MethodGet, "/objects/3937316334313964643630393333313334336465", nil)
	w2 := httptest.NewRecorder()

	vars := map[string]string{
		"hash": "3937316334313964643630393333313334336465",
	}

	req2 = mux.SetURLVars(req2, vars)

	SetGetRoute(w2, req2, kademlia)
	res2 := w2.Result()
	defer res2.Body.Close()
	data2, err2 := ioutil.ReadAll(res2.Body)

	if err2 != nil {
		fmt.Println(err2)
		t.Fail()
	}
	if string(data2) != "FOUND file" {
		t.Fail()
	}

	return
}

func TestGetRouteFailure(t *testing.T) {
	nodeID := NewRandomKademliaID()
	contact := NewContact(nodeID, "localhost:8081")
	network := NewNetwork(&contact)
	kademlia := NewKademliaStruct(network)

	go GetRoute("localhost", 8081, kademlia)
	time.Sleep(1 * time.Second)

	req2, _ := http.NewRequest(http.MethodGet, "/objects/0000000000000000000000000000000000000000", nil)
	w2 := httptest.NewRecorder()

	vars := map[string]string{
		"hash": "0000000000000000000000000000000000000000",
	}

	req2 = mux.SetURLVars(req2, vars)

	SetGetRoute(w2, req2, kademlia)
	res2 := w2.Result()
	defer res2.Body.Close()
	data2, err2 := ioutil.ReadAll(res2.Body)

	if err2 != nil {
		fmt.Println(err2)
		t.Fail()
	}
	if string(data2) != "NO SUCH VALUE FOUND" {
		t.Fail()
	}

	return
}

func TestPostRouteTTLFailure(t *testing.T) {
	nodeID := NewRandomKademliaID()
	contact := NewContact(nodeID, "localhost:8082")
	network := NewNetwork(&contact)
	kademlia := NewKademliaStruct(network)

	go GetRoute("localhost", 8082, kademlia)
	time.Sleep(1 * time.Second)

	postBody, _ := json.Marshal(map[string]string{
		"value": "file",
		"ttl":   "gdijgdamdk",
	})

	responseBody := bytes.NewBuffer(postBody)

	req := httptest.NewRequest(http.MethodPost, "/object", responseBody)
	w := httptest.NewRecorder()
	SetPostRoute(w, req, kademlia, "localhost:8082")
	res := w.Result()
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)

	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	fmt.Println(string(data))
	if string(data) != "FAILED TO PARSE TTL" {
		t.Fail()
	}

	return
}

func TestPostRouteJSONFailure(t *testing.T) {
	nodeID := NewRandomKademliaID()
	contact := NewContact(nodeID, "localhost:8083")
	network := NewNetwork(&contact)
	kademlia := NewKademliaStruct(network)

	go GetRoute("localhost", 8083, kademlia)
	time.Sleep(1 * time.Second)

	// responseBody := bytes.NewBuffer(postBody)

	req := httptest.NewRequest(http.MethodPost, "/object", nil)
	w := httptest.NewRecorder()
	SetPostRoute(w, req, kademlia, "localhost:8083")
	res := w.Result()
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)

	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	fmt.Println(string(data))
	if string(data) != "COULD NOT DECODE JSON" {
		t.Fail()
	}

	return
}
