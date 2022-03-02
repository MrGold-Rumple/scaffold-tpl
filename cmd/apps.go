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

func generateApp(apiDir, name, modName string) {

	appName := strings.ToLower(name)
	appTitle := strings.Title(appName)

	appDir := path.Join(apiDir, appName)
	checkDir(appDir)

	// files := []string{"controller.go", "data_language.go", "model.go", "structure.go"}

	appParam := tpl.AppParam{
		AppName:    appName,
		AppTitle:   appTitle,
		BQ:         "`",
		ModuleName: modName,
	}

	tasks := []Task{{
		Name:     "DL",
		Filename: path.Join(appDir, "data_language.go"),
		Tpl:      tpl.GoOnlyPkgFile,
		Data: tpl.GoPkgFileParam{
			PkgName: appName,
			Comment: "// define your DQL DML",
		},
	}, {
		Name:     "model",
		Filename: path.Join(appDir, "model.go"),
		Tpl:      tpl.AppModelGO,
		Data:     appParam,
	}, {
		Name:     "structure",
		Filename: path.Join(appDir, "structure.go"),
		Tpl:      tpl.AppStructureGO,
		Data:     appParam,
	}, {
		Name:     "controller",
		Filename: path.Join(appDir, "controller.go"),
		Tpl:      tpl.AppControllerGo,
		Data:     appParam,
	}}

	generateAllTemplateFiles(tasks)

}

func NewApps() error {
	return nil
}
