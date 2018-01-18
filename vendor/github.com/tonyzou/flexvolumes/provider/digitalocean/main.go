package main

import (
	"io/ioutil"
	"strings"
	"time"

	"github.com/digitalocean/godo"
	f "github.com/tonyzou/flexvolumes"
)

const CRED_FILE string = "/tmp/do_creds"
const DEVICE_PREFIX string = "/dev/disk/by-id/scsi-0DO_Volume_"

type DigitalOceanOptions struct {
	f.DefaultOptions
	VolumeId string `json:"volumeId"`
}

type DigitalOceanPlugin struct{}

func (DigitalOceanPlugin) NewOptions() interface{} {
	return &DigitalOceanOptions{}
}

func (DigitalOceanPlugin) Init() f.Result {
	return f.Succeed()
}

func (DigitalOceanPlugin) Attach(opts interface{}) f.Result {
	opt := opts.(*DigitalOceanOptions)
	client := getClient(opt.ApiKey)
	if client == nil {
		return f.Fail("Could not create client")
	}

	vol, _, err := client.Storage.GetVolume(opt.VolumeId)
	if err != nil {
		return f.Fail("Could not get volume \"", opt.VolumeId, "\": ", err.Error())
	}

	droplet, err := detectDroplet(client)
	if err != nil {
		return f.Fail(err.Error())
	}

	action, _, err := client.StorageActions.Attach(vol.ID, droplet.ID)
	if err != nil {
		return f.Fail(err.Error())
	}

	if res := awaitAction(client, vol.ID, action); res != nil {
		return *res
	}

	if err = ioutil.WriteFile(CRED_FILE, []byte(opt.ApiKey), 0600); err != nil {
		return f.Fail("Could not save credentials: ", err.Error())
	}

	return f.Result{
		Status: "Success",
		Device: DEVICE_PREFIX + vol.Name,
	}
}

func (DigitalOceanPlugin) Detach(device string) f.Result {
	credBytes, err := ioutil.ReadFile(CRED_FILE)
	if err != nil {
		return f.Fail("Failed to read credentials: ", err.Error())
	}

	client := getClient(string(credBytes))
	if client == nil {
		return f.Fail("Could not create client")
	}

	if !strings.HasPrefix(device, DEVICE_PREFIX) {
		return f.Fail("Expected path starting with ", DEVICE_PREFIX, "; instead got ", device)
	}

	diskName := device[len(DEVICE_PREFIX):]

	vols, _, err := client.Storage.ListVolumes(nil)
	if err != nil {
		return f.Fail("Failed to list volumes: ", err.Error())
	}

	for _, vol := range vols {
		if vol.Name == diskName {
			action, _, err := client.StorageActions.Detach(vol.ID)
			if err != nil {
				return f.Fail("Failed to detach: ", err.Error())
			}

			if res := awaitAction(client, vol.ID, action); res != nil {
				return *res
			}
			return f.Succeed()
		}
	}
	return f.Fail("Could not find volume attached at ", device)
}

func awaitAction(client *godo.Client, volId string, action *godo.Action) *f.Result {
	start := time.Now()
	actionId := action.ID
	for {
		switch action.Status {
		case "errored":
			res := f.Fail("attach failed: ", action.String())
			return &res

		case "completed":
			return nil

		case "in-progress":

		default:
			res := f.Fail("Unknown action status ", action.Status)
			return &res
		}
		if 10*time.Minute < time.Since(start) {
			res := f.Fail("Took too long to attach volume")
			return &res
		}

		time.Sleep(250 * time.Millisecond)

		var err error
		action, _, err = client.StorageActions.Get(volId, actionId)
		if err != nil {
			res := f.Fail(err.Error())
			return &res
		}
	}
	return nil
}

func (DigitalOceanPlugin) Mount(mountDir string, device string, opt interface{}) f.Result {
	return f.Mount(mountDir, device, opt.(*DigitalOceanOptions).DefaultOptions)
}

func (DigitalOceanPlugin) Unmount(mountDir string) f.Result {
	return f.Unmount(mountDir)
}

func main() {
	f.RunPlugin(&DigitalOceanPlugin{})
}
