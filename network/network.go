package network

import (
	"github.com/rackspace/gophercloud/openstack/networking/v2/networks"
	"github.com/rackspace/gophercloud"
	"fmt"
	"regexp"
	"log"
	"strings"
)

const (
	port    = "9696"
	version = "v2.0"
	api     = "networks"
)

func CreateNetwork(client *gophercloud.ServiceClient) {
	opts := networks.CreateOpts{Name: "main_network", AdminStateUp: networks.Up}

	network, err := create(client, opts).Extract()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v", network.ID)

}

func create(c *gophercloud.ServiceClient, opts networks.CreateOptsBuilder) networks.CreateResult {
	var res networks.CreateResult

	reqBody, err := opts.ToNetworkCreateMap()
	if err != nil {
		res.Err = err
		return res
	}

	_, res.Err = c.Post(createURL(c, port, version, api), reqBody, &res.Body, nil)
	return res
}

func createURL(client *gophercloud.ServiceClient, parts ...string) string {
	domainRegexp := regexp.MustCompile(`^http://[1-2]{0,1}[0-9]{0,1}[0-9]{1}.[1-2]{0,1}[0-9]{0,1}[0-9]{1}.[1-2]{0,1}[0-9]{0,1}[0-9]{1}.[1-2]{0,1}[0-9]{0,1}[0-9]{1}:`)
	domainIP := domainRegexp.FindString(client.Endpoint)
	//domainIP := strings.Split(client.Endpoint, ":")[0:2]
	//fmt.Println(domainIP)
	return domainIP + strings.Join(parts, "/")
}
