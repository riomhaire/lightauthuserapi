// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"

	"github.com/riomhaire/lightauthuserapi/frameworks/application/lightauthuserapi/bootstrap"
	"github.com/spf13/cobra"
)

// versionCmd represents the user command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Returns application version",
	Long:  ` `,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(bootstrap.VERSION)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)

}
