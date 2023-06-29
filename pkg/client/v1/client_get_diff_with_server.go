package v1

import (
	"path/filepath"

	"github.com/lazybark/go-cloud-sync/pkg/fse"
)

func (c *FSWClient) GetDiffListWithServer(locObjs []fse.FSObject, srvObjs []fse.FSObject) (dld []fse.FSObject, crtd []fse.FSObject, updtd []fse.FSObject, err error) {
	temp := make(map[string]fse.FSObject)

	for _, o := range locObjs {
		temp[filepath.Join(o.Path, o.Name)] = o
	}
	var key string
	for _, serv := range srvObjs {
		key = filepath.Join(serv.Path, serv.Name)
		if local, ok := temp[key]; ok {
			if local.Hash != serv.Hash && !local.IsDir {
				if serv.UpdatedAt.After(local.UpdatedAt) {
					dld = append(dld, local)
				}
				//We ask server to recieve file when we have newer of time-equal copy.
				//This conflict should be solved by server alone. Client always keeps local copy
				//over server copy.
				if local.UpdatedAt.After(serv.UpdatedAt) || local.UpdatedAt.Equal(serv.UpdatedAt) {
					updtd = append(updtd, local)
				}
			}
			/*if local.IsDir {
				if serv.UpdatedAt.After(local.UpdatedAt) {
					dld = append(dld, local)
				}
				if local.UpdatedAt.After(serv.UpdatedAt) || local.UpdatedAt.Equal(serv.UpdatedAt) {
					updtd = append(updtd, local)
				}
			}*/
			delete(temp, key)
		} else {
			dld = append(dld, serv)
		}
	}
	//Now we check if we have something server doesn't
	if len(temp) > 0 {
		//TO DO: think how to check if element was renamed while fs watcher was stopped?
		for _, local := range temp {
			crtd = append(crtd, local)
		}
	}

	return
}
