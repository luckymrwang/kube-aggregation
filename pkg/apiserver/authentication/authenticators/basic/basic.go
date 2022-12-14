/*
Copyright 2020 The KubeAggregation Authors.

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

package basic

import (
	"context"

	"kube-aggregation/pkg/apiserver/authentication/request/basictoken"
	"kube-aggregation/pkg/models/auth"

	"k8s.io/apiserver/pkg/authentication/authenticator"
	"k8s.io/apiserver/pkg/authentication/user"
)

// TokenAuthenticator implements kubernetes token authenticate interface with our custom logic.
// TokenAuthenticator will retrieve user info from cache by given token. If empty or invalid token
// was given, authenticator will still give passed response at the condition user will be user.Anonymous
// and group from user.AllUnauthenticated. This helps requests be passed along the handler chain,
// because some resources are public accessible.
type basicAuthenticator struct {
	authenticator auth.PasswordAuthenticator
}

func NewBasicAuthenticator(authenticator auth.PasswordAuthenticator) basictoken.Password {
	return &basicAuthenticator{
		authenticator: authenticator,
	}
}

func (t *basicAuthenticator) AuthenticatePassword(ctx context.Context, username, password string) (*authenticator.Response, bool, error) {
	authenticated, _, err := t.authenticator.Authenticate(ctx, username, password)
	if err != nil {
		return nil, false, err
	}
	return &authenticator.Response{
		User: &user.DefaultInfo{
			Name:   authenticated.GetName(),
			Groups: append(authenticated.GetGroups(), user.AllAuthenticated),
		},
	}, true, nil
}
