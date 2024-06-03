package main

import (
	"context"
	"github.com/caarlos0/env/v10"
	"github.com/romsar/hlsoc"
	"github.com/romsar/hlsoc/bcrypt"
	"github.com/romsar/hlsoc/grpc"
	"github.com/romsar/hlsoc/jwt"
	"github.com/romsar/hlsoc/postgres"
	"github.com/romsar/hlsoc/rabbitmq"
	"github.com/romsar/hlsoc/redis"
	"golang.org/x/sync/errgroup"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	err := run()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}

func run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	errWg, ctx := errgroup.WithContext(ctx)

	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	var tokenizer hlsoc.Tokenizer
	var userRepository hlsoc.UserRepository
	var postRepository hlsoc.PostRepository
	var postBroker hlsoc.PostBroker
	var passwordHasher hlsoc.PasswordHasher

	// jwt
	tokenizer = jwt.New(cfg.JWT.Secret)

	// postgres
	{
		slog.Info("opening postgres connection")

		var db *postgres.DB
		db, err = postgres.Open(cfg.Postgres.DSN)
		if err != nil {
			return err
		}

		errWg.Go(func() error {
			<-ctx.Done()

			slog.Info("terminating postgres connection")
			return db.Close()
		})

		userRepository = db
		postRepository = db
	}

	// rabbitmq
	{
		slog.Info("opening rabbitmq connection")

		var rab *rabbitmq.Client
		rab, err = rabbitmq.New(
			cfg.AMQP.Addr,
		)
		if err != nil {
			return err
		}

		errWg.Go(func() error {
			<-ctx.Done()

			slog.Info("terminating rabbitmq connection")
			return rab.Close()
		})

		postBroker = rab
	}

	// redis
	{
		slog.Info("opening redis connection")

		var red *redis.Client
		red, err = redis.New(
			cfg.Redis.Addr,
			cfg.Redis.Password,
			redis.WithUserRepository(userRepository),
			redis.WithPostRepository(postRepository),
			redis.WithPostBroker(postBroker),
		)
		if err != nil {
			return err
		}

		errWg.Go(func() error {
			<-ctx.Done()

			slog.Info("terminating redis connection")
			return red.Close()
		})

		postRepository = red
	}

	// bcrypt
	{
		slog.Info("creating bcrypt password hasher")
		passwordHasher = bcrypt.New(14)
	}

	// grpc
	{
		opts := []grpc.Option{
			grpc.WithTokenizer(tokenizer),
			grpc.WithUserRepository(userRepository),
			grpc.WithPasswordHasher(passwordHasher),
			grpc.WithPostRepository(postRepository),
			grpc.WithPostBroker(postBroker),
		}

		s := grpc.New(cfg.GRPC.Addr, opts...)
		errWg.Go(func() error {
			slog.Info("starting grpc server")
			return s.Start()
		})
		errWg.Go(func() error {
			<-ctx.Done()

			slog.Info("stopping grpc server")
			s.Stop()

			return nil
		})
	}

	err = errWg.Wait()
	if err != nil {
		return err
	}

	slog.Info("app was gracefully stopped")

	return nil
}

func loadConfig() (*Config, error) {
	var c Config

	err := env.Parse(&c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

type Config struct {
	JWT      JWTConfig      `envPrefix:"JWT_"`
	GRPC     GRPCConfig     `envPrefix:"GRPC_"`
	Postgres PostgresConfig `envPrefix:"PG_"`
	Redis    RedisConfig    `envPrefix:"REDIS_"`
	AMQP     AMQPConfig     `envPrefix:"AMQP_"`
}

type JWTConfig struct {
	Secret string `env:"SECRET,required"`
}

type GRPCConfig struct {
	Addr string `env:"ADDR,required" envDefault:":9090"`
}

type PostgresConfig struct {
	DSN string `env:"DSN,required"`
}

type RedisConfig struct {
	Addr     string `env:"ADDR,required"`
	Password string `env:"PASSWORD,required"`
}

type AMQPConfig struct {
	Addr string `env:"ADDR,required"`
}
