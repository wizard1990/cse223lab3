package triblab

import (
    "sort"
    "trib"
)

type ByTime []*trib.Trib

func (a ByTime) Len() int {
    return len(a)
}

func (a ByTime) Swap(i, j int) {
    a[i], a[j] = a[j], a[i]
}

func (a ByTime) Less(i, j int) bool {
    return compare(a[i], a[j])
}

func compare(t1, t2 *trib.Trib) {
    if t1.Clock < t2.Clock {
        return true
    }
    else if t1.Clock == t2.Clock {
        if t1.Time.Before(t2.Time) {
            return true
        }
        else if t1.Time.Equal(t2.Time) {
            if t1.User < t2.User {
                return true
            }
            else if t1.User == t2.User {
                return t1.Message < t2.Message
            }
        }
    }
    return false
}

func tribSort(tribs []*trib.Trib){
    sort.Sort(ByTime(tribs))
}