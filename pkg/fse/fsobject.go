package fse

//FSObject represents and filesystem object that can be watched by IFilesystemWatcher
type FSObject struct {
	//Path holds full filesystem path from root to the file itself
	Path string

	//IsDir = true points that the object is a directory
	IsDir bool

	//Hash describes the object internals. Objects with same hash + type should be considered
	//identical.
	Hash string
}
