package commands

import (
	"os"

	"github.com/eris-ltd/eris/config"
	"github.com/eris-ltd/eris/initialize"
	"github.com/eris-ltd/eris/log"
	"github.com/eris-ltd/eris/util"

	"github.com/spf13/cobra"
)

var Init = &cobra.Command{
	Use:   "init",
	Short: "initialize your work space for smart contract glory",
	Long:  `create the Eris root ` + util.Tilde(config.ErisRoot) + ` directory and associated sub-directories`,
	Run: func(cmd *cobra.Command, args []string) {
		Router(cmd, args)
	},
}

func buildInitCommand() {
	addInitFlags()
}

func addInitFlags() {
	Init.Flags().BoolVarP(&do.Pull, "pull-images", "", true, "by default, pulls and/or update latest primary images. use flag to skip pulling/updating of images.")
	Init.Flags().BoolVarP(&do.Yes, "yes", "y", false, "over-ride command-line prompts")
	Init.Flags().BoolVarP(&do.Quiet, "testing", "", false, "DO NOT USE (for testing only)")
}

func Router(cmd *cobra.Command, args []string) {
	err := initialize.Initialize(do)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
}
