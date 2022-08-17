/*
Copyright 2019 The KubeAggregation Authors.

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

package jwt

import (
	"context"

	"k8s.io/apiserver/pkg/authentication/authenticator"
	"k8s.io/apiserver/pkg/authentication/user"
	"k8s.io/klog"

	"kube-aggregation/pkg/models/auth"
)

// TokenAuthenticator implements kubernetes token authenticate interface with our custom logic.
// TokenAuthenticator will retrieve user info from cache by given token. If empty or invalid token
// was given, authenticator will still give passed response at the condition user will be user.Anonymous
// and group from user.AllUnauthenticated. This helps requests be passed along the handler chain,
// because some resources are public accessible.
type tokenAuthenticator struct {
	tokenOperator auth.TokenManagementInterface
}

func NewTokenAuthenticator(tokenOperator auth.TokenManagementInterface) authenticator.Token {
	return &tokenAuthenticator{
		tokenOperator: tokenOperator,
	}
}

func (t *tokenAuthenticator) AuthenticateToken(ctx context.Context, token string) (*authenticator.Response, bool, error) {
	verified, err := t.tokenOperator.Verify(token)
	if err != nil {
		klog.Warning(err)
		return nil, false, err
	}

	// AuthLimitExceeded state should be ignored
	return &authenticator.Response{
		User: &user.DefaultInfo{
			Name:   verified.User.GetName(),
			Groups: append(verified.User.GetGroups(), user.AllAuthenticated),
		},
	}, true, nil
}
