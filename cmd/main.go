package main

import (
	"context"
	"time"

	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"newsplatform/internal/agent"
	"newsplatform/internal/database"
	"newsplatform/internal/server"
)

type PlatformEnvs struct {
	Endpoint     string        `required:"true" split_words:"true" desc:"Newsapi API endpoint"`
	Token        string        `required:"true" split_words:"true" desc:"Newsapi client token"`
	Q            string        `split_words:"true" desc:"Key word query" default:"bitcoin"`
	PageSize     int           `split_words:"true" desc:"Key word query" default:"2"`
	MaxPage      int           `split_words:"true" desc:"Key word query" default:"5"`
	TimeInterval time.Duration `required:"true" split_words:"true" desc:"Key word query" default:"10s"`
}

func main() {
	ctx := context.Background()

	// Read mandatory variables from env
	var envs PlatformEnvs
	err := envconfig.Process("", &envs)
	if err != nil {
		log.Fatal("can not read required env variables", err)
	}

	// Setup and connect to database
	err = database.ConnectDB()
	if err != nil {
		log.Fatalf("Database can not be connected: %v", err)
	}

	// Run new collector agent as a goroutine
	a := agent.NewAgent(agent.AgentOpts{
		Endpoint:     envs.Endpoint,
		Token:        envs.Token,
		Q:            envs.Q,
		PageSie:      envs.PageSize,
		MaxPage:      envs.MaxPage,
		TimeInterval: envs.TimeInterval,
	})
	go a.Run(ctx)

	// Start RESTAPI server
	server.SetupServer().Run(":3000")
}
