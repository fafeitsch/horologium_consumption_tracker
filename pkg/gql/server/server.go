package standaloneServer

import (
	orm "github.com/fafeitsch/Horologium/pkg/persistance"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/handler"
	"github.com/fafeitsch/Horologium/pkg/gql"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	http.Handle("/", handler.Playground("GraphQL playground", "/query"))
	db, _ := orm.CreateInMemoryDb()
	seriesService := orm.NewSeriesService(db)
	planService := orm.NewPricingPlanService(db)
	resolver := gql.NewResolver(seriesService, planService)
	http.Handle("/query", handler.GraphQL(gql.NewExecutableSchema(gql.Config{Resolvers: resolver})))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
