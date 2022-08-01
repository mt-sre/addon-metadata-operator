package validator

import "golang.stackrox.io/kube-linter/pkg/lintcontext"

type LintContextImpl struct {
	objects []lintcontext.Object
}

// Objects returns the (valid) objects loaded from this LintContext.
func (l *LintContextImpl) Objects() []lintcontext.Object {
	return l.objects
}

// addObject adds a valid object to this LintContext
func (l *LintContextImpl) AddObjects(objs ...lintcontext.Object) {
	l.objects = append(l.objects, objs...)
}

// InvalidObjects returns any objects that we attempted to load, but which were invalid.
func (l *LintContextImpl) InvalidObjects() []lintcontext.InvalidObject {
	return nil
}
