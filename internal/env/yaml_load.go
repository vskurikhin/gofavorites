/*
 * This file was last modified at 2024-08-06 18:25 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * yaml_load.go
 * $Id$
 */
//!+

// Package env работа с настройками и окружением.
package env

import (
	"os"

	"github.com/spf13/viper"
)

// LoadConfig - TODO.
// Example:
//
// favorites:
//
//	enabled: true
//	cache:
//	  enabled: true
//	  expire_ms: 1000
//	  gc_interval_sec: 10
//	db:
//	  enabled: true
//	  name: db
//	  host: localhost
//	  port: 5432
//	  username: dbuser
//	  password: password
//	external:
//	  asset_grpc_address: 127.0.0.1:8444
//	  auth_grpc_address: 127.0.0.1:8444
//	  request_timeout_interval_ms: 911
//	grpc:
//	  address: localhost
//	  enabled: true
//	  port: 8442
//	  proto: tcp
//	  tls:
//	    enabled: true
//	    ca_file: cert/grpc-ca-cert.pem
//	    cert_file: cert/grpc-server-cert.pem
//	    key_file: cert/grpc-server-key.pem
//	http:
//	  address: localhost
//	  enabled: true
//	  port: 8443
//	  tls:
//	    enabled: true
//	    ca_file: cert/http-ca-cert.pem
//	    cert_file: cert/http-server-cert.pem
//	    key_file: cert/http-server-key.pem
//	jwt:
//	  jwt_secret: TzzVGdLUJGcYKaf5he4zeLW5QdSJws9UoUug3Q3kCMeLVijBSjPY3k0pNu2XWhB
//	  jwt_expired_in: 60m
//	  jwt_max_age_sec: 3600
//	mongo:
//	  enabled: true
//	  name: db
//	  host: localhost
//	  port: 27017
//	  username: mongouser
//	  password: password
//	token: '$2a$11$ZTzzVGdLUJGcYKJws9UoUug3Q3kCMELVziajBSJPY3k0pNu2XWHBy'
//	upk:
//	  rsa_private_key_file: cert/upk-private-key.pem
//	  rsa_public_key_file: cert/upk-public-key.pem
//	  secret: g16Ug0b1zVaCYQzxD45C6p99fUxMkaSL2npjmi5qBjRMAy6kjpZP0/zahKE4zQGvTlp7lKavV4z3RWIm9Uch1pBgaYLZ/pAZDbgr8roqVc/QEzQnsaLqoe7ZzOcPsj7NzbrXz/l+rWVAGdyAkLGs7NIZ3GgNlyZ5lrjglAIRdHA6PpW0jBzbcKb5Z5Y5U80N75+wrenlWPFUKTrN8exuUhzLK6FHWpAzuivD+pg42bZFvdSLE/0oXd0U1W+SxSBXv3RxEkRMquYG+9/VHpT745BzF+QQlR+CicLC5XaUusAZKtqFf3LokISPY1kxjP32gW3SqtZThZa/4pPMpesrXA==
func LoadConfig(path string) (cfg Config, err error) {

	if os.Getenv("GO_FAVORITES_SKIP_LOAD_CONFIG") != "" {
		return &config{}, err
	}
	viper.SetConfigName("go-favorites.yaml")  // мя файла yamlConfig
	viper.SetConfigType("yaml")               // REQUIRED если файл yamlConfig не имеет расширения в имени
	viper.AddConfigPath("/etc/go-favorites/") // путь для поиска файла yamlConfig
	viper.AddConfigPath("$HOME/go-favorites") // несколько раз, чтобы добавить несколько путей поиска
	viper.AddConfigPath(path)

	err = viper.ReadInConfig() // Find and read the yamlConfig file
	if err != nil {            // Handle errors reading the yamlConfig file
		return
	}
	var c config
	err = viper.Unmarshal(&c)
	if err != nil {
		return
	}
	cfg = &c

	return
}

//!-
/* vim: set tabstop=4 softtabstop=4 shiftwidth=4 noexpandtab: */
