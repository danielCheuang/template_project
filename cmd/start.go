package cmd

import (
	"template_project/config"
	"template_project/db/mysql"
	"template_project/db/redis"
	"template_project/logger"
	"template_project/server"
	"fmt"

	"github.com/spf13/cobra"
)

var configFile *string

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "api(.exe) start",
	Long:  "api(.exe) start -c ./build/app.json",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("start service")

		cfg := config.Init(configFile)

		if cfg.MySQL.Enable {
			mysql.Init()
		}

		if cfg.Redis.Enable {
			redis.Init()
		}

		logger.Init()

		logger.Log.Info("Config:", cfg)

		api, err := server.New(&cfg)
		if err != nil {
			panic(err)
		}

		if err := api.Start(); err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
	configFile = startCmd.Flags().StringP("config", "c", "", "start config file (required)")
	err := startCmd.MarkFlagRequired("config")
	if err != nil {
		fmt.Println(err)
	}
}
