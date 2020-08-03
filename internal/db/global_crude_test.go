package db

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

//TestRun runs all tests
func TestRun(t *testing.T) {
	suite.Run(t, new(CategoryTestSuite))
	suite.Run(t, new(ProductTestSuite))
}
