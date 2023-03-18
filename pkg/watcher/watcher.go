package watcher

import (
	"github.com/lazybark/go-cloud-sync/pkg/fse"
	v1 "github.com/lazybark/go-cloud-sync/pkg/watcher/v1"
)

// NewV1 returns v1 implementation of IFilesystemWatcher interface
func NewV1() fse.IFilesystemWatcher {
	return v1.NewWatcher()
}
