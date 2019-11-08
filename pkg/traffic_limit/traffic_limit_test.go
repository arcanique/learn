package traffic_limit

import (
	"testing"
	"time"
	"fmt"
)

func TestTraffic_FixedTimeWindow(t *testing.T) {
	t1 := time.Now().Unix()
	time.Sleep(time.Second)
	g := time.Now().Unix() - t1
	print(g)
	tra := &Traffic{
		start: time.Now().Unix(),
		counter:0,
	}
	stopCh := make(chan struct{},1)
	tic := time.Tick(time.Second)
	t_times := 0
	o_times := t_times
	go func() {
		for {
			<- tic

			fmt.Println("==============one second================",t_times - o_times)
			o_times = t_times
		}
	}()
	for {
		reject := tra.FixedTimeWindow()
		//fmt.Printf("%+v\n",reject)
		if !reject {
			t_times++
		}
	}
	<- stopCh
}

func TestTraffic_SlidingTimeWindow (t *testing.T) {
	t1 := time.Now().Unix()
	time.Sleep(time.Second)
	g := time.Now().Unix() - t1
	print(g)
	tra := &Traffic{
		start: time.Now().Unix(),
		requestTimeStamp:make([]int64,0),
		counter:0,
	}
	stopCh := make(chan struct{},1)
	tic := time.Tick(time.Second)
	t_times := 0
	o_times := t_times
	tra.requestTimeStamp = append(tra.requestTimeStamp,tra.start)
	go func() {
		for {
			<- tic

			fmt.Println("==============one second================",t_times - o_times)
			o_times = t_times
		}
	}()
	for {
		reject := tra.SlidingTimeWindow()
		//fmt.Printf("%+v\n",reject)
		if !reject {
			t_times++
		}
	}
	<- stopCh
}