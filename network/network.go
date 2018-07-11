package network

import (
	"github.com/rackspace/gophercloud/openstack/networking/v2/networks"
	"github.com/rackspace/gophercloud"
	"fmt"
	"regexp"
	"log"
	"strings"
	"github.com/rackspace/gophercloud/pagination"
)

type request struct {
	pager  *pagination.Pager
	client *gophercloud.ServiceClient
}

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

func ListNetwork(client *gophercloud.ServiceClient) {
	shared := bool(false)
	opts := networks.ListOpts{Shared: &shared}
	pager := networks.List(client, opts)
	req := request{pager: &pager, client: client}

	//test := networks.List(client, opts)
	//err := test.EachPage(func(page pagination.Page) (bool, error) {
	//	networkList, err := networks.ExtractNetworks(page)
	//
	//	for _, n := range networkList {
	//		// "n" will be a networks.Network
	//		fmt.Printf("%v/n", n)
	//	}
	//	return false, err
	//})

	err := req.EachPage(createURL(client, port, version, api), func(page pagination.Page) (bool, error) {
		networkList, err := networks.ExtractNetworks(page)

		for _, n := range networkList {
			// "n" will be a networks.Network
			fmt.Println(n)
		}
		return false, err
	})
	if err != nil {
		log.Fatal(err)
	}
}

func (r request) EachPage(url string, handler func(pagination.Page) (bool, error)) error {
	if r.pager.Err != nil {
		return r.pager.Err
	}
	currentURL := url
	for {
		currentPage, err := r.fetchNextPage(currentURL)
		if err != nil {
			return err
		}

		empty, err := currentPage.IsEmpty()
		if err != nil {
			return err
		}
		if empty {
			return nil
		}

		ok, err := handler(currentPage)
		if err != nil {
			return err
		}
		if !ok {
			return nil
		}

		currentURL, err = currentPage.NextPageURL()
		if err != nil {
			return err
		}
		if currentURL == "" {
			return nil
		}
	}
}

func (r request) fetchNextPage(url string) (pagination.Page, error) {
	resp, err := pagination.Request(r.client, r.pager.Headers, url)
	if err != nil {
		return nil, err
	}

	remembered, err := pagination.PageResultFrom(resp)
	if err != nil {
		return nil, err
	}

	return networks.NetworkPage{pagination.LinkedPageBase{PageResult: remembered}}, nil
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
