package digitalocean

import (
	"context"

	"github.com/digitalocean/godo"
	. "github.com/pharmer/flexvolumes/cloud"
)

type VolumeManager struct {
	ctx    context.Context
	client *godo.Client
}

var _ Interface = &VolumeManager{}

const (
	UID           = "digitalocean"
	DEVICE_PREFIX = "/dev/disk/by-id/scsi-0DO_Volume_"
)

func init() {
	RegisterCloudManager(UID, func(ctx context.Context) (Interface, error) { return New(ctx), nil })

}

func New(ctx context.Context) Interface {
	return &VolumeManager{ctx: ctx}
}
