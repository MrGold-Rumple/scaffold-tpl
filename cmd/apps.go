/*
* @Author: Rumple
* @Email: wrp357711589@163.com
* @DateTime: 2022/3/1 15:41
 */

package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/MrGold-Rumple/scaffold-tpl/tpl"
	"github.com/spf13/cobra"
)

var appCmd = &cobra.Command{
	Use:   "add",
	Short: "generate a new apps",
	RunE: func(cmd *cobra.Command, args []string) error {
		return NewApps()
	},
}

var (
	aFlagDir  string
	aFlagApps []string
)

func init() {
	appCmd.PersistentFlags().StringVarP(&aFlagDir, "dir", "c", "", "project dir")
	appCmd.PersistentFlags().StringSliceVarP(&aFlagApps, "apps", "a", nil, "-a app1,app2,app3")
	rootCmd.AddCommand(appCmd)
}

// modelDir /apps/model
// apiDir /apps/api
// name app‘s name
// modName go mod name
func generateApp(appsDir, name, modName string) {

	appName := strings.ToLower(name)
	appTitle := Title(appName)

	apiDir := path.Join(appsDir, "api")
	dalDir := path.Join(appsDir, "dal")
	modelDir := path.Join(appsDir, "model")
	appDir := path.Join(apiDir, appName)
	checkDir(appDir)

	// files := []string{"controller.go", "protocol.go"}
	appParam := tpl.AppParam{
		AppName:    appName,
		AppTitle:   appTitle,
		BQ:         "`",
		ModuleName: modName,
	}

	tasks := []GenTask{{
		Name:     "dal",
		Filename: path.Join(dalDir, appName+".go"),
		Tpl:      tpl.DalGo,
		// Filename: path.Join(appDir, "dal.go"),
		// Tpl:      tpl.GoOnlyPkgFile,
		// Data: tpl.GoPkgFileParam{
		// 	PkgName: appName,
		// 	Comment: "// define your DQL DML",
		// },
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

func getModNameFromFile(mod string) (result string, err error) {

	var f *os.File

	if f, err = os.Open(mod); err != nil {
		return
	}

	defer func() {
		_ = f.Close()
	}()

	reader := bufio.NewReader(f)

	for {
		var b []byte

		if b, _, err = reader.ReadLine(); err != nil {
			return
		}
		line := string(b)
		prefix := "module "

		if strings.HasPrefix(line, prefix) {
			result = strings.TrimLeft(line, prefix)
			return
		}
	}
}

// NewApps 新增app
func NewApps() error {

	if len(aFlagApps) == 0 {
		return errors.New("nothing todo")
	}

	// 如果不传则默认找寻当前目录
	goMod := "go.mod"
	if aFlagDir == "" {
		aFlagDir = "."
	}

	goMod = path.Join(aFlagDir, goMod)
	modName, err := getModNameFromFile(goMod)
	if err != nil {
		return fmt.Errorf("get go mod name error:%s", err)
	}

	for _, a := range aFlagApps {

		modDir := path.Join(aFlagDir, "apps/model")

		generateApp(path.Join(aFlagDir, "apps"), strings.ToLower(a), modName)

		content := fmt.Sprintf("        new(%s),\n", Title(a))

		if err = InsertStringToFile(path.Join(modDir, "model.go"), content, 28); err != nil {
			return fmt.Errorf("add app model error:%s", err)
		}
	}
	return nil
}
