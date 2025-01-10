// List of functions to quickly unwrap objects with possible errors
// This is done when we are already in an error state so we don't care about new errors happening at this point
package noerror

import "io"

func ReadAll(r io.Reader) []byte {
	b, err := io.ReadAll(r)
	if err != nil {
		panic(err)
	}
	return b
}
