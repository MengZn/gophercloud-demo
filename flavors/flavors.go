package flavors

import (
	"github.com/rackspace/gophercloud"
	"fmt"
	"errors"
)

//func GetFlavors(client *gophercloud.ServiceClient) {
//	opts := flavors.ListOpts{}
//	// Retrieve a pager (i.e. a paginated collection)
//	pager := flavors.ListDetail(client, opts)
//	pager.EachPage(func(page pagination.Page) (bool, error) {
//		flavorList, err := flavors.ExtractFlavors(page)
//
//		if err != nil {
//			fmt.Errorf("Fatal error Extract Images:  %s \n", err)
//		}
//		for _, i := range flavorList {
//			// "i" will be a images.Image
//
//			fmt.Printf("images is %v \n", i)
//		}
//		return false, err
//	})
//}
func TestF(client *gophercloud.ServiceClient) {
	opts := CreateOpts{
		Name:  "test_flavor",
		Ram:   "1024",
		Vcpus: "2",
		Disk:  "10",
		//Id:          "auto",
		//RxTxFactor:  "1",
		//Swap:        "0",
		//Description: "test description",
	}
	NewFlavor(client, opts)
}

type CreateOptsBuilder interface {
	FlavorCreateMap() (map[string]interface{}, error)
}

type CreateOpts struct {
	Name        string
	Ram         string
	Vcpus       string
	Disk        string
	Id          string
	Swap        string
	RxTxFactor  string
	Description string
}

func (opts CreateOpts) FlavorCreateMap() (map[string]interface{}, error) {
	if opts.Name == "" {
		return nil, errors.New("Missing field required for flavors creation: Name")
	} else if opts.Ram == "" {
		return nil, errors.New("Missing field required for flavors creation: Ram")
	} else if opts.Vcpus == "" {
		return nil, errors.New("Missing field required for flavors creation: Vcpus")
	} else if opts.Disk == "" {
		return nil, errors.New("Missing field required for flavors creation: Disk")
	}
	flavors := make(map[string]interface{})
	flavors["flavor"] = map[string]interface{}{
		"name":  opts.Name,
		"ram":   opts.Ram,
		"vcpus": opts.Vcpus,
		"disk":  opts.Disk,
	}
	return flavors, nil
}

type CreateResult struct {
	flavorResult
}

type flavorResult struct {
	gophercloud.Result
}

func postFlavor(client *gophercloud.ServiceClient, opts CreateOptsBuilder) CreateResult {
	var res CreateResult
	url := client.ServiceURL("flavors")
	reqBody, _ := opts.FlavorCreateMap()
	_, res.Err = client.Post(url, reqBody, &res.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})

	return res
}

func NewFlavor(client *gophercloud.ServiceClient, opt CreateOpts) {
	result := postFlavor(client, opt)
	fmt.Printf("%v", result)
}
