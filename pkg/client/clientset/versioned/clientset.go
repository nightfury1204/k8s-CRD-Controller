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

package versioned

import (
	glog "github.com/golang/glog"
	podwatchersv1alpha1 "k8s-crd-controller/pkg/client/clientset/versioned/typed/nahid.try.com/v1alpha1"
	discovery "k8s.io/client-go/discovery"
	rest "k8s.io/client-go/rest"
	flowcontrol "k8s.io/client-go/util/flowcontrol"
)

type Interface interface {
	Discovery() discovery.DiscoveryInterface
	PodwatchersV1alpha1() podwatchersv1alpha1.PodwatchersV1alpha1Interface
	// Deprecated: please explicitly pick a version if possible.
	Podwatchers() podwatchersv1alpha1.PodwatchersV1alpha1Interface
}

// Clientset contains the clients for groups. Each group has exactly one
// version included in a Clientset.
type Clientset struct {
	*discovery.DiscoveryClient
	podwatchersV1alpha1 *podwatchersv1alpha1.PodwatchersV1alpha1Client
}

// PodwatchersV1alpha1 retrieves the PodwatchersV1alpha1Client
func (c *Clientset) PodwatchersV1alpha1() podwatchersv1alpha1.PodwatchersV1alpha1Interface {
	return c.podwatchersV1alpha1
}

// Deprecated: Podwatchers retrieves the default version of PodwatchersClient.
// Please explicitly pick a version.
func (c *Clientset) Podwatchers() podwatchersv1alpha1.PodwatchersV1alpha1Interface {
	return c.podwatchersV1alpha1
}

// Discovery retrieves the DiscoveryClient
func (c *Clientset) Discovery() discovery.DiscoveryInterface {
	if c == nil {
		return nil
	}
	return c.DiscoveryClient
}

// NewForConfig creates a new Clientset for the given config.
func NewForConfig(c *rest.Config) (*Clientset, error) {
	configShallowCopy := *c
	if configShallowCopy.RateLimiter == nil && configShallowCopy.QPS > 0 {
		configShallowCopy.RateLimiter = flowcontrol.NewTokenBucketRateLimiter(configShallowCopy.QPS, configShallowCopy.Burst)
	}
	var cs Clientset
	var err error
	cs.podwatchersV1alpha1, err = podwatchersv1alpha1.NewForConfig(&configShallowCopy)
	if err != nil {
		return nil, err
	}

	cs.DiscoveryClient, err = discovery.NewDiscoveryClientForConfig(&configShallowCopy)
	if err != nil {
		glog.Errorf("failed to create the DiscoveryClient: %v", err)
		return nil, err
	}
	return &cs, nil
}

// NewForConfigOrDie creates a new Clientset for the given config and
// panics if there is an error in the config.
func NewForConfigOrDie(c *rest.Config) *Clientset {
	var cs Clientset
	cs.podwatchersV1alpha1 = podwatchersv1alpha1.NewForConfigOrDie(c)

	cs.DiscoveryClient = discovery.NewDiscoveryClientForConfigOrDie(c)
	return &cs
}

// New creates a new Clientset for the given RESTClient.
func New(c rest.Interface) *Clientset {
	var cs Clientset
	cs.podwatchersV1alpha1 = podwatchersv1alpha1.New(c)

	cs.DiscoveryClient = discovery.NewDiscoveryClient(c)
	return &cs
}
