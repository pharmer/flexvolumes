package cmds

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/pharmer/flexvolumes/cloud"
	"k8s.io/kubernetes/pkg/volume/flexvolume"
)

const (
	StatusFailure = "Failure"
)

type DriverOutput struct {
	flexvolume.DriverStatus
	Options interface{} `json:"options,omitempty"`
}

func Success() DriverOutput {
	return DriverOutput{
		DriverStatus: flexvolume.DriverStatus{
			Status:  flexvolume.StatusSuccess,
			Message: "Flex driver initialized",
		},
	}
}

func Device(device string) DriverOutput {
	return DriverOutput{
		DriverStatus: flexvolume.DriverStatus{
			Status:     flexvolume.StatusSuccess,
			DevicePath: device,
		},
	}
}

func Error(err error) DriverOutput {
	output := DriverOutput{
		DriverStatus: flexvolume.DriverStatus{
			Status:  StatusFailure,
			Message: err.Error(),
		},
	}
	if err == cloud.ErrNotSupported {
		output.Status = flexvolume.StatusNotSupported
	}
	return output
}

func (do DriverOutput) Capability() DriverOutput {
	do.Capabilities = &flexvolume.DriverCapabilities{
		Attach:         true,
		SELinuxRelabel: true,
	}
	return do
}

func (do DriverOutput) Print() {
	b, _ := json.Marshal(do)
	fmt.Printf("%s\n", string(b))
	code := 1
	if do.Status == flexvolume.StatusSuccess {
		code = 0
	}
	os.Exit(code)
}
