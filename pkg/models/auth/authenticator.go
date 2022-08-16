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

package auth

import (
	"context"
	"fmt"
	authuser "k8s.io/apiserver/pkg/authentication/user"
	"net/http"
)

var (
	RateLimitExceededError  = fmt.Errorf("auth rate limit exceeded")
	IncorrectPasswordError  = fmt.Errorf("incorrect password")
	AccountIsNotActiveError = fmt.Errorf("account is not active")
)

// PasswordAuthenticator is an interface implemented by authenticator which take a
// username and password.
type PasswordAuthenticator interface {
	Authenticate(ctx context.Context, username, password string) (authuser.Info, string, error)
}

type OAuthAuthenticator interface {
	Authenticate(ctx context.Context, provider string, req *http.Request) (authuser.Info, string, error)
}
