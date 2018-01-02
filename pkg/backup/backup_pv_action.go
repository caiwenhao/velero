/*
Copyright 2017 the Heptio Ark contributors.

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

package backup

import (
	"github.com/sirupsen/logrus"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/heptio/ark/pkg/apis/ark/v1"
	"github.com/heptio/ark/pkg/util/collections"
)

// backupPVAction inspects a PersistentVolumeClaim for the PersistentVolume
// that it references and backs it up
type backupPVAction struct {
	log logrus.FieldLogger
}

func NewBackupPVAction(log logrus.FieldLogger) ItemAction {
	return &backupPVAction{log: log}
}

var pvGroupResource = schema.GroupResource{Group: "", Resource: "persistentvolumes"}

func (a *backupPVAction) AppliesTo() (ResourceSelector, error) {
	return ResourceSelector{
		IncludedResources: []string{"persistentvolumeclaims"},
	}, nil
}

// Execute finds the PersistentVolume bound by the provided
// PersistentVolumeClaim, if any, and backs it up
func (a *backupPVAction) Execute(item runtime.Unstructured, backup *v1.Backup) (runtime.Unstructured, []ResourceIdentifier, error) {
	a.log.Info("Executing backupPVAction")

	var additionalItems []ResourceIdentifier

	pvc := item.UnstructuredContent()

	volumeName, err := collections.GetString(pvc, "spec.volumeName")
	// if there's no volume name, it's not an error, since it's possible
	// for the PVC not be bound; don't return an additional PV item to
	// back up.
	if err != nil || volumeName == "" {
		a.log.Info("No spec.volumeName found for PersistentVolumeClaim")
		return nil, nil, nil
	}

	additionalItems = append(additionalItems, ResourceIdentifier{
		GroupResource: pvGroupResource,
		Name:          volumeName,
	})

	return item, additionalItems, nil
}
