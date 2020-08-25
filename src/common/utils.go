package common

func AbsInt32(val int32) int32 {
	if val < 0 {
		return -val
	} else {
		return val
	}
}

func AbsInt64(val int64) int64 {
	if val < 0 {
		return -val
	} else {
		return val
	}
}
