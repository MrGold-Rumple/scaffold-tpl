/*
* @Author: Rumple
* @Email: wrp357711589@163.com
* @DateTime: 2022/2/21 17:56
 */

package main

import (
	"os"

	"github.com/MrGold-Rumple/scaffold-tpl/cmd"
	"github.com/MrGold-Rumple/scaffold-tpl/console"
)

func main() {
	if err := cmd.Execute(); err != nil {
		console.Error("execute error:%s", err)
		os.Exit(1)
	}
}
