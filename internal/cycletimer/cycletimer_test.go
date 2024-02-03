package cycletimer

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestCycleTimerCustom(t *testing.T) {
	start := time.Now()
	//start := origin

	for _, tc := range []struct {
		name string
		when time.Time
		want Cycle
	}{
		{
			name: "start",
			when: start,
			want: Cycle{ID: 0, Number: 1, Phase: Rest, PhaseDuration: 10 * time.Minute, PhaseRemaining: 10 * time.Minute},
		},
		{
			name: "10 minutes in",
			when: start.Add(10 * time.Minute),
			want: Cycle{ID: 0, Number: 1, Phase: Work, PhaseDuration: 30 * time.Minute, PhaseRemaining: 30 * time.Minute},
		},
		{
			name: "13 minutes in",
			when: start.Add(13*time.Minute + 33*time.Second),
			want: Cycle{ID: 0, Number: 1, Phase: Work, PhaseDuration: 30 * time.Minute, PhaseRemaining: 26*time.Minute + 27*time.Second},
		},
		{
			name: "20 minutes in",
			when: start.Add(20 * time.Minute),
			want: Cycle{ID: 0, Number: 1, Phase: Work, PhaseDuration: 30 * time.Minute, PhaseRemaining: 20 * time.Minute},
		},
		{
			name: "40 minutes in",
			when: start.Add(40 * time.Minute),
			want: Cycle{ID: 1, Number: 2, Phase: Rest, PhaseDuration: 10 * time.Minute, PhaseRemaining: 10 * time.Minute},
		},
		{
			name: "way in",
			when: start.Add(1033 * time.Minute),
			want: Cycle{ID: 25, Number: 6, Phase: Done},
		},
		{
			name: "Done",
			when: start.Add(4 * time.Hour),
			want: Cycle{ID: 6, Number: 6, Phase: Done},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			timer := NewCustom(30*time.Minute, 10*time.Minute, start, start, 6)

			got := timer.CycleAt(tc.when)

			got.ID = 0
			tc.want.ID = 0

			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("\nwant: %+v\n got: %+v", tc.want, got)
			}
		})
	}
}

func TestTicker(t *testing.T) {
	start := time.Now()
	timer := NewCustom(30*time.Minute, 10*time.Minute, start, start, 6)
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
