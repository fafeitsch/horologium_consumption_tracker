package main

import (
	"flag"
	"github.com/99designs/gqlgen/handler"
	"github.com/fafeitsch/Horologium/pkg/gql"
	"github.com/fafeitsch/Horologium/pkg/http"
	orm "github.com/fafeitsch/Horologium/pkg/persistance"
	"log"
)

func main() {
	arguments := make(map[string]*string)
	for _, setting := range getFlags() {
		var value string
		flag.StringVar(&value, setting.name, setting.defaultValue, setting.usage)
		flag.StringVar(&value, setting.shortOption, setting.defaultValue, setting.usage+" (short for "+setting.name+")")
		arguments[setting.name] = &value
	}
	flag.Parse()
	config, err := http.ConvertToConfiguration(arguments)
	if err != nil {
		log.Fatalf("could not parse arguments: %v", err)
	} else {
		db, err := orm.ConnectToFileDb(*arguments["home"] + "/horologium.sqlite")
		if err != nil {
			log.Fatalf("could not connect to database: %v", err)
		}
		seriesService := orm.NewSeriesService(db)
		resolver := gql.NewResolver(seriesService)
		apiHandler := handler.GraphQL(gql.NewExecutableSchema(gql.Config{Resolvers: resolver})).ServeHTTP
		server, e := http.NewServer(config, apiHandler)
		if e != nil {
			log.Fatalf("could not create server: %v", err)
		}
		err = server.StartServer()
		if err != nil {
			log.Fatalf("could not start server: %v", err)
		}
	}
}

type flagSetting struct {
	name         string
	shortOption  string
	defaultValue string
	usage        string
}

func getFlags() []flagSetting {
	return []flagSetting{
		{name: "bindAddress", shortOption: "a", defaultValue: "127.0.0.1", usage: "The address the server should bind to. Set to 0.0.0.0 to bind to all interfaces."},
		{name: "port", shortOption: "p", defaultValue: "9551", usage: "The port the server is running on."},
		{name: "jwtKey", shortOption: "k", defaultValue: "", usage: "The key used to sign the JWT-Keys. If it is empty then a random key is picked at runtime."},
		{name: "home", shortOption: "h", defaultValue: "./", usage: "The working directory of the application. Should be writable."},
	}
}
