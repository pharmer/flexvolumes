package main

import (
	"testing"
	"fmt"
	"github.com/digitalocean/godo"
)

func TestVolume(t *testing.T) {
	client := getClient("") //put token here
	_, _, err := client.Storage.CreateVolume(&godo.VolumeCreateRequest{
		Name: "flexvolume",
		Region: "nyc3",
		SizeGigaBytes:int64(10),
	})
	fmt.Println(err)
	//vol, _, err := client.Storage.GetVolume("371679e3-f3")
	//fmt.Println(vol.Name, err)
}
