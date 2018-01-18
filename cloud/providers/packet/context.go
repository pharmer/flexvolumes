package main

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
	UID = "packet"
)

func init() {
	RegisterCloudManager(UID, func(ctx context.Context) (Interface, error) { return New(ctx), nil })

}

func New(ctx context.Context) Interface {
	return &VolumeManager{ctx: ctx}
}
