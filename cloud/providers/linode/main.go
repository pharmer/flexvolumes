package linode

import (
	"fmt"
	"strconv"

	. "github.com/pharmer/flexvolumes/cloud"
	"github.com/taoh/linodego"
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
	var vol *linodego.Volume
	var err error
	if opt.VolumeID == "" {
		vName := opt.VolumeName
		if vName == "" {
			vName = opt.PVorVolumeName
		}

		vol, err = getVolumeFromName(v.client, vName)
		if err != nil {
			return "", err
		}
	} else {
		volID, err := strconv.Atoi(opt.VolumeID)
		if err != nil {
			return "", err
		}
		resp, err := v.client.Volume.List(volID)
		if err != nil {
			return "", err
		}
		if len(resp.Volume) == 1 {
			vol = &resp.Volume[0]
		} else {
			return "", fmt.Errorf("no volume found with id %v", opt.VolumeID)
		}
	}

	linode, err := linodeByName(v.client, nodeName)
	if err != nil {
		return "", err
	}

	isAttached := false
	if vol.LinodeId == linode.LinodeId {
		isAttached = true
	}

	if !isAttached {
		args := map[string]string{
			"LinodeID": strconv.Itoa(linode.LinodeId),
		}
		_, err := v.client.Volume.Update(vol.VolumeId, args)
		if err != nil {
			return "", err
		}
		if err = awaitAction(v.client, linode.LinodeId); err != nil {
			return "", err
		}
	}

	return DEVICE_PREFIX + vol.Label.String(), nil
}

func (v *VolumeManager) Detach(device, nodeName string) error {
	vol, err := getVolumeFromName(v.client, device)
	if err != nil {
		return err
	}
	linode, err := linodeByName(v.client, nodeName)
	if err != nil {
		return err
	}
	isDetached := true
	if vol.LinodeId == linode.LinodeId {
		isDetached = false
	}

	if !isDetached {
		args := map[string]string{
			"LinodeID": "0",
		}
		_, err := v.client.Volume.Update(vol.VolumeId, args)
		if err != nil {
			return err
		}
		if err = awaitAction(v.client, linode.LinodeId); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("could not find volume attached at %v", device)
}

func (v *VolumeManager) MountDevice(mountDir string, device string, options interface{}) error {
	return ErrNotSupported
}

func (v *VolumeManager) Mount(mountDir string, options interface{}) error {
	return ErrNotSupported
}

func (v *VolumeManager) Unmount(mountDir string) error {
	return ErrNotSupported
}
