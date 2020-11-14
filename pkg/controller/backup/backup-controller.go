package backup

//import (
//	"fmt"
//	kubeinformers "k8s.io/client-go/informers"
//	"k8s.io/client-go/kubernetes"
//	"k8s.io/kubernetes/pkg/controller"
//	"time"
//
//	corev1 "k8s.io/api/core/v1"
//	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
//	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
//	"k8s.io/apimachinery/pkg/util/wait"
//	"k8s.io/client-go/tools/cache"
//	"k8s.io/client-go/tools/record"
//	"k8s.io/client-go/util/workqueue"
//	"k8s.io/klog"
//
//	mydumperv1alpha1 "github.com/backup-operator/pkg/apis/mydumper/v1alpha1"
//	"github.com/backup-operator/pkg/client/clientset/versioned"
//	clientset "github.com/backup-operator/pkg/client/clientset/versioned"
//	informers "github.com/backup-operator/pkg/client/informers/externalversions"
//	mydumperInformer "github.com/backup-operator/pkg/client/informers/externalversions/mydumper/v1alpha1"
//	mydumperlister "github.com/backup-operator/pkg/client/listers/mydumper/v1alpha1"
//	"k8s.io/client-go/kubernetes/scheme"
//	eventv1 "k8s.io/client-go/kubernetes/typed/core/v1"
//	"sigs.k8s.io/controller-runtime/pkg/client"
//)
//
//type BackupController struct {
//	kubeclient  kubernetes.Interface
//	mydumperClient  	clientset.Interface
//
//	mydumperLister        mydumperlister.mydumperLister
//	mydumperInformer     mydumperInformer.mydumperInformer
//	mydumperSynced       cache.InformerSynced
//
//	workqueue workqueue.RateLimitingInterface
//
//	recorder record.EventRecorder
//	cli versioned.Interface
//
//	syncHandler   func(jobKey string) (bool, error)
//
//	expectations controller.ControllerExpectationsInterface
//}
//
//func NewBackupController(
//	kubeclientset kubernetes.Interface,
//	cli versioned.Interface,
//	genericCli client.Client,
//	informerFactory informers.SharedInformerFactory,
//	kubeInformerFactory kubeinformers.SharedInformerFactory) *MydumperController {
//
//	utilruntime.Must(backupscheme.AddToScheme(scheme.Scheme))
//	klog.Info("Creating event broadcaster")
//	eventBroadcaster := record.NewBroadcasterWithCorrelatorOptions(record.CorrelatorOptions{QPS: 1})
//	eventBroadcaster.StartLogging(klog.V(2).Infof)
//	eventBroadcaster.StartRecordingToSink(&eventv1.EventSinkImpl{
//		Interface: eventv1.New(kubeclientset.CoreV1().RESTClient()).Events("")})
//	recorder := eventBroadcaster.NewRecorder(mydumperv1alpha1.Scheme, corev1.EventSource{Component: "backup-controller"})
//
//	nbInformer := informerFactory.Backup().V1alpha1().Backups()
//
//	controller := &BackupController{
//		kubeclient:     				kubeclientset,
//		mydumperInformer:   			nbInformer,
//		mydumperLister:        			nbInformer.Lister(),
//		mydumperSynced:        			nbInformer.Informer().HasSynced,
//		workqueue:         				workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "Backups"),
//		recorder:          				recorder,
//	}
//
//	controller.syncHandler = controller.syncHandler
//
//	klog.Info("Setting up event handlers")
//
//	nbInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
//		AddFunc: controller.createBackupCluster,
//		UpdateFunc: controller.updateBackupCluster,
//		DeleteFunc: controller.deleteBackupCluster,
//	})
//
//	//podInformer := kubeInformerFactory.Core().V1().Pods()
//	//podInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
//	//	AddFunc: controller.createBackupCluster,
//	//	UpdateFunc: controller.updateBackupCluster,
//	//	DeleteFunc: controller.deleteBackupCluster,
//	//})
//
//	return controller
//}
//
//func (c *BackupController) createBackupCluster(obj interface{}) {
//	backup, ok := obj.(*mydumperv1alpha1.Backup)
//
//	if !ok {
//		return
//	}
//
//	return
//}
//
//func (c *BackupController) deleteBackup(obj interface{}) {
//
//	backup := obj.(*mydumperv1alpha1.Backup)
//
//}
//
//
//func (c *BackupController) enqueueBackup(obj interface{}) {
//	var key string
//	var err error
//	if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
//		utilruntime.HandleError(err)
//		return
//	}
//	c.workqueue.Add(key)
//}
//
//func (c *BackupController) handleObject(obj interface{}) {
//	var object metav1.Object
//	var ok bool
//	if object, ok = obj.(metav1.Object); !ok {
//		tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
//		if !ok {
//			utilruntime.HandleError(fmt.Errorf("error decoding object, invalid type"))
//			return
//		}
//		object, ok = tombstone.Obj.(metav1.Object)
//		if !ok {
//			utilruntime.HandleError(fmt.Errorf("error decoding object tombstone, invalid type"))
//			return
//		}
//		klog.Info("Recovered deleted object '%s' from tombstone", object.GetName())
//	}
//	klog.Info("Processing object: %s", object.GetName())
//	if ownerRef := metav1.GetControllerOf(object); ownerRef != nil {
//		// If this object is not owned by a Foo, we should not do anything more
//		// with it.
//		if ownerRef.Kind != "backup" {
//			return
//		}
//
//		backup, err := c.mydumperLister.Backup(object.GetNamespace()).Get(ownerRef.Name)
//		if err != nil {
//			klog.Info("ignoring orphaned object '%s' of Backup '%s'", object.GetSelfLink(), ownerRef.Name)
//			return
//		}
//
//		c.enqueueBackup(backup)
//		return
//	}
//}
//
//func (c *BackupController) Run(workers int, stopCh <-chan struct{}) {
//	defer utilruntime.HandleCrash()
//	defer c.workqueue.ShutDown()
//
//	klog.Info("Starting Backuo controller")
//	defer klog.Info("Shutting down Backup controller")
//
//	if ok := cache.WaitForCacheSync(stopCh, c.backupSynced); !ok {
//		klog.Error("Wait For Cache Sync Failed")
//		return
//	}
//
//	klog.Info("Starting workers")
//	for i := 0; i < workers; i++ {
//		go wait.Until(c.worker, time.Second, stopCh)
//	}
//
//	<-stopCh
//}
//
//
//func (c *BackupController) worker() {
//	for c.processNextWorkItem() {
//	}
//}
//
//func (c *BackupController) processNextWorkItem() bool {
//	klog.Info("processNextWorkItem")
//	key, quit := c.workqueue.Get()
//	if quit {
//		return false
//	}
//	defer c.workqueue.Done(key)
//
//	forget, err := c.syncHandler(key.(string))
//	if err == nil {
//		if forget {
//			c.workqueue.Forget(key)
//		}
//		return true
//	}
//
//	utilruntime.HandleError(fmt.Errorf("Error syncing job: %v", err))
//	c.workqueue.AddRateLimited(key)
//
//	return true
//}
//
//func (c *BackupController) syncBackup(key string) (bool, error) {
//	klog.Info("Sync Backup: %v", key)
//	return true, nil
//}
//
//func (c *BackupController) updateBackup(old, new interface{}) {
//	c.enqueueBackup(new)
//}