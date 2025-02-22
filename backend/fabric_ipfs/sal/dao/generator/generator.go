package main

import (
	"fabric_ipfs/sal/dao/generator/model"
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
)

func main() {
	modelGenerator()
	queryGenerator()
}

func modelGenerator() {
	g := gen.NewGenerator(gen.Config{
		FieldNullable: true,
		OutPath:       "./sal/dao/generator/model",
		Mode:          gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface,
	})

	db, _ := gorm.Open(mysql.Open("root:mysql_Grm7Rm@tcp(172.24.79.196:3306)/fabric_ebl?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	g.UseDB(db)

	g.GenerateModelAs("users", "UserPO", gen.FieldGenType("deleted_at", "gorm.DeletedAt"))
	g.GenerateModelAs("companys", "CompanyPO", gen.FieldGenType("deleted_at", "gorm.DeletedAt"))

	g.Execute()
}

func queryGenerator() {
	g := gen.NewGenerator(gen.Config{
		FieldNullable: true,
		OutPath:       "./sal/dao/generator/query",
		Mode:          gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface,
	})

	db, _ := gorm.Open(mysql.Open("root:mysql_Grm7Rm@tcp(172.24.79.196:3306)/fabric_ebl?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	g.UseDB(db)

	g.ApplyBasic(model.UserPO{})
	g.ApplyBasic(model.CompanyPO{})
	g.Execute()
}
