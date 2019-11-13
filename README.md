# go-podman

`go-podman` is minimal abstraction for iopodman varlink package.

It operates using Kuberentes pod Objects as inputs. 

## usage

```go
import (
    "github.com/mjudeikis/go-podman/pkg/podman"
    corev1 "k8s.io/api/core/v1"
)

client, err := podman.New(context.Background(), nil)
	if err != nil {
		panic(err)
    }
    

    err = client.Create(context.Background(), corev1.Pod{})
    err = client.List(context.Background(), corev1.Pod{})
    ...
```

## roadmap

1. Extend Kube pod spec traslation into podmanPod spec better. Currently we just running plain container with no secrets, volumes, etc.
2. Add unit testing
3. Move to k8s errors 
