package main

import (
	"testing"
	"fmt"
	//"github.com/digitalocean/godo"
	"strings"
)

func TestVolume(t *testing.T) {
	client := getClient("") //put token here
	/*_, _, err := client.Storage.CreateVolume(&godo.VolumeCreateRequest{
		Name: "flexvolume",
		Region: "nyc3",
		SizeGigaBytes:int64(10),
	})
	fmt.Println(err)*/
	vol, _, err := client.Storage.GetVolume("cc2fbaa1-fa9e-11e7-8d31-")
	fmt.Println(vol.Name, err)
}

func TestStr(t *testing.T) {
	credentials := "\t Hello\n"
	fmt.Print(credentials,"..")
	cred := strings.TrimSpace(credentials)
	fmt.Print(cred,"...")
}
