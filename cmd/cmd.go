package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "authonomy",
	Short: "authonomy service",
	Long:  "authonomy service build over ssi service",
}

func init() {
	cobra.OnInitialize(getConfig)
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().BoolVarP(&resetFlag, "reset", "r", false, "Reset the service")
}

func getConfig() {
	// Set the base name of the config file, without the file extension.
	viper.SetConfigName("config")
	// Set the path to look for the config file in.
	viper.AddConfigPath(".")
	// Read in environment variables that match
	viper.AutomaticEnv()
	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		fmt.Println("Error reading config file:", err)
	}
}

func servicePort() string {
	port := viper.GetString("service.port")
	if port == "" {
		port = "8081"
	}
	return ":" + port
}

var resetFlag bool
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the authonomy service",
	Run: func(cmd *cobra.Command, args []string) {
		dbPath := viper.GetString("service.badger_path")
		secret := viper.GetString("service.db_encryption_key")
		apiKey := viper.GetString("api.x-api-key")
		ssiUrl := viper.GetString("service.ssi_service_url")
		Start(dbPath, secret, apiKey, servicePort(), ssiUrl, resetFlag)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
