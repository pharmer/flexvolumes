package linode

import  (
	"context"
	"github.com/taoh/linodego"
	. "github.com/pharmer/flexvolumes/cloud"
)
type VolumeManager struct {
	ctx    context.Context
	client *linodego.Client
}

var _ Interface = &VolumeManager{}

const (
	UID           = "linode"
	DEVICE_PREFIX = "/dev/disk/by-id/scsi-0"
)

func init() {
	RegisterCloudManager(UID, func(ctx context.Context) (Interface, error) { return New(ctx), nil })

}

func New(ctx context.Context) Interface {
	return &VolumeManager{ctx: ctx}
}
