/*
* @Author: Rumple
* @Email: wrp357711589@163.com
* @DateTime: 2022/2/21 17:56
 */

package cmd

import (
	"github.com/MrGold-Rumple/scaffold-tpl/console"
	"github.com/spf13/cobra"
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
