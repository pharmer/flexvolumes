package linode

import (
	"fmt"

	. "github.com/pharmer/flexvolumes/cloud"
	"github.com/taoh/linodego"
	"k8s.io/apimachinery/pkg/util/wait"
)

func getVolumeFromName(client *linodego.Client, volName string) (*linodego.Volume, error) {
	volumes, err := client.Volume.List(0)
	if err != nil {
		return nil, err
	}
	for _, v := range volumes.Volume {
		if v.Label.String() == volName {
			return &v, nil
		}
	}
	return nil, fmt.Errorf("no volume found with %v volName", volName)
}

func linodeByName(client *linodego.Client, nodeName string) (*linodego.Linode, error) {
	linodes, err := client.Linode.List(0)
	if err != nil {
		return nil, err
	}

	for _, linode := range linodes.Linodes {
		if linode.Label.String() == string(nodeName) {
			return &linode, nil
		}
	}
	return nil, fmt.Errorf("no linode found with name %v", nodeName)
}

func awaitAction(client *linodego.Client, id int) error {
	attempt := 0
	return wait.PollImmediate(RetryInterval, RetryTimeout, func() (bool, error) {
		attempt++

		resp, err := client.Job.List(id, 0, true)
		if err != nil {
			return false, nil
		}
		if len(resp.Jobs) == 0 {
			return true, nil
		}
		return false, nil
	})
}
