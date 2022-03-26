package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"text/template"

	"github.com/codeskyblue/go-sh"
	"github.com/wuruipeng404/scaffold-tpl/console"
)

type GenTask struct {
	Name     string
	Filename string
	Tpl      string
	Data     any
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

	// go mod 名称等于空时 默认等于文件夹名称
	if flagModName == "" {
		flagModName = flagProjectDir
	}

	// 检察数据库参数
	switch flagDb {
	case _PostgresFlag, _Mysql, _Postgres:
	default:
		return errors.New("db only support pg(postgres) and mysql")
	}
	// 需要初始化的app
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

func generateAllTemplateFiles(tasks []GenTask) {

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
		result = append(result, fmt.Sprintf("model.%s", strings.Title(n)))
	}
	return result
}

func swaggerInit(session *sh.Session) (err error) {

	if err = session.Command("swag", "-v").Run(); err != nil {
		if err = session.Command("go", "install", "github.com/swaggo/swag/cmd/swag@latest").Run(); err != nil {
			return fmt.Errorf("install swag cmd error:%s", err)
		}
	}

	if err = session.Command("swag", "init", "--parseDependency", "--parseInternal", "--parseDepth", "3").Run(); err != nil {
		return fmt.Errorf("generate swagger file error:%s", err)
	}
	return nil
}
