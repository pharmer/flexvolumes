package lightsail

import (
	"fmt"
	. "github.com/pharmer/flexvolumes/cloud"
	//"github.com/aws/aws-sdk-go/service/lightsail"
	"github.com/aws/aws-sdk-go/service/lightsail"
)

type LightsailOptions struct {
	DefaultOptions
}

func (v *VolumeManager) NewOptions() interface{} {
	return &LightsailOptions{}
}

func (v *VolumeManager) Initialize() error {
	token, err := getCredential()
	if err != nil {
		return err
	}
	v.client, err = token.getClient()
	return err
}

func (v *VolumeManager) Init() error {
	return nil
}

func (v *VolumeManager) Attach(options interface{}, nodeName string) (string, error) {
	opt := options.(*LightsailOptions)
	name := opt.VolumeName
	if name == "" {
		name = opt.VolumeID
	}

	instance, err := instanceByName(v.client, nodeName)
	if err != nil {
		return "", err
	}

	disk, err := getDiskByName(v.client, name)
	if err != nil {
		return "", err
	}

	isAttached := false

	if disk.AttachedTo == instance.Name {
		isAttached = true
	}
	if !isAttached {
		v.client.AttachDisk(&lightsail.AttachDiskInput{
			DiskName: disk.Name,
			InstanceName: instance.Name,
		})
		//v.client.attachdi
	}





	return DEVICE_PREFIX + opt.VolumeID, nil
}

func (v *VolumeManager) Detach(device, nodeName string) error {
	return fmt.Errorf("could not find volume attached at %v", device)
}

func (v *VolumeManager) MountDevice(mountDir string, device string, options interface{}) error {
	opt := options.(*LightsailOptions)
	return Mount(mountDir, device, opt.DefaultOptions)
}

func (v *VolumeManager) Mount(mountDir string, options interface{}) error {
	return ErrNotSupported
}

func (v *VolumeManager) Unmount(mountDir string) error {
	return Unmount(mountDir)
}
