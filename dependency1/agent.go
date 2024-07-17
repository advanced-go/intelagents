package dependency1

type dependency struct {
}

func newDependencyManager() *dependency {
	d := new(dependency)
	return d
}
