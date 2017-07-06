package runtime

import "fmt"

type ReferenceError struct {
	ref string
}

func (self *ReferenceError) Error() string {
	return fmt.Sprintf("ReferenceError: %s is not defined", self.ref)
}
