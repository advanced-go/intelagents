package guidance

import "time"

func IngressInterval() time.Duration {
	return time.Second * 1
}

func EgressInterval() time.Duration {
	return time.Second * 2
}
