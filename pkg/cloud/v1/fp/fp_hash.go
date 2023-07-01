package fp

import (
	"time"

	"github.com/lazybark/go-cloud-sync/pkg/synclink/v1/proto"
	"github.com/lazybark/go-helpers/hasher"
)

func (fp *FileProcessor) CheckFileHash(obj proto.FSObject) string {
	var sleep int
	var hash string
	var err error
	for {
		hash, err = hasher.HashFilePath(obj.Path, hasher.SHA256, 8192)
		if err != nil {
			if sleep >= 3 {
				//If object isn't readable - it's just ignored until next action.
				//Deprecating rescanBuffer for now. Seems useless as we still recieve info about new action
				//and we can create object after it's modified.
				break
			}
			time.Sleep(time.Second * 1)
			sleep++
		} else {
			break
		}
	}

	return hash
}

func (fp *FileProcessor) CheckDirHash(obj proto.FSObject) string {
	return ""
}
