package main

import (
	"context"
	"database/sql"
	"expvar"
	"flag"
	"log"
	"log/slog"
	"os"
	"runtime"
	"sync"
	"time"

	db "github.com/Hopertz/rmgmt/db/sqlc"
	"github.com/Hopertz/rmgmt/pkg/beem"
	_ "github.com/lib/pq"
	"gopkg.in/go-playground/validator.v9"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}

	phone string

	beemSmS struct {
		apiKey    string
		secretKey string
	}
}

type envelope map[string]interface{}

type application struct {
	config    config
	beem      *beem.Beem
	wg        sync.WaitGroup
	store     db.Store
	validator *validator.Validate
}

func init() {

	var programLevel = new(slog.LevelVar)

	h := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: programLevel})
	slog.SetDefault(slog.New(h))

}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 5040, "API server port")
	flag.StringVar(&cfg.env, "env", os.Getenv("ENV_STAGE"), "Environment (development|Staging|production")
	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("DB_DSN"), "PostgreSQL DSN")

	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max ilde connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection  connections")

	flag.StringVar(&cfg.beemSmS.apiKey, "beem-apikey", os.Getenv("BEEM_APIKEY"), "beem-apikey")
	flag.StringVar(&cfg.beemSmS.secretKey, "beem-secretkey", os.Getenv("BEEM_SECRETKEY"), "beem-secretkey")

	flag.StringVar(&cfg.phone, "admin-phone", os.Getenv("ADMIN_PHONE"), "admin phone number")

	flag.Parse()

	dbConn, err := openDB(cfg)
	if err != nil {
		log.Fatal("error opening db", err)
	}

	defer dbConn.Close()

	slog.Info("database connection established")

	expvar.NewString("version").Set(version)

	expvar.Publish("goroutines", expvar.Func(func() interface{} {
		return runtime.NumGoroutine()
	}))

	expvar.Publish("database", expvar.Func(func() interface{} {
		return dbConn.Stats()
	}))

	expvar.Publish("timestamp", expvar.Func(func() interface{} {
		return time.Now().Unix()
	}))

	app := &application{
		config:    cfg,
		beem:      beem.New(cfg.beemSmS.apiKey, cfg.beemSmS.secretKey),
		store:     db.NewStore(dbConn),
		validator: validator.New(),
	}

	err = app.serve()
	if err != nil {
		log.Fatal("error starting server", err)
	}

}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)

	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	db.SetMaxIdleConns(cfg.db.maxIdleConns)

	duration, err := time.ParseDuration(cfg.db.maxIdleTime)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)

	if err != nil {
		return nil, err
	}

	return db, nil
}
