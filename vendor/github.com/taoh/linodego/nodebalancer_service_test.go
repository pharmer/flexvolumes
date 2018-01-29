package linodego

import (
	"testing"
	"fmt"
)

func TestListNB(t *testing.T) {
	client := NewClient("", nil)
	resp, err := client.NodeBalancer.List(0)
	fmt.Println(err)
	fmt.Println(resp.NodeBalancer)
	for _, lb := range resp.NodeBalancer {
		fmt.Println(lb.Label.String(), lb.Address4)
	}
}

func TestCreateNB(t *testing.T) {
	client := NewClient("", nil)
	args := map[string]string{}
	resp, err := client.NodeBalancer.Create(3, "aa623189901ae11e8b0473231561571e", args)
	fmt.Println(err)
	fmt.Println(resp.Errors, resp.NodeBalancerId)
}
