package tool

import (
	"bufio"
	"fmt"
	"os"
)

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

func Input(hint string, len int, a ...interface{}) error {
	print(hint + InputStr(len))
	_, err := fmt.Scan(a...)
	return err
}

func Stop() {
	fmt.Printf("输入任意键继续...")
	ClearIOBuffer()
	b := make([]byte, 1)
	// 不知道为什么清空缓冲区后，还是有残留一个ascii为10的字符。。。但是goland里面调试时好的。。
	_, err := os.Stdin.Read(b)
	_, err = os.Stdin.Read(b)
	if err != nil {
		return
	}
}

func ClearIOBuffer() {
	myReader := bufio.NewReader(nil)
	myReader.Reset(os.Stdin)
}
