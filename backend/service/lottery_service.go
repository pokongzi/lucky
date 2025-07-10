package service

import (
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"time"
)

func RandomSSQ() string {
	rand.Seed(time.Now().UnixNano())
	red := rand.Perm(33)[:6]
	for i := range red {
		red[i]++
	}
	sort.Ints(red)
	blue := rand.Intn(16) + 1
	redStrs := make([]string, 6)
	for i, v := range red {
		redStrs[i] = fmt.Sprintf("%02d", v)
	}
	return fmt.Sprintf("%s|%02d", strings.Join(redStrs, ","), blue)
}

func RandomDLT() string {
	rand.Seed(time.Now().UnixNano())
	front := rand.Perm(35)[:5]
	for i := range front {
		front[i]++
	}
	sort.Ints(front)
	back := rand.Perm(12)[:2]
	for i := range back {
		back[i]++
	}
	sort.Ints(back)
	frontStrs := make([]string, 5)
	for i, v := range front {
		frontStrs[i] = fmt.Sprintf("%02d", v)
	}
	backStrs := make([]string, 2)
	for i, v := range back {
		backStrs[i] = fmt.Sprintf("%02d", v)
	}
	return fmt.Sprintf("%s|%s", strings.Join(frontStrs, ","), strings.Join(backStrs, ","))
}
