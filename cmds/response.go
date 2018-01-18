package cmds

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/urfave/cli"
)

const (
	StatusSuccess      = "Success"
	StatusFailure      = "Failure"
	StatusNotSupported = "Not supported"
)

var ErrNotSupported = errors.New("Not Supported")
var ErrIncorrectArgNumber = errors.New("Incorrect number of args")

type DriverOutput struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Device  string `json:"device"`
}

func Success() DriverOutput {
	return DriverOutput{
		Status: StatusSuccess,
	}
}

func Device(device string) DriverOutput {
	return DriverOutput{
		Status: StatusSuccess,
		Device: device,
	}
}

func Error(err error) DriverOutput {
	output := DriverOutput{
		Status:  StatusFailure,
		Message: err.Error(),
	}
	if err == ErrNotSupported {
		output.Status = StatusNotSupported
	}
	return output
}

func handleErr(f func(*cli.Context) error) func(*cli.Context) error {
	return func(c *cli.Context) error {
		err := f(c)
		if err != nil {
			Error(err).Print()
		}
		return err
	}
}

func Options(options map[string]interface{}) RancherDriverOutput {
	return RancherDriverOutput{
		DriverOutput: DriverOutput{
			Status: StatusSuccess,
		},
		Options: options,
	}
}

type RancherDriverOutput struct {
	DriverOutput
	Options map[string]interface{} `json:"options"`
}

func (rdo RancherDriverOutput) Print() {
	b, _ := json.Marshal(rdo)
	fmt.Printf("%s\n", string(b))
}

func (do DriverOutput) Print() {
	b, _ := json.Marshal(do)
	fmt.Printf("%s\n", string(b))
	code := 1
	if do.Status == StatusSuccess {
		code = 0
	}
	os.Exit(code)
}
