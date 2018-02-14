package main

/*
example input:
list: 24 27 30 33 36 39 42 45 69 93 111 117 141
target: 176
 */

import (
	"flag"
	"sort"
	"fmt"
	"math"
	"strings"
	"strconv"
	"os"
	"errors"
)

var (
	targetThickness = flag.Int("t",176,"Target thickness")
	targetMargin = flag.Int("m",1,"Margin on target")
	shimList = flag.String("shimlist","24,27,30,33,36,39,42,45,69,93,111,117,141","Comma seperated list of shims")
	maxIterations = flag.Int("M",6,"Maximums iteration depht (max shims in one set)")
)

type ShimList struct {
	Shims		[]int        `json:"shims"`
}

func newShimList(s []int) (sl ShimList){
	sort.Sort(sort.Reverse(sort.IntSlice(s)))
	sl.Shims = s
	return
}

func newShimListFromString(s string)(sl ShimList, err error){
	sShims := strings.Split(s,",")
	for _,s := range sShims {
		if i,e := strconv.Atoi(s); e == nil {
			sl.Shims = append(sl.Shims,i)
		} else {
			err = errors.New("Error: could not convert shimlist input to integers")
			return
		}
	}
	return
}

type ResultSet struct {
	Thickness int        `json:"thikness"`
	Shims     []int      `json:"shims"`
}

type ResultList struct {
	Results 	[]ResultSet        `json:"results"`
}

func T(s []int) (r int) {
	for _,t := range s {
		r = r + t
	}
	return
}

func (rs ResultSet) String() string{
	return fmt.Sprintf("%v  -> %v",rs.Thickness,rs.Shims)
}

func (rl ResultList) String() (str string){
	for _,r := range rl.Results {
		str += fmt.Sprintf("%s\n",r)
	}
	return
}

func GenArrays(l int,sl ShimList) (lists []ResultSet){
	if l == 1 {
		for _,shim := range sl.Shims {
			t := []int{shim}
			rs := ResultSet{
				Thickness: 	shim,
				Shims:		t,
			}
			lists = append(lists,rs)
		}
	} else {
		sublists := GenArrays( l - 1 , sl )
		for _,shim := range sl.Shims {
			for _,rs := range sublists {
				t := make([]int,l-1)
				copy(t,rs.Shims)
				t = append(t,shim)

				rsn := ResultSet{
					Thickness:	rs.Thickness + shim,
					Shims:		t,
				}
				lists = append(lists,rsn)
			}
		}
	}
	return
}

func (r *ResultList) Generate(l int,shims ShimList) {
	r.Results = GenArrays(l,shims)
}

func Calculate(target, tolerance int, shims ShimList) (oklist ResultList){
	stop := false
	for i := 1; stop == false; i++ {
		rl := ResultList{}
		rl.Generate(i,shims)
		stop = true
		for _,rs := range rl.Results {
			stop = false
			if rs.Thickness <= target + tolerance {
				if math.Abs(float64(rs.Thickness - target)) <= float64(tolerance) {
					oklist.Results = append(oklist.Results, rs)
				}
			}
		}
		if i >= *maxIterations {
			stop = true
			fmt.Printf("Maximum of %v iterations reached",*maxIterations)
		}
	}
	return
}

func main() {
	flag.Parse()

	shims,err := newShimListFromString(*shimList)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	target := *targetThickness
	tolerance := *targetMargin

	fmt.Printf("Available shims: %v\n",shims)
	fmt.Printf("Target: %v  (tolerance: %v)\n",target,tolerance)

	result := Calculate(target,tolerance,shims)
	fmt.Printf("%s",result)
}
