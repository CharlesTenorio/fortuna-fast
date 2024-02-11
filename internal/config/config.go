package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/jwtauth"
)

const (
	DEVELOPER    = "developer"
	HOMOLOGATION = "homologation"
	PRODUCTION   = "production"
)

type Config struct {
	PORT           string `json:"port"`
	Mode           string `json:"mode"`
	MongoDBConfig  `json:"mongo_config"`
	RedisConfig    RedisDBConfig `json:"redis_config"`
	RMQConfig      RMQConfig     `json:"rmq_config"`
	ConsumerConfig `json:"cconfig"`
	JWTSecretKey   string `json:"jwt_secret_key"`
	JWTTokenExp    int    `json:"jwt_token_exp"`
	TokenAuth      *jwtauth.JWTAuth
	DataInicial    time.Time
	DataFinal      time.Time
}

type MongoDBConfig struct {
	MDB_URI              string            `json:"mdb_uri"`
	MDB_NAME             string            `json:"mdb_name"`
	MDB_CLIENT           string            `json:"mdb_client"`
	MDB_DELIVERY_ADDRESS string            `json:"mdb_delivery_address"`
	MDB_GRIFE            string            `json:"mdb_grife"`
	MDB_ORDER            string            `json:"mdb_order"`
	MDB_DET_ORDER        string            `json:"mdb_det_order"`
	MDB_PAYMENT          string            `json:"mdb_payment"`
	MDB_PRODUTC          string            `json:"mdb_product"`
	MDB_SUPPLIER         string            `json:"mdb_supplier"`
	MDB_USER             string            `json:"mdb_user"`
	MDB_COLLECTIONS      map[string]string `json:"mdb_collections"`
}

type RedisDBConfig struct {
	RDB_HOST string `json:"rdb_host"`
	RDB_PORT string `json:"rdb_port"`
	RDB_USER string `json:"rdb_user"`
	RDB_PASS string `json:"rdb_pass"`
	RDB_DB   int64  `json:"rdb_db"`
	RDB_DSN  string `json:"-"`
}

type RMQConfig struct {
	RMQ_URI                  string `json:"rmq_uri"`
	RMQ_MAXX_RECONNECT_TIMES int    `json:"rmq_maxx_reconnect_times"`
}

type ConsumerConfig struct {
	ExchangeName  string `json:"exchange_name"`
	ExchangeType  string `json:"exchange_type"`
	RoutingKey    string `json:"routing_key"`
	QueueName     string `json:"queue_name"`
	ConsumerName  string `json:"consumer_name"`
	ConsumerCount int    `json:"consumer_count"`
	PrefetchCount int    `json:"prefetch_count"`
	Reconnect     struct {
		MaxAttempt int `json:"max_attempt"`
		Interval   int `json:"interval"`
	}
}

func NewConfig() *Config {
	conf := defaultConf()

	SRV_PORT := os.Getenv("SRV_PORT")
	if SRV_PORT != "" {
		conf.PORT = SRV_PORT
	}

	SRV_MODE := os.Getenv("SRV_MODE")
	if SRV_MODE != "" {
		conf.Mode = SRV_MODE
	}

	SRV_RDB_HOST := os.Getenv("SRV_RDB_HOST")
	if SRV_RDB_HOST != "" {
		conf.RedisConfig.RDB_HOST = SRV_RDB_HOST
	}

	SRV_RDB_PORT := os.Getenv("SRV_RDB_PORT")
	if SRV_RDB_PORT != "" {
		conf.RedisConfig.RDB_PORT = SRV_RDB_PORT
	}

	SRV_RDB_USER := os.Getenv("SRV_RDB_USER")
	if SRV_RDB_USER != "" {
		conf.RedisConfig.RDB_USER = SRV_RDB_USER
	}

	SRV_RDB_PASS := os.Getenv("SRV_R_PASS")
	if SRV_RDB_PASS != "" {
		conf.RedisConfig.RDB_PASS = SRV_RDB_PASS
	}

	SRV_RDB_DB := os.Getenv("SRV_R_DB")
	if SRV_RDB_DB != "" {
		conf.RedisConfig.RDB_DB, _ = strconv.ParseInt(SRV_RDB_DB, 10, 64)
	}

	SRV_MDB_URI := os.Getenv("SRV_MDB_URI")
	if SRV_MDB_URI != "" {
		conf.MDB_URI = SRV_MDB_URI
	}

	SRV_MDB_NAME := os.Getenv("SRV_MDB_NAME")
	if SRV_MDB_NAME != "" {
		conf.MDB_NAME = SRV_MDB_NAME
	}

	SRV_MDB_COLLECTIONS := os.Getenv("SRV_MDB_COLLECTIONS")
	if SRV_MDB_COLLECTIONS != "" {
		collectionsMap := parseCollectionsString(SRV_MDB_COLLECTIONS)
		conf.MDB_COLLECTIONS = collectionsMap
	}

	SRV_RDB_DSN := os.Getenv("SRV_RDB_DSN")
	if SRV_RDB_DSN != "" {
		conf.RedisConfig.RDB_DSN = SRV_RDB_DSN
	}

	if len(conf.RedisConfig.RDB_HOST) > 3 {

		// "redis://<user>:<pass>@localhost:6379/<db>"
		// https://redis.uptrace.dev/guide/go-redis.html#connecting-to-redis-server

		conf.RedisConfig.RDB_DSN = fmt.Sprintf("redis://%s:%s@%s:%s/%v", conf.RedisConfig.RDB_USER, conf.RedisConfig.RDB_PASS, conf.RedisConfig.RDB_HOST, conf.RedisConfig.RDB_PORT, conf.RedisConfig.RDB_DB)
	}

	SRV_RMQ_URI := os.Getenv("SRV_RMQ_URI")
	if SRV_RMQ_URI != "" {
		conf.RMQConfig.RMQ_URI = SRV_RMQ_URI
	}

	CC_EX_NAME := os.Getenv("CC_EX_NAME")
	if CC_EX_NAME != "" {
		conf.ConsumerConfig.ExchangeName = CC_EX_NAME
	}

	CC_EX_TYPE := os.Getenv("CC_EX_TYPE")
	if CC_EX_TYPE != "" {
		conf.ConsumerConfig.ExchangeType = CC_EX_TYPE
	}

	CC_RT_KEY := os.Getenv("CC_RT_KEY")
	if CC_RT_KEY != "" {
		conf.ConsumerConfig.RoutingKey = CC_RT_KEY
	}

	CC_QU_NAME := os.Getenv("CC_QU_NAME")
	if CC_QU_NAME != "" {
		conf.ConsumerConfig.QueueName = CC_QU_NAME
	}

	CC_C_NAME := os.Getenv("CC_C_NAME")
	if CC_C_NAME != "" {
		conf.ConsumerConfig.ConsumerName = CC_C_NAME
	}

	CC_C_COUNT := os.Getenv("CC_C_COUNT")
	if CC_C_COUNT != "" {
		conf.ConsumerConfig.ConsumerCount, _ = strconv.Atoi(CC_C_COUNT)
	}

	CC_PREF_COUNT := os.Getenv("CC_PREF_COUNT")
	if CC_PREF_COUNT != "" {
		conf.ConsumerConfig.PrefetchCount, _ = strconv.Atoi(CC_PREF_COUNT)
	}

	CC_MAX_ATTEMPT := os.Getenv("CC_MAX_ATTEMPT")
	if CC_MAX_ATTEMPT != "" {
		conf.ConsumerConfig.Reconnect.MaxAttempt, _ = strconv.Atoi(CC_MAX_ATTEMPT)
	}

	CC_INTERVAL := os.Getenv("CC_INTERVAL")
	if CC_INTERVAL != "" {
		conf.ConsumerConfig.Reconnect.Interval, _ = strconv.Atoi(CC_INTERVAL)
	}

	SRV_JWT_SECRET_KEY := os.Getenv("SRV_JWT_SECRET_KEY")
	if SRV_JWT_SECRET_KEY != "" {
		conf.JWTSecretKey = SRV_JWT_SECRET_KEY
	}

	SRV_JWT_TOKEN_EXP := os.Getenv("SRV_JWT_TOKEN_EXP")
	if SRV_JWT_SECRET_KEY != "" {
		conf.JWTTokenExp, _ = strconv.Atoi(SRV_JWT_TOKEN_EXP)
	}

	SRV_DATA_INICIAL := os.Getenv("SRV_DATA_INICIAL")
	if SRV_DATA_INICIAL != "" {
		dataInicial, err := time.Parse(time.RFC3339, SRV_DATA_INICIAL)
		if err != nil {
			fmt.Println("Erro ao converter a string para time.Time:", err)

		}
		conf.DataInicial = dataInicial
	}

	SRV_DATA_FINAL := os.Getenv("SRV_DATA_INICIAL")
	if SRV_DATA_FINAL != "" {
		dataFinal, err := time.Parse(time.RFC3339, SRV_DATA_FINAL)
		if err != nil {
			fmt.Println("Erro ao converter a string para time.Time:", err)

		}
		conf.DataFinal = dataFinal
	}

	return conf
}

func defaultConf() *Config {

	default_conf := Config{
		PORT:         "8080",
		Mode:         DEVELOPER,
		JWTSecretKey: "RgUkXp2s5v8y/B?E(H+KbPeShVmYq3t6", // "----your-256-bit-secret-here----" length 32
		JWTTokenExp:  300,
		// 15m
		MongoDBConfig: MongoDBConfig{
			MDB_URI:         "mongodb://admin:supersenha@localhost:27017/",
			MDB_NAME:        "teste_db",
			MDB_COLLECTIONS: make(map[string]string),
		},

		RedisConfig: RedisDBConfig{
			RDB_HOST: "localhost",
			RDB_PORT: "6379",
			RDB_DB:   0,
			RDB_DSN:  "redis://localhost:6379/0",
		},
		RMQConfig: RMQConfig{
			RMQ_URI: "amqp://admin:supersenha@localhost:5672/",
		},

		ConsumerConfig: ConsumerConfig{
			ExchangeName:  "message_cotafacil",
			ExchangeType:  "direct",
			RoutingKey:    "create",
			QueueName:     "SEND_MESSAGE_COTAFACIL",
			ConsumerName:  "CONSUMER_MESSAGE_COTAFACIL",
			ConsumerCount: 3,
			PrefetchCount: 1,
		},
	}
	// Adicione as coleções padrão ao mapa MDB_COLLECTIONS
	defaultCollections := "cfStore, usuarios"
	collectionsMap := parseCollectionsString(defaultCollections)
	default_conf.MongoDBConfig.MDB_COLLECTIONS = collectionsMap
	default_conf.TokenAuth = jwtauth.New("HS256", []byte(default_conf.JWTSecretKey), nil)

	return &default_conf
}

func parseCollectionsString(collectionsString string) map[string]string {
	collections := make(map[string]string)

	collectionNames := strings.Split(collectionsString, ",")

	for _, name := range collectionNames {

		name = strings.TrimSpace(name)
		collections[name] = name
	}

	return collections
}
