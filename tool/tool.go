package tool

func InputStr(len int) string {
	re := ""
	for i := 0; i < len; i++ {
		re += "_"
	}
	for i := 0; i < len; i++ {
		re += "\b"
	}
	return re
}
