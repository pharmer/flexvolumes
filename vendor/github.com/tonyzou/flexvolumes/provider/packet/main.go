package main

import (
	// "github.com/packethost/packngo"
	f "github.com/tonyzou/flexvolumes"
)

type PacketOptions struct {
	f.DefaultOptions
	VolumeId string `json:"volumeId"`
}

type PacketPlugin struct{}

func (PacketPlugin) NewOptions() interface{} {
	return &PacketOptions{}
}

func (PacketPlugin) Init() f.Result {
	return f.Succeed()
}

func (PacketPlugin) Attach(opts interface{}) f.Result {
	return f.Fail("Not implemented")
}

func (PacketPlugin) Detach(device string) f.Result {
	return f.Fail("Not implemented")
}

func (PacketPlugin) Mount(mountDir string, device string, opt interface{}) f.Result {
	return f.Mount(mountDir, device, opt.(*PacketOptions).DefaultOptions)
}

func (PacketPlugin) Unmount(mountDir string) f.Result {
	return f.Unmount(mountDir)
}

func main() {
	f.RunPlugin(&PacketPlugin{})
}
