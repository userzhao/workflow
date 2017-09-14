package utils

import (
	"time"
	"strconv"
	"project/workflow/models"
)

func TimeFormat(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

func StrToInt(str string) int {
	res, _ := strconv.Atoi(str)
	return res
}

func Exist(obj int, slice []int) bool {
	for _, v := range slice {
		if v == obj {
			return true
		}
	}
	return false
}


func Merge(s ...[]*models.Instance) (slice []*models.Instance) {
	switch len(s) {
	case 0:
		break
	case 1:
		slice = s[0]
		break
	default:
		s1 := s[0]
		s2 := Merge(s[1:]...)//...将数组元素打散
		slice = make([]*models.Instance, len(s1)+len(s2))
		copy(slice, s1)
		copy(slice[len(s1):], s2)
		break
	}
	return slice
}
