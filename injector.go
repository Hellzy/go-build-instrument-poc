package main

// Injector visits Command types
type Injector interface {
	InjectCompile(*compileCommand)
	InjectLink(*linkCommand)
}
