/*
* @Author: Rumple
* @Email: ruipeng.wu@cyclone-robotics.com
* @DateTime: 2022/2/21 17:56
 */

package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"strings"
	"text/template"

	"github.com/codeskyblue/go-sh"
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

// 检查文件夹是否存在, 不存在则创建, 存在且不为空则报错.
func checkDir(cd string) {
	fs, err := ioutil.ReadDir(cd)
	if err == nil {
		if len(fs) > 0 {
			console.Fatal("%s already exists and not empty", cd)
		}
	} else {
		if err = os.MkdirAll(cd, defaultFilePerm); err != nil {
			console.Fatal("create dir %s error:%s", cd, err)
		}
	}
}

func _ParamCheck() error {
	checkDir(flagProjectDir)

	if flagModName == "" {
		flagModName = flagProjectDir
	}

	switch flagDb {
	case _PostgresFlag, _Mysql, _Postgres:
	default:
		return errors.New("db only support pg(postgres) and mysql")
	}

	var lowerApps []string
	for _, i := range flagApps {
		lowerApps = append(lowerApps, strings.ToLower(i))
	}
	flagApps = lowerApps
	return nil
}

func newSession(dir string) *sh.Session {

	session := sh.NewSession()
	session.PipeFail = true
	session.PipeStdErrors = true
	session.SetDir(dir)

	return session
}

func newFile(filename string, content string) {
	if err := ioutil.WriteFile(filename, []byte(content), defaultFilePerm); err != nil {
		console.Fatal("create file %s error %s", filename, err)
	}
}

func generateFile(tplName, filename, tpl string, data interface{}) {
	var err error
	defer func() {
		if err != nil {
			console.Fatal(err.Error())
		}
	}()

	tmpl, err := template.New(tplName).Parse(tpl)
	if err != nil {
		err = fmt.Errorf("parse template %s error %s", filename, err)
		return
	}

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, defaultFilePerm)
	if err != nil {
		err = fmt.Errorf("open file %s error %s", filename, err)
		return
	}

	if err = tmpl.Execute(file, data); err != nil {
		err = fmt.Errorf("render template %s error %s", filename, err)
		return
	}
	return
}

type Task struct {
	Name     string
	Filename string
	Tpl      string
	Data     interface{}
}

func generateAllTemplateFiles(tasks []Task) {

	for _, t := range tasks {
		console.Info("generate template", t.Name)
		generateFile(t.Name, t.Filename, t.Tpl, t.Data)
	}

}

func makeAllDir(dirs []string) {
	var err error

	for _, d := range dirs {
		if err = os.MkdirAll(d, defaultFilePerm); err != nil {
			console.Fatal("make dir %s error %s", d, err)
		}
	}
}

func goModInit(session *sh.Session) {
	var err error
	console.Info("start init git repo and go mod")
	if err = session.Command("go", "mod", "init", flagModName).Run(); err != nil {
		console.Fatal("go mod init error:%s", err)
	}

	if err = session.Command("git", "init").Run(); err != nil {
		console.Fatal("git init error:%s", err)
	}
}

func goModTidy(session *sh.Session) {
	console.Info("start download golang package")

	if err := session.Command("go", "mod", "tidy").Run(); err != nil {
		console.Fatal("go mod tidy error:%s", err)
	}
	if err := session.Command("go", "mod", "vendor").Run(); err != nil {
		console.Fatal("go mod vendor error:%s", err)
	}
}

func parseAppName() []string {
	var result []string
	for _, n := range flagApps {
		result = append(result, fmt.Sprintf("%s/api/%s", flagModName, n))
	}
	return result
}

func getRenderAppModel() []string {
	var result []string
	for _, n := range flagApps {
		result = append(result, fmt.Sprintf("%s.%s", n, strings.Title(n)))
	}
	return result
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
	apiBase := path.Join(apiDir, "base")
	// middleware := path.Join(apiDir, "middleware.go")
	// routers := path.Join(apiDir, "routers.go")
	confDir := path.Join(flagProjectDir, "config")
	docsDir := path.Join(flagProjectDir, "docs")
	internal := path.Join(flagProjectDir, "internal")
	internalDb := path.Join(internal, "db")

	needMkDirs := []string{apiBase, confDir, docsDir, internalDb}

	console.Info(fmt.Sprintf("start create all dir %+v", needMkDirs))
	makeAllDir(needMkDirs)

	// newFile(path.Join(confDir, "config.yaml"), "")

	scriptSubfix := "ps1"

	if runtime.GOOS != "windows" {
		scriptSubfix = "sh"
	}

	tasks := []Task{{
		Name:     "Dockerfile",
		Filename: path.Join(flagProjectDir, "Dockerfile"),
		Tpl:      tpl.DockerFile,
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
		Tpl:      tpl.BuildScript,
		Data: tpl.DockerBuildParam{
			ContainerName: flagProjectDir + "-crt",
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
	tasks = append(tasks, Task{
		Name:     "main",
		Filename: path.Join(flagProjectDir, "main.go"),
		Tpl:      tpl.MainGo,
		Data: tpl.MainGoParam{
			ModuleName:   flagModName,
			AppList:      parseAppName(),
			LogFileName:  flagProjectDir,
			AppModelList: getRenderAppModel(),
		},
	})

	// api/routers.go
	tasks = append(tasks, Task{
		Name:     "routers",
		Filename: path.Join(apiDir, "routers.go"),
		Tpl:      tpl.RoutersGo,
		Data: tpl.RouterParam{
			ImportApps: parseAppName(),
			Apps:       flagApps,
			ModuleName: flagModName,
		},
	})

	// api/base/model.go
	tasks = append(tasks, Task{
		Name:     "baseModel",
		Filename: path.Join(apiBase, "model.go"),
		Tpl:      tpl.BaseModelGo,
		Data:     tpl.BackQuote{BQ: "`"},
	})

	// api/base/request.go
	tasks = append(tasks, Task{
		Name:     "baseRequest",
		Filename: path.Join(apiBase, "request.go"),
		Tpl:      tpl.BaseRequestGo,
		Data:     tpl.BackQuote{BQ: "`"},
	})

	// db
	var (
		dbTask Task
		cp     = &tpl.ConfigParam{
			DB:      strings.Title(_Postgres),
			DbExist: true,
			BQ:      "`",
			JsonTag: strings.ToLower(_Postgres),
			DbType:  "_Postgres",
		}
	)

	if flagDb == _Mysql {
		dbTask = Task{
			Name:     _Mysql,
			Filename: path.Join(internalDb, "mysql.go"),
			Tpl:      tpl.DbMysql,
			Data:     tpl.ModuleParam{ModuleName: flagModName},
		}

		cp.DB = strings.Title(_Mysql)
		cp.JsonTag = _Mysql
		cp.DbType = "_Mysql"

	} else {
		dbTask = Task{
			Name:     _Postgres,
			Filename: path.Join(internalDb, "postgres.go"),
			Tpl:      tpl.DbPostgres,
			Data:     tpl.ModuleParam{ModuleName: flagModName},
		}
	}
	tasks = append(tasks, Task{
		Name:     "yamlConfig",
		Filename: path.Join(confDir, "config.yaml"),
		Tpl:      tpl.ConfigYaml,
		Data: tpl.ConfigYamlParam{
			Db:     cp.JsonTag,
			DbName: flagProjectDir,
		},
	})
	tasks = append(tasks, dbTask)

	// internal/config.go
	tasks = append(tasks, Task{
		Name:     "config",
		Filename: path.Join(internal, "config.go"),
		Tpl:      tpl.ConfigGo,
		Data:     cp,
	})

	generateAllTemplateFiles(tasks)

	for _, i := range flagApps {
		appName := strings.ToLower(i)
		generateApp(apiDir, appName, flagModName)
	}

	// goModTidy(session)
	//
	// swaggerInit(session)

	console.Info("Done!")
	return nil
}

func swaggerInit(session *sh.Session) {
	if err := session.Command("swag", "init", "--parseDependency", "--parseInternal", "--parseDepth", "3").Run(); err != nil {
		console.Warn(fmt.Sprintf("generate swagger file error:%s", err))
	}
}
