package util

import (
	"encoding/json"
	"fmt"
	clientver "k8s-crd-controller/pkg/client/clientset/versioned"
	"time"

	"github.com/golang/glog"
	nahidtrycomv1alpha1 "k8s-crd-controller/pkg/apis/nahid.try.com/v1alpha1"
	apiv1 "k8s.io/api/core/v1"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/strategicpatch"
	"k8s.io/apimachinery/pkg/util/wait"
)

func CreateOrPatchPodWatch(c clientver.Clientset, meta metav1.ObjectMeta, transform func(watch *nahidtrycomv1alpha1.PodWatch) *nahidtrycomv1alpha1.PodWatch) (result *nahidtrycomv1alpha1.PodWatch, err error) {
	cur, err := c.NahidV1alpha1().PodWatchs(meta.Namespace).Get(meta.Name, metav1.GetOptions{})
	if kerr.IsNotFound(err) {
		glog.V(3).Infof("Creating PodWatch %s/%s.", meta.Namespace, meta.Name)
		out, err := c.NahidV1alpha1().PodWatchs(meta.Namespace).Create(transform(&nahidtrycomv1alpha1.PodWatch{
			TypeMeta: metav1.TypeMeta{
				Kind:       "PodWatch",
				APIVersion: nahidtrycomv1alpha1.SchemeGroupVersion.String(),
			},
			ObjectMeta: meta,
		}))
		return out, err
	} else if err != nil {
		return nil, err
	}
	return PatchPodWatch(c, cur, transform)
}

func PatchPodWatch(c clientver.Clientset, cur *nahidtrycomv1alpha1.PodWatch, transform func(*nahidtrycomv1alpha1.PodWatch) *nahidtrycomv1alpha1.PodWatch) (*nahidtrycomv1alpha1.PodWatch, error) {
	curJson, err := json.Marshal(cur)
	if err != nil {
		return nil, err
	}
	modJson, err := json.Marshal(transform(cur.DeepCopy()))
	if err != nil {
		return nil, err
	}

	patch, err := strategicpatch.CreateTwoWayMergePatch(curJson, modJson, nahidtrycomv1alpha1.PodWatch{})
	if err != nil {
		return nil, err
	}
	if len(patch) == 0 || string(patch) == "{}" {
		return cur, nil
	}
	glog.V(3).Infof("Patching Podwatch %s/%s with %s.", cur.Namespace, cur.Name, string(patch))
	return c.NahidV1alpha1().PodWatchs(apiv1.NamespaceDefault).Patch(cur.Name, types.StrategicMergePatchType, patch)
}

func tryPatchPodWatch(c clientver.Clientset, meta metav1.ObjectMeta, transform func(watch *nahidtrycomv1alpha1.PodWatch) *nahidtrycomv1alpha1.PodWatch) (result *nahidtrycomv1alpha1.PodWatch, err error) {
	attempt := 0
	err = wait.PollImmediate(1*time.Second, 5*time.Second, func() (bool, error) {
		attempt++
		cur, e2 := c.NahidV1alpha1().PodWatchs(apiv1.NamespaceDefault).Get(meta.Name, metav1.GetOptions{})
		if kerr.IsNotFound(e2) {
			return false, e2
		} else if e2 == nil {
			result, e2 = PatchPodWatch(c, cur, transform)
			return e2 == nil, nil
		}
		glog.Errorf("Attempt %d failed to patch Service %s/%s due to %v.", attempt, cur.Namespace, cur.Name, e2)
		return false, nil
	})

	if err != nil {
		err = fmt.Errorf("failed to patch Service %s/%s after %d attempts due to %v", meta.Namespace, meta.Name, attempt, err)
	}
	return
}
