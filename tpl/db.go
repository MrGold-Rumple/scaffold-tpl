/*
* @Author: Rumple
* @Email: ruipeng.wu@cyclone-robotics.com
* @DateTime: 2022/3/1 16:06
 */

package tpl

// const DbPostgres = `
// package db
//
// import (
// 	"fmt"
// 	"log"
// 	"time"
//
// 	"gorm.io/driver/postgres"
// 	"gorm.io/gorm"
// 	"gorm.io/gorm/logger"
// 	"{{.ModuleName}}/internal"
// )
//
// var (
// 	orm *gorm.DB
// )
//
// func init() {
// 	initDB()
// }
//
// func initDB() {
// 	var err error
//
// 	pc := internal.Config().Postgres
//
// 	dsn := fmt.Sprintf(
// 		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Shanghai",
// 		pc.Host, pc.Port, pc.User, pc.Pass, pc.DB,
// 	)
//
// 	if orm, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
// 		PrepareStmt: true,
// 		Logger:      logger.Default.LogMode(logger.Silent),
// 	}); err != nil {
// 		log.Fatalf("create orm connect error:%s", err)
// 	}
//
// 	sqlDb, err := orm.DB()
// 	if err != nil {
// 		log.Fatalf("get sql error:%s", err)
// 	}
//
// 	sqlDb.SetMaxIdleConns(10)
// 	sqlDb.SetMaxOpenConns(100)
// 	sqlDb.SetConnMaxLifetime(time.Hour)
//
// 	log.Println("connect postgres success")
// }
//
// func ORM() *gorm.DB {
// 	return orm
// }
// `
//
// const DbMysql = `
// package db
//
// import (
// 	"fmt"
// 	"log"
// 	"time"
//
// 	"gorm.io/driver/mysql"
// 	"gorm.io/gorm"
// 	"gorm.io/gorm/logger"
// 	"{{.ModuleName}}/internal"
// )
//
// var (
// 	orm *gorm.DB
// )
//
// func init() {
// 	initDB()
// }
//
// func initDB() {
// 	var err error
//
// 	mc := internal.Config().Mysql
//
// 	dsn := fmt.Sprintf(
// 		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
// 		mc.User, mc.Pass, mc.Host, mc.Port, mc.DB,
// 	)
//
// 	if orm, err = gorm.Open(mysql.New(mysql.Config{
// 		DSN:                     dsn,
// 		DefaultStringSize:       256,
// 		DontSupportRenameIndex:  true,
// 		DontSupportRenameColumn: true,
// 	}), &gorm.Config{
// 		PrepareStmt: true,
// 		Logger:      logger.Default.LogMode(logger.Silent),
// 	}); err != nil {
// 		log.Fatalf("create orm connect error:%s", err)
// 	}
//
// 	sqlDb, err := orm.DB()
// 	if err != nil {
// 		log.Fatalf("get sql error:%s", err)
// 	}
// 	sqlDb.SetMaxIdleConns(10)
// 	sqlDb.SetMaxOpenConns(100)
// 	sqlDb.SetConnMaxLifetime(time.Hour)
// 	log.Println("connect mysql success")
// }
//
// func ORM() *gorm.DB {
// 	return orm
// }
// `

const DbRedis = `

`
const DbMongo = `

`
