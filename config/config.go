package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	ServerConfig     `json:"server"`
	PostgresConfig   `json:"postgres"`
	ClickhouseConfig `json:"clickhouse"`
	RedisConfig      `json:"redis"`
	NatsConfig       `json:"nats"`
	LogFilePath      string `json:"logFilePath"`
}

type ServerConfig struct {
	Addr         string `json:"addr"`
	ReadTimeout  int    `json:"read_timeout"`
	WriteTimeout int    `json:"write_timeout"`
}

type PostgresConfig struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	DbName   string `json:"db_name"`
	SslMode  string `json:"sslMode"`
}

type ClickhouseConfig struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	DbName   string `json:"db_name"`
}

type RedisConfig struct {
	Addr     string `json:"addr"`
	Password string `json:"password"`
	Db       int    `json:"db"`
	TTL      int    `json:"ttl"`
}

type NatsConfig struct {
	Url string `json:"Url"`
}

func NewConfig() (Config, error) {
	var config Config

	data, err := os.ReadFile("./config/config.json")
	if err != nil {
		return config, fmt.Errorf("error while trying to read config file: %w", err)
	}

	if err = json.Unmarshal(data, &config); err != nil {
		return config, fmt.Errorf("error while trying to unmarshall config json: %w", err)
	}

	return config, nil
}

func (pg PostgresConfig) ConnectionString() string {
	return fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=%s", pg.Host, pg.Port, pg.Username, pg.DbName, pg.Password, pg.SslMode)
}

func (ch ClickhouseConfig) ConnectionString() string {
	return fmt.Sprintf("http://%s:%d/%s?username=%s&password=%s", ch.Host, ch.Port, ch.DbName, ch.Username, ch.Password)
}
