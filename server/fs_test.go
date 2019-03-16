package main

import (
	"path"
	"testing"
)

func Test_file(t *testing.T) {
	t.Log(path.Join(nginxSslPath, "xxx"))
	t.Log(path.Join(nginxSslPath, "../"))
}
