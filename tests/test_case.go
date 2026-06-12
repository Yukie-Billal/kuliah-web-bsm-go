package tests

import (
	"github.com/goravel/framework/testing"

	"kuliah-web-bsm-go/bootstrap"
)

func init() {
	bootstrap.Boot()
}

type TestCase struct {
	testing.TestCase
}
