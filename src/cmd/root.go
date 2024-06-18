package cmd

import (
	"fmt"
	"github.com/mjedari/health-checker/app/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"strings"
)

var rootCmd = &cobra.Command{
	Use:   "Health Checker",
	Short: "api health checker application based on clean architecture",
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here

	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
func init() {
	rootCmd.PersistentFlags().StringP("author", "a", "Mahdi Jedari", "i.jedari@gmail.com")
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath("../config")

	viper.AutomaticEnv()
	viper.SetEnvPrefix("health")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		logrus.Fatalf("Fatal error config file: %s \n", err)
	}
	viper.Unmarshal(&config.Config)
	logrus.Println("configuration initialized! (Notice: configurations may be initialised from OS ENV)")
}
