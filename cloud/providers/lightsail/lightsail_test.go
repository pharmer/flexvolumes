package lightsail

import (
	"fmt"
	"testing"

	_aws "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lightsail"
	//"strconv"
)

func TestVloume(t *testing.T) {
	tkn := &TokenSource{
		AccessKeyID:     "",
		SecretAccessKey: "",
	}
	client, err := tkn.getClient()
	fmt.Println(err)
	resp, err := client.GetDisk(&lightsail.GetDiskInput{
		DiskName: _aws.String("flextest"),
	})
	fmt.Println(err)
	fmt.Println(*resp.Disk)

	ins, err := instanceByName(client, "medium-1-0-pool-d5oi7g")

	fmt.Println(*ins, err)
	fmt.Println(*resp.Disk.AttachedTo == *ins.Name)
	/*re, err := client.AttachDisk(&lightsail.AttachDiskInput{
		DiskName: resp.Disk.Name,
		InstanceName: ins.Name,
		DiskPath: _aws.String("/dev/xvdf"),
	})
	fmt.Println(*re, err)*/

	r, e := client.DetachDisk(&lightsail.DetachDiskInput{
		DiskName: _aws.String("flextest"),
	})
	fmt.Println(*r, e)
}

func TestPath(t *testing.T) {
	deviceMappings := make(map[string]bool, 0)
	deviceMappings["/dev/xvdf"] = true
	for i := 102; i <= 122; i++ {
		path := DEVICE_PREFIX + string((i))
		if _, found := deviceMappings[path]; !found {
			fmt.Println(path)
		}
	}

}
