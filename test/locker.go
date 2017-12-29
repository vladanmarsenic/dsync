package test


type Locker interface {
	Lock(resource string) error
	Release(resource string) error
	Status()
}