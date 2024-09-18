package redirect1

func updatePercentage(curr int) int {
	next := 0
	switch curr {
	case 0:
		next = 10
	case 10:
		next = 20
	case 20:
		next = 40
	case 40:
		next = 70
	case 70:
		next = 100
	default:
	}
	return next
}
