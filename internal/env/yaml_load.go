/*
 * This file was last modified at 2024-07-29 15:17 by Victor N. Skurikhin.
 * This is free and unencumbered software released into the public domain.
 * For more information, please refer to <http://unlicense.org>
 * yaml_load.go
 * $Id$
 */
//!+

// Package env работа с настройками и окружением.
package env

import (
	"github.com/spf13/viper"
	"os"
)

// LoadConfig - TODO.
// Example:
//
// favorites:
//
//	Enabled: true
//	db:
//	    Enabled: true
//	    name: db
//	    host: localhost
//	    port: 5432
//	    username: dbuser
//	    password: password
//	grpc:
//	    address: localhost
//	    Enabled: true
//	    port: 8443
//	    proto: tcp
//	    tls:
//	        Enabled: true
//	        ca_file: cert/test_ca-cert.pem
//	        cert_file: cert/test_server-cert.pem
//	        key_file: cert/test_server-key.pem
//	http:
//	    address: localhost
//	    Enabled: true
//	    port: 443
//	    tls:
//	        Enabled: true
//	        ca_file: cert/test_ca-cert.pem
//	        cert_file: cert/test_server-cert.pem
//	        key_file: cert/test_server-key.pem
//	token: 89h3f98hbwf987h3f98wenf89ehf
func LoadConfig(path string) (cfg Config, err error) {

	if os.Getenv("GO_FAVORITES_SKIP_LOAD_CONFIG") != "" {
		return &config{}, err
	}
	viper.SetConfigName("go-favorites.yaml")  // name of yamlConfig file (without extension)
	viper.SetConfigType("yaml")               // REQUIRED if the yamlConfig file does not have the extension in the name
	viper.AddConfigPath("/etc/go-favorites/") // path to look for the yamlConfig file in
	viper.AddConfigPath("$HOME/go-favorites") // call multiple times to add many search paths
	viper.AddConfigPath(path)                 // optionally look for yamlConfig in the working directory

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
