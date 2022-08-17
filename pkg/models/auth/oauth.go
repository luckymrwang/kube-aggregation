/*

 Copyright 2021 The KubeAggregation Authors.

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

package auth

import (
	"context"
	"net/http"

	kubesphere "kube-aggregation/pkg/client/clientset/versioned"

	"kube-aggregation/pkg/apiserver/authentication"

	authuser "k8s.io/apiserver/pkg/authentication/user"
	"k8s.io/klog"
)

type oauthAuthenticator struct {
	ksClient kubesphere.Interface
	options  *authentication.Options
}

func NewOAuthAuthenticator(ksClient kubesphere.Interface,
	options *authentication.Options) OAuthAuthenticator {
	authenticator := &oauthAuthenticator{
		ksClient: ksClient,
		options:  options,
	}
	return authenticator
}

func (o *oauthAuthenticator) Authenticate(_ context.Context, provider string, req *http.Request) (authuser.Info, string, error) {
	providerOptions, err := o.options.OAuthOptions.IdentityProviderOptions(provider)
	// identity provider not registered
	if err != nil {
		klog.Error(err)
		return nil, "", err
	}

	return &authuser.DefaultInfo{Name: "admin"}, providerOptions.Name, nil
}
