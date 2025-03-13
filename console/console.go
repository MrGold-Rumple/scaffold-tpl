/*
* @Author: Rumple
* @Email: wrp357711589@163.com
* @DateTime: 2022/2/21 18:02
 */

package console

import (
	"fmt"
	"os"

	"github.com/gookit/color"
)

func Warn(args ...any) {
	na := []any{"【 WARN 】"}
	na = append(na, args...)
	color.Warnln(na...)
}

func Info(args ...any) {
	na := []any{"【 INFO 】"}
	na = append(na, args...)
	color.Greenln(na...)
}

func Error(t string, args ...any) {
	na := []any{"【 ERROR 】"}
	na = append(na, fmt.Sprintf(t, args...))
	color.Redln(na...)
}

func Fatal(t string, args ...any) {
	Error(t, args...)
	os.Exit(1)
}
