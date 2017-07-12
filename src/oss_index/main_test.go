package oss_index

import (
	"testing"
)

func TestMain2(t *testing.T) {
	Main()
}

func TestAllPath(t *testing.T) {
	getAllPath(getBucket())
}