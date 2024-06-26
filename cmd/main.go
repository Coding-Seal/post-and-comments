package main

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/99designs/gqlgen/graphql/playground"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/99designs/gqlgen/graphql/handler"
	"post-and-comments/internal/adapters/graph"
	"post-and-comments/internal/adapters/graph/generated"
	inmemory "post-and-comments/internal/adapters/repos/inmemory"
	"post-and-comments/internal/adapters/repos/postgres"
	inports "post-and-comments/internal/inner-ports"
	"post-and-comments/internal/services"
)

const dbUrl = "postgres://postgres:mysecretpassword@localhost:5432/postgres?sslmode=disable"

func main() {
	var useInMemoryStorage bool

	flag.BoolVar(&useInMemoryStorage, "memory", false, "Set true if you want to use in memory storage")

	var postStore inports.PostStore

	var userStore inports.UserStore

	var commentStore inports.CommentStore

	if useInMemoryStorage {
		postStore = inmemory.NewPostStore()
		userStore = inmemory.NewUserStore()
		commentStore = inmemory.NewCommentStore()
	} else {
		db, err := sql.Open("pgx", dbUrl)
		if err != nil {
			log.Fatal(err)
		}
		err = db.Ping()
		if err != nil {
			log.Fatal(err)
		}

		defer db.Close()
		postStore = postgres.NewPostStore(db)
		userStore = postgres.NewUserStore(db)
		commentStore = postgres.NewCommentStore(db)
	}

	resolver := &graph.Resolver{
		PostSrv:    services.Post{PostStore: postStore},
		CommentSrv: services.Comment{CommentStore: commentStore},
		UserSrv:    services.User{UserStore: userStore},
	}
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{
		Resolvers:  resolver,
		Directives: generated.DirectiveRoot{},
		Complexity: generated.ComplexityRoot{},
	}))
	srv.SetRecoverFunc(func(ctx context.Context, err any) (userMessage error) {
		log.Print(err)
		debug.PrintStack()

		return errors.New("user message on panic")
	})

	http.Handle("/", playground.Handler("Comments", "/query"))
	http.Handle("/query", srv)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
