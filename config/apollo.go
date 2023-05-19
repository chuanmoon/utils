package config

import (
	"os"

	"github.com/apolloconfig/agollo/v4"
	"github.com/apolloconfig/agollo/v4/env/config"
)

var client agollo.Client
var namespace string

func init() {
	appId := os.Getenv("CHUANMOON_APOLLO_APP_ID")
	cluster := os.Getenv("CHUANMOON_APOLLO_CLUSTER")
	ip := os.Getenv("CHUANMOON_APOLLO_IP")
	secret := os.Getenv("CHUANMOON_APOLLO_SECRET")
	namespace = os.Getenv("CHUANMOON_APOLLO_NAMESPACE")

	if appId == "" || cluster == "" || ip == "" || secret == "" || namespace == "" {
		panic("Apollo config not set, please check your environment variables: CHUANMOON_APOLLO_APP_ID, CHUANMOON_APOLLO_CLUSTER, CHUANMOON_APOLLO_IP, CHUANMOON_APOLLO_SECRET, CHUANMOON_APOLLO_NAMESPACE")
	}

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
