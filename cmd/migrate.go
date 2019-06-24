package cmd

import (
	"template_project/config"
	"template_project/db/mysql"
	"template_project/model"
	"fmt"

	"github.com/spf13/cobra"
)

var cfgFile *string

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "api(.exe) migrate",
	Long:  "api(.exe) migrate -c ./build/app.json",
	Run: func(cmd *cobra.Command, args []string) {
		config.Init(cfgFile)
		mysql.Init()
		migrate()
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
	cfgFile = migrateCmd.Flags().StringP("config", "c", "", "start config file (required)")
	err := migrateCmd.MarkFlagRequired("config")
	if err != nil {
		fmt.Println(err)
	}
}

func migrate()  {
	var tables []interface{}
	tables = append(tables, &model.User{},)
	mysql.DB.RegistTables(tables)
}
