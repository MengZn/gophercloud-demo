package network

import (
	"github.com/rackspace/gophercloud/openstack/networking/v2/networks"
	"github.com/rackspace/gophercloud"
	"fmt"
	"regexp"
	"log"
	"strings"
	"github.com/rackspace/gophercloud/pagination"
	"net/http"
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

func CreateNetwork(client *gophercloud.ServiceClient, opts *networks.CreateOpts) {
	network, err := create(client, opts).Extract()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v \n", network.ID)
}
func DeleteNetwork(client *gophercloud.ServiceClient, id string) {
	resp, _ := delete(client, id)
	if resp.StatusCode == 404 {
		fmt.Printf("Network %v could not be found.\n", id)
	} else if resp.StatusCode == 204 {
		fmt.Printf("Network %v is deleted.\n", id)
	}
}

func delete(client *gophercloud.ServiceClient, id string) (*http.Response, error) {
	var responseBody interface{}
	opts := &gophercloud.RequestOpts{JSONResponse: &responseBody}
	return client.Request("DELETE", createURL(client, port, version, api, id), *opts)
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

func GetNetworkDetails(client *gophercloud.ServiceClient) {
	//network, err := networks.Get(client, "1e83aecc-a0c1-489b-b918-5e9bc49b26ec").Extract()
	//fmt.Println(err)
	//fmt.Println(network)

	var responseBody interface{}
	//var response struct {
	//	Networks struct {
	//		Network *networks.Network `json:"network"`
	//	}
	//}
	//var response map[string]interface{
	//
	//}
	//var response struct {
	//	Network *networks.Network `json:"network"`
	//}
	opts := &gophercloud.RequestOpts{JSONResponse: &responseBody}
	client.Request("GET", createURL(client, port, version, api), *opts)
	//fmt.Printf("%v \n", responseBody)
	//a, err := redis.StringMap(responseBody, nil)
	//err := mapstructure.Decode(&responseBody, &response)
	//network, err := networks.Get(client, "1e83aecc-a0c1-489b-b918-5e9bc49b26ec").Extract()
	//if err != nil {
	//	log.Fatal(err)
	//}
	fmt.Println("-----------------")
	fmt.Println(responseBody)
	fmt.Println("-----------------")
	assertResponseBody(&responseBody)
	//fmt.Println(a)
	//fmt.Println(err)
}

func assertResponseBody(response *interface{}) {
	m := (*response).(map[string]interface{})
	for k, v := range m {
		switch vv := v.(type) {
		case string:
			fmt.Println(k, "is string", vv)
		case float64:
			fmt.Println(k, "is float64", vv)
		case []interface{}:
			fmt.Println(k, "is an array:")
			for i, u := range vv {
				fmt.Println(i, u)
			}

		default:
			fmt.Println(k)
			fmt.Println(v)
			fmt.Println(k, "is of a type I don't know how to handle")
		}
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
