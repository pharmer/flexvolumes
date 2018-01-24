package lightsail

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"github.com/aws/aws-sdk-go/service/lightsail"
	_aws "github.com/aws/aws-sdk-go/aws"
)


func instanceByName(client *lightsail.Lightsail, nodeName string) (*lightsail.Instance, error) {
	host, err := client.GetInstance(&lightsail.GetInstanceInput{
		InstanceName: _aws.String(nodeName),
	})
	if err != nil {
		return nil, err
	}
	if host.Instance != nil {
		return host.Instance, nil
	}
	return nil, fmt.Errorf("no instance found with name %v", nodeName)

}

func getDiskByName(client *lightsail.Lightsail, name string)(*lightsail.Disk, error)  {
	resp, err := client.GetDisk(&lightsail.GetDiskInput{
		DiskName: _aws.String("flextest"),
	})
	if err != nil {
		return nil, err
	}
	if resp.Disk != nil {
		return nil, fmt.Errorf("no volume found with %v volName", name)
	}
	return resp.Disk, nil

}


func getMountDevicePath(ins *lightsail.Instance) (string, error) {
	deviceMappings := make(map[string]bool, 0)

	for _, disk := range ins.Hardware.Disks {
		deviceMappings[_aws.StringValue(disk.Path)] = true
	}
	for i:= 103; i<= 123; i++ {

	}

	return "", nil

}


func getRegion() (string, error) {
	zone, err := getAvailabilityZone()
	if err != nil {
		return "", err
	}
	region, err := azToRegion(zone)
	if err != nil {
		return "", err
	}

	return region, nil
}

// Derives the region from a valid az name.
// Returns an error if the az is known invalid (empty)
func azToRegion(az string) (string, error) {
	if len(az) < 1 {
		return "", fmt.Errorf("invalid (empty) AZ")
	}
	region := az[:len(az)-1]
	return region, nil
}

func getAvailabilityZone() (string, error) {
	zone := "placement/availability-zone"
	return GetMetadata(zone)
}


func GetMetadata(path string) (string, error) {
	resp, err := http.Get(metadataURL + path)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	return string(body), err
}
