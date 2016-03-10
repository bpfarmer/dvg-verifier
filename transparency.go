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
	fmt.Println(k[0].R)
	var m []*merkle.Node
	for n := 0; n < 2; n++ {
		m = append(m, &merkle.Node{})
	}
	k = merkle.AppendLeaves(l, m)
	walk(k[0])
	fmt.Println(k[0].R)
	fmt.Println(len(k[0].Leaves()))
	fmt.Println(len(k[0].Leaves()))
}

func walk(n *merkle.Node) {
	if n.IsLeaf() {
		fmt.Println("LEAF")
	}
	if n.L != nil {
		//fmt.Println("Walkin' Left")
		walk(n.L)
	}
	if n.R != nil {
		//fmt.Println("Walkin' Right")
		walk(n.R)
	}
}
