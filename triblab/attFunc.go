package triblab

import (
	"fmt"
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

func DeleteClockList(l *trib.List) *trib.List {
	return nil
}
func FindLargestClock(l1 *trib.List, l2 *trib.List, l3 *trib.List) (uint64, *trib.List, string) {
	L := make([]*trib.List, 3)
	L[0] = l1
	L[1] = l2
	L[2] = l3

	var max_count uint64
	max_count = 0
	max_collect := make(map[uint64]string)
	for i, _ := range L {
		for _, v := range L[i].L {
			c, _ := SplitClock(v)
			if c >= max_count {
				max_count = c
			}

			//Check it if in the max_collection
			value, ok := max_collect[c]
			if !ok {
				max_collect[c] = v
			} else {
				if value != v {
					fmt.Println("Error!", value, " != ", v)
				}
			}
		}
	}

	//rebuild the tribList
	new_triblist := trib.List{}
	new_triblist.L = make([]string, len(max_collect))

	j := 0
	for _, v := range max_collect {
		//fmt.Println(i, v)
		new_triblist.L[j] = v
		j++
	}

	_, message := SplitClock(max_collect[max_count])
	return max_count, &new_triblist, message
}
