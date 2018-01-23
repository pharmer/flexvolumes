package lightsail


import (
	"context"


	. "github.com/pharmer/flexvolumes/cloud"
	"github.com/aws/aws-sdk-go/service/lightsail"
)

type VolumeManager struct {
	ctx    context.Context
	client *lightsail.Lightsail
}

var _ Interface = &VolumeManager{}

const (
	UID           = "lightsail"
	DEVICE_PREFIX = "/dev/disk/by-id/scsi-0DO_Volume_"
	metadataURL = "http://169.254.169.254/latest/meta-data/"
)

func init() {
	RegisterCloudManager(UID, func(ctx context.Context) (Interface, error) { return New(ctx), nil })

}

func New(ctx context.Context) Interface {
	return &VolumeManager{ctx: ctx}
}

