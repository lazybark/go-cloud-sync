package v1

import (
	"fmt"

	"github.com/lazybark/go-cloud-sync/pkg/fselink/v1/proto"
)

func (c *FSWClient) GetServerObjList() (l []proto.FSObject, err error) {
	l, err = c.link.GetObjList()
	if err != nil {
		return l, fmt.Errorf("[GetServerObjList]: %w", err)
	}

	return
}

func (c *FSWClient) GetLocalObjects() (objs []proto.FSObject, err error) {
	objs, err = c.fp.ProcessDirectory(c.cfg.Root)
	if err != nil {
		err = fmt.Errorf("[DiffListWithServer]: %w", err)
		return
	}
	/*for _, o := range objs {
		fmt.Println(o.Path, o.Name)
	}*/
	return
}
