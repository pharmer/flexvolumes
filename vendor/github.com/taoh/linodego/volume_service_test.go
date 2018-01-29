package linodego

import (
	"testing"
	"fmt"
)

func TestListVolume(t *testing.T) {
	client := NewClient(APIKey, nil)
	resp, err := client.Volume.List(0)
	fmt.Println(err)
	fmt.Println(resp.Volume)
}

func TestCreateVolume(t *testing.T) {
	client := NewClient(APIKey, nil)
	args := map[string]string{
		"DatacenterID": "3",
	}
	resp, err := client.Volume.Create(10,"apitest",args)
	fmt.Println(err, resp)

}

func TestDeleteVolume(t *testing.T) {
	client := NewClient(APIKey, nil)
	resp, err := client.Volume.Delete(3789)
	fmt.Println(err, resp)

}

func TestVolumeUpdate(t *testing.T) {
	client := NewClient(APIKey, nil)
	args := map[string]string{
		"LinodeID": "0",
	}
	resp, err := client.Volume.Update(3791, args)
	fmt.Println(err, resp)
}
