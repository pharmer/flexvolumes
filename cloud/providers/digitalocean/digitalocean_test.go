package digitalocean

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/digitalocean/godo"
	"github.com/pharmer/flexvolumes/cloud"
	"golang.org/x/sys/unix"
)

type fakeDropletService struct {
	listFunc           func(ctx context.Context, opt *godo.ListOptions) ([]godo.Droplet, *godo.Response, error)
	listByTagFunc      func(ctx context.Context, tag string, opt *godo.ListOptions) ([]godo.Droplet, *godo.Response, error)
	getFunc            func(ctx context.Context, dropletID int) (*godo.Droplet, *godo.Response, error)
	createFunc         func(ctx context.Context, createRequest *godo.DropletCreateRequest) (*godo.Droplet, *godo.Response, error)
	createMultipleFunc func(ctx context.Context, createRequest *godo.DropletMultiCreateRequest) ([]godo.Droplet, *godo.Response, error)
	deleteFunc         func(ctx context.Context, dropletID int) (*godo.Response, error)
	deleteByTagFunc    func(ctx context.Context, tag string) (*godo.Response, error)
	kernelsFunc        func(ctx context.Context, dropletID int, opt *godo.ListOptions) ([]godo.Kernel, *godo.Response, error)
	snapshotsFunc      func(ctx context.Context, dropletID int, opt *godo.ListOptions) ([]godo.Image, *godo.Response, error)
	backupsFunc        func(ctx context.Context, dropletID int, opt *godo.ListOptions) ([]godo.Image, *godo.Response, error)
	actionsFunc        func(ctx context.Context, dropletID int, opt *godo.ListOptions) ([]godo.Action, *godo.Response, error)
	neighborsFunc      func(cxt context.Context, dropletID int) ([]godo.Droplet, *godo.Response, error)
}

func (f *fakeDropletService) List(ctx context.Context, opt *godo.ListOptions) ([]godo.Droplet, *godo.Response, error) {
	return f.listFunc(ctx, opt)
}

func (f *fakeDropletService) ListByTag(ctx context.Context, tag string, opt *godo.ListOptions) ([]godo.Droplet, *godo.Response, error) {
	return f.listByTagFunc(ctx, tag, opt)
}

func (f *fakeDropletService) Get(ctx context.Context, dropletID int) (*godo.Droplet, *godo.Response, error) {
	return f.getFunc(ctx, dropletID)
}

func (f *fakeDropletService) Create(ctx context.Context, createRequest *godo.DropletCreateRequest) (*godo.Droplet, *godo.Response, error) {
	return f.createFunc(ctx, createRequest)
}

func (f *fakeDropletService) CreateMultiple(ctx context.Context, createRequest *godo.DropletMultiCreateRequest) ([]godo.Droplet, *godo.Response, error) {
	return f.createMultipleFunc(ctx, createRequest)
}

func (f *fakeDropletService) Delete(ctx context.Context, dropletID int) (*godo.Response, error) {
	return f.deleteFunc(ctx, dropletID)
}

func (f *fakeDropletService) DeleteByTag(ctx context.Context, tag string) (*godo.Response, error) {
	return f.deleteByTagFunc(ctx, tag)
}

func (f *fakeDropletService) Kernels(ctx context.Context, dropletID int, opt *godo.ListOptions) ([]godo.Kernel, *godo.Response, error) {
	return f.kernelsFunc(ctx, dropletID, opt)
}

func (f *fakeDropletService) Snapshots(ctx context.Context, dropletID int, opt *godo.ListOptions) ([]godo.Image, *godo.Response, error) {
	return f.snapshotsFunc(ctx, dropletID, opt)
}

func (f *fakeDropletService) Backups(ctx context.Context, dropletID int, opt *godo.ListOptions) ([]godo.Image, *godo.Response, error) {
	return f.backupsFunc(ctx, dropletID, opt)
}

func (f *fakeDropletService) Actions(ctx context.Context, dropletID int, opt *godo.ListOptions) ([]godo.Action, *godo.Response, error) {
	return f.actionsFunc(ctx, dropletID, opt)
}

func (f *fakeDropletService) Neighbors(ctx context.Context, dropletID int) ([]godo.Droplet, *godo.Response, error) {
	return f.neighborsFunc(ctx, dropletID)
}

type fakeStorageService struct {
	listVolumesFn    func(context.Context, *godo.ListVolumeParams) ([]godo.Volume, *godo.Response, error)
	getVolumeFn      func(context.Context, string) (*godo.Volume, *godo.Response, error)
	createVolumeFn   func(context.Context, *godo.VolumeCreateRequest) (*godo.Volume, *godo.Response, error)
	deleteVolumeFn   func(context.Context, string) (*godo.Response, error)
	listSnapshotsFn  func(ctx context.Context, volumeID string, opts *godo.ListOptions) ([]godo.Snapshot, *godo.Response, error)
	getSnapshotFn    func(context.Context, string) (*godo.Snapshot, *godo.Response, error)
	createSnapshotFn func(context.Context, *godo.SnapshotCreateRequest) (*godo.Snapshot, *godo.Response, error)
	deleteSnapshotFn func(context.Context, string) (*godo.Response, error)
}

func (f *fakeStorageService) ListVolumes(ctx context.Context, params *godo.ListVolumeParams) ([]godo.Volume, *godo.Response, error) {
	return f.listVolumesFn(ctx, params)
}

func (f *fakeStorageService) CreateVolume(ctx context.Context, createRequest *godo.VolumeCreateRequest) (*godo.Volume, *godo.Response, error) {
	return f.createVolumeFn(ctx, createRequest)
}

func (f *fakeStorageService) GetVolume(ctx context.Context, id string) (*godo.Volume, *godo.Response, error) {
	return f.getVolumeFn(ctx, id)
}

func (f *fakeStorageService) DeleteVolume(ctx context.Context, id string) (*godo.Response, error) {
	return f.deleteVolumeFn(ctx, id)
}

func (f *fakeStorageService) ListSnapshots(ctx context.Context, volumeID string, opt *godo.ListOptions) ([]godo.Snapshot, *godo.Response, error) {
	return f.listSnapshotsFn(ctx, volumeID, opt)
}

func (f *fakeStorageService) CreateSnapshot(ctx context.Context, createRequest *godo.SnapshotCreateRequest) (*godo.Snapshot, *godo.Response, error) {
	return f.createSnapshotFn(ctx, createRequest)
}

func (f *fakeStorageService) GetSnapshot(ctx context.Context, id string) (*godo.Snapshot, *godo.Response, error) {
	return f.getSnapshotFn(ctx, id)
}

func (f *fakeStorageService) DeleteSnapshot(ctx context.Context, id string) (*godo.Response, error) {
	return f.deleteSnapshotFn(ctx, id)
}

type fakeStorageActionsService struct {
	attachFn            func(ctx context.Context, volumeID string, dropletID int) (*godo.Action, *godo.Response, error)
	detachByDropletIDFn func(ctx context.Context, volumeID string, dropletID int) (*godo.Action, *godo.Response, error)
	getFn               func(ctx context.Context, volumeID string, actionID int) (*godo.Action, *godo.Response, error)
	listFn              func(ctx context.Context, volumeID string, opt *godo.ListOptions) ([]godo.Action, *godo.Response, error)
	resizeFn            func(ctx context.Context, volumeID string, sizeGigabytes int, regionSlug string) (*godo.Action, *godo.Response, error)
}

func (f *fakeStorageActionsService) Attach(ctx context.Context, volumeID string, dropletID int) (*godo.Action, *godo.Response, error) {
	return f.attachFn(ctx, volumeID, dropletID)
}

func (f *fakeStorageActionsService) DetachByDropletID(ctx context.Context, volumeID string, dropletID int) (*godo.Action, *godo.Response, error) {
	return f.detachByDropletIDFn(ctx, volumeID, dropletID)
}

func (f *fakeStorageActionsService) Get(ctx context.Context, volumeID string, actionID int) (*godo.Action, *godo.Response, error) {
	return f.getFn(ctx, volumeID, actionID)
}

func (f *fakeStorageActionsService) List(ctx context.Context, volumeID string, opt *godo.ListOptions) ([]godo.Action, *godo.Response, error) {
	return f.listFn(ctx, volumeID, opt)
}

func (f *fakeStorageActionsService) Resize(ctx context.Context, volumeID string, sizeGigabytes int, regionSlug string) (*godo.Action, *godo.Response, error) {
	return f.resizeFn(ctx, volumeID, sizeGigabytes, regionSlug)
}
func newFakeClient(fakeStorage *fakeStorageService, fakeStorageAction *fakeStorageActionsService, fakeDroplet *fakeDropletService) *godo.Client {
	client := godo.NewClient(nil)
	client.Droplets = fakeDroplet
	client.Storage = fakeStorage
	client.StorageActions = fakeStorageAction
	return client
}

func newFakeOKResponse() *godo.Response {
	return &godo.Response{
		Response: &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBufferString("test")),
		},
	}
}

func newFakeNotOKResponse() *godo.Response {
	return &godo.Response{
		Response: &http.Response{
			StatusCode: http.StatusInternalServerError,
			Body:       ioutil.NopCloser(bytes.NewBufferString("test")),
		},
	}
}

func newFakeDroplet() *godo.Droplet {
	return &godo.Droplet{
		ID:       123,
		Name:     "test-droplet",
		SizeSlug: "2gb",
		Networks: &godo.Networks{
			V4: []godo.NetworkV4{
				{
					IPAddress: "10.0.0.0",
					Type:      "private",
				},
				{
					IPAddress: "99.99.99.99",
					Type:      "public",
				},
			},
		},
		Region: &godo.Region{
			Name: "test-region",
			Slug: "test1",
		},
	}
}

func newFakeVolume() *godo.Volume {
	return &godo.Volume{
		Region:        &godo.Region{Slug: "nyc3"},
		ID:            "80d414c6-295e-4e3a-ac58-eb9456c1e1d1",
		Name:          "my volume",
		Description:   "my description",
		SizeGigaBytes: 100,
		DropletIDs:    []int{10},
		CreatedAt:     time.Date(2002, 10, 02, 15, 00, 00, 50000000, time.UTC),
	}
}

func Test_Attach(t *testing.T) {
	getVolumeFn := func(ctx context.Context, volumeId string) (*godo.Volume, *godo.Response, error) {
		vol := newFakeVolume()
		if vol.ID == volumeId {
			return vol, newFakeOKResponse(), nil
		}
		return &godo.Volume{}, newFakeNotOKResponse(), fmt.Errorf("volume with id %v not exists", volumeId)
	}

	listDropletFn := func(ctx context.Context, opt *godo.ListOptions) ([]godo.Droplet, *godo.Response, error) {
		droplet := newFakeDroplet()
		droplets := []godo.Droplet{*droplet}

		resp := newFakeOKResponse()
		return droplets, resp, nil
	}
	getDropletFunc := func(ctx context.Context, dropletID int) (*godo.Droplet, *godo.Response, error) {
		droplet := newFakeDroplet()
		if droplet.ID == dropletID {
			resp := newFakeOKResponse()
			return droplet, resp, nil
		}
		return nil, newFakeNotOKResponse(), fmt.Errorf("no droplet found with id %v", dropletID)
	}
	attachStorageActions := func(ctx context.Context, volumeID string, dropletID int) (*godo.Action, *godo.Response, error) {
		return &godo.Action{
			Status: godo.ActionCompleted,
		}, newFakeOKResponse(), nil
	}
	getStorageActions := func(ctx context.Context, volumeID string, actionID int) (*godo.Action, *godo.Response, error) {
		return &godo.Action{
			Status: godo.ActionCompleted,
		}, newFakeOKResponse(), nil
	}

	testcases := []struct {
		name        string
		getVolumeFn func(context.Context, string) (*godo.Volume, *godo.Response, error)
		listDFunc   func(ctx context.Context, opt *godo.ListOptions) ([]godo.Droplet, *godo.Response, error)
		getDFunc    func(ctx context.Context, dropletID int) (*godo.Droplet, *godo.Response, error)
		attachFn    func(ctx context.Context, volumeID string, dropletID int) (*godo.Action, *godo.Response, error)
		getSAFn     func(ctx context.Context, volumeID string, actionID int) (*godo.Action, *godo.Response, error)
		options     *DigitalOceanOptions
		nodeName    string
		device      string
		err         error
	}{
		{
			"volume attach",
			getVolumeFn,
			listDropletFn,
			getDropletFunc,
			attachStorageActions,
			getStorageActions,
			&DigitalOceanOptions{
				DefaultOptions: cloud.DefaultOptions{
					VolumeID: "80d414c6-295e-4e3a-ac58-eb9456c1e1d1",
				},
			},
			"test-droplet",
			DEVICE_PREFIX,
			nil,
		},
		{
			"volume does not exist",
			getVolumeFn,
			listDropletFn,
			getDropletFunc,
			attachStorageActions,
			getStorageActions,
			&DigitalOceanOptions{
				DefaultOptions: cloud.DefaultOptions{
					VolumeID: "test-volume",
				},
			},
			"test-droplet",
			"",
			fmt.Errorf("volume with id test-volume not exists"),
		},
		{
			"node does not exist",
			getVolumeFn,
			listDropletFn,
			getDropletFunc,
			attachStorageActions,
			getStorageActions,
			&DigitalOceanOptions{
				DefaultOptions: cloud.DefaultOptions{
					VolumeID: "80d414c6-295e-4e3a-ac58-eb9456c1e1d1",
				},
			},
			"node-1",
			"",
			fmt.Errorf("no droplet found with node-1 name"),
		},
		{
			"volume attach fail",
			getVolumeFn,
			listDropletFn,
			getDropletFunc,
			func(ctx context.Context, volumeID string, dropletID int) (*godo.Action, *godo.Response, error) {
				return &godo.Action{
					Status: godo.ActionInProgress,
				}, newFakeOKResponse(), nil
			},
			func(ctx context.Context, volumeID string, actionID int) (*godo.Action, *godo.Response, error) {
				return &godo.Action{
					Status: "errored",
				}, newFakeOKResponse(), nil
			},
			&DigitalOceanOptions{
				DefaultOptions: cloud.DefaultOptions{
					VolumeID: "80d414c6-295e-4e3a-ac58-eb9456c1e1d1",
				},
			},
			"test-droplet",
			"",
			fmt.Errorf(`attach failed: godo.Action{ID:0, Status:"errored", Type:"", ResourceID:0, ResourceType:"", RegionSlug:""}`),
		},
	}

	for _, test := range testcases {
		t.Run(test.name, func(t *testing.T) {
			fakeD := &fakeDropletService{
				getFunc:  test.getDFunc,
				listFunc: test.listDFunc,
			}
			fakeS := &fakeStorageService{
				getVolumeFn: test.getVolumeFn,
			}
			fakeSA := &fakeStorageActionsService{
				attachFn: test.attachFn,
				getFn:    test.getSAFn,
			}

			fakeClient := newFakeClient(fakeS, fakeSA, fakeD)
			v := &VolumeManager{ctx: context.Background(), client: fakeClient}

			device, err := v.Attach(test.options, test.nodeName)
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
	listDropletFn := func(ctx context.Context, opt *godo.ListOptions) ([]godo.Droplet, *godo.Response, error) {
		droplet := newFakeDroplet()
		vol := newFakeVolume()
		droplet.VolumeIDs = []string{vol.ID}
		droplets := []godo.Droplet{*droplet}

		resp := newFakeOKResponse()
		return droplets, resp, nil
	}
	getDropletFunc := func(ctx context.Context, dropletID int) (*godo.Droplet, *godo.Response, error) {
		droplet := newFakeDroplet()
		if droplet.ID == dropletID {
			vol := newFakeVolume()
			droplet.VolumeIDs = []string{vol.ID}
			resp := newFakeOKResponse()
			return droplet, resp, nil
		}
		return nil, newFakeNotOKResponse(), fmt.Errorf("no droplet found with id %v", dropletID)
	}
	listVolumeFn := func(context.Context, *godo.ListVolumeParams) ([]godo.Volume, *godo.Response, error) {
		vol := newFakeVolume()
		return []godo.Volume{*vol}, newFakeOKResponse(), nil

	}
	detachStorageActions := func(ctx context.Context, volumeID string, dropletID int) (*godo.Action, *godo.Response, error) {
		return &godo.Action{
			Status: godo.ActionCompleted,
		}, newFakeOKResponse(), nil
	}
	getStorageActions := func(ctx context.Context, volumeID string, actionID int) (*godo.Action, *godo.Response, error) {
		return &godo.Action{
			Status: godo.ActionCompleted,
		}, newFakeOKResponse(), nil
	}

	testcases := []struct {
		name                string
		listDFunc           func(ctx context.Context, opt *godo.ListOptions) ([]godo.Droplet, *godo.Response, error)
		getDFunc            func(ctx context.Context, dropletID int) (*godo.Droplet, *godo.Response, error)
		listVolumesFn       func(context.Context, *godo.ListVolumeParams) ([]godo.Volume, *godo.Response, error)
		detachByDropletIDFn func(ctx context.Context, volumeID string, dropletID int) (*godo.Action, *godo.Response, error)
		getSAFn             func(ctx context.Context, volumeID string, actionID int) (*godo.Action, *godo.Response, error)
		nodeName            string
		device              string
		err                 error
	}{
		{
			"volume detach",
			listDropletFn,
			getDropletFunc,
			listVolumeFn,
			detachStorageActions,
			getStorageActions,
			"test-droplet",
			"test-volume",
			nil,
		},
		{
			"volume not attached",
			listDropletFn,
			func(ctx context.Context, dropletID int) (*godo.Droplet, *godo.Response, error) {
				droplet := newFakeDroplet()
				if droplet.ID == dropletID {
					resp := newFakeOKResponse()
					return droplet, resp, nil
				}
				return nil, newFakeNotOKResponse(), fmt.Errorf("no droplet found with id %v", dropletID)
			},
			listVolumeFn,
			detachStorageActions,
			getStorageActions,
			"test-droplet",
			"test-volume",
			fmt.Errorf("could not find volume attached at test-volume"),
		},
		{
			"volume detach fail",
			listDropletFn,
			getDropletFunc,
			listVolumeFn,
			func(ctx context.Context, volumeID string, dropletID int) (*godo.Action, *godo.Response, error) {
				return &godo.Action{
					Status: godo.ActionInProgress,
				}, newFakeOKResponse(), nil
			},
			func(ctx context.Context, volumeID string, actionID int) (*godo.Action, *godo.Response, error) {
				return &godo.Action{
					Status: "errored",
				}, newFakeOKResponse(), nil
			},
			"test-droplet",
			"test-volume",
			fmt.Errorf(`attach failed: godo.Action{ID:0, Status:"errored", Type:"", ResourceID:0, ResourceType:"", RegionSlug:""}`),
		},
	}
	for _, test := range testcases {
		t.Run(test.name, func(t *testing.T) {
			fakeD := &fakeDropletService{
				getFunc:  test.getDFunc,
				listFunc: test.listDFunc,
			}
			fakeS := &fakeStorageService{
				listVolumesFn: test.listVolumesFn,
			}
			fakeSA := &fakeStorageActionsService{
				detachByDropletIDFn: test.detachByDropletIDFn,
				getFn:               test.getSAFn,
			}

			fakeClient := newFakeClient(fakeS, fakeSA, fakeD)
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
	opt := &DigitalOceanOptions{
		DefaultOptions: cloud.DefaultOptions{
			VolumeID:   "80d414c6-295e-4e3a-ac58-eb9456c1e1d1",
			FsType:     "ext4",
			RW:         "rw",
			VolumeName: "test-volume",
		},
	}
	testcases := []struct {
		name     string
		options  *DigitalOceanOptions
		mountDir string
		device   string
		err      error
	}{
		{
			"mount device",
			opt,
			"/tmp/mount",
			"test-volume",
			cloud.ErrNotSupported,
		},
		{
			"fs type not specified",
			&DigitalOceanOptions{
				DefaultOptions: cloud.DefaultOptions{
					VolumeID:   "80d414c6-295e-4e3a-ac58-eb9456c1e1d1",
					RW:         "rw",
					VolumeName: "test-volume",
				},
			},
			"/tmp/mount",
			"test-volume",
			cloud.ErrNotSupported,
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
			cloud.ErrNotSupported,
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
