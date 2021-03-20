package userNotice

import "sort"

type Notice struct {
	sort.Interface
}

func (notice Notice) Less(i, j int) bool {
	if i < j {
		return true
	} else {
		return false
	}
}
