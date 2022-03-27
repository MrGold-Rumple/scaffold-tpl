/*
* @Author: Rumple
* @Email: ruipeng.wu@cyclone-robotics.com
* @DateTime: 2022/2/28 14:44
 */

package tpl

type MainGoParam struct {
	ModuleName  string `json:"module_name"`
	LogFileName string `json:"log_file_name"`
}

const MainGo = `
package main

import (
	"github.com/wuruipeng404/scaffold"
	"github.com/wuruipeng404/scaffold/logger"
	"{{.ModuleName}}/api"
)

func init() {
	logger.InitLogger("log/{{.LogFileName}}.log")
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
}

const ConfigGo = `
package internal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"

	"github.com/wuruipeng404/scaffold/orm"
	"github.com/wuruipeng404/scaffold/util"
	"gopkg.in/yaml.v2"
)

var sc *ServerConfig

type ServerConfig struct {
	Database *_Database {{.BQ}}yaml:"database" json:"database"{{.BQ}}
}
type _Database struct {
	Type orm.DbType {{.BQ}}yaml:"type" json:"type"{{.BQ}}
	User string     {{.BQ}}yaml:"user" json:"user"{{.BQ}}
	Pass string     {{.BQ}}yaml:"pass" json:"pass"{{.BQ}}
	Host string     {{.BQ}}yaml:"host" json:"host"{{.BQ}}
	Port int        {{.BQ}}yaml:"port" json:"port"{{.BQ}}
	DB   string     {{.BQ}}yaml:"db" json:"db"{{.BQ}}
}

func (c *ServerConfig) String() string {
	b, _ := json.MarshalIndent(c, "", "  ")
	return string(b)
}

func init() {

	yf, err := ioutil.ReadFile(util.Env("CONFIG_FILE", "./config/config.yaml"))

	if err == nil {
		_ = yaml.Unmarshal(yf, &sc)
	}

	if sc == nil {
		sc = &ServerConfig{}
	}

	if sc.Database == nil {
		sc.Database = &_Database{}
	}

	dbType := util.Env("DB_TYPE", string(orm.MySQL))
	dbUser := util.Env("DB_USER", "")
	dbPass := util.Env("DB_PASS", "")
	dbHost := util.Env("DB_HOST", "")
	dbPort := util.Env("DB_PORT", "")
	dbDb := util.Env("DB_NAME", "")

	// environment variables take precedence over configuration files if env is not null
	if dbType != "" {
		sc.Database.Type = orm.DbType(dbType)
	}

	if dbUser != "" {
		sc.Database.User = dbUser
	}

	if dbPass != "" {
		sc.Database.Pass = dbPass
	}

	if dbHost != "" {
		sc.Database.Host = dbHost
	}

	if dbPort != "" {
		dp, err := strconv.Atoi(dbPort)
		if err != nil {
			log.Fatalf("get db port error:%s", err)
		}
		sc.Database.Port = dp
	}

	if dbDb != "" {
		sc.Database.DB = dbDb
	}

	log.Println(fmt.Sprintf("reading config :%+v", sc))
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
	{{- range $i,$v := .ImportApps}}
	"{{$v -}}"
	{{- end}}
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

func _ApiRouters() {
	api := server.Group("/api")
	v1 := api.Group("/v1")

	{{- range $i,$v := .Apps}}
	_{{$v -}}Routers(v1)
	{{- end}}
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

type ModelParam struct {
	ModuleName   string   `json:"module_name"`
	AppModelList []string `json:"app_model_list"`
	DbType       string   `json:"db_type"`
}

const ModelGo = `
package model

import (
	"log"

	"github.com/wuruipeng404/scaffold/orm"
	"{{.ModuleName}}/internal"
)

func init() {

	dc := internal.Config().Database

	orm.Init(&orm.InitOption{
		{{- if eq .DbType "mysql"}}
		Type:   orm.MySQL,
		{{- else}}
		Type:   orm.Postgres,
		{{- end}}
		User:   dc.User,
		Pass:   dc.Pass,
		Host:   dc.Host,
		Port:   dc.Port,
		DbName: dc.DB,
	})

	migrate()
}

func migrate() {
	if err := orm.C().AutoMigrate(
		{{- range $i,$v := .AppModelList}}
		new({{$v -}}),
		{{- end}}
	); err != nil {
		log.Fatalf("database auto migrate error:%s", err)
	}
	log.Println("database auto migrate success ~~")
}
`
