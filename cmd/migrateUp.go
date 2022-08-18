package cmd

import (
	"audioPhile/database"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var up = &cobra.Command{
	Use:   "up",
	Short: "migration is up",
	Long:  " ",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Running migrate up command xyz")
		err := database.ConnectAndMigrate("localhost", "5432", "audiophile", "local", "local", database.SSLModeDisable)
		if err != nil {
			logrus.Print(err)
		}
		fmt.Println("success")
	},
}

func init() {
	rootCmd.AddCommand(up)
}
