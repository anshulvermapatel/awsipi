/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"awsipi/awsutil"

	"github.com/spf13/cobra"
)

var homeDir, err = os.UserHomeDir()

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		// The below err will check if the 'os.UserHomeDir()' for 'homeDir' var threw any error.
		checkErr(err)

		// This will validate if the command arguments have been given properly.
		err := validate(args)
		checkErr(err)

		if _, err := os.Stat(homeDir + "/.awsipi"); os.IsNotExist(err) {
			err := os.Mkdir(homeDir+"/.awsipi", 0740)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
		if err := clusterExists(args[0]); err != nil {
			fmt.Println(err)
			os.Exit(1)
		} else {

			// getting the AWS zone and the Instance Ids of the nodes of the cluster to save in the Cluster file
			zone, instanceIDs := awsutil.GetZoneAndInstanceId()
			clusterFile, err := os.Create(homeDir + "/.awsipi/" + args[0])
			checkErr(err)
			// defer statement to close this file after execution of this function.
			defer clusterFile.Close()

			clusterContent := new(bytes.Buffer)
			_, err = clusterContent.WriteString(zone + "\n")
			checkErr(err)

			for _, instanceID := range instanceIDs {
				_, err = clusterContent.WriteString(*instanceID + "\n")
				checkErr(err)
			}
			_, err = clusterContent.WriteTo(clusterFile)
			checkErr(err)
		}
		fmt.Printf("Cluster has been initiated and a file got create in " + homeDir + "/.awsipi/ directory")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// This function will valiadate if the argument to the sub-command is provided properly or if the provided cluster file already exist.
func validate(args []string) error {
	if len(args) < 1 {
		return errors.New("Please specify the cluster name")
	} else if len(args) > 1 {
		return errors.New("Only one argument is accepted")
	} else {
		return nil
	}
}

func clusterExists(cluster string) error {
	clusterFiles, err := ioutil.ReadDir(homeDir + "/.awsipi")
	checkErr(err)

	for _, clusterFile := range clusterFiles {
		if clusterFile.Name() == cluster {
			return errors.New("a Cluster with this name has already been initialized")
		}
	}
	return nil
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
