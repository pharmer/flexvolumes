package linode

import (
	"fmt"
	"testing"

	"github.com/taoh/linodego"
)

func TestVol(t *testing.T) {
	client := linodego.NewClient("", nil)
	resp, err := client.Linode.List(0)
	fmt.Println(err)
	fmt.Println(resp.Linodes[1])

	vol, err := client.Volume.List(0)
	fmt.Println(err, vol.Volume, vol.Volume[0].LinodeId)

	job, err := client.Job.List(resp.Linodes[1].LinodeId, 0, true)
	fmt.Println(job.Jobs)
	//disk, err := client.Disk.List(0 , 0)

	/*r, e := client.Account.Info()
	fmt.Println(e)
	fmt.Println(string(r.Response.RawData))*/
}
