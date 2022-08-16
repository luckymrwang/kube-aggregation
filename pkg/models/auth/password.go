/*

 Copyright 2021 The KubeSphere Authors.

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
	kubesphere "kube-aggregation/pkg/client/clientset/versioned"

	"kube-aggregation/pkg/apiserver/authentication"

	"golang.org/x/crypto/bcrypt"
	authuser "k8s.io/apiserver/pkg/authentication/user"
)

type passwordAuthenticator struct {
	ksClient    kubesphere.Interface
	authOptions *authentication.Options
}

func NewPasswordAuthenticator(ksClient kubesphere.Interface,
	options *authentication.Options) PasswordAuthenticator {
	passwordAuthenticator := &passwordAuthenticator{
		ksClient:    ksClient,
		authOptions: options,
	}
	return passwordAuthenticator
}

func (p *passwordAuthenticator) Authenticate(_ context.Context, username, password string) (authuser.Info, string, error) {
	// empty username or password are not allowed
	if username == "" || password == "" {
		return nil, "", IncorrectPasswordError
	}

	// if the password is not empty, means that the password has been reset, even if the user was mapping from IDP
	u := &authuser.DefaultInfo{
		Name: username,
	}

	return u, "", nil
}

func PasswordVerify(encryptedPassword, password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(encryptedPassword), []byte(password)); err != nil {
		return IncorrectPasswordError
	}
	return nil
}
