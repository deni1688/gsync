package cli

import (
	"deni1688/gsync/domain"
	"github.com/spf13/cobra"
	"log"
)

type Runtime struct {
	rootCmd *cobra.Command
}

func New(service domain.GsyncServiceContract) *Runtime {
	rootCmd := &cobra.Command{
		Use:   "gsync",
		Short: "gsync is a tool to sync files between a remote service and local directory",
	}

	rootCmd.AddCommand(pullCmd(service))
	rootCmd.AddCommand(pushCmd(service))
	rootCmd.AddCommand(syncCmd(service))

	return &Runtime{rootCmd}
}

func (c *Runtime) Execute() error {
	return c.rootCmd.Execute()
}

func pullCmd(service domain.GsyncServiceContract) *cobra.Command {
	return &cobra.Command{
		Use:   "pull",
		Short: "pull files from remote service to local directory",
		Run: func(cmd *cobra.Command, args []string) {
			if err := service.Pull(domain.FileInfo{Name: "Gsync"}); err != nil {
				log.Println("Failed to pull: ", err.Error())
			}
		},
	}
}

func pushCmd(service domain.GsyncServiceContract) *cobra.Command {
	return &cobra.Command{
		Use:   "push",
		Short: "push files from local directory to remote service",
		Run: func(cmd *cobra.Command, args []string) {
			if err := service.Push(); err != nil {
				log.Println("Failed to push: ", err.Error())
			}
		},
	}
}

func syncCmd(service domain.GsyncServiceContract) *cobra.Command {
	return &cobra.Command{
		Use:   "sync",
		Short: "sync files between a remote service and local directory",
		Run: func(cmd *cobra.Command, args []string) {
			if err := service.Sync(); err != nil {
				log.Println("Failed to sync: ", err.Error())
			}
		},
	}
}
