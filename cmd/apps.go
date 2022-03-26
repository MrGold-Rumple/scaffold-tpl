/*
* @Author: Rumple
* @Email: ruipeng.wu@cyclone-robotics.com
* @DateTime: 2022/3/1 15:41
 */

package cmd

import (
	"path"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wuruipeng404/scaffold-tpl/tpl"
)

var appCmd = &cobra.Command{
	Use:   "add",
	Short: "generate a new apps",
	RunE: func(cmd *cobra.Command, args []string) error {
		return NewApps()
	},
}

var (
	flagDir string
)

func init() {
	appCmd.PersistentFlags().StringVarP(&flagDir, "dir", "c", "", "project dir")
	// rootCmd.AddCommand(appCmd)
}

func generateApp(modelDir, apiDir, name, modName string) {

	appName := strings.ToLower(name)
	appTitle := strings.Title(appName)

	appDir := path.Join(apiDir, appName)
	checkDir(appDir)

	// files := []string{"controller.go", "dal.go", "protocol.go"}

	appParam := tpl.AppParam{
		AppName:    appName,
		AppTitle:   appTitle,
		BQ:         "`",
		ModuleName: modName,
	}

	tasks := []GenTask{{
		Name:     "Dal",
		Filename: path.Join(appDir, "dal.go"),
		Tpl:      tpl.GoOnlyPkgFile,
		Data: tpl.GoPkgFileParam{
			PkgName: appName,
			Comment: "// define your DQL DML",
		},
	}, {
		Name:     "protocol",
		Filename: path.Join(appDir, "protocol.go"),
		Tpl:      tpl.AppProtocolGO,
		Data:     appParam,
	}, {
		Name:     "controller",
		Filename: path.Join(appDir, "controller.go"),
		Tpl:      tpl.AppControllerGo,
		Data:     appParam,
	}, {
		Name:     "model",
		Filename: path.Join(modelDir, appName+".go"),
		Tpl:      tpl.AppModelGO,
		Data:     appParam,
	}}

	generateAllTemplateFiles(tasks)

}

func NewApps() error {
	return nil
}
