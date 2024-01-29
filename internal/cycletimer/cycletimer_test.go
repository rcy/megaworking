package cycles

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestCycleTimer(t *testing.T) {
	start := origin

	for _, tc := range []struct {
		when time.Time
		want Cycle
	}{
		{
			when: start,
			want: Cycle{ID: 0, State: Rest, Remaining: 10 * time.Minute},
		},
		{
			when: start.Add(10 * time.Minute),
			want: Cycle{ID: 0, State: Work, Remaining: 30 * time.Minute},
		},
		{
			when: start.Add(13*time.Minute + 33*time.Second),
			want: Cycle{ID: 0, State: Work, Remaining: 26*time.Minute + 27*time.Second},
		},
		{
			when: start.Add(20 * time.Minute),
			want: Cycle{ID: 0, State: Work, Remaining: 20 * time.Minute},
		},
		{
			when: start.Add(40 * time.Minute),
			want: Cycle{ID: 1, State: Rest, Remaining: 10 * time.Minute},
		},
		{
			when: start.Add(1033 * time.Minute),
			want: Cycle{ID: 25, State: Work, Remaining: 7 * time.Minute},
		},
	} {
		t.Run(fmt.Sprint(tc.when), func(t *testing.T) {
			timer := New()

			got := timer.CycleAt(tc.when)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("\nwant: %v\n got: %v", tc.want, got)
			}
		})
	}
}

func TestTicker(t *testing.T) {
	timer := New()
	ticker := timer.NewTicker(time.Millisecond * 100)

	counter := 0

	go func() {
		for cycle := range ticker.C {
			counter += 1
			fmt.Println(cycle)
		}
	}()

	fmt.Println("start sleep")
	time.Sleep(550 * time.Millisecond)

	ticker.Stop()
	if counter != 6 {
		t.Errorf("wanted 6 got %d", counter)
	}

	fmt.Println("stopped ticker")
}
