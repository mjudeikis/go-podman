package e2e

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	_ "github.com/mjudeikis/go-podman/test/e2e/client"
	_ "github.com/mjudeikis/go-podman/test/e2e/specs"
)

func TestE2E(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "e2e tests")
}
