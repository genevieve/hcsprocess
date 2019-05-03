package hcsprocess_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestHcsprocess(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Hcsprocess Suite")
}
