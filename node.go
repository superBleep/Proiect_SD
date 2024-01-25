package main

import (
	"flag"
	"fmt"
)

type data struct {
	id        string
	neighbors []string
	holder    string
	using     bool
	requestQ  []string
	asked     bool
}

func main() {
	init := flag.Bool("i", false, "marks the node as the initiator of the algorithm")
	flag.Parse()
	args := flag.Args()

	nodeData := data{id: args[0]}

	for i := 1; i < len(args); i++ {
		nodeData.neighbors = append(nodeData.neighbors, args[i])
	}

	fmt.Println(nodeData)
	fmt.Println(*init)
}
