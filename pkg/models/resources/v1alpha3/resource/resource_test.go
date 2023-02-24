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

package resource

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	fakek8s "k8s.io/client-go/kubernetes/fake"

	"github.com/clusterpedia-io/clusterpedia/pkg/api"
	"github.com/clusterpedia-io/clusterpedia/pkg/apiserver/query"
	fakeks "github.com/clusterpedia-io/clusterpedia/pkg/client/clientset/versioned/fake"
	"github.com/clusterpedia-io/clusterpedia/pkg/informers"
)

func TestResourceGetter(t *testing.T) {

	resource := prepare()

	tests := []struct {
		Name           string
		Resource       string
		Namespace      string
		Query          *query.Query
		ExpectError    error
		ExpectResponse *api.ListResult
	}{
		{
			Name:      "normal case",
			Resource:  "namespaces",
			Namespace: "",
			Query: &query.Query{
				Pagination: &query.Pagination{
					Limit:  10,
					Offset: 0,
				},
				SortBy:    query.FieldName,
				Ascending: false,
				Filters:   map[query.Field]query.Value{},
			},
			ExpectError: nil,
			ExpectResponse: &api.ListResult{
				Items:      []interface{}{foo2, foo1, bar1},
				TotalItems: 3,
			},
		},
	}

	for _, test := range tests {

		result, err := resource.List(test.Resource, test.Namespace, test.Query)

		if err != test.ExpectError {
			t.Errorf("expected error: %s, got: %s", test.ExpectError, err)
		}
		if diff := cmp.Diff(test.ExpectResponse, result); diff != "" {
			t.Errorf(diff)
		}
	}
}

var (
	foo1 = &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "foo1",
			Namespace: "bar",
		},
	}

	foo2 = &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "foo2",
			Namespace: "bar",
		},
	}
	bar1 = &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "bar1",
			Namespace: "bar",
		},
	}

	namespaces = []interface{}{foo1, foo2, bar1}
)

func prepare() *ResourceGetter {

	ksClient := fakeks.NewSimpleClientset()
	k8sClient := fakek8s.NewSimpleClientset()
	fakeInformerFactory := informers.NewInformerFactories(k8sClient, ksClient)

	for _, namespace := range namespaces {
		fakeInformerFactory.KubernetesSharedInformerFactory().Core().V1().
			Namespaces().Informer().GetIndexer().Add(namespace)
	}

	return NewResourceGetter(fakeInformerFactory)
}
