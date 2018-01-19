package packet

import (
	"context"

	"github.com/packethost/packngo"
	. "github.com/pharmer/flexvolumes/cloud"
)

type VolumeManager struct {
	ctx    context.Context
	client *packngo.Client
}

var _ Interface = &VolumeManager{}

const (
	UID           = "packet"
	DEVICE_PREFIX = "/dev/mapper/"
)

func init() {
	RegisterCloudManager(UID, func(ctx context.Context) (Interface, error) { return New(ctx), nil })

}

func New(ctx context.Context) Interface {
	return &VolumeManager{ctx: ctx}
}
