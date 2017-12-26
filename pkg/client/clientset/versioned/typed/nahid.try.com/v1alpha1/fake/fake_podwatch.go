/*
Copyright 2017 The Kubernetes Authors.

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

package fake

import (
	v1alpha1 "k8s-crd-controller/pkg/apis/nahid.try.com/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakePodWatchs implements PodWatchInterface
type FakePodWatchs struct {
	Fake *FakeNahidV1alpha1
	ns   string
}

var podwatchsResource = schema.GroupVersionResource{Group: "nahid.try.com", Version: "v1alpha1", Resource: "podwatchs"}

var podwatchsKind = schema.GroupVersionKind{Group: "nahid.try.com", Version: "v1alpha1", Kind: "PodWatch"}

// Get takes name of the podWatch, and returns the corresponding podWatch object, and an error if there is any.
func (c *FakePodWatchs) Get(name string, options v1.GetOptions) (result *v1alpha1.PodWatch, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(podwatchsResource, c.ns, name), &v1alpha1.PodWatch{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.PodWatch), err
}

// List takes label and field selectors, and returns the list of PodWatchs that match those selectors.
func (c *FakePodWatchs) List(opts v1.ListOptions) (result *v1alpha1.PodWatchList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(podwatchsResource, podwatchsKind, c.ns, opts), &v1alpha1.PodWatchList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.PodWatchList{}
	for _, item := range obj.(*v1alpha1.PodWatchList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested podWatchs.
func (c *FakePodWatchs) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(podwatchsResource, c.ns, opts))

}

// Create takes the representation of a podWatch and creates it.  Returns the server's representation of the podWatch, and an error, if there is any.
func (c *FakePodWatchs) Create(podWatch *v1alpha1.PodWatch) (result *v1alpha1.PodWatch, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(podwatchsResource, c.ns, podWatch), &v1alpha1.PodWatch{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.PodWatch), err
}

// Update takes the representation of a podWatch and updates it. Returns the server's representation of the podWatch, and an error, if there is any.
func (c *FakePodWatchs) Update(podWatch *v1alpha1.PodWatch) (result *v1alpha1.PodWatch, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(podwatchsResource, c.ns, podWatch), &v1alpha1.PodWatch{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.PodWatch), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakePodWatchs) UpdateStatus(podWatch *v1alpha1.PodWatch) (*v1alpha1.PodWatch, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(podwatchsResource, "status", c.ns, podWatch), &v1alpha1.PodWatch{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.PodWatch), err
}

// Delete takes name of the podWatch and deletes it. Returns an error if one occurs.
func (c *FakePodWatchs) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(podwatchsResource, c.ns, name), &v1alpha1.PodWatch{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakePodWatchs) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(podwatchsResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &v1alpha1.PodWatchList{})
	return err
}

// Patch applies the patch and returns the patched podWatch.
func (c *FakePodWatchs) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.PodWatch, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(podwatchsResource, c.ns, name, data, subresources...), &v1alpha1.PodWatch{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.PodWatch), err
}
