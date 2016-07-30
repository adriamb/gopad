package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/adriamb/gopad/server"
	"github.com/adriamb/gopad/server/config"
	"github.com/adriamb/gopad/web"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/user"
)

var cfgFile string
var C config.Config

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "gopad",
	Short: "Minimal go notepad",
	Long:  `A minimal markdown personal notepad written in go`,
	Run: func(cmd *cobra.Command, args []string) {
		server.NewServer(C)
		web.InitWeb()
		server.StartServer()
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)

	}

}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gopad.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigType("yaml")
	viper.SetConfigName(".gopad") // name of config file (without extension)
	viper.AddConfigPath("$HOME")  // adding home directory as first search path
	viper.SetEnvPrefix("gopad")   // so viper.AutomaticEnv will get matching envvars starting with O2M_
	viper.AutomaticEnv()          // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		panic(fmt.Sprintf("Config file not found: %v", err))
	}
	if err := viper.Unmarshal(&C); err != nil {
		panic(err)
	}

	if C.DataDir == "" {
		usr, err := user.Current()
		if err != nil {
			log.Fatal(err)
		}
		C.DataDir = usr.HomeDir + "/.gopad"
	}

	if C.TmpDir == "" {
		C.TmpDir = "/tmp/gopad/temp"
	}

	if C.CacheDir == "" {
		C.CacheDir = "/tmp/gopad/cache"
	}

	json, _ := json.MarshalIndent(C, "", "  ")
	fmt.Println("Efective configuration: " + string(json))

}
