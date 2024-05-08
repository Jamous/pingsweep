package main
/*
Make maxMask in ignoreSubnet a passed variabel

*/

import (
	"fmt"
	"net"
)

func pingDriver() error{
	fmt.Println("Welcome to pingDriver")
	
	//Get list of ipv4 addresses
	//addressList, err := getInterface()
	//if err != nil {
	//	return fmt.Errorf("Could not get interfaces. %s", err)
	//}

	//fmt.Println(addressList)

	//Generate list of address to ping
	generateAddresses()
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
				if ! ignoreSubnet(ignoreSubnets, ipAddr) {
					addressList = append(addressList, addr)
				}			
			}
		}
	}

	//return addressList
	return addressList, nil
}

func ignoreSubnet(ignoreSubnets []*net.IPNet, ipAddr *net.IPNet) bool {
	//Set a max mask size. Default is /20, We dont want to wait all day to ping 2048 unnesicary addresses.
	maxMask := 21

	//Iterate through each subnet
	for _, subnet := range ignoreSubnets {
		//Check if subnet is in the existing subnet
		if subnet.Contains(ipAddr.IP) {
			return true
		}

		//Check if subnet is larger than max mask. Return true if it is
		mask, _ := ipAddr.Mask.Size()
		if mask < maxMask {
			return true
		}
	}

	return false
}

func generateAddresses() {
	// Define your subnet in CIDR notation
	subnet := "192.168.1.0/24"

	// Parse the CIDR notation to get the IP address and the subnet mask
	ip, ipnet, err := net.ParseCIDR(subnet)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	//networkAddr := ip.Mask(ipnet.Mask)

	for ipnet.Contains(ip) {
		fmt.Println(ip)
		for j := len(ip) - 1; j >= 0; j-- {
			ip[j]++
			if ip[j] > 0 {
				break
			}
		}
	}
}