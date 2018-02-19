package packet

import (
	"fmt"
	"os/exec"
	"strings"

	. "github.com/pharmer/flexvolumes/cloud"
	"k8s.io/apimachinery/pkg/util/wait"
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

	projectID, err := getProjectID()
	if err != nil {
		return "", err
	}

	volume := opt.VolumeID
	if volume == "" {
		volume, err = getVolumeId(v.client, projectID, opt.PVorVolumeName)
		if err != nil {
			return "", err
		}
	}
	vol, _, err := v.client.Volumes.Get(volume)
	if err != nil {
		return "", err
	}
	device, err := getDevice(v.client, projectID, nodeName)
	if err != nil {
		return "", err
	}

	isAttached := false
	for _, v := range vol.Attachments {
		if v.Device.ID == device.ID {
			isAttached = true
		}
	}

	if !isAttached {
		_, _, err := v.client.VolumeAttachments.Create(vol.ID, device.ID)
		if err != nil {
			return "", err
		}
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
		//Href:"/storage/870d1ba4-b705-4c0f-974d-92ee73b010db"
		href := strings.Split(vid.Href, "/")
		if len(href) > 0 {
			id := href[len(href)-1]
			volume, _, err := v.client.Volumes.Get(id)
			if err != nil {
				return err
			}
			if volume.Description == device {
				isDetached = false
				attachmentID = volume.Attachments[0].ID
				break
			}
		}
	}
	if !isDetached {
		_, err := v.client.VolumeAttachments.Delete(attachmentID)
		if err != nil {
			return err
		}
		if err := wait.PollImmediate(RetryInterval, RetryTimeout, func() (bool, error) {
			at, _, _ := v.client.VolumeAttachments.Get(attachmentID)
			if at == nil {
				return true, nil
			}
			return false, nil
		}); err != nil {
			return err
		}
		return exec.Command("/sbin/multipath", "-W").Run()

	}
	return fmt.Errorf("could not find volume attached at %v", device)
}

func (v *VolumeManager) MountDevice(mountDir string, device string, options interface{}) error {
	opt := options.(*PacketOptions)
	projectID, err := getProjectID()
	if err != nil {
		return err
	}
	volume := opt.VolumeID
	if volume == "" {
		volume, err = getVolumeId(v.client, projectID, opt.PVorVolumeName)
		if err != nil {
			return err
		}
	}
	vol, _, err := v.client.Volumes.Get(volume)
	if err != nil {
		return err
	}
	cmd := exec.Command("packet-block-storage-attach", "-m", "queue", vol.Name)
	if err := cmd.Run(); err != nil {
		return err
	}
	return Mount(mountDir, device, opt.DefaultOptions)
}

func (v *VolumeManager) Mount(mountDir string, options interface{}) error {
	return ErrNotSupported
}

func (v *VolumeManager) Unmount(mountDir string) error {
	if err := Unmount(mountDir); err != nil {
		return err
	}
	cmd := exec.Command("packet-block-storage-detach")
	return cmd.Run()
}
