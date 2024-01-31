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
	state     string
}

type message struct {
	srcId   string
	destId  string
	message string
}

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

func SelectInitiator(nodes []string) (string, error) {
	var initiator string
	fmt.Print("Available nodes: ")
	for _, node := range nodes {
		fmt.Print(node, " ")
	}
	fmt.Print("\nSelect initiator: ")

	_, err := fmt.Scanln(&initiator)
	if err != nil {
		return "", err
	}

	for !contains(nodes, initiator) {
		fmt.Println("\nThe selected node is not valid.")
		fmt.Print("Available nodes: ")
		for _, node := range nodes {
			fmt.Print(node, " ")
		}
		fmt.Print("\nSelect initiator: ")
		_, err := fmt.Scanln(&initiator)
		if err != nil {
			return "", err
		}
	}

	return initiator, nil
}

func contains(nodes []string, node string) bool {
	for _, n := range nodes {
		if n == node {
			return true
		}
	}
	return false
}

func processMessage(nodeData *data, msg message, channels map[string]chan message) {
	switch nodeData.state {
	case "Initializing":
		if msg.destId == nodeData.id && msg.message == "INITIALIZE" {
			nodeData.holder = msg.srcId
			for _, neighbor := range nodeData.neighbors {
				if neighbor != msg.srcId {
					channels[neighbor] <- message{nodeData.id, neighbor, "INITIALIZE"}
					time.Sleep(100 * time.Millisecond) // Simulate delay
				}
			}
			fmt.Printf("%s received INITIALIZE message from %s.\n", nodeData.id, msg.srcId)
			if nodeData.id != nodeData.holder {
				nodeData.asked = true
				channels[nodeData.holder] <- message{nodeData.id, nodeData.holder, "REQUEST"}
				fmt.Printf("%s requested the token from %s.\n", nodeData.id, nodeData.holder)
				time.Sleep(100 * time.Millisecond) // Simulate delay
			}
			nodeData.state = "Waiting"
		}

	case "Waiting":
		if msg.destId == nodeData.id && msg.message == "TOKEN" {
			// enterCriticalSection(nodeData, channels)
		} else if msg.destId == nodeData.id && msg.message == "REQUEST" {
			nodeData.requestQ = append(nodeData.requestQ, msg.srcId)
			attemptToSendToken(nodeData, channels)
		}

	case "InCriticalSection":
		// Node is in the critical section or has terminated.
	}
}

func enterCriticalSection(nodeData *data, channels map[string]chan message) {
	nodeData.holder = nodeData.id
	nodeData.using = true
	nodeData.asked = false
	fmt.Printf("%s is in the critical section.\n", nodeData.id)
	time.Sleep(20 * time.Second) // Simulate critical section work
	exitCriticalSection(nodeData, channels)
}

func exitCriticalSection(nodeData *data, channels map[string]chan message) {
	nodeData.using = false
	attemptToSendToken(nodeData, channels)
	nodeData.state = "InCriticalSection" // Update state after critical section work is done
}

func attemptToSendToken(nodeData *data, channels map[string]chan message) {
	fmt.Printf("%s is attempting to send TOKEN.\n", nodeData.id)
	if !nodeData.using && len(nodeData.requestQ) > 0 {
		next := nodeData.requestQ[0]
		nodeData.requestQ = nodeData.requestQ[1:]
		channels[next] <- message{nodeData.id, next, "TOKEN"}
		fmt.Printf("%s sent TOKEN to %s.\n", nodeData.id, next)
	}
}

func TreeNode(id string, neighbors []string, initiator bool, channels map[string]chan message, wg *sync.WaitGroup) {
	defer wg.Done()

	nodeData := data{id: id, neighbors: neighbors, state: "Initializing"}
	myChannel := channels[id]

	if initiator {
		nodeData.holder = id // Initiator holds the token initially
		for _, neighbor := range nodeData.neighbors {
			channels[neighbor] <- message{nodeData.id, neighbor, "INITIALIZE"}
			time.Sleep(100 * time.Millisecond) // Simulate delay
		}
		nodeData.state = "Waiting"
	}

	timer := time.NewTimer(20 * time.Second) // Start a timer for 20 seconds

	for {
		select {
		case msg := <-myChannel:
			processMessage(&nodeData, msg, channels)
		case <-timer.C: // Wait for the timer to expire
			fmt.Println("Exiting loop")
			return // Exit the loop
		}
	}
}

func main() {
	var wg sync.WaitGroup

	tree, err := ReadTree()
	if err != nil {
		log.Fatalln(err)
	}

	channels := make(map[string]chan message)
	for _, arr := range tree {
		channels[arr[0]] = make(chan message, len(arr)-1)
	}

	var nodes []string
	for _, arr := range tree {
		nodes = append(nodes, arr[0])
	}

	initiator, err := SelectInitiator(nodes)
	if err != nil {
		log.Fatalln(err)
	}

	for i := 0; i < len(tree); i++ {
		wg.Add(1)
		go TreeNode(tree[i][0], tree[i][1:], tree[i][0] == initiator, channels, &wg)
		time.Sleep(time.Second)
	}

	wg.Wait()

	for _, ch := range channels {
		close(ch)
	}
}
