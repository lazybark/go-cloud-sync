//fse means filesystem events. This package stores interface for filesystem watcher that should
//watch specified directory for changes and report back.

package fse

import (
	"encoding/json"
	"fmt"
)

// FSEvent represents and event that can occur in the filesystem. It contains an object data
// and event type
type FSEvent struct {
	//Object is the object that was the target of the action
	Object FSObject

	//Action is the Action that was made over the object
	Action FSAction
}

func (e FSEvent) JSON() ([]byte, error) {
	ebg, err := json.Marshal(e)
	if err != nil {
		return ebg, fmt.Errorf("[FSEvent->JSON] %w", err)
	}

	return ebg, nil
}
