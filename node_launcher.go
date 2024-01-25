package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

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

func main() {
	var initiator string
	isInNodes := false

	// Read the tree structure from a file
	tree, err := ReadTree()
	if err != nil {
		log.Fatalln(err)
	}

	nodes := []string{}
	for _, arr := range tree {
		nodes = append(nodes, arr[0])
	}

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

	// Prepare coomand-line arguments
	// and launch each node in the tree
	for _, arr := range tree {
		builder := strings.Builder{}
		builder.WriteString("run ./node.go")

		if arr[0] == initiator {
			builder.WriteString(" -i")
		}

		for _, s := range arr {
			builder.WriteString(fmt.Sprintf(" %s", s))
		}

		args := strings.Split(builder.String(), " ")
		cmd := exec.Command("go", args...)

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err := cmd.Run()
		if err != nil {
			fmt.Println("Error!", err)
			return
		}
	}
}
