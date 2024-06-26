package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/fahmifan/devkit/pkg/logs"
	"github.com/fahmifan/devkit/pkg/mailer/smtp"
	"github.com/gomodule/redigo/redis"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		logs.Warn(".env file not found")
		return
	}

	logs.Info("load .env file to os env")
}

// Port ..
func Port() string {
	return os.Getenv("PORT")
}

// Env ..
func Env() string {
	if val, ok := os.LookupEnv("ENV"); ok {
		return val
	}

	return "development"
}

// JWTKey ..
func JWTKey() string {
	val, ok := os.LookupEnv("JWT_SECRET")
	if !ok {
		log.Fatal("JWT_SECRET not provided")
	}
	return val
}

// BaseURL ..
func BaseURL() string {
	if val, ok := os.LookupEnv("BASE_URL"); ok {
		return val
	}

	return fmt.Sprintf("http://localhost:%s", os.Getenv("PORT"))
}

func WebBaseURL() string {
	if val, ok := os.LookupEnv("WEB_BASE_URL"); ok {
		return val
	}

	return fmt.Sprintf("http://localhost:%s", os.Getenv("PORT"))
}

// PostgresDSN ..
func PostgresDSN() string {
	port := os.Getenv("DB_PORT")
	dbname := os.Getenv("DB_NAME")
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	sslmode := os.Getenv("DB_SSLMODE")

	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		user,
		password,
		host,
		port,
		dbname,
		sslmode)
}

// WorkerNamespace ..
func WorkerNamespace() string {
	return "autograd_worker"
}

// WorkerConcurrency ..
func WorkerConcurrency() uint {
	return 5
}

// RedisWorkerHost ..
func RedisWorkerHost() string {
	return os.Getenv("REDIS_WORKER_HOST")
}

// NewRedisPool ..
func NewRedisPool(host string) *redis.Pool {
	return &redis.Pool{
		MaxActive: 5,
		MaxIdle:   5,
		Wait:      true,
		Dial: func() (redis.Conn, error) {
			return redis.DialURL(host)
		},
	}
}

// FileUploadPath ..
func FileUploadPath() string {
	val, ok := os.LookupEnv("FILE_UPLOAD_PATH")
	if ok {
		return val
	}

	return "file_upload_path"
}

func AutogradAuthToken() string {
	val, _ := os.LookupEnv("AUTOGRAD_AUTH_TOKEN")
	return val
}

func AutogradServerURL() string {
	val, _ := os.LookupEnv("AUTOGRAD_SERVER_URL")
	return val
}

func SenderEmail() string {
	return os.Getenv("SENDER_EMAIL")
}

func SMTPConfig() smtp.Config {
	return smtp.Config{
		Host:     os.Getenv("SMTP_HOST"),
		Port:     parseInt(os.Getenv("SMTP_PORT")),
		Username: os.Getenv("SMTP_USERNAME"),
		Password: os.Getenv("SMTP_PASSWORD"),
	}
}

func Debug() bool {
	val, _ := os.LookupEnv("DEBUG")
	return val == "true"
}

func parseInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}
