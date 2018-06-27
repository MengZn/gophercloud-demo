package flavors

import (
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack/compute/v2/flavors"
	"github.com/rackspace/gophercloud/pagination"
	"fmt"
)

func GetFlavors(client *gophercloud.ServiceClient) {
	opts := flavors.ListOpts{}
	// Retrieve a pager (i.e. a paginated collection)
	pager := flavors.ListDetail(client, opts)
	pager.EachPage(func(page pagination.Page) (bool, error) {
		flavorList, err := flavors.ExtractFlavors(page)

		if err != nil {
			fmt.Errorf("Fatal error Extract Images:  %s \n", err)
		}
		for _, i := range flavorList {
			// "i" will be a images.Image

			fmt.Printf("images is %v \n", i)
		}
		return false, err
	})
}
