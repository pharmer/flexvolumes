package vultr

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	vultr "github.com/JamesClonk/vultr/lib"
	"github.com/pharmer/flexvolumes/cloud"
)

//adding volume to node for first time
func TestFirstVolumeAttach(t *testing.T) {
	v := VolumeManager{
		ctx:    context.Background(),
		client: vultr.NewClient(tGetToken(), nil),
	}
	opts := cloud.DefaultOptions{
		VolumeName: "flextest",
		VolumeID:   "13105873",
	}
	res, err := v.Attach(&VultrOptions{
		opts,
	}, "94-pool-omtfjj")
	if err != nil {
		t.Error(err)
	}

	if res != "/dev/vdb" {
		t.Error("Expected /dev/vdb, found %v", res)
	}
}

func TestVolumeDetach(t *testing.T) {
	v := VolumeManager{
		ctx:    context.Background(),
		client: vultr.NewClient(tGetToken(), nil),
	}
	opts := cloud.DefaultOptions{
		VolumeName: "flextest",
		VolumeID:   "13105873",
	}
	err := v.Detach(opts.VolumeName, "94-pool-omtfjj")
	if err != nil {
		t.Error(err)
	}
}

func TestNextDeviceName(t *testing.T) {
	v := VolumeManager{
		ctx:    context.Background(),
		client: vultr.NewClient(tGetToken(), nil),
	}

	serverId, err := getServerID(v.client, "94-pool-omtfjj")
	if err != nil {
		t.Error(err)
	}

	name, err := getNextDeviceName(v.client, serverId)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(name)
}

func TestVolumeResp(t *testing.T) {
	v := VolumeManager{
		ctx:    context.Background(),
		client: vultr.NewClient(tGetToken(), nil),
	}
	blockStorageList, err := v.client.GetBlockStorages()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(blockStorageList)
}

func tGetToken() string {
	b, _ := tReadFile("/home/ac/Downloads/cred/vultr.json")
	v := struct {
		Token string `json:"token"`
	}{}
	fmt.Println(json.Unmarshal(b, &v))
	//fmt.Println(v)
	return v.Token
}

func tReadFile(name string) ([]byte, error) {
	crtBytes, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, fmt.Errorf("failed to read `%s`.Reason: %v", name, err)
	}
	return crtBytes, nil
}

