# electricityHelper

示例：

```go
package main

import (
	"fmt"
	"github.com/Mmx233/electricityHelper"
)

func main() {
	d, err := elec.GetInfo(10000) //寝室号
	if err != nil {
		panic(err)
	}

	fmt.Println(d)
}
```