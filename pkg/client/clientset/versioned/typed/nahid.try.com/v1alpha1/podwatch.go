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

package v1alpha1

import (
	v1alpha1 "k8s-crd-controller/pkg/apis/nahid.try.com/v1alpha1"
	scheme "k8s-crd-controller/pkg/client/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// PodWatchsGetter has a method to return a PodWatchInterface.
// A group's client should implement this interface.
type PodWatchsGetter interface {
	PodWatchs(namespace string) PodWatchInterface
}

// PodWatchInterface has methods to work with PodWatch resources.
type PodWatchInterface interface {
	Create(*v1alpha1.PodWatch) (*v1alpha1.PodWatch, error)
	Update(*v1alpha1.PodWatch) (*v1alpha1.PodWatch, error)
	UpdateStatus(*v1alpha1.PodWatch) (*v1alpha1.PodWatch, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha1.PodWatch, error)
	List(opts v1.ListOptions) (*v1alpha1.PodWatchList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.PodWatch, err error)
	PodWatchExpansion
}

// podWatchs implements PodWatchInterface
type podWatchs struct {
	client rest.Interface
	ns     string
}

// newPodWatchs returns a PodWatchs
func newPodWatchs(c *NahidV1alpha1Client, namespace string) *podWatchs {
	return &podWatchs{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the podWatch, and returns the corresponding podWatch object, and an error if there is any.
func (c *podWatchs) Get(name string, options v1.GetOptions) (result *v1alpha1.PodWatch, err error) {
	result = &v1alpha1.PodWatch{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("podwatchs").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of PodWatchs that match those selectors.
func (c *podWatchs) List(opts v1.ListOptions) (result *v1alpha1.PodWatchList, err error) {
	result = &v1alpha1.PodWatchList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("podwatchs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested podWatchs.
func (c *podWatchs) Watch(opts v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("podwatchs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a podWatch and creates it.  Returns the server's representation of the podWatch, and an error, if there is any.
func (c *podWatchs) Create(podWatch *v1alpha1.PodWatch) (result *v1alpha1.PodWatch, err error) {
	result = &v1alpha1.PodWatch{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("podwatchs").
		Body(podWatch).
		Do().
		Into(result)
	return
}

// Update takes the representation of a podWatch and updates it. Returns the server's representation of the podWatch, and an error, if there is any.
func (c *podWatchs) Update(podWatch *v1alpha1.PodWatch) (result *v1alpha1.PodWatch, err error) {
	result = &v1alpha1.PodWatch{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("podwatchs").
		Name(podWatch.Name).
		Body(podWatch).
		Do().
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().

func (c *podWatchs) UpdateStatus(podWatch *v1alpha1.PodWatch) (result *v1alpha1.PodWatch, err error) {
	result = &v1alpha1.PodWatch{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("podwatchs").
		Name(podWatch.Name).
		SubResource("status").
		Body(podWatch).
		Do().
		Into(result)
	return
}

// Delete takes name of the podWatch and deletes it. Returns an error if one occurs.
func (c *podWatchs) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("podwatchs").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *podWatchs) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("podwatchs").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched podWatch.
func (c *podWatchs) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.PodWatch, err error) {
	result = &v1alpha1.PodWatch{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("podwatchs").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
