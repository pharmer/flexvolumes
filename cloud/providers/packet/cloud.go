package packet

import (
	"fmt"

	"github.com/packethost/packngo"
)

func getDevice(client *packngo.Client, projectID string, nodeName string) (*packngo.Device, error) {
	devices, _, err := client.Devices.List(projectID)
	if err != nil {
		return nil, err
	}

	for _, device := range devices {
		if device.Hostname == nodeName {
			return &device, nil
		}
	}
	return nil, fmt.Errorf("no device found with %v name", nodeName)
}
