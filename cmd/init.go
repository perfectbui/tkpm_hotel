package main

import (
	"TKPM/common/flags"
	"TKPM/configs"
	"TKPM/internals/delivery"
	"TKPM/internals/domain"
	"TKPM/internals/models"
	"TKPM/internals/repository"
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go/config"
	"github.com/urfave/cli"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

var (
	app          = NewApp()
	server       = new(srv)
	startCommand = cli.Command{
		Action:      flags.MigrateFlags(Start),
		Name:        "start",
		Usage:       "start server",
		ArgsUsage:   "<genesisPath>",
		Flags:       []cli.Flag{},
		Description: `start server`,
	}
)

type srv struct {
	cfg    *configs.Config
	logger logr.Logger

	mgoDB            *mongo.Database
	mgoClient        *mongo.Client
	mgoClientOptions *options.ClientOptions

	roomDomain     domain.Room
	accountDomain  domain.Account
	contractDomain domain.Contract

	roomDelivery     delivery.RoomDelivery
	accountDelivery  delivery.AccountDelivery
	contractDelivery delivery.ContractDelivery

	tracer opentracing.Tracer
}

func init() {
	app.Action = cli.ShowAppHelp
	app.Commands = []cli.Command{
		startCommand,
	}
	app.Flags = []cli.Flag{
		flags.ServerHostFlag,
		flags.ServerPortFlag,
		flags.ServerNameFlag,

		flags.MongoDatabaseNameFlag,
		flags.MongoHostFlag,
		flags.MongoPortFlag,

		flags.StorageAccessKeyFlag,
		flags.StorageSecretKeyFlag,
		flags.StorageRegionFlag,
		flags.StorageNameFlag,

		flags.JaegerHostFlag,
		flags.JaegerPortFlag,
	}
}

func (s *srv) loadConfig(ctx *cli.Context) error {
	server.cfg = &configs.Config{
		HTTP: configs.ConnAddress{
			Host: ctx.GlobalString(flags.ServerHostFlag.GetName()),
			Port: ctx.GlobalString(flags.ServerPortFlag.GetName()),
		},
		Mongo: configs.Mongo{
			Host:     ctx.GlobalString(flags.MongoHostFlag.GetName()),
			Port:     ctx.GlobalString(flags.MongoPortFlag.GetName()),
			Database: ctx.GlobalString(flags.MongoDatabaseNameFlag.GetName()),
		},
		Storage: configs.Storage{
			AccessKey:  ctx.GlobalString(flags.StorageAccessKeyFlag.GetName()),
			SecretKey:  ctx.GlobalString(flags.StorageSecretKeyFlag.GetName()),
			BucketName: ctx.GlobalString(flags.StorageNameFlag.GetName()),
			Region:     ctx.GlobalString(flags.StorageRegionFlag.GetName()),
		},
		Tracer: configs.ConnAddress{
			Host: ctx.GlobalString(flags.JaegerHostFlag.GetName()),
			Port: ctx.GlobalString(flags.JaegerPortFlag.GetName()),
		},
		ServiceName: ctx.GlobalString(flags.ServerNameFlag.GetName()),
	}
	return nil
}

func (s *srv) loadLogger() error {
	zapLog, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	s.logger = zapr.NewLogger(zapLog)
	return nil
}

func (s *srv) connectMongo() error {
	s.mgoClientOptions = options.Client().ApplyURI("mongodb://my_database:27017")

	// connect to mongoDb
	var err error
	s.mgoClient, err = mongo.Connect(context.TODO(), s.mgoClientOptions)
	if err != nil {
		s.logger.Error(err, "fail to connect mongodb")
		return err
	}
	s.mgoDB = s.mgoClient.Database(s.cfg.Mongo.Database)
	s.logger.Info("connect mongodb successfull ")
	return nil
}

func (s *srv) loadTracing() error {
	cfg := config.Configuration{
		ServiceName: s.cfg.ServiceName,
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:            false,
			BufferFlushInterval: 1 * time.Second,
			LocalAgentHostPort:  s.cfg.Tracer.Host + ":" + s.cfg.Tracer.Port,
		},
	}
	tracer, _, err := cfg.NewTracer()
	if err != nil {
		s.logger.Error(err, "fail to load tracing")
		return err
	}

	opentracing.SetGlobalTracer(tracer)

	s.tracer = tracer
	s.logger.Info("load tracing successfull ")
	return nil
}

func (s *srv) loadDomain() error {
	s.accountDomain = domain.NewAccountDomain(repository.NewAccountRepository(s.mgoDB.Collection(models.AccountCollection)))
	s.contractDomain = domain.NewContractDomain(repository.NewContractRepository(s.mgoDB.Collection(models.ContractCollection)))
	s.roomDomain = domain.NewRoomDomain(repository.NewRoomRepository(s.mgoDB.Collection(models.RoomCollection)))
	s.logger.Info("load domain successfull ")
	return nil
}

func (s *srv) loadDelivery() error {
	s.accountDelivery = delivery.NewAccountDelivery(s.accountDomain)
	s.contractDelivery = delivery.NewContractDelivery(s.contractDomain, s.roomDomain)
	s.roomDelivery = delivery.NewRoomDelivery(s.roomDomain, s.cfg.Storage)
	s.logger.Info("load delivery successfull ")
	return nil
}

func (s *srv) startHTTPServer() {
	handler := delivery.NewHTTPHandler(s.roomDelivery, s.accountDelivery, s.contractDelivery, s.tracer)
	server := &http.Server{
		Addr:    s.cfg.HTTP.Host + ":" + s.cfg.HTTP.Port,
		Handler: delivery.AllowCORS(handler),
	}
	s.logger.Info(fmt.Sprintf("start http server at port %v\n", s.cfg.HTTP.Port))
	log.Fatal(server.ListenAndServe())
}

// NewApp creates an app with sane defaults.
func NewApp() *cli.App {
	app := cli.NewApp()
	app.Action = cli.ShowAppHelp
	app.Name = "Hotel management"
	app.Author = "Bui Hoan Hao"
	app.Email = "haopro@gmail.com"
	app.Usage = "Server API"
	return app
}

// Start ...
func Start(ctx *cli.Context) error {
	if err := server.loadConfig(ctx); err != nil {
		return err
	}

	if err := server.loadLogger(); err != nil {
		return err
	}

	if err := server.connectMongo(); err != nil {
		return err
	}

	if err := server.loadTracing(); err != nil {
		return err
	}

	if err := server.loadDomain(); err != nil {
		return err
	}

	if err := server.loadDelivery(); err != nil {
		return err
	}

	server.startHTTPServer()
	return nil
}
