package packet

import (
	"github.com/packethost/packngo"
	"github.com/pkg/errors"
)

func getDevice(client *packngo.Client, projectID string, nodeName string) (*packngo.Device, error) {
	devices, _, err := client.Devices.List(projectID, nil)
	if err != nil {
		return nil, err
	}

	for _, device := range devices {
		if device.Hostname == nodeName {
			return &device, nil
		}
	}
	return nil, errors.Errorf("no device found with %v name", nodeName)
}

func getVolumeId(client *packngo.Client, projectID, pvName string) (string, error) {
	vols, _, err := client.Volumes.List(projectID, nil)
	if err != nil {
		return "", err
	}
	for _, v := range vols {
		if v.Description == pvName {
			return v.ID, nil
		}
	}
	return "", errors.Errorf("no volume found with description %s", pvName)

}
