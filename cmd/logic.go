package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"text/template"

	"github.com/codeskyblue/go-sh"
	"github.com/wuruipeng404/scaffold-tpl/console"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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
	checkDir(nFlagDir)

	// go mod 名称等于空时 默认等于文件夹名称
	if nFlagModName == "" {
		nFlagModName = nFlagDir
	}

	// 检察数据库参数
	switch nFlagDb {
	case _PostgresFlag, _Mysql, _Postgres:
	default:
		return errors.New("db only support pg(postgres) and mysql")
	}
	// 需要初始化的app
	var lowerApps []string
	for _, i := range nFlagApps {
		lowerApps = append(lowerApps, strings.ToLower(i))
	}
	nFlagApps = lowerApps
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

func generateFile(tplName, filename, tpl string, data any) {
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
	if err = session.Command("go", "mod", "init", nFlagModName).Run(); err != nil {
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
	for _, n := range nFlagApps {
		result = append(result, fmt.Sprintf("%s/apps/api/%s", nFlagModName, n))
	}
	return result
}

func getRenderAppModel() []string {
	var result []string
	for _, n := range nFlagApps {
		result = append(result, Title(n))
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

func Title(v string) string {
	c := cases.Title(language.Und)
	return c.String(v)
}

func File2lines(filePath string) ([]string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = f.Close()
	}()
	return LinesFromReader(f)
}

func LinesFromReader(r io.Reader) ([]string, error) {
	var lines []string
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

func InsertStringToFile(path, str string, index int) error {
	lines, err := File2lines(path)
	if err != nil {
		return err
	}

	fileContent := ""
	for i, line := range lines {
		if i == index {
			fileContent += str
		}
		fileContent += line
		fileContent += "\n"
	}

	return ioutil.WriteFile(path, []byte(fileContent), defaultFilePerm)
}
