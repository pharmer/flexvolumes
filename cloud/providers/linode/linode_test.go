package linode

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"strings"
	"testing"

	"github.com/pharmer/flexvolumes/cloud"
	"github.com/taoh/linodego"
	"golang.org/x/sys/unix"
)

type fakeLinodeVolumeService struct {
	cloneFn  func(int, string) (*linodego.LinodeVolumeResponse, error)
	createFn func(int, string, map[string]string) (*linodego.LinodeVolumeResponse, error)
	deleteFn func(int) (*linodego.LinodeVolumeResponse, error)
	listFn   func(int) (*linodego.LinodeVolumeListResponse, error)
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

func newFakeVolume() linodego.Volume {
	label := linodego.CustomString{}
	err := label.UnmarshalJSON([]byte("test-volume"))
	if err != nil {
		return linodego.Volume{}
	}
	return linodego.Volume{
		VolumeId: 987,
		Label:    label,
		Size:     50,
		Status:   "active",
	}
}

func Test_Attach(t *testing.T) {
	listVFn := func(i int) (*linodego.LinodeVolumeListResponse, error) {
		volume := newFakeVolume()
		if i != 0 && volume.VolumeId != i {
			return &linodego.LinodeVolumeListResponse{
				newFakeOKResponse("volume.list"),
				[]linodego.Volume{},
			}, nil
		}
		return &linodego.LinodeVolumeListResponse{
			newFakeOKResponse("volume.list"),
			[]linodego.Volume{volume},
		}, nil
	}
	updateVFn := func(i int, strings map[string]string) (*linodego.LinodeVolumeResponse, error) {
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

	listLFunc := func(i int) (*linodego.LinodesListResponse, error) {
		linode := newFakeLinode()
		linodes := []linodego.Linode{linode}
		return &linodego.LinodesListResponse{
			newFakeOKResponse("linode.list"),
			linodes,
		}, nil
	}

	listJFn := func(i int, i2 int, b bool) (*linodego.LinodesJobListResponse, error) {
		return &linodego.LinodesJobListResponse{
			newFakeOKResponse("linode.job.list"),
			[]linodego.Job{},
		}, nil
	}

	testcases := []struct {
		name      string
		listVFn   func(int) (*linodego.LinodeVolumeListResponse, error)
		updateVFn func(int, map[string]string) (*linodego.LinodeVolumeResponse, error)
		listLFn   func(int) (*linodego.LinodesListResponse, error)
		listJFn   func(int, int, bool) (*linodego.LinodesJobListResponse, error)
		options   *LinodeOptions
		device    string
		err       error
	}{
		{
			"volume attach with id",
			listVFn,
			updateVFn,
			listLFunc,
			listJFn,
			&LinodeOptions{
				DefaultOptions: cloud.DefaultOptions{
					VolumeID: "987",
				},
			},
			DEVICE_PREFIX,
			nil,
		},
		{
			"volume attach with name",
			listVFn,
			updateVFn,
			listLFunc,
			listJFn,
			&LinodeOptions{
				DefaultOptions: cloud.DefaultOptions{
					VolumeName: "test-volume",
				},
			},
			DEVICE_PREFIX,
			nil,
		},
		{
			"volume not exists",
			listVFn,
			updateVFn,
			listLFunc,
			listJFn,
			&LinodeOptions{
				DefaultOptions: cloud.DefaultOptions{
					VolumeID: "1111",
				},
			},
			"",
			fmt.Errorf("no volume found with id %v", 1111),
		},
	}

	for _, test := range testcases {
		t.Run(test.name, func(t *testing.T) {
			fakeV := &fakeLinodeVolumeService{
				listFn:   test.listVFn,
				updateFn: test.updateVFn,
			}
			fakeL := &fakeLinodeService{
				listFunc: test.listLFn,
			}
			fakeJ := &fakeLinodeJobService{
				listFn: test.listJFn,
			}
			fakeClient := newFakeClient(fakeV, fakeL, fakeJ)
			v := &VolumeManager{ctx: context.Background(), client: fakeClient}

			device, err := v.Attach(test.options, "test-linode")
			if !reflect.DeepEqual(err, test.err) {
				t.Errorf("unexpected err, expected nil. got: %v", err)
			}
			if !strings.HasPrefix(device, test.device) {
				t.Errorf("unexpected device prefix, expected %s, got %s", DEVICE_PREFIX, device)
			}
		})
	}

}

func Test_Detach(t *testing.T) {
	listVFn := func(i int) (*linodego.LinodeVolumeListResponse, error) {
		volume := newFakeVolume()
		if i != 0 && volume.VolumeId != i {
			return &linodego.LinodeVolumeListResponse{
				newFakeOKResponse("volume.list"),
				[]linodego.Volume{},
			}, nil
		}
		volume.LinodeId = 1234
		return &linodego.LinodeVolumeListResponse{
			newFakeOKResponse("volume.list"),
			[]linodego.Volume{volume},
		}, nil
	}
	updateVFn := func(i int, strings map[string]string) (*linodego.LinodeVolumeResponse, error) {
		vol := newFakeVolume()
		if i != vol.VolumeId {
			return nil, fmt.Errorf("volume not found")
		}
		return &linodego.LinodeVolumeResponse{
			newFakeOKResponse("volume.update"),
			linodego.VolumeId{i},
		}, nil
	}
	listLFunc := func(i int) (*linodego.LinodesListResponse, error) {
		linode := newFakeLinode()
		if i != 0 && linode.LinodeId != i {
			return &linodego.LinodesListResponse{
				newFakeOKResponse("linode.list"),
				[]linodego.Linode{},
			}, nil
		}
		linodes := []linodego.Linode{linode}
		return &linodego.LinodesListResponse{
			newFakeOKResponse("linode.list"),
			linodes,
		}, nil
	}
	listJFn := func(i int, i2 int, b bool) (*linodego.LinodesJobListResponse, error) {
		return &linodego.LinodesJobListResponse{
			newFakeOKResponse("linode.job.list"),
			[]linodego.Job{},
		}, nil
	}
	testcases := []struct {
		name      string
		listVFn   func(int) (*linodego.LinodeVolumeListResponse, error)
		updateVFn func(int, map[string]string) (*linodego.LinodeVolumeResponse, error)
		listLFn   func(int) (*linodego.LinodesListResponse, error)
		listJFn   func(int, int, bool) (*linodego.LinodesJobListResponse, error)
		device    string
		nodeName  string
		err       error
	}{
		{
			"volume detach",
			listVFn,
			updateVFn,
			listLFunc,
			listJFn,
			"test-volume",
			"test-linode",
			nil,
		},
		{
			"volume not attached",
			func(i int) (*linodego.LinodeVolumeListResponse, error) {
				volume := newFakeVolume()
				return &linodego.LinodeVolumeListResponse{
					newFakeOKResponse("volume.list"),
					[]linodego.Volume{volume},
				}, nil
			},
			updateVFn,
			listLFunc,
			listJFn,
			"test-volume",
			"test-linode",
			fmt.Errorf("could not find volume attached at %v", "test-volume"),
		},
	}
	for _, test := range testcases {
		t.Run(test.name, func(t *testing.T) {
			fakeV := &fakeLinodeVolumeService{
				listFn:   test.listVFn,
				updateFn: test.updateVFn,
			}
			fakeL := &fakeLinodeService{
				listFunc: test.listLFn,
			}
			fakeJ := &fakeLinodeJobService{
				listFn: test.listJFn,
			}
			fakeClient := newFakeClient(fakeV, fakeL, fakeJ)
			v := &VolumeManager{ctx: context.Background(), client: fakeClient}

			err := v.Detach(test.device, test.nodeName)
			if !reflect.DeepEqual(err, test.err) {
				t.Errorf("unexpected err, expected nil. got: %v", err)
			}
		})
	}

}

// https://npf.io/2015/06/testing-exec-command/
func fakeExecCommand(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	//testing: warning: no tests to run
	return cmd
}

func fakeUnixStat(path string, stat *unix.Stat_t) error {
	stat.Mode = unix.S_IFBLK
	return nil
}

func TestMount(t *testing.T) {
	opt := &LinodeOptions{
		DefaultOptions: cloud.DefaultOptions{
			VolumeID:   "987",
			FsType:     "ext4",
			RW:         "rw",
			VolumeName: "test-volume",
		},
	}
	testcases := []struct {
		name     string
		options  *LinodeOptions
		mountDir string
		device   string
		err      error
	}{
		{
			"mount device",
			opt,
			"/tmp/mount",
			"test-volume",
			nil,
		},
		{
			"fs type not specified",
			&LinodeOptions{
				DefaultOptions: cloud.DefaultOptions{
					VolumeID:   "987",
					RW:         "rw",
					VolumeName: "test-volume",
				},
			},
			"/tmp/mount",
			"test-volume",
			fmt.Errorf("No filesystem type specified"),
		},
	}

	for _, test := range testcases {
		t.Run(test.name, func(t *testing.T) {
			cloud.ExecCommand = fakeExecCommand
			cloud.UnixStat = fakeUnixStat
			fakeClient := newFakeClient(nil, nil, nil)
			v := &VolumeManager{ctx: context.Background(), client: fakeClient}

			err := v.MountDevice(test.mountDir, test.device, test.options)
			if !reflect.DeepEqual(err, test.err) {
				t.Errorf("unexpected err, expected %v. got: %v", test.err, err)
			}
		})
	}

}

func Test_Unmount(t *testing.T) {
	testcases := []struct {
		name     string
		mountDir string
		err      error
	}{
		{
			"mount device",
			"/tmp/mount",
			nil,
		},
	}

	for _, test := range testcases {
		t.Run(test.name, func(t *testing.T) {
			cloud.ExecCommand = fakeExecCommand
			fakeClient := newFakeClient(nil, nil, nil)
			v := &VolumeManager{ctx: context.Background(), client: fakeClient}

			err := v.Unmount(test.mountDir)
			if !reflect.DeepEqual(err, test.err) {
				t.Errorf("unexpected err, expected nil. got: %v", err)
			}
		})
	}
}
