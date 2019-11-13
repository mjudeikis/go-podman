package specs

import (
	"context"

	"github.com/mjudeikis/go-podman/test/e2e/client"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Create Update test", func() {
	podCreate := kubePod.DeepCopy()
	podCreate.Name = "test-pod-update"
	It("should create pod validate podman pod", func() {
		err := client.PodmanClient.Create(context.Background(), podCreate)
		if err != nil {
			Expect(err).ToNot(HaveOccurred())
		}

		pod, err := client.PodmanClient.Get(context.Background(), podCreate)
		if err != nil {
			Expect(err).ToNot(HaveOccurred())
		}

		newPod := pod.DeepCopy()
		newPod.Spec.Containers[0].Name = "new-container"

		err = client.PodmanClient.Update(context.Background(), newPod)
		if err != nil {
			Expect(err).ToNot(HaveOccurred())
		}

		pod, err = client.PodmanClient.Get(context.Background(), podCreate)
		if err != nil {
			Expect(err).ToNot(HaveOccurred())
		}
		if pod.Spec.Containers[0].Name != "new-container" {
			Fail("pod update failed")
		}
	})

	AfterEach(func() {
		client.PodmanClient.Delete(context.Background(), podCreate)
	})
})
