package main

import (
	"log"
	"fmt"
)

func main() {
	//Setup new default PSconfig
	//psconfig := NewPSconfig()

	//Build PSconfig
	psconfig := PSconfig{UseDefaultNetwork: true, MaxSubnetSize: 21}

	//Call pingDriver
	subnetList, err := PingDriver(psconfig)
	if err != nil {
		log.Fatalf("%s\n", err)
	}

	fmt.Println(subnetList)
}