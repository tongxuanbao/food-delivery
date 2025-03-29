//go:build !production

package assert

func Assert(condition bool, message string) {
	if !condition {
		panic(message)
	}
}
