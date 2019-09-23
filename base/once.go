package main

import "fmt"
import "sync"

var a string
var once sync.once

func setup() {
	a = "Hello, world"
}

func doprint() {
	once.Do(setup)
	fmt.Print(a)
}

func main() {
	go doprint()
	go doprint()
}
