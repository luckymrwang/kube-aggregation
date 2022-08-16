/*
Copyright 2020 KubeSphere Authors

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

package options

import (
	"crypto/tls"
	"flag"
	"fmt"

	"kube-aggregation/pkg/apiserver/authentication/token"

	"k8s.io/client-go/kubernetes/scheme"
	cliflag "k8s.io/component-base/cli/flag"
	"k8s.io/klog"
	runtimecache "sigs.k8s.io/controller-runtime/pkg/cache"
	runtimeclient "sigs.k8s.io/controller-runtime/pkg/client"

	"kube-aggregation/pkg/apis"
	"kube-aggregation/pkg/apiserver"
	apiserverconfig "kube-aggregation/pkg/apiserver/config"
	"kube-aggregation/pkg/informers"
	genericoptions "kube-aggregation/pkg/server/options"
	"kube-aggregation/pkg/simple/client/cache"

	"net/http"
	"strings"

	"kube-aggregation/pkg/simple/client/k8s"
)

type ServerRunOptions struct {
	ConfigFile              string
	GenericServerRunOptions *genericoptions.ServerRunOptions
	*apiserverconfig.Config

	//
	DebugMode bool

	// Enable gops or not.
	GOPSEnabled bool
}

func NewServerRunOptions() *ServerRunOptions {
	s := &ServerRunOptions{
		GenericServerRunOptions: genericoptions.NewServerRunOptions(),
		Config:                  apiserverconfig.New(),
	}

	return s
}

func (s *ServerRunOptions) Flags() (fss cliflag.NamedFlagSets) {
	fs := fss.FlagSet("generic")
	fs.BoolVar(&s.DebugMode, "debug", false, "Don't enable this if you don't know what it means.")
	fs.BoolVar(&s.GOPSEnabled, "gops", false, "Whether to enable gops or not. When enabled this option, "+
		"ks-apiserver will listen on a random port on 127.0.0.1, then you can use the gops tool to list and diagnose the ks-apiserver currently running.")
	s.GenericServerRunOptions.AddFlags(fs, s.GenericServerRunOptions)
	s.KubernetesOptions.AddFlags(fss.FlagSet("kubernetes"), s.KubernetesOptions)
	s.AuthenticationOptions.AddFlags(fss.FlagSet("authentication"), s.AuthenticationOptions)
	s.AuthorizationOptions.AddFlags(fss.FlagSet("authorization"), s.AuthorizationOptions)
	s.RedisOptions.AddFlags(fss.FlagSet("redis"), s.RedisOptions)
	s.LoggingOptions.AddFlags(fss.FlagSet("logging"), s.LoggingOptions)
	s.MultiClusterOptions.AddFlags(fss.FlagSet("multicluster"), s.MultiClusterOptions)

	fs = fss.FlagSet("klog")
	local := flag.NewFlagSet("klog", flag.ExitOnError)
	klog.InitFlags(local)
	local.VisitAll(func(fl *flag.Flag) {
		fl.Name = strings.Replace(fl.Name, "_", "-", -1)
		fs.AddGoFlag(fl)
	})

	return fss
}

const fakeInterface string = "FAKE"

// NewAPIServer creates an APIServer instance using given options
func (s *ServerRunOptions) NewAPIServer(stopCh <-chan struct{}) (*apiserver.APIServer, error) {
	apiServer := &apiserver.APIServer{
		Config: s.Config,
	}

	kubernetesClient, err := k8s.NewKubernetesClient(s.KubernetesOptions)
	if err != nil {
		return nil, err
	}
	apiServer.KubernetesClient = kubernetesClient

	informerFactory := informers.NewInformerFactories(kubernetesClient.Kubernetes(), kubernetesClient.KubeSphere())
	apiServer.InformerFactory = informerFactory

	var cacheClient cache.Interface
	if s.RedisOptions != nil && len(s.RedisOptions.Host) != 0 {
		if s.RedisOptions.Host == fakeInterface && s.DebugMode {
			apiServer.CacheClient = cache.NewSimpleCache()
		} else {
			cacheClient, err = cache.NewRedisClient(s.RedisOptions, stopCh)
			if err != nil {
				return nil, fmt.Errorf("failed to connect to redis service, please check redis status, error: %v", err)
			}
			apiServer.CacheClient = cacheClient
		}
	} else {
		klog.Warning("ks-apiserver starts without redis provided, it will use in memory cache. " +
			"This may cause inconsistencies when running ks-apiserver with multiple replicas.")
		apiServer.CacheClient = cache.NewSimpleCache()
	}

	server := &http.Server{
		Addr: fmt.Sprintf(":%d", s.GenericServerRunOptions.InsecurePort),
	}

	if s.GenericServerRunOptions.SecurePort != 0 {
		certificate, err := tls.LoadX509KeyPair(s.GenericServerRunOptions.TlsCertFile, s.GenericServerRunOptions.TlsPrivateKey)
		if err != nil {
			return nil, err
		}

		server.TLSConfig = &tls.Config{
			Certificates: []tls.Certificate{certificate},
		}
		server.Addr = fmt.Sprintf(":%d", s.GenericServerRunOptions.SecurePort)
	}

	sch := scheme.Scheme
	if err := apis.AddToScheme(sch); err != nil {
		klog.Fatalf("unable add APIs to scheme: %v", err)
	}

	apiServer.RuntimeCache, err = runtimecache.New(apiServer.KubernetesClient.Config(), runtimecache.Options{Scheme: sch})
	if err != nil {
		klog.Fatalf("unable to create controller runtime cache: %v", err)
	}

	apiServer.RuntimeClient, err = runtimeclient.New(apiServer.KubernetesClient.Config(), runtimeclient.Options{Scheme: sch})
	if err != nil {
		klog.Fatalf("unable to create controller runtime client: %v", err)
	}

	apiServer.Issuer, err = token.NewIssuer(s.AuthenticationOptions)
	if err != nil {
		klog.Fatalf("unable to create issuer: %v", err)
	}

	apiServer.Server = server

	return apiServer, nil
}
