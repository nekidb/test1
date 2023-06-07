package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/nekidb/test1/internal/config"
	"github.com/nekidb/test1/internal/server"
	"github.com/nekidb/test1/internal/shortener"
	"github.com/nekidb/test1/internal/storage"
	"golang.org/x/sync/errgroup"
)

func main() {
	config, err := config.Get(os.DirFS("."), "config.json")
	if err != nil {
		log.Fatal(err)
	}

	storage, err := storage.NewBoltStorage(config.DB)
	if err != nil {
		log.Fatal(err)
	}
	defer storage.Close()

	shortener, err := shortener.NewShortenerService(storage)
	if err != nil {
		log.Fatal(err)
	}

	server := server.NewServer(shortener)

	ln, err := net.Listen("tcp", config.Host+config.Port)
	if err != nil {
		log.Fatal(err)
	}

	appContext, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	g, gCtx := errgroup.WithContext(appContext)

	g.Go(func() error {
		return server.Serve(ln)
	})

	g.Go(func() error {
		<-gCtx.Done()

		return server.Shutdown()
	})

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}
