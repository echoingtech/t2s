package main

import (
	"flag"
	"fmt"
	"g.echo.tech/dev/go-micro/hakot/config"
	"g.echo.tech/dev/go-micro/hakot/mysql"
	"g.echo.tech/dev/go-micro/log4go"
	"g.echo.tech/dev/t2s/src"
	"gopkg.in/yaml.v2"
)

const (
	Template                = "%s:%s@tcp(%s:%s)/%s?charset=utf8"
	ConfigTypeDSN    string = "dsn"
	ConfigTypeApollo string = "apollo"
)

func main() {

	var (
		configType         string
		dbName, tableNames string
		user, pswd         string
		host, port         string
		outPath            string
		packageName        string
	)

	flag.StringVar(&user, "u", user, "mysql username (dsn必传模式)")

	flag.StringVar(&pswd, "p", pswd, "mysql password (dsn必传模式)")

	flag.StringVar(&host, "H", host, "mysql host (dsn模式必传)")

	flag.StringVar(&port, "P", port, "mysql port (dsn模式必传)")

	flag.StringVar(&dbName, "db", dbName, "数据库名")

	flag.StringVar(&tableNames, "t", tableNames, "表名, 多表按照 ',' 隔开")

	flag.StringVar(&outPath, "out", outPath, "生成文件地址,不指定则不生成")

	flag.StringVar(&packageName, "package", "po", "生成文件的packageName")

	flag.StringVar(&configType, "c", ConfigTypeApollo, "配置类型 apollo (从apollo读取mysql配置)  dsn (自定义mysql配置)")

	flag.Parse()

	var (
		err error
	)
	defer func() {
		if err != nil {
			flag.Usage()
		}
	}()

	if len(dbName) <= 0 || len(tableNames) <= 0 {
		err = log4go.Error("必须指定dbName, tableNames")
		return
	}

	var (
		mysqlStr string
	)
	switch configType {
	case ConfigTypeApollo:
		mysqlStr, err = LoadFromApollo(dbName)
		if err != nil {
			panic(err)
		}
	case ConfigTypeDSN:
		if len(host) <= 0 {
			err = log4go.Error("自定义模式必须指定 mysql host")
			return
		}
		if len(port) <= 0 {
			err = log4go.Error("自定义模式必须指定 mysql port")
			return
		}
		if len(user) <= 0 {
			err = log4go.Error("自定义模式必须指定 mysql user")
			return
		}
		if len(pswd) <= 0 {
			err = log4go.Error("自定义模式必须指定 mysql password")
			return
		}
		mysqlStr = LoadFromDSN(user, pswd, host, port, dbName)
	default:
		flag.Usage()
		return
	}

	fmt.Println("使用 ", configType, " 模式 convert table")

	func() {

		t2s := src.NewTable2Struct()
		if len(outPath) > 0 {
			t2s.SavePath(outPath)
		}

		if len(packageName) > 0 {
			t2s.PackageName(packageName)
		}

		err := t2s.Dsn(mysqlStr).
			Tables(tableNames).
			PackageName(packageName).
			Run()

		if err != nil {
			_ = log4go.Error("t2s error: %v", err)
		}

	}()

}

func LoadFromApollo(mysqlName string) (string, error) {

	client, err := config.NewConfigurator("", "", config.WithApolloAppIDAndDefaultPrefix("conn.mysql", "mysql"))
	if err != nil {
		panic(err)
	}

	res, err := client.Get(mysqlName)
	if err != nil {
		_ = log4go.Error("get apollo conf error: %v", err.Error())
		return "", err
	}

	var conf mysql.MysqlConfig
	err = yaml.Unmarshal(res, &conf)
	if err != nil {
		_ = log4go.Error("unmarshal mysql config error: %v", err.Error())
		return "", err
	}

	if len(conf.Instances) <= 0 {
		_ = log4go.Error("no instances find in : %v", mysqlName)
		return "", err
	}

	db := conf.Instances[0]

	dsn := fmt.Sprintf(Template, db.UserID, db.Password, db.Server, db.Port, db.DB)
	return dsn, nil

}

func LoadFromDSN(user, pswd, host, port, dbname string) string {
	server := host + ":" + port
	return fmt.Sprintf(Template, user, pswd, server, dbname)
}
