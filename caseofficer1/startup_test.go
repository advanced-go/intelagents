package caseofficer1

import (
	"fmt"
	"github.com/advanced-go/stdlib/host"
)

func ExampleStartupPing() {
	status := host.Ping(PkgPath)
	fmt.Printf("test: Ping() -> [status:%v]\n", status)

	//Output:
	//test: Ping() -> [status:OK]

}
