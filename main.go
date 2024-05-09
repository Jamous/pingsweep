package main

import "log"

func main() {
	//Setup new PSconfig
	psconfig := newPSconfig()

	//Call pingDriver
	err := pingDriver(psconfig)
	if err != nil {
		log.Fatalf("%s\n", err)
	}
}