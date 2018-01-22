package packet

import (
	"fmt"
	"testing"

	"github.com/packethost/packngo"
)

func TestVol(t *testing.T) {
	client := getClient()

	vol, _, err := client.Volumes.Get("")
	if err != nil {
		fmt.Println(err)
	}

	device, err := getDevice(client, "", "")
	fmt.Println(device.ID, err)

	for _, v := range vol.Attachments {
		fmt.Println(v.Device.ID)
	}
	/*device, err := getDevice(client, "", "baremetal-0-pool-mjbwrc")
	if err != nil {
		fmt.Println(err)
	}

	//fmt.Println(device)

	fmt.Println("device id = ", device.ID)
	droplet, _, err := client.Devices.Get(device.ID)
	fmt.Println(err)
	//fmt.Println(device)
	for _, vid := range droplet.Volumes {
		fmt.Println("..", vid.Attachments)
	}

	for _, v := range vol.Attachments {
		fmt.Println(v.ID, v.Href, v.Device.ID, v.Volume)
	}
	fmt.Println(vol.Name, vol.ID, vol.Description)*/
}

func getClient() *packngo.Client {
	return packngo.NewClient("", "", nil)
}
