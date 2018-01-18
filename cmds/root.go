package cmds

import (
	"flag"
	//	"log"
	"os"
	"strings"

	v "github.com/appscode/go/version"
	"github.com/jpillora/go-ogle-analytics"
	_ "github.com/pharmer/flexvolumes/cloud/providers"
	"github.com/spf13/cobra"
	//"github.com/spf13/pflag"
)

const (
	gaTrackingCode = "UA-62096468-20"
)

func NewRootCmd(version string) *cobra.Command {
	var (
		enableAnalytics = true
	)
	rootCmd := &cobra.Command{
		Use:               "flexvolume [command]",
		Short:             `Pharm flexvolume by Appscode - Start farms`,
		DisableAutoGenTag: true,
		PersistentPreRun: func(c *cobra.Command, args []string) {
			/*c.Flags().VisitAll(func(flag *pflag.Flag) {
				log.Printf("FLAG: --%s=%q", flag.Name, flag.Value)
			})*/
			if enableAnalytics && gaTrackingCode != "" {
				if client, err := ga.NewClient(gaTrackingCode); err == nil {
					parts := strings.Split(c.CommandPath(), " ")
					client.Send(ga.NewEvent(parts[0], strings.Join(parts[1:], "/")).Label(version))
				}
			}
		},
	}
	rootCmd.PersistentFlags().AddGoFlagSet(flag.CommandLine)
	// ref: https://github.com/kubernetes/kubernetes/issues/17162#issuecomment-225596212
	flag.CommandLine.Parse([]string{})
	rootCmd.PersistentFlags().BoolVar(&enableAnalytics, "analytics", enableAnalytics, "Send analytical events to Google Analytics")

	if len(os.Args) > 1 {
		checkSupported(os.Args[1])
	}
	rootCmd.AddCommand(NewCmdInit())
	rootCmd.AddCommand(NewCmdAttach())
	rootCmd.AddCommand(NewCmdMount())
	rootCmd.AddCommand(NewCmdDetach())
	rootCmd.AddCommand(NewCmdUnmount())

	rootCmd.AddCommand(v.NewCmdVersion())

	return rootCmd
}

func checkSupported(cmd string) bool {
	supported := []string{"init", "attach", "detach", "mount", "unmount"}
	for _, s := range supported {
		if s == cmd {
			return true
		}
	}
	Error(ErrNotSupported).Print()
	return false
}
