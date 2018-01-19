package cmds

import (
	"context"
	"errors"

	"github.com/pharmer/flexvolumes/cloud"
	"github.com/pharmer/flexvolumes/cmds/options"
	"github.com/spf13/cobra"
)

func NewCmdDetach() *cobra.Command {
	cfg := options.NewConfig()
	cmd := &cobra.Command{
		Use:               "detach",
		Short:             "Detach the volume from the Kubelet node",
		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				Error(cloud.ErrIncorrectArgNumber).Print()
			}
			cloud, err := cloud.GetCloudManager(cfg.Provider, context.Background())
			if err != nil {
				Error(err).Print()
			}

			device := args[0]
			nodeName := args[1]
			if nodeName == "" {
				Error(errors.New("node name not found")).Print()
			}

			if err = cloud.Initialize(); err != nil {
				Error(err).Print()
			}

			if err := cloud.Detach(device, nodeName); err != nil {
				Error(err).Print()
			}
			Success().Print()
		},
	}
	return cmd
}
