/*
Copyright 2020 The KubeSphere Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"gopkg.in/yaml.v2"

	"kube-aggregation/pkg/apiserver/authentication"
	"kube-aggregation/pkg/apiserver/authentication/oauth"
	"kube-aggregation/pkg/apiserver/authorization"
	"kube-aggregation/pkg/simple/client/cache"
	"kube-aggregation/pkg/simple/client/k8s"
	"kube-aggregation/pkg/simple/client/ldap"
	"kube-aggregation/pkg/simple/client/logging"
	"kube-aggregation/pkg/simple/client/metering"
	"kube-aggregation/pkg/simple/client/monitoring/prometheus"
	"kube-aggregation/pkg/simple/client/multicluster"
	"kube-aggregation/pkg/simple/client/notification"
)

func newTestConfig() (*Config, error) {
	var conf = &Config{
		KubernetesOptions: &k8s.KubernetesOptions{
			KubeConfig: "/Users/zry/.kube/config",
			Master:     "https://127.0.0.1:6443",
			QPS:        1e6,
			Burst:      1e6,
		},
		LdapOptions: &ldap.Options{
			Host:            "http://openldap.kubesphere-system.svc",
			ManagerDN:       "cn=admin,dc=example,dc=org",
			ManagerPassword: "P@88w0rd",
			UserSearchBase:  "ou=Users,dc=example,dc=org",
			GroupSearchBase: "ou=Groups,dc=example,dc=org",
			InitialCap:      10,
			MaxCap:          100,
			PoolName:        "ldap",
		},
		RedisOptions: &cache.Options{
			Host:     "localhost",
			Port:     6379,
			Password: "KUBESPHERE_REDIS_PASSWORD",
			DB:       0,
		},
		MonitoringOptions: &prometheus.Options{
			Endpoint: "http://prometheus.kubesphere-monitoring-system.svc",
		},
		LoggingOptions: &logging.Options{
			Host:        "http://elasticsearch-logging.kubesphere-logging-system.svc:9200",
			IndexPrefix: "elk",
			Version:     "6",
		},
		NotificationOptions: &notification.Options{
			Endpoint: "http://notification.kubesphere-alerting-system.svc:9200",
		},
		AuthorizationOptions: authorization.NewOptions(),
		AuthenticationOptions: &authentication.Options{
			AuthenticateRateLimiterMaxTries: 5,
			AuthenticateRateLimiterDuration: 30 * time.Minute,
			JwtSecret:                       "xxxxxx",
			LoginHistoryMaximumEntries:      100,
			MultipleLogin:                   false,
			OAuthOptions: &oauth.Options{
				Issuer:            oauth.DefaultIssuer,
				IdentityProviders: []oauth.IdentityProviderOptions{},
				Clients: []oauth.Client{{
					Name:                         "kubesphere-console-client",
					Secret:                       "xxxxxx-xxxxxx-xxxxxx",
					RespondWithChallenges:        true,
					RedirectURIs:                 []string{"http://ks-console.kubesphere-system.svc/oauth/token/implicit"},
					GrantMethod:                  oauth.GrantHandlerAuto,
					AccessTokenInactivityTimeout: nil,
				}},
				AccessTokenMaxAge:            time.Hour * 24,
				AccessTokenInactivityTimeout: 0,
			},
		},
		MultiClusterOptions: multicluster.NewOptions(),
		MeteringOptions: &metering.Options{
			RetentionDay: "7d",
		},
	}
	return conf, nil
}

func saveTestConfig(t *testing.T, conf *Config) {
	content, err := yaml.Marshal(conf)
	if err != nil {
		t.Fatalf("error marshal config. %v", err)
	}
	err = ioutil.WriteFile(fmt.Sprintf("%s.yaml", defaultConfigurationName), content, 0640)
	if err != nil {
		t.Fatalf("error write configuration file, %v", err)
	}
}

func cleanTestConfig(t *testing.T) {
	file := fmt.Sprintf("%s.yaml", defaultConfigurationName)
	if _, err := os.Stat(file); os.IsNotExist(err) {
		t.Log("file not exists, skipping")
		return
	}

	err := os.Remove(file)
	if err != nil {
		t.Fatalf("remove %s file failed", file)
	}

}

func TestGet(t *testing.T) {
	conf, err := newTestConfig()
	if err != nil {
		t.Fatal(err)
	}
	saveTestConfig(t, conf)
	defer cleanTestConfig(t)

	conf.RedisOptions.Password = "P@88w0rd"
	os.Setenv("KUBESPHERE_REDIS_PASSWORD", "P@88w0rd")

	conf2, err := TryLoadFromDisk()
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(conf, conf2); diff != "" {
		t.Fatal(diff)
	}
}

func TestStripEmptyOptions(t *testing.T) {
	var config Config

	config.RedisOptions = &cache.Options{Host: ""}
	config.MonitoringOptions = &prometheus.Options{Endpoint: ""}
	config.LdapOptions = &ldap.Options{Host: ""}
	config.LoggingOptions = &logging.Options{Host: ""}
	config.NotificationOptions = &notification.Options{Endpoint: ""}
	config.MultiClusterOptions = &multicluster.Options{Enable: false}

	config.stripEmptyOptions()

	if config.RedisOptions != nil ||
		config.MonitoringOptions != nil ||
		config.LdapOptions != nil ||
		config.LoggingOptions != nil ||
		config.NotificationOptions != nil ||
		config.MultiClusterOptions != nil {
		t.Fatal("config stripEmptyOptions failed")
	}
}
