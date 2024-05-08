package main

import (
	"fmt"
	"net"
)

func pingDriver() error{
	fmt.Println("Welcome to pingDriver")
	
	//Get list of ipv4 addresses
	addressList, err := getInterface()
	if err != nil {
		return fmt.Errorf("Could not get interfaces. %s", err)
	}

	fmt.Println(addressList)

	/*
	//Get the current address
	host, _ := os.Hostname()
	fmt.Println(net.LookupIP(host))
	fmt.Println(host)
	//myip := net.ParseIP("216.252.192.167")
	
	fmt.Println(net.IPv4allrouter.DefaultMask())

	fmt.Println(net.InterfaceAddrs())
	*/

	return nil
}

func getInterface() ([]net.Addr, error) {
	//Get interfaces
	interfaces, err := net.InterfaceAddrs()
	if err != nil {
		return nil, fmt.Errorf("Could not get interfaces. %s", err)
	}
	
	//setup addressList and ignoreSubnets
	addressList := []net.Addr{}
	ignoreSubnets := []*net.IPNet{}
	ignoreList := []string{"169.254.0.0/16", "127.0.0.0/8"}

	//Generate ignoreSubnets from ignoreList
	for _, ignore := range ignoreList {
		_, ing, err := net.ParseCIDR(ignore)
		if err != nil {
			//Skip if invalid
			continue
		}
		ignoreSubnets = append(ignoreSubnets, ing)
	}

	//Lets look only at v4, there is no need to pingsweep IPv6 thatnks to neighbor solicitations.
	for _, addr := range interfaces {
		//Convert net.Addr to *net.IPNet to check if v4 or not. This returns ok (bool), only look at addresses that return ok.
		ipAddr, ok := addr.(*net.IPNet)
		if ok {
			//If IPv4 add to addressList
			if ipAddr.IP.To4() != nil {
				//Ignore networks in unwanted subnets
				if ! inSubnet(ignoreSubnets, ipAddr) {
					addressList = append(addressList, addr)
				}			
			}
		}
	}

	//return addressList
	return addressList, nil
}

func inSubnet(ignoreSubnets []*net.IPNet, ipAddr *net.IPNet) bool {
	//Iterate through each subnet
	for _, subnet := range ignoreSubnets {
		if subnet.Contains(ipAddr.IP) {
			return true
		}
	}

	return false
}