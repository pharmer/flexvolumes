package vultr

import (
	"fmt"

	. "github.com/pharmer/flexvolumes/cloud"
	"k8s.io/apimachinery/pkg/util/wait"
)

type VultrOptions struct {
	DefaultOptions
}

func (v *VolumeManager) NewOptions() interface{} {
	return &VultrOptions{}
}

func (v *VolumeManager) Initialize() error {
	key, err := getCredential()
	if err != nil {
		return err
	}
	v.client = key.getClient()
	return nil
}

func (v *VolumeManager) Init() error {
	return nil
}

func (v *VolumeManager) Attach(options interface{}, nodeName string) (string, error) {
	opt := options.(*VultrOptions)

	vol, err := v.client.GetBlockStorage(opt.VolumeID)
	if err != nil {
		return "", err
	}

	serverID, err := getServerID(v.client, nodeName)
	if err != nil {
		return "", err
	}

	isAttached := false
	if vol.AttachedTo == serverID {
		isAttached = true
	}

	var deviceName string
	if !isAttached {
		deviceName, err = getNextDeviceName(v.client, serverID)
		if err != nil {
			return "", err
		}
		err = v.client.AttachBlockStorage(opt.VolumeID, serverID)
		if err != nil {
			return "", err
		}
		err = writeDeviceName(vol.ID, deviceName)
		if err != nil {
			return "", err
		}

	} else {
		deviceName, err = getDeviceName(vol.ID)
		if err != nil {
			return "", err
		}
	}

	return deviceName, nil
}

func (v *VolumeManager) Detach(device, nodeName string) error {
	serverID, err := getServerID(v.client, nodeName)
	if err != nil {
		return err
	}
	volID, found, err := getVolumeId(v.client, device, serverID)
	if err != nil {
		return err
	}

	isDetached := !found
	if !isDetached {
		err = v.client.DetachBlockStorage(volID)
		if err != nil {
			return err
		}
		return wait.PollImmediate(RetryInterval, RetryTimeout, func() (bool, error) {
			vol, err := v.client.GetBlockStorage(volID)
			if err != nil {
				return false, err
			} else {
				if vol.AttachedTo != serverID {
					return true, nil
				}
				return false, nil
			}
		})
	}
	return fmt.Errorf("could not find volume attached at %v", nodeName)
}

func (v *VolumeManager) MountDevice(mountDir string, device string, options interface{}) error {
	opt := options.(*VultrOptions)
	return Mount(mountDir, device, opt.DefaultOptions)
}

func (v *VolumeManager) Mount(mountDir string, options interface{}) error {
	return ErrNotSupported
}

func (v *VolumeManager) Unmount(mountDir string) error {
	return Unmount(mountDir)
}
