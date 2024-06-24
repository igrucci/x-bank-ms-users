package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"x-bank-users/config"
	"x-bank-users/core/web"
	"x-bank-users/infra/gmail"
	"x-bank-users/infra/hasher"
	"x-bank-users/infra/random"
	"x-bank-users/infra/redis"
	"x-bank-users/infra/swissknife"
	"x-bank-users/transport/http"
	"x-bank-users/transport/http/jwt"
)

var (
	addr       = flag.String("addr", ":8080", "")
	configFile = flag.String("config", "config.json", "")
)

func main() {
	flag.Parse()

	conf, err := config.Read(*configFile)
	if err != nil {
		log.Fatal(err)
	}

	knife := swissknife.NewService()
	passwordHasher := hasher.NewService()

	//jwtHs512, err := jwt.NewHS512(conf.Hs512SecretKey)
	//if err != nil {
	//	log.Fatal(err)
	//}
	jwtRs256, err := jwt.NewRS256(conf.Rs256PrivateKey, conf.Rs256PublicKey)
	if err != nil {
		log.Fatal(err)
	}

	redisService, err := redis.NewService(conf.Redis.Password, conf.Redis.Host, conf.Redis.Port, conf.Redis.Database, conf.Redis.MaxCons)
	if err != nil {
		log.Fatal(err)
	}
	gmailService := gmail.NewService(conf.Gmail.Host, conf.Gmail.SenderName, conf.Gmail.SenderEmail, conf.Gmail.Login, conf.Gmail.Password, conf.Gmail.UrlToActivate, conf.Gmail.UrlToRestore)

	randomGenerator := random.NewService()

	service := web.NewService(&knife, &randomGenerator, &redisService, &gmailService, &passwordHasher, &redisService, &redisService, &knife, &redisService)

	transport := http.NewTransport(service, &jwtRs256)

	errCh := transport.Start(*addr)
	interruptsCh := make(chan os.Signal, 1)
	signal.Notify(interruptsCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-errCh:
		log.Fatal(err)
	case <-interruptsCh:
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer shutdownCancel()
		err = transport.Stop(shutdownCtx)
		if err != nil {
			log.Fatal(err)
		}
	}
}
