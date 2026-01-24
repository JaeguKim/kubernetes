/*
Copyright 2026 The Kubernetes Authors.

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

package configmap

import (
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
	genericapirequest "k8s.io/apiserver/pkg/endpoints/request"
	apitesting "k8s.io/kubernetes/pkg/api/testing"
	api "k8s.io/kubernetes/pkg/apis/core"
)

func TestDeclarativeValidateForDeclarative(t *testing.T) {
	ctx := genericapirequest.WithRequestInfo(genericapirequest.NewDefaultContext(), &genericapirequest.RequestInfo{
		APIGroup:   "",
		APIVersion: "v1",
	})
	testCases := map[string]struct {
		input        api.ConfigMap
		expectedErrs field.ErrorList
	}{
		// baseline
		"valid resource": {
			input: mkValidConfigMap(),
		},
		// metadata.name
		"name: empty": {
			input: mkValidConfigMap(func(cm *api.ConfigMap) {
				cm.Name = ""
			}),
			expectedErrs: field.ErrorList{
				field.Required(field.NewPath("metadata.name"), ""),
			},
		},
		// TODO: Add more test cases as validation rules are migrated to declarative.
	}
	for k, tc := range testCases {
		t.Run(k, func(t *testing.T) {
			apitesting.VerifyValidationEquivalence(t, ctx, &tc.input, Strategy.Validate, tc.expectedErrs)
		})
	}
}

func TestValidateUpdateForDeclarative(t *testing.T) {
	ctx := genericapirequest.WithRequestInfo(genericapirequest.NewDefaultContext(), &genericapirequest.RequestInfo{
		APIGroup:   "",
		APIVersion: "v1",
	})
	testCases := map[string]struct {
		old          api.ConfigMap
		update       api.ConfigMap
		expectedErrs field.ErrorList
	}{
		// baseline
		"no change": {
			old:    mkValidConfigMap(),
			update: mkValidConfigMap(),
		},
		// metadata.name
		"name: changed": {
			old: mkValidConfigMap(),
			update: mkValidConfigMap(func(cm *api.ConfigMap) {
				cm.Name += "x"
			}),
			expectedErrs: field.ErrorList{
				field.Invalid(field.NewPath("metadata.name"), nil, ""),
			},
		},
		// TODO: Add more test cases as validation rules are migrated to declarative.
	}
	for k, tc := range testCases {
		t.Run(k, func(t *testing.T) {
			tc.old.ObjectMeta.ResourceVersion = "1"
			tc.update.ObjectMeta.ResourceVersion = "1"
			apitesting.VerifyUpdateValidationEquivalence(t, ctx, &tc.update, &tc.old, Strategy.ValidateUpdate, tc.expectedErrs)
		})
	}
}

// mkValidConfigMap produces a ConfigMap which passes validation with no tweaks.
func mkValidConfigMap(tweaks ...func(cm *api.ConfigMap)) api.ConfigMap {
	cm := api.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{Name: "abc", Namespace: metav1.NamespaceDefault},
		Data:       map[string]string{"key": "value"},
	}
	for _, tweak := range tweaks {
		tweak(&cm)
	}
	return cm
}
