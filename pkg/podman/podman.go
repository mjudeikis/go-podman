package podman

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/varlink/go/varlink"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"

	"github.com/mjudeikis/go-podman/pkg/converter"
	"github.com/mjudeikis/go-podman/pkg/iopodman"
)

var (
	// Provider configuration defaults.
	defaultSocket = "unix:/run/podman/io.podman"
	defaultSleep  = time.Millisecond * 100
)

// Config defines podman configurables
type Config struct {
	Socket *string
	Log    *zap.SugaredLogger
}

type podman struct {
	c   *varlink.Connection
	log *zap.SugaredLogger
}

// Podman is an simplified interface to interfact with
// podman varlink api
type Podman interface {
	Create(ctx context.Context, pod *corev1.Pod) error
	CreateOrUpdate(ctx context.Context, pod *corev1.Pod) error
	Delete(ctx context.Context, pod *corev1.Pod) error
	Update(ctx context.Context, pod *corev1.Pod) error
	Get(ctx context.Context, pod *corev1.Pod) (*corev1.Pod, error)
	GetByName(ctx context.Context, name string) (*corev1.Pod, error)
	List(ctx context.Context) (*corev1.PodList, error)
}

// New created new instance of podman interface
func New(ctx context.Context, c *Config) (Podman, error) {
	podman := podman{}
	cfg := getConfig(c)
	var err error
	podman.c, err = varlink.NewConnection(ctx, *cfg.Socket)
	if err != nil {
		return nil, err
	}
	podman.log = cfg.Log

	return podman, nil
}

func getConfig(c *Config) *Config {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	log := logger.Sugar()

	if c != nil {
		if c.Socket == nil {
			c.Socket = &defaultSocket
		}
		if c.Log == nil {
			c.Log = log
		}
		return c
	}

	return &Config{
		Socket: &defaultSocket,
		Log:    log,
	}
}

// Create creates podman pod and containers within
func (p podman) Create(ctx context.Context, pod *corev1.Pod) error {
	if pod == nil {
		return fmt.Errorf("create pod can't be nil")
	}

	key := converter.BuildKey(pod)
	podmanPod, err := converter.GetPodmanPod(key, pod)
	if err != nil {
		p.log.Error("getPodmanPod failed", "err", err.Error())
		return err
	}

	podmanPodName, err := iopodman.CreatePod().Call(ctx, p.c, *podmanPod)
	if err != nil {
		p.log.Error("create pod failed", "err", err.Error())
		return err
	}
	p.log.Info("pod created ", "podName ", podmanPodName)

	time.Sleep(defaultSleep)

	// add containers in the pod
	for _, c := range pod.Spec.Containers {
		p.log.Info("create container ", "pod ", podmanPodName, " container ", c.Name)
		container := converter.KubeSpecToPodmanContainer(c, podmanPodName)
		_, err := iopodman.CreateContainer().Call(ctx, p.c, container)
		if err != nil {
			p.log.Error("error createContainer", "err", err.Error())
			return err
		}
	}

	time.Sleep(defaultSleep)

	// start pod
	_, err = iopodman.StartPod().Call(ctx, p.c, podmanPodName)
	if err != nil {
		p.log.Error("error startPod", "err", err.Error())
		return err
	}

	// Podman uses bolt db. Which is laggy
	// https://github.com/containers/libpod/issues/4005
	time.Sleep(defaultSleep)

	// check pod status
	retry := 1
	for retry < 5 {
		retry++
		podmanPodStatus, err := iopodman.InspectPod().Call(ctx, p.c, podmanPodName)
		if err != nil {
			p.log.Error("error GetPod.InspectPod ", "err ", err.Error())
			return err
		}
		var status PodmanPod
		err = json.Unmarshal([]byte(podmanPodStatus), &status)
		if err != nil {
			return err
		}
		healty := true
		for _, c := range status.Containers {
			if c.State != "running" {
				healty = false
			}
		}
		if healty {
			continue
		}
	}

	return nil
}

func (p podman) CreateOrUpdate(ctx context.Context, pod *corev1.Pod) error {
	if pod == nil {
		return fmt.Errorf("create pod can't be nil")
	}

	// for logging only
	key := converter.BuildKey(pod)

	pp, err := p.Get(ctx, pod)
	if err != nil {
		if _, ok := err.(*iopodman.PodNotFound); ok {
			p.log.Debugf("pod not found, creating", " pod ", key)
			return p.Create(ctx, pod)
		}
		if pp != nil && err == nil {
			p.log.Debugf("pod exist, update", " pod ", key)
			return p.Update(ctx, pod)
		}
	}

	return nil
}

func (p podman) Delete(ctx context.Context, pod *corev1.Pod) error {
	if pod == nil {
		p.log.Error("pod can't be nil")
		return fmt.Errorf("pod can't be nil")
	}

	key := converter.BuildKey(pod)

	_, err := iopodman.RemovePod().Call(ctx, p.c, key, true)
	if err != nil {
		p.log.Error("error while deleting pod", " pod ", key, " err ", err.Error())
		return err
	}

	time.Sleep(defaultSleep)

	return nil
}

func (p podman) Update(ctx context.Context, pod *corev1.Pod) error {
	err := p.Delete(ctx, pod)
	if err != nil {
		return err
	}
	return p.Create(ctx, pod)
}

func (p podman) Get(ctx context.Context, input *corev1.Pod) (pod *v1.Pod, err error) {
	key := converter.BuildKey(input)
	return p.GetByName(ctx, key)
}

func (p podman) GetByName(ctx context.Context, name string) (pod *v1.Pod, err error) {
	_, err = iopodman.GetPod().Call(ctx, p.c, name)
	if err != nil {
		return nil, err
	}

	ppod, err := iopodman.InspectPod().Call(ctx, p.c, name)
	if err != nil {
		return nil, err
	}

	kpod, err := converter.GetKubePod(ppod)
	if err != nil {
		return nil, err
	}

	return kpod, nil
}

func (p podman) List(ctx context.Context) (podList *corev1.PodList, err error) {
	ppods, err := iopodman.ListPods().Call(ctx, p.c)
	if err != nil {
		return nil, err
	}

	kpodsList := &corev1.PodList{}
	for _, podData := range ppods {
		kpod, err := p.GetByName(ctx, podData.Name)
		if err != nil {
			return nil, err
		}
		kpodsList.Items = append(kpodsList.Items, *kpod)
	}

	return kpodsList, nil
}
