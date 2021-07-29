// Copyright (c) 2021 Terminus, Inc.
//
// This program is free software: you can use, redistribute, and/or modify
// it under the terms of the GNU Affero General Public License, version 3
// or later ("AGPL"), as published by the Free Software Foundation.
//
// This program is distributed in the hope that it will be useful, but WITHOUT
// ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or
// FITNESS FOR A PARTICULAR PURPOSE.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package crd

import (
	"context"

	apiextensionv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apiextension "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
)

const (
	CRDKind         string = "Dice"
	CRDSingular     string = "dice"
	CRDPlural       string = "dices"
	CRDGroup        string = "dice.terminus.io"
	CRDVersion      string = "v1beta1"
	CRDFullname     string = CRDPlural + "." + CRDGroup
	CRDGroupVersion string = CRDGroup + "/" + CRDVersion
)

// CreateCRD create 'dice' crd if not exists yet
func CreateCRD(config *rest.Config) error {
	client, err := apiextension.NewForConfig(config)
	if err != nil {
		return err
	}
	crd := apiextensionv1beta1.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{Name: CRDFullname},
		Spec: apiextensionv1beta1.CustomResourceDefinitionSpec{
			Group: CRDGroup,
			Versions: []apiextensionv1beta1.CustomResourceDefinitionVersion{{
				Name:    CRDVersion,
				Served:  true,
				Storage: true,
			}},
			Scope: apiextensionv1beta1.NamespaceScoped,
			Names: apiextensionv1beta1.CustomResourceDefinitionNames{
				Plural:   CRDPlural,
				Singular: CRDSingular,
				Kind:     CRDKind,
			},
			AdditionalPrinterColumns: []apiextensionv1beta1.CustomResourceColumnDefinition{
				{
					Name:        "Status",
					Type:        "string",
					Description: "Dice cluster current status",
					JSONPath:    ".status.phase",
				},
				{
					Name:        "LastMessage",
					Type:        "string",
					Description: "last message",
					JSONPath:    ".status.conditions[0].reason",
				},
			},
			Subresources: &apiextensionv1beta1.CustomResourceSubresources{
				Status: &apiextensionv1beta1.CustomResourceSubresourceStatus{},
			},
		},
	}
	_, err = client.ApiextensionsV1beta1().CustomResourceDefinitions().Create(context.Background(), &crd, metav1.CreateOptions{})
	if err != nil && apierrors.IsAlreadyExists(err) {
		return nil
	}
	return err
}
