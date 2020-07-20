package util

import (
	"fmt"
	"io"
	"log"
)

// 三元表达式
func If(flag bool, a, b interface{}) interface{} {
	if flag {
		return a
	}
	return b
}

// 判断是否报错
func IsNilError(info string, err error) bool {
	if err != nil {
		fmt.Println(info, err)
		log.Println(info, err)
		panic(err)
		// os.Exit(1)
		return false
	}
	return true
}

// 判断是否已读取到末尾
func IsEofError(info string, err error) bool {
	if err != io.EOF {
		fmt.Println(err, info)
		log.Println(err, info)
		return false
	}
	return true
}

func IsHttpNilError(info string, err error) bool {
	if err != nil && err != io.EOF {
		fmt.Println(info, err)
		log.Println(info, err)
		panic(err)
		// os.Exit(1)
		return false
	}
	return true
}
