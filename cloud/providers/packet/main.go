package main

import (
	. "github.com/pharmer/flexvolumes/cloud"
)

type PacketOptions struct {
	DefaultOptions
}

func (v *VolumeManager) NewOptions() interface{} {
	return &PacketOptions{}
}

func (v *VolumeManager) Initialize() error {
	return ErrNotSupported
}

func (v *VolumeManager) Init() error {
	return nil
}

func (v *VolumeManager) Attach(options interface{}, nodeName string) (string, error) {
	return "", ErrNotSupported
}

func (v *VolumeManager) Detach(device, nodeName string) error {
	return ErrNotSupported
}

func (v *VolumeManager) Mount(mountDir string, device string, options interface{}) error {
	return ErrNotSupported
}

func (v *VolumeManager) Unmount(mountDir string) error {
	return Unmount(mountDir)
}
