package client

import (
	"context"

	"github.com/mjudeikis/go-podman/pkg/podman"
)

var PodmanClient podman.Podman

func init() {
	var err error
	PodmanClient, err = podman.New(context.Background(), nil)
	if err != nil {
		panic(err)
	}
}
