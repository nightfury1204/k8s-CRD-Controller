package controller

import (
	"fmt"
	"time"

	"github.com/golang/glog"

	clientver "k8s-crd-controller/pkg/client/clientset/versioned"
	myutil "k8s-crd-controller/pkg/util"

	nahidtrycomv1alpha1 "k8s-crd-controller/pkg/apis/nahid.try.com/v1alpha1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime2 "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/client-go/kubernetes"
	"math/rand"
	"strconv"
	"strings"
)

const (
	inf int = 100000000
	AllPods int = inf
	Processing string = "processing"
	Ready string = "ready"
)

type Controller struct {
	kubeClientset    kubernetes.Clientset
	podIndexer  cache.Indexer
	podInformer cache.Controller
	podQueue    workqueue.RateLimitingInterface

	//for PodWatch
	clientset        clientver.Clientset
	podWatchIndexer  cache.Indexer
	podWatchInformer cache.Controller
	podWatchQueue    workqueue.RateLimitingInterface

	//to find pod assigned to which podwatch
	PodMapToPodWatch map[string]string
	//Current status of pod
	PodStatus map[string]string
}

func customListWatcherForPodWatch(clientset clientver.Clientset) *cache.ListWatch {
	return &cache.ListWatch{
		ListFunc: func(opts metav1.ListOptions) (runtime2.Object, error) {
			return clientset.NahidV1alpha1().PodWatchs(apiv1.NamespaceDefault).List(metav1.ListOptions{IncludeUninitialized:true})
		},
		WatchFunc: func(opts metav1.ListOptions) (watch.Interface, error) {
			return clientset.NahidV1alpha1().PodWatchs(apiv1.NamespaceDefault).Watch(metav1.ListOptions{IncludeUninitialized:true})
		},
	}
}

func customListWatcherForPod(clientset kubernetes.Clientset) *cache.ListWatch {
	return &cache.ListWatch{
		ListFunc: func(opts metav1.ListOptions) (runtime2.Object, error) {
			return clientset.CoreV1().Pods(apiv1.NamespaceDefault).List(metav1.ListOptions{})
		},
		WatchFunc: func(opts metav1.ListOptions) (watch.Interface, error) {
			return clientset.CoreV1().Pods(apiv1.NamespaceDefault).Watch(metav1.ListOptions{})
		},
	}
}

func NewController(clientset clientver.Clientset, kubeClientset kubernetes.Clientset) *Controller {
	//for podWatch
	podWatchQueue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())

	//create watcher for podWatch
	podWatchLW := customListWatcherForPodWatch(clientset)

	//for podWatch
	podWatchIndexer, podWatchInformer := cache.NewIndexerInformer(podWatchLW, &nahidtrycomv1alpha1.PodWatch{}, time.Second*30, cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(obj)
			if err == nil {
				podWatchQueue.Add(key)
			}
		},
		UpdateFunc: func(old interface{}, new interface{}) {
			oldPW := old.(*nahidtrycomv1alpha1.PodWatch)
			newPW := new.(*nahidtrycomv1alpha1.PodWatch)

			if oldPW != newPW {
				key, err := cache.MetaNamespaceKeyFunc(new)
				if err == nil {
					podWatchQueue.Add(key)
				}
			}
		},
		DeleteFunc: func(obj interface{}) {
			key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
			if err == nil {
				podWatchQueue.Add(key)
			}
		},
	}, cache.Indexers{})

	//for pod
	podQueue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())

	//create watcher for pod
	podLW := customListWatcherForPod(kubeClientset)

	//for pod
	podIndexer, podInformer := cache.NewIndexerInformer(podLW, &apiv1.Pod{}, 0, cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(obj)
			if err == nil {
				podQueue.Add(key)
			}
		},
		UpdateFunc: func(old interface{}, new interface{}) {
			oldPW := old.(*apiv1.Pod)
			newPW := new.(*apiv1.Pod)

			if oldPW != newPW {
				key, err := cache.MetaNamespaceKeyFunc(new)
				if err == nil {
					podQueue.Add(key)
				}
			}
		},
		DeleteFunc: func(obj interface{}) {
			key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
			if err == nil {
				podQueue.Add(key)
			}
		},
	}, cache.Indexers{})


	return &Controller{
		kubeClientset:    kubeClientset,
		podInformer: podInformer,
		podIndexer:  podIndexer,
		podQueue:    podQueue,

		clientset:        clientset,
		podWatchInformer: podWatchInformer,
		podWatchIndexer:  podWatchIndexer,
		podWatchQueue:    podWatchQueue,
		PodMapToPodWatch: map[string]string{},
		PodStatus: map[string]string{},
	}
}

func (c *Controller) processNextItemForPodWatch() bool {
	// Wait until there is a new item in the working podWatchQueue
	key, quit := c.podWatchQueue.Get()
	if quit {
		return false
	}
	// Tell the podWatchQueue that we are done with processing this key. This unblocks the key for other workers
	// This allows safe parallel processing because two pods with the same key are never processed in
	// parallel.
	defer c.podWatchQueue.Done(key)

	// Invoke the method containing the business logic
	err := c.performOpAccrodingToPodWatch(key.(string))
	// Handle the error if something went wrong during the execution of the business logic
	c.handleErrForPodWatch(err, key)
	return true
}

func (c *Controller) processNextItemForPod() bool {
	// Wait until there is a new item in the working podQueue
	key, quit := c.podQueue.Get()
	if quit {
		return false
	}
	// Tell the podQueue that we are done with processing this key. This unblocks the key for other workers
	// This allows safe parallel processing because two pods with the same key are never processed in
	// parallel.
	defer c.podQueue.Done(key)

	// Invoke the method containing the business logic
	err := c.performOpAccrodingToPod(key.(string))
	// Handle the error if something went wrong during the execution of the business logic
	c.handleErrForPod(err, key)
	return true
}

// performOpAccrodingToPodWatch is the business logic of the controller. In this controller it simply prints
// information about the pod to stdout. In case an error happened, it has to simply return the error.
// The retry logic should not be part of the business logic.
func (c *Controller) performOpAccrodingToPodWatch(key string) error {
	obj, exists, err := c.podWatchIndexer.GetByKey(key)
	if err != nil {
		glog.Errorf("Fetching object with key %s from store failed with %v", key, err)
		return err
	}

	if !exists {
		// Below we will warm up our cache with a Pod, so that we will see a delete for one pod
		fmt.Printf("PodWatcher %s deleted\n", key)
	} else {
		// Note that you also have to check the uid if you have a local controlled resource, which
		// is dependent on the actual instance, to detect that a Pod was recreated with the same name
		podWatch := obj.(*nahidtrycomv1alpha1.PodWatch)

		if podWatch.DeletionTimestamp != nil {
			// finalizer
			if myutil.HasFinalizer(podWatch.ObjectMeta,"finalizer.podwatch.nahid.try.com") {
				//delete pods for podwatch
				err := c.deletePod(podWatch.Spec.LabelSelector.MatchLabels, AllPods)
				if err!=nil {
					fmt.Println("failed to delete pod")
					return err
				}
				_,err = myutil.PatchPodWatch(c.clientset,podWatch, func(podW *nahidtrycomv1alpha1.PodWatch) *nahidtrycomv1alpha1.PodWatch {
					podW.ObjectMeta = myutil.RemoveFinalizer(podW.ObjectMeta,"finalizer.podwatch.nahid.try.com")
					return podW
				})
				if err!=nil {
					return err
				}
			}

		} else if podWatch.GetInitializers() != nil {
			//initializer
			pendingInitializers := podWatch.GetInitializers().Pending
			if pendingInitializers[0].Name == "podwatchinit.nahid.try.com" {
				//add finalizer
				_,err := myutil.PatchPodWatch(c.clientset,podWatch, func(podW *nahidtrycomv1alpha1.PodWatch) *nahidtrycomv1alpha1.PodWatch {
					podW.ObjectMeta = myutil.AddFinalizer(podW.ObjectMeta, "finalizer.podwatch.nahid.try.com")
					podW.ObjectMeta = myutil.RemoveInitializer(podW.ObjectMeta)
					return podW
				})
				if err!=nil {
					return err
				}
			}
		} else {
			fmt.Printf("PodWatcher %s( avialable | required | current): %v | %v | %v\n", podWatch.GetName(),podWatch.Status.AvailabelReplicas,podWatch.Spec.Replicas,podWatch.Status.CurrentlyProcessing)

			if podWatch.Status.AvailabelReplicas+podWatch.Status.CurrentlyProcessing == podWatch.Spec.Replicas {
				//ignore
			} else if podWatch.Status.AvailabelReplicas+podWatch.Status.CurrentlyProcessing < podWatch.Spec.Replicas {

				err := c.createPod(podWatch.Spec.Template, podWatch.GetName())
				if err!=nil {
					fmt.Println("Failed to create pod")
					return err
				}

				err = c.updateStatusForPodWatch(podWatch, podWatch.Status.AvailabelReplicas, podWatch.Status.CurrentlyProcessing+1)
				if err!=nil {
					fmt.Println("Failed to update podWatch")
					return err
				}
			} else {
				//delete pod
				//currently not required
				err := c.deletePod(podWatch.Spec.LabelSelector.MatchLabels, int(podWatch.Status.AvailabelReplicas+podWatch.Status.CurrentlyProcessing-podWatch.Spec.Replicas))
				if err!=nil {
					fmt.Println("failed to delete pod")
					return err
				}
			}
		}

	}
	return nil
}

// performOpAccrodingToPod is the business logic of the controller. In this controller it simply prints
// information about the pod to stdout. In case an error happened, it has to simply return the error.
// The retry logic should not be part of the business logic.
func (c *Controller) performOpAccrodingToPod(key string) error {
	obj, exists, err := c.podIndexer.GetByKey(key)
	if err != nil {
		glog.Errorf("Fetching object with key %s from store failed with %v", key, err)
		return err
	}

	if !exists {
		// Below we will warm up our cache with a Pod, so that we will see a delete for one pod
		fmt.Printf("Pod %s deleted\n", key)
		values := strings.Split(key,"/")
		podName := values[len(values)-1]
		podWatchName,ok := c.PodMapToPodWatch[podName]
		if ok {
			podWatch,err := c.clientset.NahidV1alpha1().PodWatchs(apiv1.NamespaceDefault).Get(podWatchName,metav1.GetOptions{})
			if err!=nil {
				return err
			}
			err = c.updateStatusForPodWatch(podWatch, podWatch.Status.AvailabelReplicas-1,podWatch.Status.CurrentlyProcessing)
			return err
		}
		delete(c.PodMapToPodWatch,podName)
		delete(c.PodStatus,podName)
	} else {
		// Note that you also have to check the uid if you have a local controlled resource, which
		// is dependent on the actual instance, to detect that a Pod was recreated with the same name
		pod := obj.(*apiv1.Pod)

		fmt.Printf("Sync/Add/Update for Pod %s\n", pod.GetName())
		status,ok := c.PodStatus[pod.GetName()]
		//pod is not assigned to podwatch
		if !ok{
			return nil
		}

		if (pod.Status.Phase== apiv1.PodRunning || pod.Status.Phase== apiv1.PodSucceeded) && status==Processing {
			podWatchName,ok := c.PodMapToPodWatch[pod.GetName()]
			if ok {
				podWatch,err := c.clientset.NahidV1alpha1().PodWatchs(apiv1.NamespaceDefault).Get(podWatchName,metav1.GetOptions{})
				if err!=nil {
					return err
				}
				err = c.updateStatusForPodWatch(podWatch, podWatch.Status.AvailabelReplicas+1,podWatch.Status.CurrentlyProcessing-1)
				if err != nil {
					return err
				}
				c.PodStatus[pod.GetName()]=Ready
			}
		}

	}
	return nil
}

// handleErrForPodWatch checks if an error happened and makes sure we will retry later.
func (c *Controller) handleErrForPodWatch(err error, key interface{}) {
	if err == nil {
		// Forget about the #AddRateLimited history of the key on every successful synchronization.
		// This ensures that future processing of updates for this key is not delayed because of
		// an outdated error history.
		c.podWatchQueue.Forget(key)
		return
	}

	// This controller retries 5 times if something goes wrong. After that, it stops trying.
	if c.podWatchQueue.NumRequeues(key) < 5 {
		glog.Infof("Error syncing podwatch %v: %v", key, err)

		// Re-enqueue the key rate limited. Based on the rate limiter on the
		// podWatchQueue and the re-enqueue history, the key will be processed later again.
		c.podWatchQueue.AddRateLimited(key)
		return
	}

	c.podWatchQueue.Forget(key)
	// Report to an external entity that, even after several retries, we could not successfully process this key
	runtime.HandleError(err)
	glog.Infof("Dropping podwatch %q out of the podWatchQueue: %v", key, err)
}

// handleErrForPod checks if an error happened and makes sure we will retry later.
func (c *Controller) handleErrForPod(err error, key interface{}) {
	if err == nil {
		// Forget about the #AddRateLimited history of the key on every successful synchronization.
		// This ensures that future processing of updates for this key is not delayed because of
		// an outdated error history.
		c.podQueue.Forget(key)
		return
	}

	// This controller retries 5 times if something goes wrong. After that, it stops trying.
	if c.podQueue.NumRequeues(key) < 5 {
		glog.Infof("Error syncing pod %v: %v", key, err)

		// Re-enqueue the key rate limited. Based on the rate limiter on the
		// podWatchQueue and the re-enqueue history, the key will be processed later again.
		c.podQueue.AddRateLimited(key)
		return
	}

	c.podQueue.Forget(key)
	// Report to an external entity that, even after several retries, we could not successfully process this key
	runtime.HandleError(err)
	glog.Infof("Dropping pod %q out of the podQueue: %v", key, err)
}

func (c *Controller) Run(threadiness int, stopCh chan struct{}) {
	defer runtime.HandleCrash()

	// Let the workers stop when we are done
	defer c.podWatchQueue.ShutDown()
	defer c.podQueue.ShutDown()
	fmt.Println("Starting PodWatch controller")

	go c.podWatchInformer.Run(stopCh)
	go c.podInformer.Run(stopCh)

	// Wait for all involved caches to be synced, before processing items from the podWatchQueue is started
	if !cache.WaitForCacheSync(stopCh, c.podWatchInformer.HasSynced, c.podInformer.HasSynced) {
		runtime.HandleError(fmt.Errorf("Timed out waiting for caches to sync"))
		return
	}

	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runWorkerForPodWatch, time.Second, stopCh)
		go wait.Until(c.runWorkerForPod, time.Second, stopCh)
	}

	<-stopCh
	fmt.Println("Stopping PodWatch controller")
}

func (c *Controller) runWorkerForPodWatch() {
	for c.processNextItemForPodWatch() {
	}
}

func (c *Controller) runWorkerForPod() {
	for c.processNextItemForPod() {
	}
}

//create pod according to podTemplate
func (c *Controller) createPod(podTemplate nahidtrycomv1alpha1.PodTemplate, podWatchName string) error {
	podClient := c.kubeClientset.CoreV1().Pods(apiv1.NamespaceDefault)
	pod := &apiv1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: podWatchName+"-"+genRandomName(),
			Labels:	podTemplate.GetObjectMeta().GetLabels(),
		},
		Spec: podTemplate.Spec,
	}

	result, err := podClient.Create(pod)
	if err==nil {
		c.PodMapToPodWatch[result.GetName()]=podWatchName
		c.PodStatus[result.GetName()]=Processing
	}
	return err
}

//delete pods whose lebel matches selector
func (c *Controller) deletePod(selector map[string]string ,podDeleteLimit int) error {
	podClient := c.kubeClientset.CoreV1().Pods(apiv1.NamespaceDefault)
	for key,val := range selector {
		podList, err := c.kubeClientset.CoreV1().Pods(apiv1.NamespaceDefault).List(metav1.ListOptions{LabelSelector:key+"="+val})
		if err!=nil {
			return err
		}
		cnt :=int(0)
		if podDeleteLimit==AllPods {
			podDeleteLimit = len(podList.Items)
		}
		for _,pod := range podList.Items {
			if cnt < podDeleteLimit {
				//clear podMapToPodwatch,podStatus
				delete(c.PodMapToPodWatch,pod.GetName())
				delete(c.PodStatus,pod.GetName())
				err = podClient.Delete(pod.GetName(),&metav1.DeleteOptions{})
				if err!=nil {
					return err
				}
				cnt++
			}
		}
	}
	return nil
}

//update status of currenlty processing pod and currently avialble pods in PodWatch
func (c *Controller) updateStatusForPodWatch(curPodWatch *nahidtrycomv1alpha1.PodWatch, numOfReplicas int32, currentlyP int32)  error {
	_, err := myutil.PatchPodWatch(c.clientset, curPodWatch, func(podWatch *nahidtrycomv1alpha1.PodWatch) *nahidtrycomv1alpha1.PodWatch {
		podWatch.Status.AvailabelReplicas = numOfReplicas
		podWatch.Status.CurrentlyProcessing = currentlyP
		return podWatch
	})

	return err
}

//generate a random integer
func genRandomName() string {
	return strconv.Itoa(rand.Int())
}