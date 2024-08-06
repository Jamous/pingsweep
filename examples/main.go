package main

import (
	"fmt"
	"log"

	"github.com/Jamous/pingsweep"
)

func main() {
	//Setup new default PSconfig
	psconfig := pingsweep.NewPSconfig()

	//Build PSconfig
	//psconfig := pingsweep.PSconfig{UseDefaultNetwork: true, MaxSubnetSize: 21, CustomSubnet: "10.0.0.0/24"}

	//Call pingDriver
	subnetList, err := pingsweep.PingDriver(psconfig)
	if err != nil {
		log.Fatalf("%s\n", err)
	}

	fmt.Println(subnetList)
}
