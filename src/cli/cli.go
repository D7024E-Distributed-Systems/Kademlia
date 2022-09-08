package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func Init() {
	go do()
}

func do() {
	reader := bufio.NewReader(os.Stdin)
	for {
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\r\n", "", -1)
		if strings.Compare("exit", text) == 0 {
			fmt.Println("Are you sure you want to exit? Y/n")
			text, _ := reader.ReadString('\n')
			text = strings.Replace(text, "\r\n", "", -1)
			if strings.Compare("Y", text) == 0 || strings.Compare("y", text) == 0 ||
				strings.Compare("Yes", text) == 0 || strings.Compare("yes", text) == 0 {
				os.Exit(0)
			}
		} else if strings.Compare("help", text) == 0 || strings.Compare("?", text) == 0 {
			fmt.Println("Available commands:")
			fmt.Println("\texit - shuts down the current node and all data will be lost")
		}
	}
}
