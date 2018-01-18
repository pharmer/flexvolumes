package digitalocean

import (
	"fmt"

	"github.com/digitalocean/godo"
	. "github.com/pharmer/flexvolumes/cloud"
	"golang.org/x/oauth2"
	"k8s.io/apimachinery/pkg/util/wait"
)

func getDroplet(client *godo.Client, nodeName string) (*godo.Droplet, error) {
	droplets, _, err := client.Droplets.List(oauth2.NoContext, &godo.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, droplet := range droplets {
		if droplet.Name == nodeName {
			d, _, err := client.Droplets.Get(oauth2.NoContext, droplet.ID)
			if err != nil {
				return nil, err
			}
			return d, nil
		}
	}

	return nil, fmt.Errorf("no droplet found with %v name", nodeName)
}

func awaitAction(client *godo.Client, volId string, action *godo.Action) error {
	actionId := action.ID
	return wait.PollImmediate(RetryInterval, RetryTimeout, func() (bool, error) {
		switch action.Status {
		case "errored":
			return false, fmt.Errorf("attach failed: %v", action.String())

		case godo.ActionCompleted:
			return true, nil

		case godo.ActionInProgress:

		default:
			return false, fmt.Errorf("unknown action status %v", action.Status)
		}
		var err error
		action, _, err = client.StorageActions.Get(oauth2.NoContext, volId, actionId)
		if err != nil {
			return false, err
		}
		return false, nil
	})
}
