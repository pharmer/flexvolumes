package linode

import (
	"fmt"
	"testing"

	"github.com/taoh/linodego"
	"context"
	"github.com/pharmer/flexvolumes/cloud"
	"strings"
)

type fakeLinodeVolumeService struct {
	cloneFn func(int, string) (*linodego.LinodeVolumeResponse, error)
	createFn func(int, string, map[string]string) (*linodego.LinodeVolumeResponse, error)
	deleteFn func(int) (*linodego.LinodeVolumeResponse, error)
	listFn func(int) (*linodego.LinodeVolumeListResponse, error)
	updateFn func(int, map[string]string) (*linodego.LinodeVolumeResponse, error)
}

func (f *fakeLinodeVolumeService) List(volumeId int) (*linodego.LinodeVolumeListResponse, error) {
	return f.listFn(volumeId)
}

func (f *fakeLinodeVolumeService) Create(size int, label string, args map[string]string) (*linodego.LinodeVolumeResponse, error) {
	return f.createFn(size, label, args)
}

func (f *fakeLinodeVolumeService) Update(volumeId int, args map[string]string) (*linodego.LinodeVolumeResponse, error) {
	return f.updateFn(volumeId, args)
}

func (f *fakeLinodeVolumeService) Delete(volumeId int) (*linodego.LinodeVolumeResponse, error) {
	return f.deleteFn(volumeId)
}

func (f *fakeLinodeVolumeService) Clone(cloneFromId int, label string) (*linodego.LinodeVolumeResponse, error) {
	return f.cloneFn(cloneFromId, label)
}

type fakeLinodeService struct {
	bootFunc     func(int, int) (*linodego.JobResponse, error)
	cloneFunc    func(int, int, int, int) (*linodego.LinodeResponse, error)
	createFunc   func(int, int, int) (*linodego.LinodeResponse, error)
	deleteFunc   func(int, bool) (*linodego.LinodeResponse, error)
	listFunc     func(int) (*linodego.LinodesListResponse, error)
	rebootFunc   func(int, int) (*linodego.JobResponse, error)
	resizeFunc   func(int, int) (*linodego.LinodeResponse, error)
	shutdownFunc func(int) (*linodego.JobResponse, error)
	updateFunc   func(int, map[string]interface{}) (*linodego.LinodeResponse, error)
}

func (f *fakeLinodeService) Boot(linodeId int, configId int) (*linodego.JobResponse, error) {
	return f.bootFunc(linodeId, configId)
}

func (f *fakeLinodeService) Clone(linodeId int, dataCenterId int, planId int, paymentTerm int) (*linodego.LinodeResponse, error) {
	return f.cloneFunc(linodeId, dataCenterId, planId, paymentTerm)
}

func (f *fakeLinodeService) Create(dataCenterId int, planId int, paymentTerm int) (*linodego.LinodeResponse, error) {
	return f.createFunc(dataCenterId, planId, paymentTerm)
}

func (f *fakeLinodeService) Delete(linodeId int, skipChecks bool) (*linodego.LinodeResponse, error) {
	return f.deleteFunc(linodeId, skipChecks)
}

func (f *fakeLinodeService) List(linodeId int) (*linodego.LinodesListResponse, error) {
	return f.listFunc(linodeId)
}

func (f *fakeLinodeService) Reboot(linodeId int, configId int) (*linodego.JobResponse, error) {
	return f.rebootFunc(linodeId, configId)
}

func (f *fakeLinodeService) Resize(linodeId int, planId int) (*linodego.LinodeResponse, error) {
	return f.resizeFunc(linodeId, planId)
}

func (f *fakeLinodeService) Shutdown(linodeId int) (*linodego.JobResponse, error) {
	return f.shutdownFunc(linodeId)
}

func (f *fakeLinodeService) Update(linodeId int, args map[string]interface{}) (*linodego.LinodeResponse, error) {
	return f.updateFunc(linodeId, args)
}

type fakeLinodeJobService struct {
	listFn func(int, int, bool) (*linodego.LinodesJobListResponse, error)
}

func (f *fakeLinodeJobService) List(linodeId int, jobId int, pendingOnly bool) (*linodego.LinodesJobListResponse, error) {
	return f.listFn(linodeId, jobId, pendingOnly)
}

func newFakeClient(fakeVolume *fakeLinodeVolumeService, fakeLinode *fakeLinodeService, fakeJob *fakeLinodeJobService) *linodego.Client {
	client := linodego.NewClient("", nil)
	client.Volume = fakeVolume
	client.Linode = fakeLinode
	client.Job = fakeJob
	return client
}

func newFakeOKResponse(action string) linodego.Response {
	return linodego.Response{
		Errors: nil,
		Action: action,
	}
}
func newFakeLinode() linodego.Linode {
	label := linodego.CustomString{}
	err := label.UnmarshalJSON([]byte("test-linode"))
	if err != nil {
		return linodego.Linode{}
	}
	return linodego.Linode{
		Label:    label,
		LinodeId: 1234,
		PlanId:   2,
	}
}

func newFakeVolume() linodego.Volume  {
	label := linodego.CustomString{}
	err := label.UnmarshalJSON([]byte("test-volume"))
	if err != nil {
		return linodego.Volume{}
	}
	return linodego.Volume{
		VolumeId: 987,
		Label: label,
		Size: 50,
		Status: "active",
	}
}

func Test_Attach(t *testing.T) {
	fakeVol := &fakeLinodeVolumeService{}
	fakeVol.listFn = func(i int) (*linodego.LinodeVolumeListResponse, error) {
		volume := newFakeVolume()
		return &linodego.LinodeVolumeListResponse{
			newFakeOKResponse("volume.list"),
			[]linodego.Volume{volume},
		}, nil
	}
	fakeVol.updateFn = func(i int, strings map[string]string) (*linodego.LinodeVolumeResponse, error) {
		vol := newFakeVolume()
		linode := strings["LinodeID"]
		if linode != "1234" {
			return nil, fmt.Errorf("linode not found")
		}
		if i != vol.VolumeId {
			return nil, fmt.Errorf("volume not found")
		}
		return &linodego.LinodeVolumeResponse{
			newFakeOKResponse("volume.update"),
			linodego.VolumeId{i},
		}, nil
	}

	fakeLinode := &fakeLinodeService{}
	fakeLinode.listFunc = func(i int) (*linodego.LinodesListResponse, error) {
		linode := newFakeLinode()
		linodes := []linodego.Linode{linode}
		return &linodego.LinodesListResponse{
			newFakeOKResponse("linode.list"),
			linodes,
		}, nil
	}

	fakeJob := &fakeLinodeJobService{}
	fakeJob.listFn = func(i int, i2 int, b bool) (*linodego.LinodesJobListResponse, error) {
		return &linodego.LinodesJobListResponse{
			newFakeOKResponse("linode.job.list"),
			[]linodego.Job{},
		}, nil
	}

	fakeClient := newFakeClient(fakeVol, fakeLinode, fakeJob)
	v := &VolumeManager{ctx:context.Background(), client:fakeClient}

	options := &LinodeOptions{
		DefaultOptions: cloud.DefaultOptions{
			VolumeID: "987",
		},
	}
	device, err := v.Attach(options, "test-linode")
	if err != nil {
		t.Errorf("unexpected err, expected nil. got: %v", err)
	}
	if !strings.HasPrefix(device, DEVICE_PREFIX) {
		t.Errorf("unexpected device prefix, expected %s, got %s", DEVICE_PREFIX, device)
	}


}
