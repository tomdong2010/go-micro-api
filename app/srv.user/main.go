/**
 * Created by Goland.
 * User: yan.wang5
 * Date: 2019/9/6
 */
package main

import (
	"demo/app/srv.user/conf"
	"demo/app/srv.user/handler"
	"demo/cmn"
	"demo/mdw"
	"demo/proto"
	"demo/utility/db"
	"fmt"
	"github.com/lestrrat-go/file-rotatelogs"
	"github.com/micro/go-micro"
	"github.com/sirupsen/logrus"
	"os"
)

func main() {
	service := micro.NewService(
		micro.Name(cmn.APP_NAME_PREFIX + cmn.APP_SRV_USER),
		micro.WrapHandler(mdw.LogMdw),
	)

	service.Init( )

	// 初始化公共配置文件
	checkErr("InitCommonConfig", cmn.InitConfig(service.Server()))
	fmt.Print("InitCommonConfig Success!!!\n")

	// 初始化app配置文件
	checkErr("InitAppConfig", conf.InitConfig(service.Server(), cmn.APP_SRV_USER))
	fmt.Print("InitAppConfig Success!!!\n")

	// 创建文件日志，按天分割，日志文件仅保留一周
	w, err := rotatelogs.New(conf.GetLogPath())
	checkErr("CreateRotateLog", err)

	// 设置日志
	logrus.SetOutput(w)
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetReportCaller(true)

	// 启动mysql
	defer db.CloseMysql()
	fmt.Print("InitMysql...\r")
	checkErr("InitMysql", db.InitMysql(cmn.GetMysqlConfig()))
	fmt.Print("InitMysql Success!!!\n")

	// 启动redis
	defer db.CloseRedis()
	fmt.Print("InitRedis...\r")
	checkErr("InitRedis", db.InitRedis(cmn.GetRedisConfig()))
	fmt.Print("InitRedis Success!!!\n")

	_ = proto.RegisterGreeterHandler(service.Server(), new(handler.Greeter))

	checkErr("server run", service.Run())
}

func checkErr(errMsg string, err error) {
	if err != nil {
		fmt.Printf("%s Error: %v\n", errMsg, err)
		os.Exit(1)
	}
}
