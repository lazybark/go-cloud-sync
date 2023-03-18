//fse means filesystem events. This package stores interface for filesystem watcher that should
//watch specified directory for changes and report back.

package fse

//IFilesystemWatcher represents watcher that uses event (evc) and error (erc) channels to report all
//changes in specified dir and its subdirs.
type IFilesystemWatcher interface {
	//Init sets initial parameters for fs watcher
	Init(root string, evc chan (FSEvent), erc chan (error)) error
	//Start launches the watcher routine or returns an error
	Start() error
	//Stop stops the watcher routine or returns an error. It also should close event & error channels,
	//which means new Start() will need new Init() with new channels
	Stop() error
}

//FSEvent represents and event that can occur in the filesystem. It contains an object data
//and event type
type FSEvent struct {
	Object FSObject
	Action FSAction
}

//FSAction is an action that was performed on a filesystem element
type FSAction int

const (
	//NoAction means represents an empty value. For example, when there is no need for
	//other party to know about event or event was unknown
	NoAction FSAction = iota
	//Create means that some object was created in the filesystem
	Create
	//Write represents any possible update to object in the filesystem
	Write
	//Remove means that some object was deleted from the filesystem
	Remove
	//Rename means that some object was renamed
	//(also can be processed like deleted->created series of actions by some systems)
	Rename
	//Chmod means that access permissions were changed
	Chmod
)

func (a FSAction) String() string {
	return [...]string{"No action", "Create", "Write", "Remove", "Rename", "Chmod"}[a]
}

//FSObject represents and filesystem object that can be watched by IFilesystemWatcher
type FSObject struct {
	Path  string
	IsDir bool
}
