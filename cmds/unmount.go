package cmds

import (
	"context"

	"github.com/pharmer/flexvolumes/cloud"
	"github.com/pharmer/flexvolumes/cmds/options"
	"github.com/spf13/cobra"
)

func NewCmdUnmount() *cobra.Command {
	cfg := options.NewConfig()
	cmd := &cobra.Command{
		Use:               "unmount",
		Short:             "Unmount the volume",
		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				Error(cloud.ErrIncorrectArgNumber).Print()
			}
			cloud, err := cloud.GetCloudManager(cfg.Provider, context.Background())
			if err != nil {
				Error(err).Print()
			}
			dir := args[0]

			if err = cloud.Initialize(); err != nil {
				Error(err).Print()
			}

			if err := cloud.Unmount(dir); err != nil {
				Error(err).Print()
			}
			Success().Print()
		},
	}
	return cmd
}
