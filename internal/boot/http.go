package boot

import (
	"log"
	"net/http"

	_ "github.com/lib/pq" // Postgres driver

	"github.com/essajiwa/teratur/internal/config"
	"github.com/essajiwa/teratur/pkg/safesql"

	userData "github.com/essajiwa/teratur/internal/data/user"
	httpsrv "github.com/essajiwa/teratur/internal/delivery/http"
	userH "github.com/essajiwa/teratur/internal/delivery/http/user"
	userSvc "github.com/essajiwa/teratur/internal/service/user"
)

// HTTP will load configuration, do Dependency Injection and then start the HTTP server
func HTTP() error {

	var (
		s   httpsrv.Server  // HTTP server object
		ud  userData.Data   // User domain data layer
		us  userSvc.Service // User domain service layer
		uh  *userH.Handler  // User domain Handler
		cfg *config.Config  // Config object
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
		log.Fatalf("init follower DB failed: %v", err)
	}

	// User domain init
	ud = userData.New(slaveDB)
	us = userSvc.New(ud)
	uh = userH.New(us)

	// Inject service used on handler here
	s = httpsrv.Server{
		User: uh,
	}

	if err := s.Serve(cfg.Server.Port); err != http.ErrServerClosed {
		// Error starting or closing listener:
		return err
	}

	return nil
}
