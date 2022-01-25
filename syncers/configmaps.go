package syncers

import (
	"github.com/loft-sh/vcluster-sdk/syncer"
	syncercontext "github.com/loft-sh/vcluster-sdk/syncer/context"
	"github.com/loft-sh/vcluster-sdk/syncer/translator"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func NewConfigMapSyncer(ctx *syncercontext.RegisterContext) syncer.Syncer {
	return &configMapSyncer{
		NamespacedTranslator: translator.NewNamespacedTranslator(ctx, "configmap", &corev1.ConfigMap{}),
	}
}

type configMapSyncer struct {
	translator.NamespacedTranslator
}

// Make sure the interface is implemented
var _ syncer.Syncer = &configMapSyncer{}

func (s *configMapSyncer) SyncDown(ctx *syncercontext.SyncContext, vObj client.Object) (ctrl.Result, error) {
	return s.SyncDownCreate(ctx, vObj, s.translate(vObj.(*corev1.ConfigMap)))
}

func (s *configMapSyncer) Sync(ctx *syncercontext.SyncContext, pObj client.Object, vObj client.Object) (ctrl.Result, error) {
	pConfigMap := pObj.(*corev1.ConfigMap)
	vConfigMap := vObj.(*corev1.ConfigMap)
	if pConfigMap.Immutable != nil && vConfigMap.Immutable != nil && *pConfigMap.Immutable && !*vConfigMap.Immutable {
		// if the ConfigMap in the host is Immutable, while ConfigMap in vcluster
		// is not Immutable, then we need to delete it from the host to reconcile
		// it into the expected state. We force requeue to trigger recreation.
		_, err := syncer.DeleteObject(ctx, pObj)
		return ctrl.Result{Requeue: true}, err
	}

	return s.SyncDownUpdate(ctx, vObj, s.translateUpdate(pConfigMap, vConfigMap))
}

func (s *configMapSyncer) translate(vObj client.Object) *corev1.ConfigMap {
	return s.TranslateMetadata(vObj).(*corev1.ConfigMap)
}

// translateUpdate returns nil if the host side ConfigMap doesn't need to be updated,
// otherwise it returns an updated ConfigMap.
// Note: the caller has to cover the case where the vObj.Immutable is true, and pObj.Immutable is false
func (s *configMapSyncer) translateUpdate(pObj, vObj *corev1.ConfigMap) *corev1.ConfigMap {
	var updated *corev1.ConfigMap

	// check if the annotations or labels have changed
	changed, updatedAnnotations, updatedLabels := s.TranslateMetadataUpdate(vObj, pObj)
	if changed {
		updated = newIfNil(updated, pObj)
		updated.Labels = updatedLabels
		updated.Annotations = updatedAnnotations
	}

	// check if the data has changed
	if !equality.Semantic.DeepEqual(vObj.Data, pObj.Data) {
		updated = newIfNil(updated, pObj)
		updated.Data = vObj.Data
	}

	// check if the binary data has changed
	if !equality.Semantic.DeepEqual(vObj.BinaryData, pObj.BinaryData) {
		updated = newIfNil(updated, pObj)
		updated.BinaryData = vObj.BinaryData
	}

	// check if the Immutable field has changed
	// Note: the caller has to cover the case where the vObj.Immutable is true, and pObj.Immutable is false
	if !equality.Semantic.DeepEqual(vObj.Immutable, pObj.Immutable) {
		updated = newIfNil(updated, pObj)
		updated.Immutable = vObj.Immutable
	}

	return updated
}

func newIfNil(updated *corev1.ConfigMap, pObj *corev1.ConfigMap) *corev1.ConfigMap {
	if updated == nil {
		return pObj.DeepCopy()
	}
	return updated
}
