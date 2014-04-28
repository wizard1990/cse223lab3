package triblab

import (
	//"fmt"
	"strconv"
	"strings"
	"trib"
)

func SplitClock(str string) (uint64, string) {
	p := strings.Index(str, ",")
	clock_str := str[:p]
	message_str := str[p+1:]

	clock, _ := strconv.ParseUint(clock_str, 10, 64)
	return clock, message_str
}
func AddClock(clock uint64, str string) string {
	clock_str := strconv.FormatUint(clock, 10)
	str = clock_str + "," + str
	return str
}
func FindLargestClock(l1 *trib.List, l2 *trib.List, l3 *trib.List) (uint64, *trib.List, string) {
	//max_collect := make(map[uint64]string)

	return 0, nil, ""
}
