package v1

type FileProcessor struct {
	escSymbol string
	root      string
	cacheRoot string
}

func NewFileProcessorV1(escSymbol, root, cacheRoot string) *FileProcessor {
	fp := FileProcessor{
		escSymbol: escSymbol,
		root:      root,
		cacheRoot: cacheRoot,
	}

	return &fp
}
