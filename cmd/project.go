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
	flagProjectDir string
	flagModName    string
	flagDb         string
	flagApps       []string
	flagNosqlList  []string
)

func init() {

	generateCmd.PersistentFlags().StringVarP(&flagProjectDir, "dir", "c", defaultDir, "the new project dir")
	generateCmd.PersistentFlags().StringVarP(&flagModName, "module", "m", "", "init go module name, (default same as \"${dir}\")")
	generateCmd.PersistentFlags().StringVarP(&flagDb, "db", "d", _PostgresFlag, "sql driver pg,mysql")
	generateCmd.PersistentFlags().StringSliceVarP(&flagApps, "apps", "a", []string{"user"}, "init apps, example -a user,file,category")
	generateCmd.PersistentFlags().StringSliceVarP(&flagNosqlList, "nosql", "n", []string{}, "init no sql driver -n=mongodb,es,redis")

	rootCmd.AddCommand(generateCmd)
}

func Generate() (err error) {
	console.Info("check param")
	if err = _ParamCheck(); err != nil {
		return err
	}
	session := newSession(flagProjectDir)
	goModInit(session)

	console.Info(fmt.Sprintf("start init project %s, db %+v, nosql %+v, module name %s, apps %+v",
		flagProjectDir, flagDb, flagNosqlList, flagModName, flagApps))
	// 目录已经生成
	apiDir := path.Join(flagProjectDir, "api")
	apiProtocol := path.Join(apiDir, "protocol")
	confDir := path.Join(flagProjectDir, "config")
	docsDir := path.Join(flagProjectDir, "docs")
	internal := path.Join(flagProjectDir, "internal")
	modelDir := path.Join(flagProjectDir, "model")
	// internalDb := path.Join(internal, "db")

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
		Filename: path.Join(flagProjectDir, "Dockerfile"),
		Tpl:      tpl.DockerFileNew,
		Data:     tpl.DockerFileParam{BinName: flagModName},
	}, {
		Name:     "gitignore",
		Filename: path.Join(flagProjectDir, ".gitignore"),
		Tpl:      tpl.GitIgnore,
	}, {
		Name:     "dockerIgnore",
		Filename: path.Join(flagProjectDir, ".dockerignore"),
		Tpl:      tpl.DockerIgnore,
	}, {
		Name:     "readme",
		Filename: path.Join(flagProjectDir, "README.md"),
		Tpl:      tpl.Readme,
	}, {
		Name:     "dockerBuildScript",
		Filename: path.Join(flagProjectDir, "t."+scriptSubfix),
		Tpl:      stpl,
		Data: tpl.DockerBuildParam{
			ContainerName: flagProjectDir + "-sc",
			ImageTag:      flagProjectDir + ":latest",
			BuildArg:      "--build-arg config=config",
			ExportPort:    "7788",
		},
	}, {
		Name:     "swaggerScript",
		Filename: path.Join(flagProjectDir, "gen-swagger."+scriptSubfix),
		Tpl:      tpl.GenSwagger,
	}}

	// main.go
	tasks = append(tasks, GenTask{
		Name:     "main",
		Filename: path.Join(flagProjectDir, "main.go"),
		Tpl:      tpl.MainGo,
		Data: tpl.MainGoParam{
			ModuleName: flagModName,
			// AppList:      parseAppName(),
			LogFileName:  flagProjectDir,
			AppModelList: getRenderAppModel(),
		},
	})

	// api/routers.go
	tasks = append(tasks, GenTask{
		Name:     "routers",
		Filename: path.Join(apiDir, "routers.go"),
		Tpl:      tpl.RoutersGo,
		Data: tpl.RouterParam{
			ImportApps: parseAppName(),
			Apps:       flagApps,
			ModuleName: flagModName,
		},
	})

	// model/base.go
	// tasks = append(tasks, GenTask{
	// 	Name:     "baseModel",
	// 	Filename: path.Join(apiBase, "model.go"),
	// 	Tpl:      tpl.BaseModelGo,
	// 	Data:     tpl.BackQuote{BQ: "`"},
	// })

	// api/base/request.go
	tasks = append(tasks, GenTask{
		Name:     "baseRequest",
		Filename: path.Join(apiProtocol, "request.go"),
		Tpl:      tpl.BaseRequestGo,
		Data:     tpl.BackQuote{BQ: "`"},
	})

	// db
	var (
		// dbTask GenTask
		cp = &tpl.ConfigParam{
			DB:      strings.Title(_Postgres),
			DbExist: true,
			BQ:      "`",
			JsonTag: strings.ToLower(_Postgres),
			DbType:  "_Postgres",
		}
	)
	if flagDb == _Mysql {
		// 	dbTask = GenTask{
		// 		Name:     _Mysql,
		// 		Filename: path.Join(internalDb, "mysql.go"),
		// 		Tpl:      tpl.DbMysql,
		// 		Data:     tpl.ModuleParam{ModuleName: flagModName},
		// 	}
		//
		cp.DB = strings.Title(_Mysql)
		cp.JsonTag = _Mysql
		cp.DbType = "_Mysql"
		//
	} else {
		// 	dbTask = GenTask{
		// 		Name:     _Postgres,
		// 		Filename: path.Join(internalDb, "postgres.go"),
		// 		Tpl:      tpl.DbPostgres,
		// 		Data:     tpl.ModuleParam{ModuleName: flagModName},
		// 	}
	}
	tasks = append(tasks, GenTask{
		Name:     "yamlConfig",
		Filename: path.Join(confDir, "config.yaml"),
		Tpl:      tpl.ConfigYaml,
		Data: tpl.ConfigYamlParam{
			Db:     cp.JsonTag,
			DbName: flagProjectDir,
		},
	})
	// tasks = append(tasks, dbTask)

	// internal/config.go
	tasks = append(tasks, GenTask{
		Name:     "config",
		Filename: path.Join(internal, "config.go"),
		Tpl:      tpl.ConfigGo,
		Data:     cp,
	})

	generateAllTemplateFiles(tasks)

	for _, i := range flagApps {
		appName := strings.ToLower(i)
		generateApp(modelDir, apiDir, appName, flagModName)
	}

	goModTidy(session)
	if err = swaggerInit(session); err != nil {
		console.Error("init swagger error:%s", err)
		return err
	}

	console.Info("Done!")
	return nil
}
