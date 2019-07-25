package main

import (
	"encoding/json"
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/yakaa/log4g"
	"io/ioutil"
	"mall-go/user/command/api/config"
	"mall-go/user/controller"
	"mall-go/user/logic"
	"mall-go/user/model"
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

	if err != nil{
		log4g.Info(err)
	}

	client := redis.NewClient(&redis.Options{
		Addr:     conf.Redis.DataSource,
		Password: conf.Redis.Auth, // no password set
	})
	userModel := model.NewUserModel(engine, client, conf.Mysql.Table.User)
	//userLogic := logic.NewUserLogic{userModel}
	userLogic := logic.NewUserLogic(userModel, client)
	userController := controller.NewUserController(userLogic)
	log4g.Init(log4g.Config{Path: "logs"})
	gin.DefaultWriter = log4g.InfoLog
	gin.DefaultErrorWriter = log4g.ErrorLog

	r := gin.Default()

	userRouteGroup := r.Group("/user")
	{
		userRouteGroup.POST("/register", userController.Register)
		userRouteGroup.POST("/login", userController.Login)
	}

	log4g.Error(r.Run(conf.Port)) // listen and serve on 0.0.0.0:8080
}