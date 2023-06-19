package fse

//FSAction is an action that was performed on a filesystem element
type FSAction int

const (
	//NoAction represents an empty value. For example, when there is no need for
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
