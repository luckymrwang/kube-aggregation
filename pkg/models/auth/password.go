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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	kubesphere "kube-aggregation/pkg/client/clientset/versioned"

	"kube-aggregation/pkg/apiserver/authentication"

	"golang.org/x/crypto/bcrypt"
	"k8s.io/apimachinery/pkg/api/errors"
	authuser "k8s.io/apiserver/pkg/authentication/user"
	"k8s.io/klog"
	"kube-aggregation/pkg/apiserver/authentication/identityprovider"
	"kube-aggregation/pkg/apiserver/authentication/oauth"
	"kube-aggregation/pkg/constants"
)

type passwordAuthenticator struct {
	ksClient    kubesphere.Interface
	userGetter  *userGetter
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
	// generic identity provider has higher priority
	for _, providerOptions := range p.authOptions.OAuthOptions.IdentityProviders {
		// the admin account in kubesphere has the highest priority
		if username == constants.AdminUserName {
			break
		}
		if genericProvider, _ := identityprovider.GetGenericProvider(providerOptions.Name); genericProvider != nil {
			authenticated, err := genericProvider.Authenticate(username, password)
			if err != nil {
				if errors.IsUnauthorized(err) {
					continue
				}
				return nil, providerOptions.Name, err
			}
			linkedAccount, err := p.userGetter.findMappedUser(providerOptions.Name, authenticated.GetUserID())
			if err != nil && !errors.IsNotFound(err) {
				klog.Error(err)
				return nil, providerOptions.Name, err
			}
			// using this method requires you to manually provision users.
			if providerOptions.MappingMethod == oauth.MappingMethodLookup && linkedAccount == nil {
				continue
			}
			// the user will automatically create and mapping when login successful.
			if linkedAccount == nil && providerOptions.MappingMethod == oauth.MappingMethodAuto {
				if !providerOptions.DisableLoginConfirmation {
					return preRegistrationUser(providerOptions.Name, authenticated), providerOptions.Name, nil
				}

				linkedAccount, err = p.ksClient.IamV1alpha2().Users().Create(context.Background(), mappedUser(providerOptions.Name, authenticated), metav1.CreateOptions{})
				if err != nil {
					return nil, providerOptions.Name, err
				}
			}
			if linkedAccount != nil {
				return &authuser.DefaultInfo{Name: linkedAccount.GetName()}, providerOptions.Name, nil
			}
		}
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
