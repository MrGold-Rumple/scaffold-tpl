/*
* @Author: Rumple
* @Email: ruipeng.wu@cyclone-robotics.com
* @DateTime: 2022/2/21 18:02
 */

package console

import (
	"fmt"
	"os"

	"github.com/gookit/color"
)

func Warn(args ...interface{}) {
	na := []interface{}{"【 WARN 】"}
	na = append(na, args...)
	color.Warnln(na...)
}

func Info(args ...interface{}) {
	na := []interface{}{"【 INFO 】"}
	na = append(na, args...)
	color.Greenln(na...)
}

func Error(t string, args ...interface{}) {
	na := []interface{}{"【 ERROR 】"}
	na = append(na, fmt.Sprintf(t, args...))
	color.Redln(na...)
}

func Fatal(t string, args ...interface{}) {
	Error(t, args...)
	os.Exit(1)
}
