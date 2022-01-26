package syncers

import (
	"testing"

	"github.com/loft-sh/vcluster-sdk/syncer"
	synccontext "github.com/loft-sh/vcluster-sdk/syncer/context"
	generictesting "github.com/loft-sh/vcluster-sdk/syncer/testing"
	"github.com/loft-sh/vcluster-sdk/syncer/translator"
	"github.com/loft-sh/vcluster-sdk/translate"
	"gotest.tools/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func TestSync(t *testing.T) {
	baseConfigMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-configmap",
			Namespace: "test",
		},
	}
	updatedConfigMap := &corev1.ConfigMap{
		ObjectMeta: baseConfigMap.ObjectMeta,
		Data: map[string]string{
			"test": "test",
		},
	}
	syncedConfigMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      translate.PhysicalName(baseConfigMap.Name, baseConfigMap.Namespace),
			Namespace: "test",
			Annotations: map[string]string{
				translator.NameAnnotation:      baseConfigMap.Name,
				translator.NamespaceAnnotation: baseConfigMap.Namespace,
			},
			Labels: map[string]string{
				translate.NamespaceLabel: baseConfigMap.Namespace,
			},
		},
	}
	updatedSyncedConfigMap := &corev1.ConfigMap{
		ObjectMeta: syncedConfigMap.ObjectMeta,
		Data:       updatedConfigMap.Data,
	}
	basePod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: baseConfigMap.Namespace,
		},
		Spec: corev1.PodSpec{
			Volumes: []corev1.Volume{
				{
					Name: "test",
					VolumeSource: corev1.VolumeSource{
						ConfigMap: &corev1.ConfigMapVolumeSource{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: baseConfigMap.Name,
							},
						},
					},
				},
			},
		},
	}

	generictesting.RunTests(t, []*generictesting.SyncTest{
		{
			Name: "Unused config map",
			InitialVirtualState: []runtime.Object{
				baseConfigMap,
			},
			ExpectedPhysicalState: map[schema.GroupVersionKind][]runtime.Object{
				corev1.SchemeGroupVersion.WithKind("ConfigMap"): {
					syncedConfigMap,
				},
			},
			Sync: func(ctx *synccontext.RegisterContext) {
				syncCtx, syncer := generictesting.FakeStartSyncer(t, ctx, newSyncer)
				_, err := syncer.(*configMapSyncer).SyncDown(syncCtx, baseConfigMap)
				assert.NilError(t, err)
			},
		},
		{
			Name: "Used config map",
			InitialVirtualState: []runtime.Object{
				baseConfigMap,
				basePod,
			},
			ExpectedPhysicalState: map[schema.GroupVersionKind][]runtime.Object{
				corev1.SchemeGroupVersion.WithKind("ConfigMap"): {
					syncedConfigMap,
				},
			},
			Sync: func(ctx *synccontext.RegisterContext) {
				syncCtx, syncer := generictesting.FakeStartSyncer(t, ctx, newSyncer)
				_, err := syncer.(*configMapSyncer).SyncDown(syncCtx, baseConfigMap)
				assert.NilError(t, err)
			},
		},
		{
			Name: "Update config map",
			InitialVirtualState: []runtime.Object{
				updatedConfigMap,
			},
			InitialPhysicalState: []runtime.Object{
				syncedConfigMap,
			},
			ExpectedPhysicalState: map[schema.GroupVersionKind][]runtime.Object{
				corev1.SchemeGroupVersion.WithKind("ConfigMap"): {
					updatedSyncedConfigMap,
				},
			},
			Sync: func(ctx *synccontext.RegisterContext) {
				syncCtx, syncer := generictesting.FakeStartSyncer(t, ctx, newSyncer)
				_, err := syncer.(*configMapSyncer).Sync(syncCtx, syncedConfigMap, updatedConfigMap)
				assert.NilError(t, err)
			},
		},
	})
}

func newSyncer(ctx *synccontext.RegisterContext) (syncer.Base, error) {
	return NewConfigMapSyncer(ctx), nil
}
