//go:build production

package assert

func Assert(condition bool, message string) {
	// No-op in production
}
