package cisco_v1

func ciscoActionExchange(action bool) int64 {
	if action {
		return 2
	}
	return 1
}
