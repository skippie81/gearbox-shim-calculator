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

func NewSets(setLengt,shims, iteration int) int {
	if setLengt == 1 {
		if iteration == 0 {
			return shims
		} else {
			return 0
		}
	} else {
		if iteration == 0 {
			t := 0
			for i := 0; i < shims; i++ {
				t = t + NewSets(setLengt-1,shims,i)
			}
			return t
		} else {
			if setLengt != 2 {
				return NewSets(setLengt, shims, iteration - 1) - NewSets(setLengt - 1, shims, iteration - 1)
			} else {
				return NewSets(setLengt, shims, iteration - 1) - 1
			}
		}
	}
}

func StartIndex(lenght, setLenght, shims, iteration int) int {
	return lenght - NewSets(setLenght,shims, iteration)
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

		for c,shim := range sl.Shims {
			for i := StartIndex(len(sublists),l,len(sl.Shims),c); i < len(sublists); i++ {
				t := make([]int,l-1)
				copy(t,sublists[i].Shims)
				t = append(t,shim)

				rsn := ResultSet{
					Thickness:	sublists[i].Thickness + shim,
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
		for _,rs := range rl.Results {
			if rs.Thickness <= target + tolerance {
				if math.Abs(float64(rs.Thickness - target)) <= float64(tolerance) {
					oklist.Results = append(oklist.Results, rs)
					if rs.Thickness == target {
						stop = true
					}
				}
			}
		}
		if i >= *maxIterations {
			stop = true
			fmt.Printf("Maximum of %v iterations reached\n",*maxIterations)
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
