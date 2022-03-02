/*
* @Author: Rumple
* @Email: ruipeng.wu@cyclone-robotics.com
* @DateTime: 2022/2/21 17:56
 */

package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "scaffold",
	Short: "scaffold is a generate standard web tool",
	Long:  ``,
}

func Execute() error {
	return rootCmd.Execute()
}
