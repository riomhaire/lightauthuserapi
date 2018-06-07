// Copyright © 2018 Gary Leeson
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
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/riomhaire/lightauthuserapi/frameworks/application/lightauthuserapi/bootstrap"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "A Simple User Service",
	Long: `Is a user service which stores users and their claims,
	       and accessed via key`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("serve called")

		application := bootstrap.Application{}

		application.Initialize(cmd, args)

		c := make(chan os.Signal, 2)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)

		go func() {
			<-c
			log.Println("Shutting Down")
			// pprof.StopCPUProfile()
			// //trace.Stop()
			// tracefile.Close()
			application.Stop()
			os.Exit(0)
		}()
		application.Run()

	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().IntP("port", "p", 3060, "Default Port to Listen to.")
	serveCmd.Flags().StringP("key", "k", "secret", "Secret needed to access api.")
	serveCmd.Flags().StringP("usersFile", "u", "users.csv", "If User File used this is the one to use - must be r/w.")
	serveCmd.Flags().StringP("rolesFile", "r", "roles.csv", "If Role File used this is the one to use - must be r/w.")

	serveCmd.Flags().StringP("consulHost", "t", "", "Host where consul resides usually something like http://consul:8500 ")
	serveCmd.Flags().BoolP("consul", "c", false, "Enable consul support")

}
