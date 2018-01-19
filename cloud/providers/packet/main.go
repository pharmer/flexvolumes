package packet

import (
	"fmt"

	. "github.com/pharmer/flexvolumes/cloud"
)

type PacketOptions struct {
	DefaultOptions
}

func (v *VolumeManager) NewOptions() interface{} {
	return &PacketOptions{}
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
	opt := options.(*PacketOptions)

	vol, _, err := v.client.Volumes.Get(opt.VolumeID)
	if err != nil {
		return "", err
	}

	projectID, err := getProjectID()
	if err != nil {
		return "", err
	}
	device, err := getDevice(v.client, projectID, nodeName)
	if err != nil {
		return "", err
	}

	isAttached := false
	for _, v := range vol.Attachments {
		if v.ID == device.ID {
			isAttached = true
		}
	}

	if !isAttached {
		_, _, err := v.client.VolumeAttachments.Create(vol.ID, device.ID)
		if err != nil {
			return "", err
		}

		//TODO(sanjid): add wait here
		/*if err = awaitAction(v.client, vol.ID, action); err != nil {
			return "", err
		}*/
	}

	return DEVICE_PREFIX + vol.Name, nil
}

func (v *VolumeManager) Detach(device, nodeName string) error {
	projectID, err := getProjectID()
	if err != nil {
		return err
	}
	droplet, err := getDevice(v.client, projectID, nodeName)
	if err != nil {
		return err
	}

	isDetached := true
	attachmentID := ""
	for _, vid := range droplet.Volumes {
		if vid.Name == device {
			isDetached = false
			attachmentID = vid.Attachments[0].ID
			break
		}
	}
	if !isDetached {
		_, err := v.client.VolumeAttachments.Delete(attachmentID)
		if err != nil {
			return err
		}

		//TODO(sanjid): add wait here
		/*if err = awaitAction(v.client, vol.ID, action); err != nil {
			return err
		}*/
		return nil
	}
	return fmt.Errorf("could not find volume attached at %v", device)
}

func (v *VolumeManager) MountDevice(mountDir string, device string, options interface{}) error {
	opt := options.(*PacketOptions)
	return Mount(mountDir, device, opt.DefaultOptions)
}

func (v *VolumeManager) Mount(mountDir string, options interface{}) error {
	return ErrNotSupported
}

func (v *VolumeManager) Unmount(mountDir string) error {
	return Unmount(mountDir)
}
