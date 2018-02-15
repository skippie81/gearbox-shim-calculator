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
	threads = flag.Int("threads",2,"Threads to use")

	dataSet = &DataSet{make(map[int]*ResultList)}
	runningThreads *int
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

type ResultSet struct {
	Thickness int        `json:"thikness"`
	Shims     []int      `json:"shims"`
}

type ResultList struct {
	Results      []ResultSet     `json:"results"`
	ResultLength int             `json:"result_length"`
}

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

type DataSet struct {
	Data 	map[int]*ResultList
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

	if dataSet.Data[l] != nil {
		fmt.Printf("Cache hit for %v\n",l)
		return dataSet.Data[l].Results
	}

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
	r.ResultLength = l
}

func (oklist *ResultList) LoadResults(target,tolerance int,rl *ResultList){
	stop := true
	for _,rs := range rl.Results {
		if rs.Thickness <= target + tolerance {
			stop = false
			if math.Abs(float64(rs.Thickness - target)) <= float64(tolerance) {
				oklist.Results = append(oklist.Results, rs)
			}
		}
	}

	if stop {
		// if all sets in current list are lager than target + tolerance
		// we just set maxitterations to 1 so there are no new calculations starting
		// as longer sets will never return smaller thikness than thinnest combination of this current set
		fmt.Printf("Set of %v lenght returned no combinations smaller than %v, stopping further calculations",rl.ResultLength,target+tolerance)
		*maxIterations = 1
	}
}

func RunCalculationThread(lenght, target, toleration int,shims ShimList,foundList *ResultList,c chan *ResultList){
	rl := ResultList{}
	rl.Generate(lenght,shims)

	c <- &rl

	foundList.LoadResults(target,toleration,&rl)
	fmt.Printf("Stop thread for setlenght %v\n",lenght)
	*runningThreads--
}

func CacheBuilder(c chan *ResultList){
	for rl := range c {
		if dataSet.Data[rl.ResultLength] == nil {
			fmt.Printf("save setlenght %v to cache\n",rl.ResultLength)
			dataSet.Data[rl.ResultLength] = rl
		}
	}
}

func Calculate(target, tolerance int, shims ShimList) (oklist ResultList){
	i :=1
	t := 0
	runningThreads = &t

	c := make(chan *ResultList)
	go CacheBuilder(c)

	for i <= *maxIterations {
		for *runningThreads < *threads -1 && i < *maxIterations {
			*runningThreads++
			fmt.Printf("Starting thread %d for setlenght %v\n",*runningThreads,i)
			go RunCalculationThread(i,target,tolerance,shims,&oklist,c)
			i++
		}
		*runningThreads++
		fmt.Printf("Starting thread %d for setlenght %v\n",*runningThreads,i)
		RunCalculationThread(i,target,tolerance,shims,&oklist,c)
		i++
	}

	close(c)

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

	if *threads >= *maxIterations {
		fmt.Printf("Setting threads to %v as this is the max set of itterations\n",*maxIterations)
		*threads = *maxIterations
	}

	fmt.Printf("Available shims: %v\n",shims)
	fmt.Printf("Target: %v  (tolerance: %v)\n",target,tolerance)
	fmt.Printf("Using %v threads for calculation\n",*threads)

	result := Calculate(target,tolerance,shims)
	fmt.Printf("%s",result)
}