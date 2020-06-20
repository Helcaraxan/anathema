package context

type Context struct{}

func Background() Context {
	return Context{}
}
