package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/Francesco99975/rosskery/cmd/boot"
	"github.com/Francesco99975/rosskery/internal/models"
	"github.com/Francesco99975/rosskery/internal/storage"
	"github.com/Francesco99975/rosskery/internal/tools"

	"github.com/stripe/stripe-go/v78"
)

func main() {
	err := boot.LoadEnvVariables()
	if err != nil {
		panic(err)
	}

	// Create a root ctx and a CancelFunc which can be used to cancel retentionMap goroutine
	rootCtx := context.Background()
	ctx, cancel := context.WithCancel(rootCtx)
	defer cancel()

	port := os.Getenv("PORT")

	models.Setup(os.Getenv("DSN"))

	go tools.GotifyQueue.ProcessQueue()

	storage.ValkeySetup(ctx)

	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

	e := createRouter(ctx)

	go func() {
		e.Logger.Fatal(e.Start(":" + port))
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
