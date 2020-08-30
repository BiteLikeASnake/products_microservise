package main

import (
	"fmt"
	"time"

	"github.com/call-me-snake/products_microservise/internal/db"
	"github.com/call-me-snake/products_microservise/internal/server"
	"github.com/jessevdk/go-flags"

	"github.com/labstack/gommon/log"
)

type Envs struct {
	ConnectionString string `long:"connection_string" env:"DB" description:"Connection string to database" default:"user=postgres password=example dbname=online_shop sslmode=disable"`
	ServerAddress    string `long:"server_address" env:"ADDRESS" description:"Address of the server" default:":8081"`

	AdminToken string `long:"admintoken" env:"ADMIN_TOKEN" description:"admin token" default:"adminpass"`
	UserToken  string `long:"usertoken" env:"USER_TOKEN" description:"user token" default:"userpass"`
}

var envs Envs
var parser = flags.NewParser(&envs, flags.Default)

func main() {
	var parser = flags.NewParser(&envs, flags.Default)
	if _, err := parser.Parse(); err != nil {
		log.Fatal(err.Error())
	}
	var database *db.Db
	var err error
	for {
		database, err = db.New(envs.ConnectionString)
		if err == nil {
			break
		}
		fmt.Println(err)
		time.Sleep(5 * time.Second)
	}
	//database, err := db.New(envs.ConnectionString)
	defer database.Close()
	fmt.Printf("got params: %v\n", envs)
	server.DefineParams(envs.AdminToken, envs.UserToken)
	conn := server.New(envs.ServerAddress)
	conn.ExecuteHandlers(database)
	fmt.Println("started")
	conn.Start()
	fmt.Println("end")
}

//cmd.exe --connection_string "user=postgres password=example dbname=online_shop sslmode=disable" --server_address ":8081" --admintoken "adminpass" --usertoken "userpass"
