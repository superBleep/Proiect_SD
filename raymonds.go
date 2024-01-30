package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

type data struct {
	id        string
	neighbors []string
	holder    string
	using     bool
	requestQ  []string
	asked     bool
}

type message struct {
	srcId   string
	destId  string
	message string
}

// ReadTree reads the tree structure from a file
func ReadTree() ([][]string, error) {
	file, err := os.Open("tree.dat")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var tree [][]string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		row := strings.Fields(line)

		tree = append(tree, row)
	}

	return tree, nil
}

// SelectInitiator returns the node which initiates the algorithm
func SelectInitiator(nodes []string) (string, error) {
	var initiator string
	isInNodes := false

	// Read the id of the initiator node from stdinput
	fmt.Print("Available nodes: ")
	for _, node := range nodes {
		fmt.Print(node, " ")
	}
	fmt.Print("\nSelect initiator: ")

	fmt.Scanln(&initiator)

	// Check the validtiy of the input value
	for isInNodes == false {
		for _, node := range nodes {
			if node == initiator {
				isInNodes = true
				break
			}
		}
		if !isInNodes {
			fmt.Println("\nThe selected node is not valid.")
			fmt.Print("Available nodes: ")
			for _, node := range nodes {
				fmt.Print(node, " ")
			}
			fmt.Print("\nSelect initiator: ")

			fmt.Scanln(&initiator)
		}
	}

	return initiator, nil
}

func TreeNode(id string, neighbors []string, initiator bool, channel chan message, wg *sync.WaitGroup) {
	defer wg.Done()

	nodeData := data{id: id, neighbors: neighbors}

	if initiator {
		for _, neighbor := range nodeData.neighbors {
			channel <- message{nodeData.id, neighbor, "INITIALIZE"}
		}
	} else {
		msg := <-channel

		if msg.destId == nodeData.id && msg.message == "INITIALIZE" {
			nodeData.holder = msg.srcId

			for _, neighbor := range nodeData.neighbors {
				if neighbor != msg.srcId {
					channel <- message{nodeData.id, neighbor, "INITIALIZE"}
				}
			}

			fmt.Println(nodeData)
		}
	}
}

func main() {
	var wg sync.WaitGroup
	var nodes []string
	channel := make(chan message)

	// Read the tree structure from a file
	tree, err := ReadTree()
	if err != nil {
		log.Fatalln(err)
	}

	// Get the names of all the nodes
	for _, arr := range tree {
		nodes = append(nodes, arr[0])
	}

	// Read the initator node from stdinput
	initiator, err := SelectInitiator(nodes)
	if err != nil {
		log.Fatalln(err)
	}

	// Simulate nodes
	for i := 0; i < len(tree); i++ {
		wg.Add(1)

		// Start nodes and mark the initiator node
		if tree[i][0] == initiator {
			go TreeNode(tree[i][0], tree[i][1:], true, channel, &wg)
		} else {
			go TreeNode(tree[i][0], tree[i][1:], false, channel, &wg)
		}

		time.Sleep(time.Second) // Simulate delay in starting nodes
	}

	close(channel) // Close the communication channel
}
