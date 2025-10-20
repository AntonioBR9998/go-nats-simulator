package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/AntonioBR9998/go-nats-simulator/utils"
	"github.com/nats-io/nats.go"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	server "github.com/AntonioBR9998/go-nats-simulator/gan/api"
	"github.com/AntonioBR9998/go-nats-simulator/gan/config"
	"github.com/AntonioBR9998/go-nats-simulator/gan/domain"
	"github.com/AntonioBR9998/go-nats-simulator/gan/repository"
	"github.com/AntonioBR9998/go-nats-simulator/gan/simulator"
)

const (
	explainedName = "{G}o {A}PI {N}ATS"
)

var (
	Version   = "dev"
	Commit    = "I'm live!"
	BuildDate = "I don't remember exactly"
)

func main() {
	app := cli.App{
		Name:        "gan",
		Usage:       explainedName,
		Description: "GAN is a microservice for manage configs and samples from IoT devices",
		Action:      startGanService,
		Version:     Version,
		Before:      BeforeFunc,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:      "config",
				Required:  true,
				TakesFile: true,
				Aliases:   []string{"c"},
				Usage:     "load configuration from `FILE`",
				EnvVars:   []string{"GAN_CONFIG_FILE"},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Printf("\n%s", err)
	}
}

func startGanService(ctx *cli.Context) error {
	log.Infoln("starting " + explainedName)

	configFilePath := ctx.String("config")
	log.Debugf("the configuration file path is: %s", configFilePath)
	log.Infof("loading configuration from file '%s'", configFilePath)
	cfg := utils.New(
		&config.Config{},
		func(cfg *config.Config) {
			cfg.ConfigPath = configFilePath
		},
	)

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{ServerName: cfg.ServerName, InsecureSkipVerify: true}

	log.Traceln("creating repository layer")
	repository := repository.NewRepository(*cfg)

	log.Traceln("creating service layer")
	natsClient, _ := nats.Connect(cfg.Nats.Host + ":" + cfg.Nats.Port)
	sensorManager := simulator.NewManager(natsClient)
	service := domain.NewService(repository, *cfg, sensorManager)

	log.Traceln("creating REST API layer")
	s := server.NewAPI(*cfg, service)

	log.Infoln("the user server is on tap now: ", cfg.API.GetURL())
	return http.ListenAndServe(cfg.API.GetRelativeURL(), s.Router())
}

func BeforeFunc(ctx *cli.Context) error {
	ch := make(chan os.Signal, 2)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-ch
		log.Infoln("shutting down GAN service!")
		os.Exit(0)
	}()

	return nil
}
