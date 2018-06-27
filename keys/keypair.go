package keys

import (
	"crypto/rsa"
	"crypto/rand"
	"fmt"
	"golang.org/x/crypto/ssh"
	"github.com/rackspace/gophercloud/openstack/compute/v2/extensions/keypairs"
	"github.com/rackspace/gophercloud"
)

func CreateKeyPairs(client *gophercloud.ServiceClient) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Errorf("Fatal GenerateKey %v \n", err)
	}
	publicKey := privateKey.PublicKey
	pub, err := ssh.NewPublicKey(&publicKey)
	pubBytes := ssh.MarshalAuthorizedKey(pub)
	pk := string(pubBytes)
	kp, err := keypairs.Create(client, keypairs.CreateOpts{
		Name:      "keypair_name",
		PublicKey: pk,
	}).Extract()
	fmt.Println("kp %v \n", kp)

}

func DeleteKeyParis(client *gophercloud.ServiceClient, name string) {
	err := keypairs.Delete(client, name).ExtractErr()
	if err != nil {
		fmt.Errorf("Fatal Delete keypairs%v \n", err)
	}
	fmt.Printf("Delete keypairs success \n")
}
func GetKeyPairs(client *gophercloud.ServiceClient, name string) {
	keyPair, err := keypairs.Get(client, name).Extract()
	if err != nil {
		fmt.Errorf("Fatal Get keypairs%v \n", err)
	}
	fmt.Printf("keypairs %v", *keyPair)
}
