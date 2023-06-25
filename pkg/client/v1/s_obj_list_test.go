package v1

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/lazybark/go-cloud-sync/pkg/fse"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetDiffListWithServer_EmptyLocal(t *testing.T) {
	now := time.Now()
	locObjs := []fse.FSObject{}
	srvObjs := []fse.FSObject{
		{Path: "1", Name: "1", Hash: "1", UpdatedAt: now},
		{Path: "2", Name: "2", Hash: "2", UpdatedAt: now},
	}
	c := &FSWClient{}

	toDownload := []fse.FSObject{
		{Path: "1", Name: "1", Hash: "1", UpdatedAt: now},
		{Path: "2", Name: "2", Hash: "2", UpdatedAt: now},
	}
	var toCreate []fse.FSObject
	var toUpdate []fse.FSObject

	download, created, updated, err := c.GetDiffListWithServer(locObjs, srvObjs)
	require.NoError(t, err)

	assert.Equal(t, true, cmp.Equal(toDownload, download))
	assert.Equal(t, true, cmp.Equal(toCreate, created))
	assert.Equal(t, true, cmp.Equal(toUpdate, updated))
}

func TestGetDiffListWithServer_EmptyServer(t *testing.T) {
	now := time.Now()
	locObjs := []fse.FSObject{
		{Path: "1", Name: "1", Hash: "1", UpdatedAt: now},
		{Path: "2", Name: "2", Hash: "2", UpdatedAt: now},
	}
	srvObjs := []fse.FSObject{}
	c := &FSWClient{}

	var toDownload []fse.FSObject
	toCreate := []fse.FSObject{
		{Path: "1", Name: "1", Hash: "1", UpdatedAt: now},
		{Path: "2", Name: "2", Hash: "2", UpdatedAt: now},
	}
	var toUpdate []fse.FSObject

	download, created, updated, err := c.GetDiffListWithServer(locObjs, srvObjs)
	require.NoError(t, err)

	assert.Equal(t, true, cmp.Equal(toDownload, download))
	assert.Equal(t, true, cmp.Equal(toCreate, created))
	assert.Equal(t, true, cmp.Equal(toUpdate, updated))
}

func TestGetDiffListWithServer_UpdatedOnBoth(t *testing.T) {
	now := time.Now()
	locObjs := []fse.FSObject{
		{Path: "1", Name: "1", Hash: "1", UpdatedAt: now},
		//Should be set to upload
		{Path: "2", Name: "2", Hash: "new_hash", UpdatedAt: now},
	}
	srvObjs := []fse.FSObject{
		//Should be set to download
		{Path: "1", Name: "1", Hash: "new_hash", UpdatedAt: now.Add(time.Second * 3)},
		{Path: "2", Name: "2", Hash: "2", UpdatedAt: now},
	}
	c := &FSWClient{}

	var toCreate []fse.FSObject
	toDownload := []fse.FSObject{
		{Path: "1", Name: "1", Hash: "1", UpdatedAt: now},
	}
	toUpdate := []fse.FSObject{
		{Path: "2", Name: "2", Hash: "new_hash", UpdatedAt: now},
	}

	download, created, updated, err := c.GetDiffListWithServer(locObjs, srvObjs)
	require.NoError(t, err)

	assert.Equal(t, true, cmp.Equal(toDownload, download))
	assert.Equal(t, true, cmp.Equal(toCreate, created))
	assert.Equal(t, true, cmp.Equal(toUpdate, updated))
}
