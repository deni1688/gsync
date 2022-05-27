package cobra

import (
	syncer2 "deni1688/gsync/domain"
	"github.com/spf13/cobra"
	"log"
)

type cobraCli struct {
	rootCmd *cobra.Command
}

func NewCLI(gs syncer2.GsyncService) *cobraCli {
	rootCmd := &cobra.Command{
		Use:   "gsync",
		Short: "gsync is a tool to sync files between a remote storage provider like Google Drive and local directory",
	}

	rootCmd.AddCommand(pull(gs))
	rootCmd.AddCommand(push(gs))
	rootCmd.AddCommand(sync(gs))

	return &cobraCli{rootCmd}
}

func (c *cobraCli) Execute() error {
	return c.rootCmd.Execute()
}

func pull(gs syncer2.GsyncService) *cobra.Command {
	return &cobra.Command{
		Use:   "pull",
		Short: "pull files from remote gs to local directory",
		Run: func(cmd *cobra.Command, args []string) {
			if err := gs.Pull(syncer2.SyncTarget{Name: "Gsync"}); err != nil {
				log.Println("Failed to pull: ", err.Error())
			}
		},
	}
}

func push(gs syncer2.GsyncService) *cobra.Command {
	return &cobra.Command{
		Use:   "push",
		Short: "push files from local directory to remote gs",
		Run: func(cmd *cobra.Command, args []string) {
			if err := gs.Push(syncer2.SyncTarget{Name: "Gsync"}); err != nil {
				log.Println("Failed to push: ", err.Error())
			}
		},
	}
}

func sync(gs syncer2.GsyncService) *cobra.Command {
	return &cobra.Command{
		Use:   "sync",
		Short: "sync files between a remote gs and local directory",
		Run: func(cmd *cobra.Command, args []string) {
			if err := gs.Sync(syncer2.SyncTarget{Name: "Gsync"}); err != nil {
				log.Println("Failed to sync: ", err.Error())
			}
		},
	}
}
