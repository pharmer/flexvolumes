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

/*
func getMountDevice(
	i *lightsail.Instance,
	volumeID *lightsail.Disk,
	assign bool) (assigned string, alreadyAttached bool, err error) {
	instanceType := i.getInstanceType()
	if instanceType == nil {
		return "", false, fmt.Errorf("could not get instance type for instance: %s", i.awsID)
	}

	deviceMappings := map[string]string{}

	for _, blockDevice := range i.Hardware.Disks {
		name := _aws.StringValue(blockDevice.DeviceName)
		if strings.HasPrefix(name, "/dev/sd") {
			name = name[7:]
		}
		if strings.HasPrefix(name, "/dev/xvd") {
			name = name[8:]
		}
		if len(name) < 1 || len(name) > 2 {
			glog.Warningf("Unexpected EBS DeviceName: %q", aws.StringValue(blockDevice.DeviceName))
		}
		deviceMappings[mountDevice(name)] = awsVolumeID(aws.StringValue(blockDevice.Ebs.VolumeId))
	}

	// We lock to prevent concurrent mounts from conflicting
	// We may still conflict if someone calls the API concurrently,
	// but the AWS API will then fail one of the two attach operations
	c.attachingMutex.Lock()
	defer c.attachingMutex.Unlock()

	for mountDevice, volume := range c.attaching[i.nodeName] {
		deviceMappings[mountDevice] = volume
	}

	// Check to see if this volume is already assigned a device on this machine
	for mountDevice, mappingVolumeID := range deviceMappings {
		if volumeID == mappingVolumeID {
			if assign {
				glog.Warningf("Got assignment call for already-assigned volume: %s@%s", mountDevice, mappingVolumeID)
			}
			return mountDevice, true, nil
		}
	}

	if !assign {
		return mountDevice(""), false, nil
	}

	// Find the next unused device name
	deviceAllocator := c.deviceAllocators[i.nodeName]
	if deviceAllocator == nil {
		// we want device names with two significant characters, starting with /dev/xvdbb
		// the allowed range is /dev/xvd[b-c][a-z]
		// http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/device_naming.html
		deviceAllocator = NewDeviceAllocator()
		c.deviceAllocators[i.nodeName] = deviceAllocator
	}
	// We need to lock deviceAllocator to prevent possible race with Deprioritize function
	deviceAllocator.Lock()
	defer deviceAllocator.Unlock()

	chosen, err := deviceAllocator.GetNext(deviceMappings)
	if err != nil {
		glog.Warningf("Could not assign a mount device.  mappings=%v, error: %v", deviceMappings, err)
		return "", false, fmt.Errorf("Too many EBS volumes attached to node %s.", i.nodeName)
	}

	attaching := c.attaching[i.nodeName]
	if attaching == nil {
		attaching = make(map[mountDevice]awsVolumeID)
		c.attaching[i.nodeName] = attaching
	}
	attaching[chosen] = volumeID
	glog.V(2).Infof("Assigned mount device %s -> volume %s", chosen, volumeID)

	return chosen, false, nil
}*/


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
