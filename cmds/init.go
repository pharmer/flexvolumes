package cmds

import (
	"context"

	"github.com/pharmer/flexvolumes/cloud"
	"github.com/pharmer/flexvolumes/cmds/options"
	"github.com/spf13/cobra"
)

func NewCmdInit() *cobra.Command {
	cfg := options.NewConfig()
	cmd := &cobra.Command{
		Use:               "init",
		Short:             "Initializes the driver.",
		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
			cloud, err := cloud.GetCloudManager(cfg.Provider, context.Background())
			if err != nil {
				Error(err).Print()
			}
			if err = cloud.Init(); err != nil {
				Error(err).Print()
			}
			Success().Capability().Print()
		},
	}
	return cmd
}
