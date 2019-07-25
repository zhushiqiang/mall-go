package main

import (
	"encoding/json"
	"flag"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/yakaa/grpcx"
	"github.com/yakaa/log4g"
	"google.golang.org/grpc"
	"io/ioutil"
	"mall-go/integral/command/rpc/config"
	"mall-go/integral/logic"
	"mall-go/integral/model"
	"mall-go/integral/protos"
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

	log4g.Init(log4g.Config{Path:"logs",Stdout:true})

	engine, err := xorm.NewEngine("mysql", conf.Mysql.DataSource)

	if err != nil{
		log4g.Info(err)
	}

	client := redis.NewClient(&redis.Options{
		Addr:     conf.Redis.DataSource,
		Password: conf.Redis.Auth, // no password set
	})
	integralModel := model.NewIntegralModel(engine, client, conf.Mysql.Table.Integral)

	integralServerLogic,err := logic.NewIntegralLogic(conf.RabbitMq.DataSource + conf.RabbitMq.VirtualHost, conf.RabbitMq.QueueName, integralModel)
	if err != nil{
		log4g.Info(err)
	}
	rpcServer, err := grpcx.MustNewGrpcxServer(conf.RpcServerConfig, func(server *grpc.Server) {
		protos.RegisterIntegralRpcServer(server, integralServerLogic)
	})
	if err != nil {
		log4g.Info(err)
	}
	//测试代码
	//integralServerLogic.PushMessage("INSERT INTO integral (user_id, integral) VALUES (22,10);")
	integralServerLogic.ConsumeMessage();

	defer integralServerLogic.CloseRabbitMqConn()
	log4g.Error(rpcServer.Run())
}