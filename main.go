package main

import (
	"log"

	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"github.com/unicok/auth/db"
	proto "github.com/unicok/auth/proto/auth"

	"github.com/unicok/auth/handler"
)

func main() {
	service := micro.NewService(
		micro.Name("com.unicok.srv.auth"),
		micro.Version("latest"),

		micro.Flags(
			cli.StringFlag{
				Name:   "mongodb_url",
				EnvVar: "MONGODB_URL",
				Usage:  "The mongodb URL e.g mongodb://username:password@127.0.0.1:27017/db",
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

	proto.RegisterAuthServiceHandler(service.Server(), new(handler.Auth))

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
