package boot

import (
	"log"
	"net/http"

	_ "github.com/lib/pq" // Postgres driver

	"github.com/essajiwa/teratur/internal/config"
	"github.com/essajiwa/teratur/pkg/errors"
	"github.com/essajiwa/teratur/pkg/httpclient"
	"github.com/essajiwa/teratur/pkg/safesql"

	shopData "github.com/essajiwa/teratur/internal/data/shop"
	userData "github.com/essajiwa/teratur/internal/data/user"
	httpsrv "github.com/essajiwa/teratur/internal/delivery/http"
	userH "github.com/essajiwa/teratur/internal/delivery/http/user"
	shopSvc "github.com/essajiwa/teratur/internal/service/shop"
	userSvc "github.com/essajiwa/teratur/internal/service/user"
)

// HTTP will load configuration, do Dependency Injection and then start the HTTP server
func HTTP() error {

	var (
		s     httpsrv.Server  // HTTP server object
		ud    userData.Data   // User domain data layer
		us    userSvc.Service // User domain service layer
		sd    shopData.Data
		ss    shopSvc.Service // Shop domain service layer
		uh    *userH.Handler  // User domain Handler
		cfg   *config.Config  // Config object
		httpc *httpclient.Client
	)

	// Config initialization
	err := config.Init()
	if err != nil {
		log.Fatalf("failed to initialize config: %v", err)
	}

	cfg = config.Get()
	// Open DB connection
	slaveDB, err := safesql.OpenSlaveDB("postgres", cfg.Database.Follower)
	if err != nil {
		log.Fatalf("init follower DB failed %v", errors.Wrap(err))
	}

	httpc = httpclient.NewClient()

	// User domain init
	ud = userData.New(slaveDB)
	us = userSvc.New(ud)
	sd = shopData.New(httpc, cfg.API.Shop)
	ss = shopSvc.New(sd)
	uh = userH.New(us, ss)

	// Inject service used on handler here
	s = httpsrv.Server{
		User: uh,
	}

	if err := s.Serve(cfg.Server.Port); err != http.ErrServerClosed {
		// Error starting or closing listener:
		return errors.Wrap(err)
	}

	return nil
}
