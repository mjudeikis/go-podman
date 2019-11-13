package specs

import (
	"context"

	"github.com/mjudeikis/go-podman/pkg/util"
	"github.com/mjudeikis/go-podman/test/e2e/client"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("List Get test", func() {
	kubepod1 := kubePod.DeepCopy()
	kubepod1.Name = "test-pod-list-1"
	kubepod2 := kubePod.DeepCopy()
	kubepod2.Name = "test-pod-list-2"

	It("should list and get podman pods", func() {
		err := client.PodmanClient.Create(context.Background(), kubepod1)
		if err != nil {
			Expect(err).ToNot(HaveOccurred())
		}

		err = client.PodmanClient.Create(context.Background(), kubepod2)
		if err != nil {
			Expect(err).ToNot(HaveOccurred())
		}

		list, err := client.PodmanClient.List(context.Background())

		if !util.PodListContainsPod(*list, *kubepod1) ||
			!util.PodListContainsPod(*list, *kubepod2) {
			Fail("kubepod1 and kubepod2 not found in running containers")
		}

	})

	AfterEach(func() {
		client.PodmanClient.Delete(context.Background(), kubepod1)
		client.PodmanClient.Delete(context.Background(), kubepod2)
	})

})
