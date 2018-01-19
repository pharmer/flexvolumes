package cmds

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/pharmer/flexvolumes/cloud"
	"github.com/pharmer/flexvolumes/cmds/options"
	"github.com/spf13/cobra"
)

func NewCmdAttach() *cobra.Command {
	cfg := options.NewConfig()
	cmd := &cobra.Command{
		Use:               "attach",
		Short:             "Attach the volume specified by the given spec on the given host",
		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				Error(cloud.ErrIncorrectArgNumber).Print()
			}
			cloud, err := cloud.GetCloudManager(cfg.Provider, context.Background())
			if err != nil {
				Error(err).Print()
			}

			opt := cloud.NewOptions()
			if err := json.Unmarshal([]byte(args[0]), opt); err != nil {
				Error(fmt.Errorf("could not parse options for attach; got %v", os.Args[0])).Print()
			}
			nodeName := args[1]
			if nodeName == "" {
				Error(errors.New("node name not found")).Print()
			}

			if err = cloud.Initialize(); err != nil {
				Error(err).Print()
			}

			device, err := cloud.Attach(opt, nodeName)
			if err != nil {
				Error(err).Print()
			}
			Device(device).Print()
		},
	}
	return cmd
}
