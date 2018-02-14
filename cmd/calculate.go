package main

import (
	"flag"
	"sort"
	"fmt"
	"math"
	"strings"
	"strconv"
	"os"
)

var (
	targetThikness = flag.Int("t",176,"Target thikness")
	targetMargin = flag.Int("m",1,"Margin on target")
	shimList = flag.String("shimlist","24,27,30,33,36,39,42,45,69,93,111,117,141","comma seperated list of shims")
	maxItterations = flag.Int("M",20,"Maximums iteration depht")
)

type ResultSet []int

func T(s []int) (r int) {
	for _,t := range s {
		r = r + t
	}
	return
}

func GenArrays(l int,shims []int) (lists [][]int){
	if l == 1 {
		for _,s := range shims {
			t := []int{s}
			lists = append(lists,t)
		}
	} else {
		sublists := GenArrays( l - 1 , shims )
		for _,s := range shims {
			for _,sl := range sublists {
				t := make([]int,len(sl))
				copy(t,sl)
				t = append(t,s)
				lists = append(lists,t)
			}
		}
	}
	return
}

func Caluculate(target, tolerance int, shims []int){
	stop := false
	for i := 1 ; stop == false; i++ {
		l := GenArrays(i,shims)
		stop = true
		for _,s := range l {
			t := T(s)
			if t <= target {
				stop = false
			}
			if math.Abs(float64(target - t)) <= float64(tolerance) {
				fmt.Printf("Possible sollution: %v (sum: %v) (tolerance: %v)\n",s,t,target - t)
			}
		}
		if i >= *maxItterations {
			stop = true
			fmt.Println("Maximum iterations reached")
		}
	}
}

func main() {
	flag.Parse()

	sShims := strings.Split(*shimList,",")
	var shims []int
	for _,s := range sShims {

		if i,e := strconv.Atoi(s); e == nil {
			shims = append(shims,i)
		} else {
			fmt.Println("Error: could not convert shimlist input to integers")
			os.Exit(1)
		}
	}

	sort.Sort(sort.Reverse(sort.IntSlice(shims)))

	target := *targetThikness
	tolerance := *targetMargin

	fmt.Printf("Available shims: %v\n",shims)
	fmt.Printf("Target: %v  (tolerance: %v)\n",target,tolerance)

	Caluculate(target,tolerance,shims)
}

/*
example input:
list: 24 27 30 33 36 39 42 45 69 93 111 117 141
target: 176
 */
