package network

import (
	"github.com/rackspace/gophercloud/openstack/networking/v2/subnets"
	"github.com/rackspace/gophercloud"
	"fmt"
	"log"
)

func CreatSubnet(client *gophercloud.ServiceClient) {
	opts := &subnets.CreateOpts{
		NetworkID: "7d7fa869-806b-47bc-8bd2-dcad292a6588",
		CIDR:      "192.168.199.0/24",
		IPVersion: subnets.IPv4,
		Name:      "my_subnet",
	}
	createSubnet(client,opts)
}

func createSubnet(client *gophercloud.ServiceClient, createOpts *subnets.CreateOpts) {
	var responseBody interface{}
	opt := &gophercloud.RequestOpts{JSONResponse: &responseBody}
	reqBody, err := createOpts.ToSubnetCreateMap()
	if err != nil {
		log.Fatalf("SubnetCreate error: %v \n", err)
	}
	opt.JSONBody = reqBody

	client.Request("POST", createURL(client, port, version, apiSubnet), *opt)
	fmt.Printf("%v \n", responseBody)
}
