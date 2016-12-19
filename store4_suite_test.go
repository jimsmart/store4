package store4_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestStore4(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Store4 Suite")
}
