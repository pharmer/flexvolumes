package cloud

import (
	"errors"
	"time"
)

type DefaultOptions struct {
	ApiKey         string `json:"kubernetes.io/secret/apiKey"`
	FsType         string `json:"kubernetes.io/fsType"`
	PVorVolumeName string `json:"kubernetes.io/pvOrVolumeName"`
	RW             string `json:"kubernetes.io/readwrite"`
	VolumeName     string `json:"volumeName,omitempty"`
	VolumeID       string `json:"volumeId,omitempty"`
}

const (
	CredentialFileEnv         = "CRED_FILE_PATH"
	CredentialDefaultLocation = "/etc/kubernetes/cloud.json"
	RetryInterval             = 5 * time.Second
	RetryTimeout              = 15 * time.Minute
)

var ErrNotSupported = errors.New("Not Supported")
var ErrIncorrectArgNumber = errors.New("Incorrect number of args")

type Interface interface {
	NewOptions() interface{}
	Initialize() error

	Init() error
	Attach(options interface{}, nodeName string) (device string, err error)
	Detach(device, nodeName string) error
	Mount(mountDir string, device string, options interface{}) error
	Unmount(mountDir string) error
}
