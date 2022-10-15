package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/D7024E-Distributed-Systems/Kademlia/src/kademlia"
	. "github.com/D7024E-Distributed-Systems/Kademlia/src/kademlia"
)

func Init(shutdownNode func(), kademlia *Kademlia) {
	do(readLine, shutdownNode, kademlia)
}

func do(readInput func() string, shutdownNode func(), kademlia *Kademlia) {
	for {
		input := readInput()
		if stringsEqual(input, "exit") {
			if shouldExit(readInput) {
				shutdownNode()
				return
			}
		} else if stringsEqual(input, "find contact") {
			findContact(readInput, kademlia.LookupContact)
		} else if stringsEqual(input, "put") {
			storeValue(readInput, kademlia.StoreValue)
		} else if stringsEqual(input, "help") {
			printHelp()
		} else if stringsEqual(input, "forget") {
			forgetHelp(kademlia, readInput, kademlia.RemoveFromKnown)
		} else if stringsEqual(input, "get") {
			getValue(readInput, kademlia)
		} else if stringsEqual(input, "table") {
			fmt.Println("table length is", len(kademlia.Network.RoutingTable.FindClosestContacts(NewRandomKademliaID(), 9999)))
		} else {
			fmt.Println("Unknown command \"" + input + "\"")
		}

	}
}

func shouldExit(readInput func() string) bool {
	fmt.Println("Are you sure you want to exit? Y/n")
	text := readInput()
	if stringsEqual(text, "Y") || stringsEqual(text, "y") ||
		stringsEqual(text, "Yes") || stringsEqual(text, "yes") {
		return true
	}
	return false
}

func printHelp() {
	fmt.Println("Available commands:")
	fmt.Println("\texit - shuts down the current node and all data will be lost.")
	fmt.Println("\tfind contact - finds the", BucketSize, "closest contacts to a given node.")
	fmt.Println("\tget - gets the value of a specific hash given.")
	fmt.Println("\tput - stores a string value.")
	fmt.Println("\tforget - forgets a value of a specific hash given.")
}

func findContact(readInput func() string, lookupContact func(*kademlia.KademliaID) ContactCandidates) {
	fmt.Println("Enter a node ID to look for")
	str := readInput()
	id := kademlia.ToKademliaID(str)
	if id == nil {
		fmt.Println("Invalid ID given")
		return
	}
	t1 := time.Now()
	contact := lookupContact(id)
	t2 := time.Now()
	fmt.Println("Found", contact.Len(), ": Closest contact", contact.GetContacts(1)[0], "-", t2.Sub(t1).Milliseconds(), "ms")
}

func storeValue(readInput func() string, StoreValue func([]byte, time.Duration) ([]*kademlia.KademliaID, string)) {
	fmt.Println("What would you like to store?")
	data := readInput()
	fmt.Println("For how long do you want to store it? (\"10s\", or \"5h30m\")")
	inp := readInput()
	ttl, err := time.ParseDuration(inp)
	if err != nil {
		fmt.Println("Invalid time, try again")
		storeValue(readInput, StoreValue)
		return
	}
	t1 := time.Now()
	storedIDs, hash := StoreValue([]byte(data), ttl)
	t2 := time.Now()
	fmt.Println("Hash of", data, "is", hash, "-", t2.Sub(t1).Milliseconds(), "ms")
	fmt.Println("Stored in nodes: ", storedIDs)
}

func forgetHelp(kademlia *kademlia.Kademlia, readInput func() string, removeFromKnown func(value string) bool) {
	fmt.Println("Which value do you want to forget?")
	text := readInput()
	success := kademlia.RemoveFromKnown(text)
	if success {
		fmt.Println("Successfully forgot value: ", text)
	} else {
		fmt.Println("Failed to forget value: ", text)
	}
}

func getValue(readInput func() string, kademlia *kademlia.Kademlia) {
	fmt.Println("Which hash do you want to get?")
	text := readInput()
	id := ToKademliaID(text)
	t1 := time.Now()
	if id == nil {
		fmt.Println("Invalid ID given try \"get\" again")
		return
	}
	value, contact := kademlia.GetValue(id)
	t2 := time.Now()
	if value == nil {
		fmt.Println("Error, value not found", "-", t2.Sub(t1).Milliseconds(), "ms")
		return
	}
	fmt.Println("\""+*value+"\" found at node", contact.ID, "-", t2.Sub(t1).Milliseconds(), "ms")

}

func stringsEqual(a, b string) bool {
	return strings.Compare(a, b) == 0
}

func readLine() string {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")
	scanner.Scan()
	return strings.Replace(scanner.Text(), "\r\n", "", -1)
}
