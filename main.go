package main

import (
	"fmt"

	"github.com/echoingtech/t2s/src"
)

func main() {

	err := src.NewTable2Struct().SavePath("./model/model.go").
		Dsn("user_im:7nJd*JL0FxZFPqfB@tcp(rm-uf6f0jy4vg5ua2bsu.mysql.rds.aliyuncs.com)/svc-im?charset=utf8").
		Table("user_syncs").
		PackageName("po").
		Run()

	fmt.Println(err)
}
