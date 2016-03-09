package main

import (
	"fmt"
	"transparency/merkle"

	_ "github.com/lib/pq"
)

func main() {
	var l []*merkle.Node
	for n := 0; n < 16; n++ {
		l = append(l, &merkle.Node{})
	}
	k := merkle.BuildTree(l)
	walk(k[0])
}

func walk(n *merkle.Node) {
	if n.L != nil {
		fmt.Println("Walkin' Left")
		walk(n.L)
	} else if n.R != nil {
		fmt.Println("Walkin' Right")
		walk(n.R)
	}
}
