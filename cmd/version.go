/*
* @Author: Rumple
* @Email: ruipeng.wu@cyclone-robotics.com
* @DateTime: 2022/2/21 17:56
 */

package cmd

import (
	"github.com/spf13/cobra"
	"github.com/wuruipeng404/scaffold-tpl/console"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "the tools version",
	Run: func(cmd *cobra.Command, args []string) {
		console.Info("release 2022-02-07 v2.4")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
