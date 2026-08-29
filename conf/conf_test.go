package conf

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestMain(m *testing.M) {
	rootPath := filepath.Join("..")
	if err := os.Chdir(rootPath); err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}

func TestConf(t *testing.T) {
	Init()
	fmt.Println(Get())
}
