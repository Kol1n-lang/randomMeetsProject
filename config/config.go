package config

import (
	"fmt"
	"log"
	"strconv"

	"github.com/BurntSushi/toml"
	"github.com/redis/go-redis/v9"
)

type Config struct {
	Database DatabaseConfig    `toml:"database"`
	Server   ServerConfig      `toml:"server"`
	Redis    RedisConfig       `toml:"redis"`
	RabbitMQ RabbitMQConfig    `toml:"rabbitmq"`
	JWT      JWTConfig         `toml:"jwt"`
	External ExternalConfig    `toml:"external"`
	Google   GoogleOAuthConfig `toml:"google"`
}

type DatabaseConfig struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	User     string `toml:"user"`
	Password string `toml:"password"`
	Name     string `toml:"name"`
}

type ServerConfig struct {
	Host string `toml:"host"`
	Port int    `toml:"port"`
}

type RedisConfig struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	Database int    `toml:"database"`
	Cache    int    `toml:"cache_ttl_minutes"`
}

type RabbitMQConfig struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	User     string `toml:"username"`
	Password string `toml:"password"`
}

type JWTConfig struct {
	AccessExpireMinutes int    `toml:"access_expire_minutes"`
	RefreshExpireDays   int    `toml:"refresh_expire_days"`
	SecretKey           string `toml:"secret_key"`
	Algorithm           string `toml:"algorithm"`
}

type ExternalConfig struct {
	Sender      string `toml:"sender"`
	AppPassword string `toml:"app_password"`
	APIKey      string `toml:"api_key"`
}

type GoogleOAuthConfig struct {
	ClientID     string `toml:"client_id"`
	ClientSecret string `toml:"client_secret"`
	RedirectURL  string `toml:"redirect_url"`
}

func LoadConfig(path string) (*Config, error) {
	var config Config
	if _, err := toml.DecodeFile(path, &config); err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &config, nil
}

func (c Config) DbUrl() string {
	cfg, _ := LoadConfig("config.toml")
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.Host,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
	)
}

func (c Config) ServerURL() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}

func (c Config) RedisClient() *redis.Client {
	redisDB := redis.NewClient(&redis.Options{
		Addr: c.Redis.Host + ":" + strconv.Itoa(c.Redis.Port),
		DB:   c.Redis.Database,
	})
	return redisDB
}

func (c Config) RabbitMQUrl() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%d", c.RabbitMQ.User, c.RabbitMQ.Password, c.RabbitMQ.Host, c.RabbitMQ.Port)
}

func (c Config) GoogleOAuthUrl() string {
	half := fmt.Sprintf("https://accounts.google.com/o/oauth2/auth?response_type=code&client_id=%s&redirect_uri=%s&", c.Google.ClientID, c.Google.RedirectURL)
	secondHalf := "access_type=offline&" + "prompt=consent&" + "scope=email%20profile"
	return half + secondHalf
}
