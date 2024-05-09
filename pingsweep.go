package main
/*
Make maxMask in ignoreSubnet a passed variabel

*/

import (
	"fmt"
	"net"
)

//Handle package config.
type PSconfig struct {
	ListAddresses     bool //Prompt users manually to use or ignore a network
	UseDefaultNetwork bool //Only use the default network, ignore all others
	MaxSubnetSize     int  //Maxinimum subnet size. Default is 21, anything longer will be ignored as a valid interface.
}

//Generates a default PSconfig.
func newPSconfig() PSconfig {
	config := PSconfig{ListAddresses: false, UseDefaultNetwork: true, MaxSubnetSize: 20}
	
	return config
}

//Driver
func pingDriver(psconfig PSconfig) error{

	fmt.Println("Welcome to pingDriver")
	
	//Get list of ipv4 addresses
	subnetList, err := getInterface(psconfig)
	if err != nil {
		return fmt.Errorf("Could not get interfaces. %s", err)
	}

	//Generate list of address to ping
	allAddresses := generateAddresses(subnetList)

	fmt.Println(allAddresses, len(allAddresses))
	_=allAddresses

	return nil
}



func getInterface(psconfig PSconfig) ([]net.Addr, error) {
	//Slice of interfaces
	var interfaces []net.Addr

	//If UseDefaultNetwork is true only add the network with the default gateway to interfaces
	if psconfig.UseDefaultNetwork {
		//Get the interface index of the default gateway
		index, err := getGateway()
		if err != nil {
			return nil, fmt.Errorf("getInterface could not get the index of the default gateway: %s", err)
		}

		// Get the network interface and addresses
		iface, err := net.InterfaceByIndex(index)
		if err != nil {
			return nil, fmt.Errorf("getInterface could not get the InterfaceByIndex of the default gateway: %s", err)
		}

		addrs, err := iface.Addrs()
		if err != nil {
			return nil, fmt.Errorf("getInterface could not get the Addrs of the default gateway: %s", err)
		}

		interfaces = append(interfaces, addrs...)
	} else {
		//Get all interfacesinterfaces
		var err error
		interfaces, err = net.InterfaceAddrs()
		if err != nil {
			return nil, fmt.Errorf("Could not get interfaces. %s", err)
		}
	}
		
	//setup addressList and ignoreSubnets
	subnetList := []net.Addr{}
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
			//If IPv4 add to subnetList
			if ipAddr.IP.To4() != nil {
				//Ignore networks in unwanted subnets
				if ! ignoreSubnet(ignoreSubnets, ipAddr, psconfig.MaxSubnetSize) {
					subnetList = append(subnetList, addr)
				}			
			}
		}
	}

	//return addressList
	return subnetList, nil
}

func getGateway() (int, error) {
    //Get the list of network interfaces
    interfaces, err := net.Interfaces()
    if err != nil {
        return 0, fmt.Errorf("getGateway could not get Interfaces: %s\n", err)
    }

    //Find the default route (gateway) among the interfaces
    for _, iface := range interfaces {
        addrs, err := iface.Addrs()
        if err != nil {
            fmt.Printf("getGateway could not cvonert iface %s: %s\n", iface, err)
            continue
        }

        for _, addr := range addrs {
            ipnet, ok := addr.(*net.IPNet)
            if !ok {
				fmt.Printf("getGateway could not cvonert addr %s: %s\n", iface, err)
                continue
            }

			//Look only at unicast interfaces
            if ipnet.IP.IsGlobalUnicast() && !ipnet.IP.IsLoopback() && !ipnet.IP.IsLinkLocalUnicast() {
				//Only accept if IPv4
				if ipnet.IP.To4() != nil {
					//fmt.Printf("Default Route (Gateway) for %s: %s, %d\n", iface.Name, ipnet.IP.String(), iface.Index)
					return iface.Index, nil
				}
				
            }
        }
    }

	//If it got this found no route was found, ignore
	return 0, nil
}

func ignoreSubnet(ignoreSubnets []*net.IPNet, ipAddr *net.IPNet, maxSubnetSize int) bool {
	//Iterate through each subnet
	for _, subnet := range ignoreSubnets {
		//Check if subnet is in the existing subnet
		if subnet.Contains(ipAddr.IP) {
			return true
		}

		//Check if subnet is larger than max mask. Return true if it is
		mask, _ := ipAddr.Mask.Size()
		if mask < maxSubnetSize {
			return true
		}
	}

	return false
}

func generateAddresses(subnetList []net.Addr) []net.IP {
	var allAddresses []net.IP

	//Iterate through every subnet
	for _, subnet := range subnetList {
		//skip over this if not ok
		ipnet, ok := subnet.(*net.IPNet)
		if ! ok {
			continue
		}

		//Set network address
		ip := ipnet.IP.Mask(ipnet.Mask)

		//Find all addresses
		for ipnet.Contains(ip) {
			//Add to allAddresses
			allAddresses = append(allAddresses, ip)

			//Increment ip
			for j := len(ip) - 1; j >= 0; j-- {
				ip[j]++
				if ip[j] > 0 {
					break
				}
			}
		}
	}

	return allAddresses
}