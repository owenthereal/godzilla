package runtime

import "fmt"

var Console = &console{}

type console struct {
}

func (c *console) Log(s string) {
	fmt.Println(s)
}
