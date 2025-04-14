package main

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/sync/errgroup"
	"log"
	"os"
	"os/signal"
	"pvz/internal/application"
	"pvz/internal/config"
	receptionRepos "pvz/internal/infrastructure/receptions"
	userRepos "pvz/internal/infrastructure/users"
	receptionServices "pvz/internal/service/receptions"
	userServices "pvz/internal/service/users"
	"syscall"
)

func main() {
	log.Println("initiating user repository...")
	url := config.GetUsersPostgresURL()
	conn, err := pgxpool.New(context.Background(), config.GetUsersPostgresURL())
	if err != nil {
		log.Fatal(err)
	}
	if err = conn.Ping(context.Background()); err != nil {
		log.Fatal(err)
	}
	log.Println(url, err)
	ur := userRepos.NewPostgresRepo(conn)
	rr := receptionRepos.NewPostgresRepo(conn)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		cancel()
	}()
	errG, gCtx := errgroup.WithContext(ctx)

	log.Println("initiating user service...")
	token := userServices.NewToken([]byte(config.GetSecretKey()))
	userService := userServices.NewUserService(ur, token)

	log.Println("initiating reception service...")
	receptionService := receptionServices.NewReceptionService(rr)

	log.Println("initializing server...")
	httpServer := application.SetupHTTPServer(userService, receptionService)

	errG.Go(func() error {
		return httpServer.ListenAndServe()
	})
	errG.Go(func() error {
		<-gCtx.Done()
		log.Println("shutting down server...")
		return httpServer.Shutdown(gCtx)
	})
	errG.Go(func() error {
		<-gCtx.Done()
		log.Println("closing database...")
		conn.Close()
		return nil
	})
	if err = errG.Wait(); err != nil {
		log.Printf("exit reason: %s \n", err)
	}
	log.Println("app shut down")
}
