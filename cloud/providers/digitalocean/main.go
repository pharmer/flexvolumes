package digitalocean

import (
	"fmt"

	"github.com/digitalocean/godo"
	. "github.com/pharmer/flexvolumes/cloud"
	"golang.org/x/oauth2"
)

type DigitalOceanOptions struct {
	DefaultOptions
}

func (v *VolumeManager) NewOptions() interface{} {
	return &DigitalOceanOptions{}
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
	opt := options.(*DigitalOceanOptions)

	vol, _, err := v.client.Storage.GetVolume(oauth2.NoContext, opt.VolumeID)
	if err != nil {
		return "", err
	}

	droplet, err := getDroplet(v.client, nodeName)
	if err != nil {
		return "", err
	}

	isAttached := false
	for _, vid := range droplet.VolumeIDs {
		if vid == vol.ID {
			isAttached = true
			break
		}
	}

	if !isAttached {
		action, _, err := v.client.StorageActions.Attach(oauth2.NoContext, vol.ID, droplet.ID)
		if err != nil {
			return "", err
		}

		if err = awaitAction(v.client, vol.ID, action); err != nil {
			return "", err
		}
	}

	return DEVICE_PREFIX + vol.Name, nil
}

func (v *VolumeManager) Detach(device, nodeName string) error {
	droplet, err := getDroplet(v.client, nodeName)
	if err != nil {
		return err
	}

	params := &godo.ListVolumeParams{
		Name:   device,
		Region: droplet.Region.Slug,
	}
	vols, _, err := v.client.Storage.ListVolumes(oauth2.NoContext, params)
	if err != nil {
		return fmt.Errorf("failed to list volumes: %v", err.Error())
	}
	vol := vols[0]

	isDetached := true
	for _, vid := range droplet.VolumeIDs {
		if vid == vol.ID {
			isDetached = false
			break
		}
	}
	if !isDetached {
		action, _, err := v.client.StorageActions.DetachByDropletID(oauth2.NoContext, vol.ID, droplet.ID)
		if err != nil {
			return err
		}

		if err = awaitAction(v.client, vol.ID, action); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("could not find volume attached at %v", device)
}

func (v *VolumeManager) MountDevice(mountDir string, device string, options interface{}) error {
	opt := options.(*DigitalOceanOptions)
	return Mount(mountDir, device, opt.DefaultOptions)
}

func (v *VolumeManager) Mount(mountDir string, options interface{}) error {
	opt := options.(*DigitalOceanOptions)

	vol, _, err := v.client.Storage.GetVolume(oauth2.NoContext, opt.VolumeID)
	if err != nil {
		return err
	}

	device := DEVICE_PREFIX + vol.Name
	return Mount(mountDir, device, opt.DefaultOptions)
}

func (v *VolumeManager) Unmount(mountDir string) error {
	return Unmount(mountDir)
}
