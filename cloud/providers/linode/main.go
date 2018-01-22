package linode

import (
"fmt"

. "github.com/pharmer/flexvolumes/cloud"
)

type LinodeOptions struct {
	DefaultOptions
}

func (v *VolumeManager) NewOptions() interface{} {
	return &LinodeOptions{}
}

func (v *VolumeManager) Initialize() error {
	token, err := getCredential()
	if err != nil {
		return err
	}
	v.client = token.getClient()
	return nil
}

func (v *VolumeManager) Init() error {
	return nil
}

func (v *VolumeManager) Attach(options interface{}, nodeName string) (string, error) {
	opt := options.(*LinodeOptions)

	/*vol, _, err := v.client.Disk.Delete()
	if err != nil {
		return "", err
	}*/

	return DEVICE_PREFIX + opt.ApiKey, nil
}

func (v *VolumeManager) Detach(device, nodeName string) error {
		return fmt.Errorf("could not find volume attached at %v", device)
}

func (v *VolumeManager) MountDevice(mountDir string, device string, options interface{}) error {
	opt := options.(*LinodeOptions)
	return Mount(mountDir, device, opt.DefaultOptions)
}

func (v *VolumeManager) Mount(mountDir string, options interface{}) error {
	return ErrNotSupported
}

func (v *VolumeManager) Unmount(mountDir string) error {
	return Unmount(mountDir)
}
