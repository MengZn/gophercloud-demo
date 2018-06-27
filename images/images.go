package images


import (
	"fmt"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
	"github.com/rackspace/gophercloud/openstack/compute/v2/images"
)
func ShowImages(client *gophercloud.ServiceClient) {
	opts := images.ListOpts{}
	pager := images.ListDetail(client, opts)
	// Define an anonymous function to be executed on each page's iteration
	pager.EachPage(func(page pagination.Page) (bool, error) {
		imageList, err := images.ExtractImages(page)
		if err != nil {
			fmt.Errorf("Fatal error Extract Images:  %s \n", err)
		}
		for _, i := range imageList {
			// "i" will be a images.Image

			fmt.Printf("images is %v \n", i)
		}
		return false, err
	})
}

