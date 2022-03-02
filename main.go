/*
* @Author: Rumple
* @Email: ruipeng.wu@cyclone-robotics.com
* @DateTime: 2022/2/21 17:56
 */

package main

import (
	"os"

	"github.com/wuruipeng404/scaffold-tpl/cmd"
	"github.com/wuruipeng404/scaffold-tpl/console"
)

func main() {
	if err := cmd.Execute(); err != nil {
		console.Error("execute error:%s", err)
		os.Exit(1)
	}
}
