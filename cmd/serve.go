package cmd

import (
	"github.com/alexmeuer/juke/internal/server"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run the web server",
	Long: `Juke is a collaborative playlist manager.
This command runs the web server, the heart of juke.`,
	Run: func(cmd *cobra.Command, args []string) {
		port, err := cmd.Flags().GetUint16("port")
		if err != nil {
			// The flag definition with init() should prevent this case.
			panic(err)
		}
		err = server.Serve(cmd.Flag("host").Value.String(), port)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	serveCmd.Flags().StringP("host", "H", "0.0.0.0", "Host to listen on")
	serveCmd.Flags().Uint16P("port", "p", 8080, "Port to listen on")
}
