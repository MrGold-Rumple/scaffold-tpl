/*
* @Author: Rumple
* @Email: ruipeng.wu@cyclone-robotics.com
* @DateTime: 2022/3/2 9:52
 */

package tpl

type AppParam struct {
	AppName    string // user
	AppTitle   string // User
	BQ         string // "`"
	ModuleName string //
}

const AppProtocolGO = `
package {{.AppName}}

import (
	"github.com/wuruipeng404/scaffold/swag"
	"{{.ModuleName}}/model"
)

type (
	_{{.AppTitle}}GetParam struct {
		Id uint {{.BQ}}form:"id" binding:"required"{{.BQ}}
	}

	_{{.AppTitle}}CreateParam struct {
	}

	_{{.AppTitle}}UpdateParam struct {
	}

	_{{.AppTitle}}DeleteParam struct {
	}
)

type (
	_SwagGet{{.AppTitle}}Resp struct {
		swag.CodeMsg
		Data model.{{.AppTitle}} {{.BQ}}json:"data"{{.BQ}}
	}

	_SwagList{{.AppTitle}}Resp struct {
		swag.CodeMsg
		swag.Page
		Data []model.{{.AppTitle}} {{.BQ}}json:"data"{{.BQ}}
	}
)

`

const AppModelGO = `
package model

import (
	"github.com/wuruipeng404/scaffold/orm"
)

type {{.AppTitle}} struct {
	orm.Model
}
`

const AppControllerGo = `
package {{.AppName}}

import (
	"github.com/gin-gonic/gin"
	"github.com/wuruipeng404/scaffold"
	"{{.ModuleName}}/api/protocol"
)

type Controller struct {
	scaffold.BeautyController
}

// Get
// @Summary get {{.AppTitle}} by id
// @Description
// @Security ApiKeyAuth
// @Tags {{.AppTitle}}
// @Accept json
// @Produce json
// @Param _{{.AppTitle}}GetParam query _{{.AppTitle}}GetParam true "get parameter"
// @Success 200 {object} _SwagGet{{.AppTitle}}Resp "request success"
// @Failure 400 {object} swag.CodeMsg "request failed"
// @Router /{{.AppName}} [get]
func (c *Controller) Get(ctx *gin.Context) {
	var (
		err error
		p   _{{.AppTitle}}GetParam
	)

	if err = ctx.ShouldBindQuery(&p); err != nil {
		c.FailedE(ctx, err)
		return
	}

}

// List
// @Summary list or search {{.AppTitle}}
// @Description list or search
// @Security ApiKeyAuth
// @Tags {{.AppTitle}}
// @Accept json
// @Produce json
// @Param protocol.SearchPageParam query protocol.SearchPageParam true "search list param"
// @Success 200 {object} _SwagList{{.AppTitle}}Resp "request success"
// @Failure 400 {object} swag.CodeMsg "request failed"
// @Router /{{.AppName}}/list [get]
func (c *Controller) List(ctx *gin.Context) {
	var (
		err error
		p   protocol.SearchPageParam
	)

	if err = ctx.ShouldBindQuery(&p); err != nil {
		c.FailedE(ctx, err)
		return
	}

}

// Create
// @Summary creates a new {{.AppTitle}} resource
// @Description
// @Security ApiKeyAuth
// @Tags {{.AppTitle}}
// @Accept json
// @Produce json
// @Param _{{.AppTitle}}CreateParam body _{{.AppTitle}}CreateParam true "create parameter"
// @Success 200 {object} _SwagGet{{.AppTitle}}Resp "request success"
// @Failure 400 {object} swag.CodeMsg "request failed"
// @Router /{{.AppName}} [post]
func (c *Controller) Create(ctx *gin.Context) {
	var (
		err error
		p   _{{.AppTitle}}CreateParam
	)

	if err = ctx.ShouldBindJSON(&p); err != nil {
		c.FailedE(ctx, err)
		return
	}

}

// Update
// @Summary update {{.AppTitle}} resource
// @Description
// @Security ApiKeyAuth
// @Tags {{.AppTitle}}
// @Accept json
// @Produce json
// @Param _{{.AppTitle}}UpdateParam body _{{.AppTitle}}UpdateParam true "update parameter"
// @Success 200 {object} _SwagGet{{.AppTitle}}Resp "request success"
// @Failure 400 {object} swag.CodeMsg "request failed"
// @Router /{{.AppName}} [put]
func (c *Controller) Update(ctx *gin.Context) {
	var (
		err error
		p   _{{.AppTitle}}UpdateParam
	)
	if err = ctx.ShouldBindJSON(&p); err != nil {
		c.FailedE(ctx, err)
		return
	}
}

// Delete
// @Summary delete {{.AppTitle}} resource
// @Description
// @Security ApiKeyAuth
// @Tags {{.AppTitle}}
// @Accept json
// @Produce json
// @Param _{{.AppTitle}}DeleteParam query _{{.AppTitle}}DeleteParam true "delete parameter"
// @Success 200 {object} swag.CodeMsg "request success"
// @Failure 400 {object} swag.CodeMsg "request failed"
// @Router /{{.AppName}} [delete]
func (c *Controller) Delete(ctx *gin.Context) {
	var (
		err error
		p   _{{.AppTitle}}DeleteParam
	)
	if err = ctx.ShouldBindQuery(&p); err != nil {
		c.FailedE(ctx, err)
		return
	}
}
`
