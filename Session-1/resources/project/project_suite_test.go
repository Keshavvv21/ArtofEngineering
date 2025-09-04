package project_test

import (
	"testing"

	. "github.com/onsi/ginkgo" // Ginkgo BDD-style testing framework
	. "github.com/onsi/gomega" // Gomega matchers for assertions
)

// TestProject is the entry point for running the test suite.
// It registers Gomega's fail handler and then starts the Ginkgo specs
// defined in this package.
//
// The "Project Suite" string is just a label for the suite name
// that will appear in the test output.
func TestProject(t *testing.T) {
	// RegisterFailHandler tells Gomega how to handle assertion failures.
	RegisterFailHandler(Fail)

	// RunSpecs runs all specs (Describe/Context/It blocks) defined
	// in *_test.go files under this package.
	RunSpecs(t, "Project Suite")
}
