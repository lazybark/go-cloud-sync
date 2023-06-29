package v1

import (
	"fmt"

	"github.com/lazybark/go-cloud-sync/pkg/fse"
)

func (c *FSWClient) GetServerObjList() (l []fse.FSObject, err error) {
	l, err = c.link.GetObjList()
	if err != nil {
		err = fmt.Errorf("[FSWATCHER][GetServerObjList]: %w", err)
		return
	}

	return
}

func (c *FSWClient) GetLocalObjects() (objs []fse.FSObject, err error) {
	objs, err = c.fp.ProcessDirectory(c.cfg.Root)
	if err != nil {
		err = fmt.Errorf("[FSWATCHER][DiffListWithServer]: %w", err)
		return
	}
	/*for _, o := range objs {
		fmt.Println(o.Path, o.Name)
	}*/
	return
}
