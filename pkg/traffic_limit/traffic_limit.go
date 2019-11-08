package traffic_limit

import "time"

const (
	time_wnd = 60
	defaultFlow = 20
)

type Traffic struct {
	start  int64
	requestTimeStamp []int64
	counter int
}

type SlidTraffic struct {
	tStart int64
	timeWnd int64 //1min
	childWnd int64 // 10,6s
	requestCounter int
	maxRequest int //
}
// fixed time window
// 无法应对突发流量，时间力度越大应对突发流量的性能越差
func (t *Traffic)FixedTimeWindow() (rejected bool){
	now := time.Now().Unix()
	if now - t.start < time_wnd {
		if t.counter <= defaultFlow {
			t.counter++
			rejected = false
			return
		}
	} else if  now - t.start >= time_wnd {
		t.start += time_wnd
		t.counter = 1
	}
	rejected = true
	return
}

func (t *SlidTraffic)SlidingTimeWindow() (rejected bool) {
	now := time.Now().Unix()
again:
	t_start := t.tStart
	if now - t_start < time_wnd {
		if len(t.requestTimeStamp) < defaultFlow {
			rejected = false
			return
		}
	} else  {
		if len(t.requestTimeStamp) <2 {
			print(t.requestTimeStamp)
		}
		t.requestTimeStamp = t.requestTimeStamp[1:]
		goto again
	}
	rejected = true
	return
}
//
//func TokenBucket() {
//	t_start := time.Now().Unix()
//	t_now := time.Now().Unix()
//	a := t_start - t_now
//	t_start = t_now
//	t_counter += a
//	if t_counter == 0 {
//		t_counter--
//		true
//		return
//	} else if counter > max {
//		counter = max
//	}
//	false
//	return
//}