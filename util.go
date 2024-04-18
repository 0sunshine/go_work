package main

import (
	"errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/Knetic/govaluate.v2"
	"strconv"
	"strings"
)

func FormatCmd(cmd string, args map[string]int) (newCmd string, err error) {

	parameters := make(map[string]interface{}, 8)
	for k, v := range args {
		parameters[k] = v
	}

	strBuilder := strings.Builder{} //生成新的

	argsBytes := []byte{}   //存参数
	cmdBytes := []byte(cmd) //转换输入进行遍历

	var prevC byte = 0 //上一个字符，判断\{转义使用

	for _, c := range cmdBytes {

		if len(argsBytes) > 0 {
			if c != '}' || prevC == '\\' {
				argsBytes = append(argsBytes, c)
			} else {
				//完整的参数
				arg := string(argsBytes[1:])
				argsBytes = []byte{}

				expression, err := govaluate.NewEvaluableExpression(arg)
				if err != nil {
					logrus.Errorf("cmd: %s, arg error, arg: %s, err: %s", cmd, arg, err)
					return newCmd, err
				}

				result, err := expression.Evaluate(parameters)
				if err != nil || result == nil {
					logrus.Errorf("cmd: %s, arg error, arg: %s, err: %s", cmd, arg, err)
					return newCmd, err
				}

				value, ok := result.(float64)
				if !ok {
					logrus.Errorf("cmd: %s, arg error, arg: %s, result not a float64, result (%T): %v", cmd, arg, result, result)
					return newCmd, errors.New("type error")
				}

				strBuilder.WriteString(strconv.Itoa(int(value)))
			}
		} else if c == '{' && prevC != '\\' {
			argsBytes = append(argsBytes, c)
		} else {
			strBuilder.WriteByte(c)
		}

		prevC = c
	}

	newCmd = strBuilder.String()
	return newCmd, err
}

func safeSplit(s string) []string {
	split := strings.Split(s, " ")

	var result []string
	var inquote string
	var block string
	for _, i := range split {
		if inquote == "" {
			if strings.HasPrefix(i, "'") || strings.HasPrefix(i, "\"") {
				inquote = string(i[0])
				block = strings.TrimPrefix(i, inquote) + " "
			} else {
				result = append(result, i)
			}
		} else {
			if !strings.HasSuffix(i, inquote) {
				block += i + " "
			} else {
				block += strings.TrimSuffix(i, inquote)
				inquote = ""
				result = append(result, block)
				block = ""
			}
		}
	}

	return result
}
