/*
* @Author: Rumple
* @Email: ruipeng.wu@cyclone-robotics.com
* @DateTime: 2022/2/28 14:44
 */

package tpl

type MainGoParam struct {
	ModuleName   string   `json:"module_name"`
	AppList      []string `json:"app_list"`
	LogFileName  string   `json:"log_file_name"`
	AppModelList []string `json:"app_model_list"`
}

const MainGo = `
package main

import (
	"log"

	"github.com/wuruipeng404/scaffold"
	"github.com/wuruipeng404/scaffold/logger"
	"github.com/wuruipeng404/scaffold/orm"
	"{{.ModuleName}}/model"
	"{{.ModuleName}}/api"
)

func init() {
	logger.InitLogger("log/{{.LogFileName}}.log")

	if err := orm.Init(&orm.InitOption{
		Type:   "",
		User:   "",
		Pass:   "",
		DbName: "",
		Host:   "",
		Port:   0,
	}); err != nil {
		log.Fatalf("init orm error:%s", err)
	}

	migrate()
}

func migrate() {
	if err := orm.C().AutoMigrate(
		{{range $i,$v := .AppModelList}}
		new({{$v}}),
		{{end}}
	); err != nil {
		logger.Fatalf("数据库迁移失败:%s", err)
	}
	logger.Info("数据库迁移成功~")
}

// @title {{.ModuleName}} API Document
// @version 1.0
// @host localhost:8000
// @BasePath /api/v1
// @securityDefinitions.apikey  ApiKeyAuth
// @in                          header
// @name                        Authorization
func main() {
	server := scaffold.NewGraceServer(":8000", api.InitEngine())
	server.Start()
}
`

type ConfigParam struct {
	DB      string // Postgres or Mysql
	DbExist bool   //
	BQ      string // "`"
	JsonTag string // lowercase(db)
	DbType  string // _Postgres  _Mysql
}

const ConfigGo = `
package internal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

var sc *ServerConfig

type ServerConfig struct {
	{{if .DbExist}}
	{{.DB}} {{.DbType}} {{.BQ}}yaml:"{{.JsonTag}}" json:"{{.JsonTag}}"{{.BQ}}
	{{end}}
}

{{if .DbExist}}
type {{.DbType}} struct {
	User string {{.BQ}}yaml:"user" json:"user"{{.BQ}}
	Pass string {{.BQ}}yaml:"pass" json:"pass"{{.BQ}}
	Host string {{.BQ}}yaml:"host" json:"host"{{.BQ}}
	Port int    {{.BQ}}yaml:"port" json:"port"{{.BQ}}
	DB   string {{.BQ}}yaml:"db" json:"db"{{.BQ}}
}
{{end}}

func (c *ServerConfig) String() string {
	b, _ := json.MarshalIndent(c, "", "  ")
	return string(b)
}

func init() {

	yf, err := ioutil.ReadFile("./config/config.yaml")
	if err != nil {
		log.Fatalf("读取配置文件失败:%s", err)
	} else {
		if err = yaml.Unmarshal(yf, &sc); err != nil {
			log.Fatalf("解析配置文件失败:%s", err)
		}
	}

	log.Println(fmt.Sprintf("读取配置成功:%+v", sc))
}

func Config() *ServerConfig {
	return sc
}
`

type ModuleParam struct {
	ModuleName string
}

type RouterParam struct {
	ImportApps []string
	Apps       []string
	ModuleName string
}

const RoutersGo = `
package api

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/wuruipeng404/scaffold"
	{{range $i,$v := .ImportApps}}
	"{{$v}}"
	{{end}}
	_ "{{.ModuleName}}/docs"
)

var server *gin.Engine

func InitEngine() *gin.Engine {
	server = gin.New()
	server.Use(gin.Recovery(), scaffold.GracefulLogger(), scaffold.Cors())
	// server.LoadHTMLGlob("template/*.html")
	// engine.Static("/static", "./static")
	// server.Static("/static", "./template/static")
	swaggerRouters()
	_ApiRouters()
	return server
}

func swaggerRouters() {
	server.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

// 接口
func _ApiRouters() {
	api := server.Group("/api")
	v1 := api.Group("/v1")

	{{range $i,$v := .Apps}}
	_{{$v}}Routers(v1)
	{{end}}
}

{{range $i,$v := .Apps}}
func _{{$v}}Routers(group *gin.RouterGroup) {

	route := group.Group("/{{$v}}")
	control := new({{$v}}.Controller)

	route.GET("/", control.Get)
	route.GET("/list", control.List)
	route.POST("/", control.Create)
	route.PUT("/", control.Update)
	route.DELETE("/", control.Delete)
}
{{end}}
`

type BackQuote struct {
	BQ string
}

const BaseRequestGo = `
package protocol

type PageParam struct {
	Page     int {{.BQ}}form:"page,default=1" binding:"gt=0"{{.BQ}}
	PageSize int {{.BQ}}form:"page_size,default=10" binding:"gt=0"{{.BQ}}
}

type SearchPageParam struct {
	PageParam
	Key string {{.BQ}}form:"key"{{.BQ}}
}
`

type GoPkgFileParam struct {
	PkgName string
	Comment string
}

const GoOnlyPkgFile = `
package {{.PkgName}}

{{.Comment}}
`
