package storage

import (
	"context"
	"os"

	valkey "github.com/Desquaredp/go-valkey"
)

var Valkey *valkey.Client

func ValkeySetup(ctx context.Context) {
	Valkey = valkey.NewClient(&valkey.Options{
		Addr:     os.Getenv("VALKEY_ADDR"),
		Password: "",
		DB:       0,
	})

	_, err := Valkey.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}

	if err := Valkey.Set(ctx, "test", "test", 0).Err(); err != nil {
		panic(err)
	}

	val, err := Valkey.Get(ctx, "test").Result()
	if err != nil || val != "test" {
		panic(err)
	}

	if err := Valkey.Del(ctx, "test").Err(); err != nil {
		panic(err)
	}

	if err := Valkey.Set(ctx, string(Online), true, 0).Err(); err != nil {
		panic(err)
	}

	if err := Valkey.Set(ctx, string(Operative), true, 0).Err(); err != nil {
		panic(err)
	}
}
