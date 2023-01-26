package config

import (
	"errors"
	"strings"
	"time"

	"github.com/AltMax/art-test/postgresql"
	"github.com/jackc/pgx"
	"github.com/spf13/viper"
)

// Config contains all configurable vars for apps.
type Config struct {
	ServerAddr        string            `mapstructure:"server_addr"`
	Postgresql        postgresql.Config `mapstructure:"postgresql"`
	LRUCacheSize      int               `mapstructure:"lru_cache_size"`
	FetchUnitsTimeout int64             `mapstructure:"fetch_units_timeout"` //seconds
}

func New() (Config, error) {
	configInstance := Config{}
	if err := viper.Unmarshal(&configInstance); err != nil {
		return Config{}, errors.New("can't load config structure")
	}
	return configInstance, nil
}

func init() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetConfigType("toml")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig() // Find and read the config file
	_ = err

	viper.SetDefault("server_addr", ":10000")

	// PostgreSQL
	viper.SetDefault("postgresql.port", 5432)
	viper.SetDefault("postgresql.host", "127.0.0.1")
	viper.SetDefault("postgresql.database", "unit_service_test")
	viper.SetDefault("postgresql.user", "postgres")
	viper.SetDefault("postgresql.password", "")
	viper.SetDefault("postgresql.secured", false)
	viper.SetDefault("postgresql.max_connections", 20)
	viper.SetDefault("postgresql.max_connection_age", time.Minute)
	viper.SetDefault("postgresql.health_check_period", time.Second*10)
	viper.SetDefault("postgresql.logger_enabled", false)
	viper.SetDefault("postgresql.log_level", pgx.LogLevelError)
	viper.SetDefault("postgresql.keep_alive", time.Second*15)

	viper.SetDefault("lru_cache_size", 500)
	viper.SetDefault("fetch_units_timeout", 60*60) //1h
}
