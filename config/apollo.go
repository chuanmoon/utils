package config

import (
	"os"

	"github.com/apolloconfig/agollo/v4"
	"github.com/apolloconfig/agollo/v4/env/config"
)

var client agollo.Client
var namespace string

func Init(namespaces []string) {
	appId := os.Getenv("apollo_app_id")
	cluster := os.Getenv("apollo_cluster")
	ip := os.Getenv("apollo_ip")
	secret := os.Getenv("apollo_secret")
	namespace = os.Getenv("apollo_namespace")

	c := &config.AppConfig{
		AppID:          appId,
		Cluster:        cluster,
		IP:             ip,
		IsBackupConfig: true,
		Secret:         secret,
		NamespaceName:  namespace,
	}

	var err error
	client, err = agollo.StartWithConfig(func() (*config.AppConfig, error) {
		return c, nil
	})

	if err != nil {
		panic(err)
	}
}

func String(key, defaultValue string) string {
	return client.GetConfig(namespace).GetStringValue(key, defaultValue)
}
