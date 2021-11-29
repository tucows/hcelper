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
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"regexp"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

type AuthCredentials struct {
	Username string `example:"jdoe"`
	Password string `example:"supersecretpassword"`
	Method   string `example:"ldap"`
}

type VaultLDAPResponse struct {
	LeaseID       string      `json:"lease_id"`
	Renewable     bool        `json:"renewable"`
	LeaseDuration int         `json:"lease_duration"`
	Data          interface{} `json:"data"`
	Auth          struct {
		ClientToken string   `json:"client_token"`
		Policies    []string `json:"policies"`
		Metadata    struct {
			Username string `json:"username"`
		} `json:"metadata"`
		LeaseDuration int  `json:"lease_duration"`
		Renewable     bool `json:"renewable"`
	} `json:"auth"`
}

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login with your organization's auth provider",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	PreRun: func(cmd *cobra.Command, args []string) {

	},
	Run: func(cmd *cobra.Command, args []string) {

		username := cmd.Flag("username").Value.String()
		method := cmd.Flag("method").Value.String()
		env := cmd.Flag("env").Value.String()

		envCheck := func(input string) error {
			_, err := regexp.MatchString(`pre|prod`, input)
			if err != nil {
				return errors.New(`env must be "pre" or "prod"`)
			}
			return nil
		}

		if env == "" {
			envPrompt := promptui.Prompt{
				Label:    "Environment (pre or prod)",
				Validate: envCheck,
			}
			env = envPrompt.Run()
		}
		os.Setenv("VAULT_ADDR", env)

		// create the password prompt
		validate := func(input string) error {
			if len(input) < 6 {
				return errors.New("password must have more than 6 characters")
			}
			return nil
		}

		passPrompt := promptui.Prompt{
			Label:    "Password",
			Validate: validate,
			Mask:     '*',
		}

		result, err := passPrompt.Run()

		if err != nil {
			fmt.Printf("Password input failed %v\n", err)
			os.Exit(1)
		}

		// prepare auth request to gateway
		authCreds := AuthCredentials{
			Username: username,
			Password: result,
			Method:   method,
		}

		body, err := json.Marshal(authCreds.Password)
		if err != nil {
			fmt.Printf("Error in parsing auth request parameters: %v\n", err)
			os.Exit(1)
		}
		if authCreds.Method != "" {
			switch authCreds.Method {
			case "ldap":
				req, err := http.NewRequest("POST", fmt.Sprintf("%s/auth/ldap/login/%s", env, authCreds.Username), bytes.NewBuffer(body))
				if err != nil {
					fmt.Printf("Error constructing LDAP login request: %v\n", err)
					os.Exit(1)
				}
				client := &http.Client{}
				resp, err := client.Do(req)
				if err != nil {
					fmt.Printf("Error logging into Vault via LDAP: %v\n", err)
					os.Exit(1)
				}
				defer resp.Body.Close()

				ldapResp := VaultLDAPResponse{}
				err = json.NewDecoder(resp.Body).Decode(&ldapResp)

				os.Setenv("VAULT_TOKEN", ldapResp.Auth.ClientToken)

			}
		}

	},
}

func init() {
	rootCmd.AddCommand(loginCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loginCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// loginCmd.Flags().StringP("gateway", "g", "", "The Mortar API gateway URL")
	loginCmd.Flags().StringP("username", "u", "", "The username for user login credentials")
	loginCmd.Flags().StringP("env", "e", "", "The environment you're logging into (pre or prod)")
	loginCmd.Flags().StringP("method", "m", "ldap", "The login method")
	loginCmd.MarkFlagRequired("username")
}
