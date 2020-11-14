package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/backup-operator/pkg/client/clientset/versioned"
	informers "github.com/backup-operator/pkg/client/informers/externalversions"
	"github.com/backup-operator/pkg/controller/backup"
	"github.com/backup-operator/pkg/scheme"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/uuid"
	"k8s.io/apimachinery/pkg/util/wait"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	typedv1 "k8s.io/client-go/kubernetes/typed/core/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

var (
	leaseDuration = 15 * time.Second
	renewDuration = 5 * time.Second
	retryPeriod   = 3 * time.Second
	waitDuration  = 5 * time.Second
)


func main() {
	cfg := config.GetConfigOrDie()

	cli, err := versioned.NewForConfig(cfg)
	if err != nil {
		klog.Fatalf("failed to create versioned client: %v", err)
	}

	kubeCli, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		klog.Fatalf("failed to create kubernetes client: %v", err)
	}

	genericCli, err := client.New(cfg, client.Options{Scheme: scheme.Scheme})
	if err != nil {
		klog.Fatalf("failed to create generic client: %v", err)
	}

	ns := os.Getenv("NAMESPACE")
	if len(ns) == 0 {
		ns = corev1.NamespaceDefault
	}

	hostname, err := os.Hostname()
	if err != nil {
		klog.Fatalf("unable to get hostname: %v", err)
	}

	id := hostname + "_" + string(uuid.NewUUID())
	recorder := createRecorder(kubeCli, "nebula-controller-manager")

	rl := &resourcelock.EndpointsLock{
		EndpointsMeta: metav1.ObjectMeta{
			Namespace: ns,
			Name:      "nebula-controller-manager",
		},
		Client: kubeCli.CoreV1(),
		LockConfig: resourcelock.ResourceLockConfig{
			Identity:      id,
			EventRecorder: recorder,
		},
	}

	controllerCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// listen for interrupts or the Linux SIGTERM signal and cancel
	// our context, which the leader election code will observe and
	// step down
	ch := make(chan os.Signal, 2)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-ch
		klog.Info("Received termination, signaling shutdown")
		cancel()
		<-ch
		os.Exit(0)
	}()

	informerFactory := informers.NewSharedInformerFactory(cli, time.Second*5)
	kubeInformerFactory := kubeinformers.NewSharedInformerFactory(kubeCli, time.Second*5)

	onStarted := func(ctx context.Context) {
		backup := backup.NewBackupController(
			kubeCli,
			cli,
			genericCli,
			informerFactory,
			kubeInformerFactory,
			false,
			2*time.Minute,
			2*time.Minute,
			2*time.Minute,
		)


		go backup.Run(4, ctx.Done())

		// Start method is non-blocking and runs all registered informers in a dedicated goroutine.
		kubeInformerFactory.Start(ctx.Done())
		informerFactory.Start(ctx.Done())
		<-ctx.Done()
	}
	onStopped := func() {
		klog.Infof("leader election lost")
		os.Exit(0)
	}

	wait.Forever(func() {
		leaderelection.RunOrDie(controllerCtx, leaderelection.LeaderElectionConfig{
			Lock:            rl,
			ReleaseOnCancel: true,
			LeaseDuration:   leaseDuration,
			RenewDeadline:   renewDuration,
			RetryPeriod:     retryPeriod,
			Name:            "nebula-controller-manager",
			Callbacks: leaderelection.LeaderCallbacks{
				OnStartedLeading: onStarted,
				OnStoppedLeading: onStopped,
			},
		})
	}, waitDuration)
}

// createRecorder creates event recorder.
func createRecorder(kubeClient kubernetes.Interface, userAgent string) record.EventRecorder {
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(klog.Infof)
	eventBroadcaster.StartRecordingToSink(&typedv1.EventSinkImpl{Interface: kubeClient.CoreV1().Events("")})
	return eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: userAgent})
}
