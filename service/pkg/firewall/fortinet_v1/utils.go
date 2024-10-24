package fortinet_v1

func actionExchange(action string) int64 {
	if action == actionAccept {
		return 2
	}
	return 1
}

func statusExchange(status string) bool {
	if status == statusEnable {
		return true
	}
	return false
}
