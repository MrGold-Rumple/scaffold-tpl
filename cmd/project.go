/*
* @Author: Rumple
* @Email: ruipeng.wu@cyclone-robotics.com
* @DateTime: 2022/2/21 17:56
 */

package cmd

import (
	"fmt"
	"path"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wuruipeng404/scaffold-tpl/console"
	"github.com/wuruipeng404/scaffold-tpl/tpl"
)

var generateCmd = &cobra.Command{
	Use:   "new",
	Short: "generate a standard web project with gin",
	RunE: func(cmd *cobra.Command, args []string) error {
		return Generate()
	},
}

const (
	_PostgresFlag = "pg"
	_Postgres     = "postgres"
	_Mysql        = "mysql"
	_Mongo        = "mongo"
	_Redis        = "redis"
	_Es           = "es"

	defaultDir      = "scaffold-demo"
	defaultFilePerm = 0755
)

var (
	nFlagDir       string
	nFlagModName   string
	nFlagDb        string
	nFlagApps      []string
	nFlagNosqlList []string
)

func init() {

	generateCmd.PersistentFlags().StringVarP(&nFlagDir, "dir", "c", defaultDir, "the new project dir")
	generateCmd.PersistentFlags().StringVarP(&nFlagModName, "module", "m", "", "init go module name, (default same as \"${dir}\")")
	generateCmd.PersistentFlags().StringVarP(&nFlagDb, "db", "d", _PostgresFlag, "sql driver pg,mysql")
	generateCmd.PersistentFlags().StringSliceVarP(&nFlagApps, "apps", "a", []string{"user"}, "init apps, example -a user,file,category")
	generateCmd.PersistentFlags().StringSliceVarP(&nFlagNosqlList, "nosql", "n", []string{}, "init no sql driver -n=mongodb,es,redis")

	rootCmd.AddCommand(generateCmd)
}

func Generate() (err error) {
	console.Info("check param")
	if err = _ParamCheck(); err != nil {
		return err
	}
	session := newSession(nFlagDir)
	goModInit(session)

	console.Info(fmt.Sprintf("start init project %s, db %+v, nosql %+v, module name %s, apps %+v",
		nFlagDir, nFlagDb, nFlagNosqlList, nFlagModName, nFlagApps))

	apiDir := path.Join(nFlagDir, "api")
	apiProtocol := path.Join(apiDir, "protocol")
	confDir := path.Join(nFlagDir, "config")
	docsDir := path.Join(nFlagDir, "docs")
	internal := path.Join(nFlagDir, "internal")
	modelDir := path.Join(nFlagDir, "model")

	needMkDirs := []string{apiProtocol, modelDir, confDir, docsDir, internal}

	console.Info(fmt.Sprintf("start create all dir %+v", needMkDirs))
	makeAllDir(needMkDirs)

	scriptSubfix := "ps1"
	stpl := tpl.PowerBuildScript

	if runtime.GOOS != "windows" {
		scriptSubfix = "sh"
		stpl = tpl.BashBuildScript
	}

	tasks := []GenTask{{
		Name:     "Dockerfile",
		Filename: path.Join(nFlagDir, "Dockerfile"),
		Tpl:      tpl.DockerFile,
		Data:     tpl.DockerFileParam{BinName: nFlagModName},
	}, {
		Name:     "gitignore",
		Filename: path.Join(nFlagDir, ".gitignore"),
		Tpl:      tpl.GitIgnore,
	}, {
		Name:     "dockerIgnore",
		Filename: path.Join(nFlagDir, ".dockerignore"),
		Tpl:      tpl.DockerIgnore,
	}, {
		Name:     "readme",
		Filename: path.Join(nFlagDir, "README.md"),
		Tpl:      tpl.Readme,
	}, {
		Name:     "dockerBuildScript",
		Filename: path.Join(nFlagDir, "t."+scriptSubfix),
		Tpl:      stpl,
		Data: tpl.DockerBuildParam{
			ContainerName: nFlagDir + "-sc",
			ImageTag:      nFlagDir + ":latest",
			BuildArg:      "--build-arg config=config",
			ExportPort:    "7788",
		},
	}, {
		Name:     "swaggerScript",
		Filename: path.Join(nFlagDir, "docs."+scriptSubfix),
		Tpl:      tpl.GenSwagger,
	}}

	// main.go
	tasks = append(tasks, GenTask{
		Name:     "main",
		Filename: path.Join(nFlagDir, "main.go"),
		Tpl:      tpl.MainGo,
		Data: tpl.MainGoParam{
			ModuleName:  nFlagModName,
			LogFileName: nFlagDir,
		},
	})

	// api/routers.go
	tasks = append(tasks, GenTask{
		Name:     "routers",
		Filename: path.Join(apiDir, "routers.go"),
		Tpl:      tpl.RoutersGo,
		Data: tpl.RouterParam{
			ImportApps: parseAppName(),
			Apps:       nFlagApps,
			ModuleName: nFlagModName,
		},
	})

	// api/base/request.go
	tasks = append(tasks, GenTask{
		Name:     "baseRequest",
		Filename: path.Join(apiProtocol, "request.go"),
		Tpl:      tpl.BaseRequestGo,
		Data:     tpl.BackQuote{BQ: "`"},
	})

	db := _Postgres
	if nFlagDb == _Mysql {
		db = _Mysql
	}

	tasks = append(tasks, GenTask{
		Name:     "yamlConfig",
		Filename: path.Join(confDir, "config.yaml"),
		Tpl:      tpl.ConfigYaml,
		Data: tpl.ConfigYamlParam{
			Db:     db,
			DbName: nFlagDir,
		},
	})

	// model/model.go
	tasks = append(tasks, GenTask{
		Name:     "modelInit",
		Filename: path.Join(modelDir, "model.go"),
		Tpl:      tpl.ModelGo,
		Data: tpl.ModelParam{
			ModuleName:   nFlagModName,
			AppModelList: getRenderAppModel(),
			DbType:       db,
		},
	})

	// internal/config.go
	tasks = append(tasks, GenTask{
		Name:     "config",
		Filename: path.Join(internal, "config.go"),
		Tpl:      tpl.ConfigGo,
		Data:     tpl.BackQuote{BQ: "`"},
	})

	generateAllTemplateFiles(tasks)

	for _, i := range nFlagApps {
		appName := strings.ToLower(i)
		generateApp(modelDir, apiDir, appName, nFlagModName)
	}

	goModTidy(session)
	if err = swaggerInit(session); err != nil {
		console.Error("init swagger error:%s", err)
		return err
	}

	console.Info("Done!")
	return nil
}
