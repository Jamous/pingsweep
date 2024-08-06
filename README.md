Pingsweep
=========

Pingsweep is a go module that will ping every device in a given subnet. It can either automatically grab the subnet, based on connected subnets, or a subnet can be specified. It does this by sending one ICMP echo request, then closing the connection.

Pingsweep is useful for network discovery and filling out IPv4 arp tables.

This module will set off some antivirus programs.

Options
-------
Pingsweep has several options, you can use the default options, or specify your own. See examples/main.go for how to do this.

* UseDefaultNetwork (bool): Only use the default network, ignore all others. This will only select the subnet associated with the default gateway.
* MaxSubnetSize (int): Maximum subnet size. Default is 21, anything longer will be ignored as a valid interface. You can change this to ping larger subnets.
* CustomSubnet (string): Specify the custom subnet to ping. Must be in CIDR notation.

Ignored networks
----------------
Pingsweep will ignore anything in the 169.254.0.0/16 and 127.0.0.0/8 subnets. These are typically not routed and do not have other hosts on them.

Versions
--------
v0.0.1 - Initial release - 08/06/24
