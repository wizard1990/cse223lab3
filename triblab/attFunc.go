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
	clock_str = fmt.Sprintf("%0104s", clock_str)
	str = clock_str + "," + str
	return str
}

func DiffList(from, to *trib.List) *trib.List {
	count := 0
	for _, m := range from.L {
		found := false
		for _, n := range to.L {
			if m == n {
				found = true
				break
			}
		}
		if !found {
			count++
		}
	}

	new_triblist := trib.List{}
	new_triblist.L = make([]string, count)

	i := 0
	for _, m := range from.L {
		found := false
		for _, n := range to.L {
			if m == n {
				found = true
				break
			}
		}
		if !found {
			new_triblist.L[i] = m
			i++
		}

	}
	return &new_triblist
	/*
		L := make([]*trib.List, 2)
		L[0] = l1
		L[1] = l2
		max_collect := make(map[uint64]string)
		for i, _ := range L {
			for _, v := range L[i].L {
				c, _ := SplitClock(v)

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
	*/
}

func list_append(l *trib.List, m string) *trib.List {
	l.L = append(l.L, m)
	return l
}
func list_remove(l *trib.List, n string) *trib.List {
	list := []string{}
	for _, m := range l.L {
		if m != n {
			list = append(list, m)
		} else {
			continue
		}
	}
	l.L = list
	return l
}
func GetDisplayList(l *trib.List) *trib.List {
	SortList(l)
	g := trib.List{[]string{}}
	for _, m := range l.L {
		if m[len(m)-6:] == "Append" {
			_, message := SplitClock(m[:len(m)-8])
			list_append(&g, message)
		}
		if m[len(m)-6:] == "Remove" {
			_, message := SplitClock(m[:len(m)-8])
			list_remove(&g, message)
		}
	}
	return &g
}

func SortList(l *trib.List) *trib.List {
	for i, m := range l.L {
		for j := i + 1; j < len(l.L); j++ {
			n := l.L[j]
			c1, m1 := SplitClock(m)
			c2, m2 := SplitClock(n)
			if c1 > c2 {
				l.L[i], l.L[j] = l.L[j], l.L[i]
			}
			if c1 == c2 {
				if len(m1) > 8 && len(m2) > 8 {
					if m2[len(m2)-6:] == "Append" {
						l.L[i], l.L[j] = l.L[j], l.L[i]
					}
				}
			}
		}
	}
	//fmt.Println(l)
	return l
}

func DeleteClockList(l *trib.List) *trib.List {

	new_triblist := trib.List{}
	new_triblist.L = make([]string, len(l.L))

	for i, v := range l.L {
		_, m := SplitClock(v)

		new_triblist.L[i] = m

	}
	return &new_triblist
}

func MergeKeyList(l1 *trib.List, l2 *trib.List, l3 *trib.List) *trib.List {
	L := make([]*trib.List, 3)
	L[0] = l1
	L[1] = l2
	L[2] = l3

	max_collect := make(map[string]string)
	for i, _ := range L {
		for _, v := range L[i].L {
			c := v

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
	return &new_triblist
}
func FindLargestClock(l1 *trib.List, l2 *trib.List, l3 *trib.List) (uint64, *trib.List, string) {
	L := make([]*trib.List, 3)
	L[0] = l1
	L[1] = l2
	L[2] = l3

	var max_count uint64
	max_count = 0
	max_collect := make(map[string]string)

	latest_message := ""
	for i, _ := range L {
		for _, v := range L[i].L {
			c, s := SplitClock(v)
			if c >= max_count {
				max_count = c
				latest_message = s
			}

			//Check it if in the max_collection
			value, ok := max_collect[v]
			if !ok {
				max_collect[v] = v
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

	if max_count == 0 {
		return 0, &new_triblist, ""
	}

	SortList(&new_triblist)
	return max_count, &new_triblist, latest_message
}
