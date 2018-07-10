package instance

import (
	"github.com/rackspace/gophercloud/openstack/compute/v2/servers"
	"fmt"
	"github.com/rackspace/gophercloud"
)

func CreateInstance(client *gophercloud.ServiceClient) {
	network := []servers.Network{
		{
			UUID: "8e08b348-ed62-47a4-a7e0-cf8ba240a70a",
		},
	}
	server, err := servers.Create(client, servers.CreateOpts{
		Name:       "GoodBoy",
		FlavorName: "m1.tiny",
		ImageName:  "cirros",
		Networks:   network,
	}).Extract()
	if err != nil {
		fmt.Println("Unable to create server: %s", err)
	}
	fmt.Println("Server ID: %s", server.ID)
}
