package nifi

import (
	"fmt"

	v1 "github.com/konpyutaika/nifikop/api/v1"

	"github.com/konpyutaika/nifikop/pkg/resources/templates"
	"github.com/konpyutaika/nifikop/pkg/util"
	nifiutil "github.com/konpyutaika/nifikop/pkg/util/nifi"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func (r *Reconciler) pvc(id int32, storage v1.StorageConfig, log zap.Logger) runtime.Object {
	return &corev1.PersistentVolumeClaim{
		ObjectMeta: templates.ObjectMetaWithGeneratedNameAndAnnotations(
			// name
			fmt.Sprintf(templates.NodeStorageTemplate, r.NifiCluster.Name, id, storage.Name),
			// labels
			util.MergeLabels(
				nifiutil.LabelsForNifi(r.NifiCluster.Name),
				storage.Metadata.Labels,
				map[string]string{
					"nodeId":                        fmt.Sprintf("%d", id),
					"storageName":                   storage.Name,
					nifiutil.NifiDataVolumeMountKey: "true",
				},
			),
			// annotations
			util.MergeAnnotations(
				storage.Metadata.Annotations,
				map[string]string{"mountPath": storage.MountPath, "storageName": storage.Name},
			),
			r.NifiCluster),
		Spec: *storage.PVCSpec,
	}
}
