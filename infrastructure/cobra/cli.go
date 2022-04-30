package cobra

import (
	services2 "deni1688/gsync/domain/syncer"
	"github.com/spf13/cobra"
	"log"
)

type cobraCli struct {
	rootCmd *cobra.Command
}

func NewCLI(service services2.GsyncService) *cobraCli {
	rootCmd := &cobra.Command{
		Use:   "gsync",
		Short: "gsync is a tool to sync files between a remote service and local directory",
	}

	rootCmd.AddCommand(pullCmd(service))
	rootCmd.AddCommand(pushCmd(service))
	rootCmd.AddCommand(syncCmd(service))

	return &cobraCli{rootCmd}
}

func (c *cobraCli) Execute() error {
	return c.rootCmd.Execute()
}

func pullCmd(service services2.GsyncService) *cobra.Command {
	return &cobra.Command{
		Use:   "pull",
		Short: "pull files from remote service to local directory",
		Run: func(cmd *cobra.Command, args []string) {
			if err := service.PullFiles(services2.SyncFile{Name: "Gsync"}); err != nil {
				log.Println("Failed to pull: ", err.Error())
			}
		},
	}
}

func pushCmd(service services2.GsyncService) *cobra.Command {
	return &cobra.Command{
		Use:   "push",
		Short: "push files from local directory to remote service",
		Run: func(cmd *cobra.Command, args []string) {
			if err := service.PushFiles(services2.SyncFile{Name: "Gsync"}); err != nil {
				log.Println("Failed to push: ", err.Error())
			}
		},
	}
}

func syncCmd(service services2.GsyncService) *cobra.Command {
	return &cobra.Command{
		Use:   "sync",
		Short: "sync files between a remote service and local directory",
		Run: func(cmd *cobra.Command, args []string) {
			if err := service.SyncFiles(services2.SyncFile{Name: "Gsync"}); err != nil {
				log.Println("Failed to sync: ", err.Error())
			}
		},
	}
}
