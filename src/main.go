package main

import (
        "crypto/tls"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	kwhhttp "github.com/slok/kubewebhook/v2/pkg/http"
	kwhlog "github.com/slok/kubewebhook/v2/pkg/log"
	kwhlogrus "github.com/slok/kubewebhook/v2/pkg/log/logrus"
	kwhmodel "github.com/slok/kubewebhook/v2/pkg/model"
	kwhmutating "github.com/slok/kubewebhook/v2/pkg/webhook/mutating"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// carmyContainer is the default carmy container from which to pull the
	// carmy binary.
	carmyContainer = "cloudarmycl/envsecrets"

	// binVolumeName is the name of the volume where the carmy binary is stored.
	binVolumeName = "injector"

	// binVolumeMountPath is the mount path where the carmy binary can be found.
	binVolumeMountPath = "/env-injector/"
)

// binInitContainer is the container that pulls the carmy binary executable
// into a shared volume mount.
var binInitContainer = corev1.Container{
	Name:            "copy-carmy-bin",
	Image:           carmyContainer,
	ImagePullPolicy: corev1.PullIfNotPresent,
	Command: []string{"sh", "-c",
		fmt.Sprintf("cp /injector/envsecrets %s", binVolumeMountPath)},
	VolumeMounts: []corev1.VolumeMount{
		{
			Name:      binVolumeName,
			MountPath: binVolumeMountPath,
		},
	},
}

// binVolume is the shared, in-memory volume where the carmy binary lives.
var binVolume = corev1.Volume{
	Name: binVolumeName,
	VolumeSource: corev1.VolumeSource{
		EmptyDir: &corev1.EmptyDirVolumeSource{
			Medium: corev1.StorageMediumMemory,
		},
	},
}

// binVolumeMount is the shared volume mount where the carmy binary lives.
var binVolumeMount = corev1.VolumeMount{
	Name:      binVolumeName,
	MountPath: binVolumeMountPath,
	ReadOnly:  true,
}

// CarmyMutator is a mutator.
type CarmyMutator struct {
	logger kwhlog.Logger
}

// Mutate implements MutateFunc and provides the top-level entrypoint for object
// mutation.
func (m *CarmyMutator) Mutate(ctx context.Context, ar *kwhmodel.AdmissionReview, obj metav1.Object) (*kwhmutating.MutatorResult, error) {
	m.logger.Infof("calling mutate")

	pod, ok := obj.(*corev1.Pod)
	if !ok {
		return &kwhmutating.MutatorResult{
			Warnings: []string{fmt.Sprintf("incoming resource is not a Pod (%T)", pod)},
		}, nil
	}

	mutated := false
	for i, c := range pod.Spec.InitContainers {
		c, didMutate := m.mutateContainer(ctx, &c)
		if didMutate {
			mutated = true
			pod.Spec.InitContainers[i] = *c
		}
	}

	for i, c := range pod.Spec.Containers {
		c, didMutate := m.mutateContainer(ctx, &c)
		if didMutate {
			mutated = true
			pod.Spec.Containers[i] = *c
		}
	}

	// If any of the containers requested carmy secrets, mount the shared volume
	// and ensure the carmy binary is available via an init container.
	if mutated {
		pod.Spec.Volumes = append(pod.Spec.Volumes, binVolume)
		pod.Spec.InitContainers = append([]corev1.Container{binInitContainer},
			pod.Spec.InitContainers...)
	}

	return &kwhmutating.MutatorResult{
		MutatedObject: pod,
	}, nil
}

// mutateContainer mutates the given container, updating the volume mounts and
// command if it contains carmy references.
func (m *CarmyMutator) mutateContainer(_ context.Context, c *corev1.Container) (*corev1.Container, bool) {

	// Carmy prepends the command from the podspec. If there's no command in the
	// podspec, there's nothing to append. Note: this is the command in the
	// podspec, not a CMD or ENTRYPOINT in a Dockerfile.
	//if len(c.Command) == 0 {
	//	m.logger.Warningf("cannot apply envsecrets to %s: container spec does not define a command", c.Name)
	//	return c, false
	//}
        
	// Add the shared volume mount
	c.VolumeMounts = append(c.VolumeMounts, binVolumeMount)
	// Prepend the command with envsecrets exec --
	//original := append(c.Command)
	///original := append(c.Command)
        /////c.Command = append([]string{binVolumeMountPath + "envsecrets", "./entrypoint.sh"}
        c.Command = []string{binVolumeMountPath + "envsecrets", "./entrypoint.sh"}
        ////c.Args = append([]string{"./entrypoint.sh"})
	//c.Command = []string{binVolumeMountPath + "envsecrets"}
	//c.Args = append([]string{"./"}, original...)
        //c.Args = append([]string{"exec", "--"}, original...)
	return c, true
}

// hasCarmyReferences parses the environment and returns true if any of the
// environment variables includes a carmy reference.
//func (m *CarmyMutator) hasCarmyReferences(env []corev1.EnvVar) bool {
//	for _, e := range env {
//		if carmy.IsReference(e.Value) {
//			return true
//		}
//	}
//	return false
//}

// webhookHandler is the http.Handler that responds to webhooks
func webhookHandler() (http.Handler, error) {
	entry := logrus.NewEntry(logrus.New())
	entry.Logger.SetLevel(logrus.DebugLevel)
	logger := kwhlogrus.NewLogrus(entry)

	mutator := &CarmyMutator{logger: logger}

	mcfg := kwhmutating.WebhookConfig{
		ID:      "carmySecrets",
		Obj:     &corev1.Pod{},
		Mutator: mutator,
		Logger:  logger,
	}

	// Create the wrapping webhook
	wh, err := kwhmutating.NewWebhook(mcfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create mutating webhook: %w", err)
	}

	// Get the handler for our webhook.
	whhandler, err := kwhhttp.HandlerFor(kwhhttp.HandlerConfig{
		Webhook: wh,
		Logger:  logger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create mutating webhook handler: %w", err)
	}
	return whhandler, nil
}

func logAt(lvl, msg string, args ...any) {
	body := map[string]any{
		"time":     time.Now().UTC().Format(time.RFC3339),
		"severity": lvl,
		"message":  fmt.Sprintf(msg, args...),
	}

	payload, err := json.Marshal(body)
	if err != nil {
		panic(fmt.Sprintf("failed to make JSON error message: %s", err))
	}
	fmt.Fprintln(os.Stderr, string(payload))
}

func logInfo(msg string, args ...any) {
	logAt("INFO", msg, args...)
}

func logError(msg string, args ...any) {
	logAt("ERROR", msg, args...)
}

func realMain() error {

	port := "8443"

	handler, err := webhookHandler()

	if err != nil {
		return fmt.Errorf("server failed to start: %w", err)
	}
        
        cert, err := tls.LoadX509KeyPair("/etc/webhook/certs/tls.crt", "/etc/webhook/certs/tls.key")

        if err != nil {
		return fmt.Errorf("server failed to load the certs: %w", err)
        }

        tlsConfig := &tls.Config{
        	Certificates: []tls.Certificate{cert},
        }

        srv := &http.Server{
        	Addr:      ":"+ port,
        	TLSConfig: tlsConfig,
		Handler: handler,
        }

	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.ListenAndServeTLS("", "")

	}()
	logInfo("server is listening on " + port)

	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, syscall.SIGINT, syscall.SIGTERM)
	<-stopCh

	ctx, done := context.WithTimeout(context.Background(), 10*time.Second)
	defer done()

	if err := srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("server failed to shutdown: %w", err)
	}

	// Wait for shutdown
	if err := <-errCh; err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("server error: %w", err)
	}

	return nil
}

func main() {
	if err := realMain(); err != nil {
		logError(err.Error())
		os.Exit(1)
	}

	logInfo("server successfully stopped")
}
