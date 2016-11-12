package main

import (
	"time"

	"github.com/micro/cli"
	"github.com/micro/go-micro"

	"github.com/unicok/auth-srv/db"
	"github.com/unicok/auth-srv/handler"
	proto "github.com/unicok/auth-srv/proto/auth"
	"github.com/unicok/misc/log"
)

func main() {
	service := micro.NewService(
		micro.Name("com.unicok.srv.auth"),
		micro.Version("latest"),
		micro.RegisterTTL(time.Minute),
		micro.RegisterInterval(time.Second*30),

		micro.Flags(
			cli.StringFlag{
				Name:   "mongodb_url",
				EnvVar: "MONGODB_URL",
				Usage:  "The mongodb URL.",
				Value:  "mongodb://username:password@127.0.0.1:27017/account",
			},
		),

		micro.Action(func(c *cli.Context) {

			if len(c.String("mongodb_url")) > 0 {
				db.Url = c.String("mongodb_url")
			}

		}),
	)

	service.Init()
	db.Init()
	log.Info("test")
	proto.RegisterAuthHandler(service.Server(), new(handler.Auth))

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
