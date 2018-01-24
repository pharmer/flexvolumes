package vultr

import (
	"context"

	vultr "github.com/JamesClonk/vultr/lib"
	. "github.com/pharmer/flexvolumes/cloud"
)

type VolumeManager struct {
	ctx    context.Context
	client *vultr.Client
}

var _ Interface = &VolumeManager{}

const (
	UID           = "vultr"
	DEVICE_PREFIX = "/dev/vd"
)

func init() {
	RegisterCloudManager(UID, func(ctx context.Context) (Interface, error) { return New(ctx), nil })
}

func New(ctx context.Context) Interface {
	return &VolumeManager{ctx: ctx}
}
