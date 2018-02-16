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
	maxIterations = flag.Int("M",10,"Maximums iteration depht (max shims in one set)")
	threads = flag.Int("threads",2,"Threads to use") // not used for now
)


// this holds the list of avialable shims
// loaded on startup

type ShimList struct {
	Shims		[]int        `json:"shims"`
}

// functions to create the shimlist from input string
// comma separated set op integers

func newShimList(s []int) (sl ShimList){
	sort.Sort(sort.Reverse(sort.IntSlice(s)))
	sl.Shims = s
	return
}

func newShimListFromString(s string)(sl ShimList, err error){
	sShims := strings.Split(s,",")
	var l []int
	for _,s := range sShims {
		if i,e := strconv.Atoi(s); e == nil {
			l = append(l,i)
		} else {
			err = errors.New("Error: could not convert shimlist input to integers")
			return
		}
	}
	sl = newShimList(l)
	return
}

// A resulstet is ONE list of shims used
// Generating a specific thikness
type ResultSet struct {
	Thickness int        `json:"thikness"`
	Shims     []int      `json:"shims"`
}

// A resultlist holds a list of resultsets
// could be all possible combinations or just a subset that matches the thickness criteria
type ResultList struct {
	Results      []ResultSet     `json:"results"`
}

// a dataset holds a map with pointers to ResultList
// is used to point to a specific resultlist of a specific lenght
type DataSet struct {
	Data 	map[int]*ResultList
}


// Implementing the String Interfac for ResultList and Result set for outpot
func (rs ResultSet) String() string{
	return fmt.Sprintf("%v  -> %v",rs.Thickness,rs.Shims)
}

func (rl ResultList) String() (str string){

	var exact string

	str += fmt.Sprintf("\nWe found the following sets for %v +- %v\n========================================================== \n\n",*targetThickness,*targetMargin)

	for _,r := range rl.Results {
		str += fmt.Sprintf("%s\n",r)
		if r.Thickness == *targetThickness {
			exact += fmt.Sprintf("%s\n",r)
		}
	}

	if exact != "" {
		str += fmt.Sprintf("\nExact Matches found: \n========================================================== \n\n%s",exact)
	} else {
		str += fmt.Sprint("\nNo exact matches found \n========================================================== \n")
	}

	return
}


// generation of the resultsets

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

func GenArrays(l int,sl ShimList) (rl ResultList){
	if l == 1 {
		for _,shim := range sl.Shims {
			t := []int{shim}
			rs := ResultSet{
				Thickness: 	shim,
				Shims:		t,
			}
			rl.Results = append(rl.Results,rs)
		}
	} else {
		sublist := GenArrays( l - 1 , sl )
		for c,shim := range sl.Shims {
			for i := StartIndex(len(sublist.Results),l,len(sl.Shims),c); i < len(sublist.Results); i++ {
				t := make([]int,l-1)
				copy(t,sublist.Results[i].Shims)
				t = append(t,shim)

				rsn := ResultSet{
					Thickness:	sublist.Results[i].Thickness + shim,
					Shims:		t,
				}
				rl.Results = append(rl.Results,rsn)
			}
		}
	}
	return
}


// check valid results
func (oklist *ResultList) LoadResults(target,tolerance int,rl *ResultList){
	for _,rs := range rl.Results {
		if rs.Thickness <= target + tolerance {
			if math.Abs(float64(rs.Thickness - target)) <= float64(tolerance) {
				oklist.Results = append(oklist.Results, rs)
			}
		}
	}
}


func Calculate(target, tolerance, maxsetlenght int, shims ShimList) (oklist ResultList){
	d := DataSet{make(map[int]*ResultList)}

	for i := maxsetlenght; i > 0; i-- {
		a := GenArrays(i,shims)
		d.Data[i] = &a
	}

	for _,rl := range d.Data{
		oklist.LoadResults(target,tolerance,rl)
	}
	return
}

func (s *ShimList) getLargestPossibleSet(target, tolerance int) int {
	smallestShim := s.Shims[len(s.Shims)-1]
	return int(( target + tolerance ) / smallestShim)
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

	fmt.Printf("Available shims: %v\n",shims.Shims)
	fmt.Printf("Target: %v  (tolerance: %v)\n",target,tolerance)

	if m := shims.getLargestPossibleSet(target,tolerance); *maxIterations > m {
		fmt.Printf("A set op %v shims is the largest possible set that could have a thikness within required range.\n",m)
		fmt.Println("Setting max set of itterations accordingly")
		*maxIterations = m
	}

	if *threads >= *maxIterations {
		fmt.Printf("Setting threads to %v as this is the max set of itterations\n",*maxIterations)
		*threads = *maxIterations
	}


	fmt.Printf("Using %v threads for calculation\n",*threads)

	result := Calculate(target,tolerance,*maxIterations,shims)
	fmt.Printf("%s",result)
}