package main

import "log"

func main() {
	err := pingDriver()
	if err != nil {
		log.Fatalf("%s\n", err)
	}
}