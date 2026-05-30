package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/palantir/go-githubapp/githubapp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	cfgPath := "config.toml"
	if len(os.Args) > 1 {
		cfgPath = os.Args[1]
	}

	cfg, err := loadConfig(cfgPath)
	if err != nil {
		log.Fatal().Err(err).Str("path", cfgPath).Msg("failed to load config")
	}

	privateKey, err := os.ReadFile(cfg.GitHub.PrivateKey)
	if err != nil {
		log.Fatal().Err(err).Str("path", cfg.GitHub.PrivateKey).Msg("failed to read private key")
	}

	ghCfg := githubapp.Config{
		V3APIURL: "https://api.github.com/",
	}
	ghCfg.App.IntegrationID = cfg.GitHub.AppID
	ghCfg.App.WebhookSecret = cfg.GitHub.WebhookSecret
	ghCfg.App.PrivateKey = string(privateKey)

	cc, err := githubapp.NewDefaultCachingClientCreator(
		ghCfg,
		githubapp.WithClientUserAgent("noai-bot/1.0"),
	)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create github client creator")
	}

	handler := &PRHandler{
		ClientCreator: cc,
		cfg:           cfg,
	}

	dispatcher := githubapp.NewEventDispatcher(
		[]githubapp.EventHandler{handler},
		cfg.GitHub.WebhookSecret,
		githubapp.WithErrorCallback(func(w http.ResponseWriter, r *http.Request, err error) {
			log.Error().Err(err).Msg("event dispatch error")
			http.Error(w, "internal error", http.StatusInternalServerError)
		}),
	)

	http.Handle("/api/github/hook", dispatcher)

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Info().Str("addr", addr).Msg("starting noai-bot")
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal().Err(err).Msg("server stopped")
	}
}
