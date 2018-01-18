package main

import (
	"errors"
	"net"

	"github.com/digitalocean/godo"
)

func detectDroplet(client *godo.Client) (*godo.Droplet, error) {
	iface, err := net.InterfaceByIndex(2)
	if err != nil {
		return nil, err
	}

	addrs, err := iface.Addrs()
	if err != nil {
		return nil, err
	}

	var ips []string

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && ipnet.IP.To4() != nil {
			ips = append(ips, ipnet.IP.String())
		}
	}

	droplets, _, err := client.Droplets.List(nil)
	if err != nil {
		return nil, err
	}

	for _, droplet := range droplets {
		dropletIP, err := droplet.PublicIPv4()
		if err != nil {
			continue
		}

		for _, ip := range ips {
			if dropletIP == ip {
				return &droplet, nil
			}
		}
	}

	return nil, errors.New("Could not detect droplet id")
}
