package specs

import (
	"context"

	"github.com/mjudeikis/go-podman/test/e2e/client"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Create Delete test", func() {
	podCreate := kubePod.DeepCopy()
	podCreate.Name = "test-pod-create"
	It("should create pod validate podman pod", func() {
		err := client.PodmanClient.Create(context.Background(), podCreate)
		if err != nil {
			Expect(err).ToNot(HaveOccurred())
		}

		err = client.PodmanClient.Delete(context.Background(), podCreate)
		if err != nil {
			Expect(err).ToNot(HaveOccurred())
		}
	})

	AfterEach(func() {
		client.PodmanClient.Delete(context.Background(), podCreate)
	})
})
