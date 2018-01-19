package cmds

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/pharmer/flexvolumes/cloud"
	"github.com/pharmer/flexvolumes/cmds/options"
	"github.com/spf13/cobra"
)

func NewCmdMountDevice() *cobra.Command {
	cfg := options.NewConfig()
	cmd := &cobra.Command{
		Use:               "mountdevice",
		Short:             "Mount device mounts the device to a global path which individual pods can then bind mount.",
		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) <= 0 && len(args) > 3 {
				Error(cloud.ErrIncorrectArgNumber).Print()
			}
			cloud, err := cloud.GetCloudManager(cfg.Provider, context.Background())
			if err != nil {
				Error(err).Print()
			}

			dir := args[0]
			device := args[1]
			opt := cloud.NewOptions()
			if err := json.Unmarshal([]byte(args[2]), opt); err != nil {
				Error(fmt.Errorf("could not parse options for attach; got %v", os.Args[2])).Print()
			}

			if err = cloud.Initialize(); err != nil {
				Error(err).Print()
			}

			if err := cloud.MountDevice(dir, device, opt); err != nil {
				Error(err).Print()
			}
			Success().Print()
		},
	}
	return cmd
}
