/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"regexp"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "hcelper",
	Short: "A helper tool for the HCE platform",
	Long: `A helper tool for the HCE platform. 

Kickstart your usage from a Vault-first perspective.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is .hce.hcl)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {

		dirname, err := os.UserHomeDir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		viper.AddConfigPath(dirname)
		viper.SetConfigName(".hce")
		viper.SetConfigType("hcl")

		// global disable TLS verify using gateway_tls_skip_verify optional setting
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: viper.GetBool("gateway_tls_skip_verify")}

	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			fmt.Println("No .hce.hcl config file found. Provide a path to one, or ensure it is in your homedir.")
			fmt.Println("See https://github.com/tucows/hcelper for more information.")
		} else {
			fmt.Sprintf("Error: %v", err)
		}
		os.Exit(1)
	}

	gatewayCheck := viper.InConfig("gateway")
	if !gatewayCheck {
		fmt.Println("Gateway not found")
		os.Exit(1)
	}
	gateway := viper.GetString("gateway")
	fqdnCheck, err := regexp.MatchString(`https\:\/\/*`, gateway)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if !fqdnCheck {
		fmt.Printf("Not a valid FQDN: %v\n", gateway)
		os.Exit(1)
	}

}
