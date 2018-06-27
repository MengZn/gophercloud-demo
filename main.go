package main

import (
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack"
	"fmt"
	"gophercloud-demo/keys"
	"gophercloud-demo/flavors"
)

func main() {
	authOpts := gophercloud.AuthOptions{
		IdentityEndpoint: "http://10.0.0.21:35357/v3",
		Username:         "admin",
		Password:         "h2u1E1mKp3jFfFKriBE55bHTtdeZYKEMVcWf0Ron",
		DomainName:       "Default",
		TenantName:       "admin",
	}
	provider, err := openstack.AuthenticatedClient(authOpts)
	if err != nil {
		fmt.Println("Fatal AuthenticatedClient")
		fmt.Errorf("Fatal error autenticating:  %s \n", err)
	}
	client, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{
		Region: "RegionOne",
	})
	if err != nil {
		fmt.Println("Fatal NewComputeV2")
		fmt.Errorf("Fatal error NewComputeV2:  %s \n", err)
	}
	keys.GetKeyPairs(client, "jason")
	flavors.GetFlavors(client)
}
