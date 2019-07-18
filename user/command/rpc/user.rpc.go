package main

import (
	"encoding/json"
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/yakaa/grpcx"
	"github.com/yakaa/log4g"
	"google.golang.org/grpc"
	"io/ioutil"
	"mall-go/user/command/rpc/config"
	"mall-go/user/logic"
	"mall-go/user/model"
	"mall-go/user/protos"
)

var configFile = flag.String("f", "config/config.json", "use config")

func main() {

	flag.Parse()
	conf := new(config.Config)
	bs,err := ioutil.ReadFile(*configFile)

	if err !=nil {
		log4g.Info(err)
	}
	if err := json.Unmarshal(bs, conf); err != nil {
		log4g.Info(err)
	}

	//if conf.Mode == gin.ReleaseMode {
	log4g.Init(log4g.Config{Path:"logs",Stdout:true})
	gin.DefaultWriter = log4g.InfoLog
	gin.DefaultErrorWriter = log4g.ErrorLog
	//}

	engine, err := xorm.NewEngine("mysql", conf.Mysql.DataSource)

	if err != err{
		log4g.Info(err)
	}

	client := redis.NewClient(&redis.Options{
		Addr:     conf.Redis.DataSource,
		Password: conf.Redis.Auth, // no password set
	})
	userModel := model.NewUserModel(engine, client, conf.Mysql.Table.User)

	userServerLogic := logic.NewUserRpcServiceLogic(userModel)

	rpcServer, err := grpcx.MustNewGrpcxServer(conf.RpcServerConfig, func(server *grpc.Server) {
		protos.RegisterUserRpcServer(server, userServerLogic)
	})
	if err != nil {
		log4g.Info(err)
	}



	log4g.Error(rpcServer.Run())
}