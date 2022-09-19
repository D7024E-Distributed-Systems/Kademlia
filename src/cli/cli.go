package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/D7024E-Distributed-Systems/Kademlia/src/kademlia"
)

func Init(shutdownNode func(), kademlia *kademlia.Kademlia) {
	go do(readLine, shutdownNode, kademlia)
}

func do(readInput func() string, shutdownNode func(), kademlia *kademlia.Kademlia) {
	for {
		input := readInput()
		if stringsEqual(input, "exit") {
			if shouldExit(readInput) {
				shutdownNode()
				return
			}
		} else if stringsEqual(input, "find contact") {
			findContact(readInput, kademlia.LookupContact)
		} else if stringsEqual(input, "help") {
			printHelp()
		} else if stringsEqual(input, "forget") {
			forgetHelp(kademlia, readInput, kademlia.RemoveFromKnown)
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
	fmt.Println("\tfind contact - finds the k closest contacts to a given node.")
}

func findContact(readInput func() string, lookupContact func(*kademlia.KademliaID) []kademlia.Contact) {
	str := readInput()
	id := kademlia.NewKademliaID(str)
	// if id == nil {
	// 	fmt.Println("Invalid kademlia ID")
	// 	return
	// }
	contact := lookupContact(id)
	fmt.Println("Found contact", contact, "from searching in CLI")
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

func stringsEqual(a, b string) bool {
	return strings.Compare(a, b) == 0
}

func readLine() string {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")
	scanner.Scan()
	return strings.Replace(scanner.Text(), "\r\n", "", -1)
}
