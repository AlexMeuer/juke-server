/*
Copyright © 2024 none
*/
package cmd

import (
	"log"

	"github.com/alexmeuer/juke/pkg/oauth"
	"github.com/spf13/cobra"
)

// authCmd represents the auth command
var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Test the Spotify OAuth flow",
	Run: func(cmd *cobra.Command, args []string) {
		if err := oauth.Debug(cmd.Context()); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(authCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// authCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// authCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
