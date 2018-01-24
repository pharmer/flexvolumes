package vultr

import (
	"fmt"

	vultr "github.com/JamesClonk/vultr/lib"
)

//returns serverID of corresponding nodeName
func getServerID(client *vultr.Client, nodeName string) (string, error) {
	serverList, err := client.GetServers()
	if err != nil {
		return "", err
	}

	for _, s := range serverList {
		if s.Name == nodeName {
			return s.ID, nil
		}
	}
	return "", fmt.Errorf("Server not found for %v", nodeName)
}

//if volName blockstorage is attached to Server, it will return volume_id,true,nil
//if volName blockstorage is not attached to Server, it will return "",false,nil
//Othewise, it will return "",false,error
func getVolumeId(client *vultr.Client, volName string, serverID string) (string, bool, error) {
	volList, err := client.GetBlockStorages()
	if err != nil {
		return "", false, err
	}
	for _, s := range volList {
		if s.Name == volName {
			if s.AttachedTo == serverID {
				return s.ID, true, nil
			} else {
				return "", false, nil
			}
		}
	}
	return "", false, fmt.Errorf("Volume not found for %v", volName)
}

//ref: https://www.vultr.com/docs/block-storage
func getNextDeviceName(client *vultr.Client, serverId string) (string, error) {
	blockStorageList, err := client.GetBlockStorages()
	if err != nil {
		return "", err
	}

	var nameSuffix byte = byte('b')
	for _, b := range blockStorageList {
		if b.AttachedTo == serverId {
			nameSuffix++
		}
	}
	return DEVICE_PREFIX + string(nameSuffix), nil
}
