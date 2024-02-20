package scheduleMerge

import (
	"sort"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

type event struct {
	StartTime time.Time
	EndTime   time.Time
	CreatedAt time.Time
	ID        int
}

func (e *event) GetStartTime() time.Time {
	return e.StartTime
}

func (e *event) GetEndTime() time.Time {
	return e.EndTime
}

func (e *event) SetStartTime(t time.Time) {
	(*e).StartTime = t
}

func (e *event) SetEndTime(t time.Time) {
	(*e).EndTime = t
}

func (e *event) Clone() Event {
	return &event{
		StartTime: e.StartTime,
		EndTime:   e.EndTime,
		CreatedAt: e.CreatedAt,
		ID:        e.ID,
	}
}

type schedule []*event

func (s schedule) SortByDesirability() {
	sort.SliceStable(s, func(i, j int) bool {
		return s[i].CreatedAt.Before(s[j].CreatedAt)
	})
}

func (s schedule) GetEvents() []Event {
	events := make([]Event, len(s))
	for i, e := range s {
		events[i] = e
	}
	return events
}

func TestEngine_Merge(t *testing.T) {
	// Event Overlap Types:
	// 1. No overlap :
	//    a: more desirable event: [----)
	//    	 less desirable event:      [----)
	//    b: more desirable event: 		[----)
	//    	 less desirable event: [----)
	// 2. Partial overlap :
	//    a: more desirable event: [----)
	//    	 less desirable event:    [----)
	//    b: more desirable event:    [----)
	//    	 less desirable event: [----)
	// 3. Full overlap :
	//    a: more desirable event: [----)
	//    	 less desirable event: [----)
	//    b: more desirable event: [------)
	//    	 less desirable event:  [----)
	//    c: more desirable event:  [----)
	//    	 less desirable event: [------)
	//    d: more desirable event: [------)
	//    	 less desirable event: [----)
	//    e: more desirable event: [------)
	//    	 less desirable event:   [----)

	tcs := []struct {
		name             string
		testSchedule     schedule
		expectedSchedule []event
		trimOverlaps     bool
	}{
		{
			// more desirable event: [----)
			// less desirable event:       [----)
			name: "2 events-[1.a]-no trim",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
			},
			trimOverlaps: false,
		},
		{
			// more desirable event: [----)
			// less desirable event:       [----)
			name: "2 events-[1.a]-trim",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
			},
			trimOverlaps: true,
		},
		{
			// more desirable event:      [----)
			// less desirable event: [----)
			name: "2 events-[1.b]-no trim",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
			},
			trimOverlaps: false,
		},
		{
			// more desirable event:      [----)
			// less desirable event: [----)
			name: "2 events-[1.b]-trim",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
			},
			trimOverlaps: true,
		},
		{
			// more desirable event: [----)
			// less desirable event:    [----)
			name: "2 events-[2.a]-no trim",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 0, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
			},
			trimOverlaps: false,
		},
		{
			// more desirable event: [----)
			// less desirable event:    [----)
			name: "2 events-[2.a]-trim",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 0, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
			},
			trimOverlaps: true,
		},
		{
			// more desirable event:    [----)
			// less desirable event: [----)
			name: "2 events-[2.b]-no trim",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
			},
			trimOverlaps: false,
		},
		{
			// more desirable event:    [----)
			// less desirable event: [----)
			name: "2 events-[2.b]-trim",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
			},
			trimOverlaps: true,
		},
		{
			// more desirable event: [----)
			// less desirable event: [----)
			name: "2 events-[3.a]-trim independent",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
			},
		},
		{
			// more desirable event: [------)
			// less desirable event:  [----)
			name: "2 events-[3.b]-trim independent",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 0, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
			},
		},
		{
			// more desirable event:  [----)
			// less desirable event: [------)
			name: "2 events-[3.c]a-no trim",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
			},
			trimOverlaps: false,
		},
		{
			// more desirable event:  [----)
			// less desirable event: [------)
			name: "2 events-[3.c]a-trim",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
			},
			trimOverlaps: true,
		},
		{
			// more desirable event: [----)
			// less desirable event: [------)
			name: "2 events-[3.c]b-no trim",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
			},
			trimOverlaps: false,
		},
		{
			// more desirable event: [----)
			// less desirable event: [------)
			name: "2 events-[3.c]b-trim",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
			},
			trimOverlaps: true,
		},
		{
			// more desirable event:   [----)
			// less desirable event: [------)
			name: "2 events-[3.c]c-no trim",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
			},
			trimOverlaps: false,
		},
		{
			// more desirable event:   [----)
			// less desirable event: [------)
			name: "2 events-[3.c]c-trim",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
			},
			trimOverlaps: true,
		},
		{
			// more desirable event: [------)
			// less desirable event: [----)
			name: "2 events-[3.d]-trim independent",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
			},
		},
		{
			// more desirable event: [------)
			// less desirable event:   [----)
			name: "2 events-[3.e]-trim independent",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 0, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
			},
		},
		{
			// more desirable event :       [----)
			// less desirable event :  [----)
			// least desirable event:            [----)
			name: "3 events-[1.a,1.b]",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
			},
			trimOverlaps: true,
		},
		{
			// more desirable event :  [----)
			// less desirable event :     [----)
			// least desirable event:          [----)
			name: "3 events-[1.a,2.a]",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
			},
			trimOverlaps: true,
		},
		{
			// more desirable event :     [----)
			// less desirable event :  [----)
			// least desirable event:          [----)
			name: "3 events-[1.a,2.b]a",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 0, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 0, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 0, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
			},
			trimOverlaps: true,
		},
		{
			// more desirable event :     [------)
			// less desirable event :  [----)
			// least desirable event:          [----)
			name: "3 events-[1.a,2.b]b",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 0, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 0, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 0, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
			},
			trimOverlaps: true,
		},
		{
			// more desirable event :     [-----------)
			// less desirable event :  [----)
			// least desirable event:          [----)
			name: "3 events-[1.a,2.b]c",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 0, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 0, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 0, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			trimOverlaps: true,
		},
		{
			// more desirable event :  [----)
			// less desirable event :  [----)
			// least desirable event:       [----)
			name: "3 events-[1.a,3.a]",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
			},
			trimOverlaps: true,
		},
		{
			// more desirable event :  [--------)
			// less desirable event :    [----)
			// least desirable event:           [----)
			name: "3 events-[1.a,3.b]a",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 0, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
			},
			trimOverlaps: true,
		},
		{
			// more desirable event :  [----------)
			// less desirable event :    [----)
			// least desirable event:           [----)
			name: "3 events-[1.a,3.b]b",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 0, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
			},
			trimOverlaps: true,
		},
		{
			// more desirable event :  [----------------)
			// less desirable event :    [----)
			// least desirable event:           [----)
			name: "3 events-[1.a,3.b]c",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 0, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			trimOverlaps: true,
		},
		{
			// more desirable event :     [----)
			// less desirable event :  [----------)
			// least desirable event:             [----)
			name: "3 events-[1.a,3.c]",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 0, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 0, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 0, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
			},
			trimOverlaps: true,
		},
		{
			// more desirable event :  [----------)
			// less desirable event :  [----)
			// least desirable event:             [----)
			name: "3 events-[1.a,3.d]a",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
			},
			trimOverlaps: true,
		},
		{
			// more desirable event :  [------------)
			// less desirable event :  [----)
			// least desirable event:             [----)
			name: "3 events-[1.a,3.d]b",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
			},
			trimOverlaps: true,
		},
		{
			// more desirable event :  [------------------)
			// less desirable event :  [----)
			// least desirable event:             [----)
			name: "3 events-[1.a,3.d]c",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 4, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 4, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			trimOverlaps: true,
		},
		{
			// more desirable event :  [----------)
			// less desirable event :        [----)
			// least desirable event:             [----)
			name: "3 events-[1.a,3.e]",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
			},
			trimOverlaps: true,
		},
		{
			// more desirable event :       [-----)
			// less desirable event :             [----)
			// least desirable event:  [----)
			name: "3 events-[1.b,1.a]a",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
			},
			trimOverlaps: true,
		},
		{
			// more desirable event :     [-------)
			// less desirable event :             [----)
			// least desirable event:  [----)
			name: "3 events-[1.b,1.a]b",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 0, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 0, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 0, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
			},
			trimOverlaps: true,
		},
		{
			// more desirable event :  [----------)
			// less desirable event :             [----)
			// least desirable event:    [----)
			name: "3 events-[1.b,1.a]c",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 0, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
			},
			trimOverlaps: true,
		},
		{
			// more desirable event :        [----)
			// less desirable event :           [----)
			// least desirable event:  [----)
			name: "3 events-[1.b,2.a]a",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
			},
			trimOverlaps: true,
		},
		{
			// more desirable event :     [-------)
			// less desirable event :           [----)
			// least desirable event:  [----)
			name: "3 events-[1.b,2.a]b",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 0, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 0, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 0, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
			},
			trimOverlaps: true,
		},
		{
			// more desirable event :  [----------)
			// less desirable event :           [----)
			// least desirable event:  [----)
			name: "3 events-[1.b,2.a]c",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
			},
			trimOverlaps: true,
		},
		{
			// more desirable event :           [----)
			// less desirable event :        [----)
			// least desirable event:  [----)
			name: "3 events-[1.b,2.b]",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			trimOverlaps: true,
		},
		{
			// more desirable event :        [----)
			// less desirable event :        [----)
			// least desirable event:  [----)
			name: "3 events-[1.b,3.a]",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			trimOverlaps: true,
		},
		{
			// more desirable event :        [--------)
			// less desirable event :          [----)
			// least desirable event:  [----)
			name: "3 events-[1.b,3.b]a",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
		},
		{
			// more desirable event :     [----------)
			// less desirable event :        [----)
			// least desirable event:  [----)
			name: "3 events-[1.b,3.b]b",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 0, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 0, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 0, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			trimOverlaps: true,
		},
		{
			// more desirable event :  [-------------)
			// less desirable event :        [----)
			// least desirable event:  [----)
			name: "3 events-[1.b,3.b]c",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 0, 00, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			trimOverlaps: true,
		},
		{
			// more desirable event :         [----)
			// less desirable event :        [------)
			// least desirable event:  [----)
			name: "3 events-[1.b,3.c]",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
			},
			trimOverlaps: true,
		},
		{
			// more desirable event :        [----)
			// less desirable event :        [------)
			// least desirable event:  [----)
			name: "3 events-[1.b,3.d]",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
			},
			trimOverlaps: true,
		},
		{
			// more desirable event :        [------)
			// less desirable event :          [----)
			// least desirable event:  [----)
			name: "3 events-[1.b,3.e]a",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			trimOverlaps: true,
		},
		{
			// more desirable event :      [--------)
			// less desirable event :          [----)
			// least desirable event:  [----)
			name: "3 events-[1.b,3.e]b",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 0, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 0, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 0, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			trimOverlaps: true,
		},
		{
			// more desirable event : [-------------)
			// less desirable event :          [----)
			// least desirable event:  [----)
			name: "3 events-[1.b,3.e]c",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 0, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 0, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 0, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			trimOverlaps: true,
		},
		{
			// more desirable event : [----)
			// less desirable event :      [----)
			// least desirable event:         [----)
			name: "3 events-[2.a,1.a]",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
			},
			trimOverlaps: true,
		},
		{
			// more desirable event :       [----)
			// less desirable event : [----)
			// least desirable event:    [-----------)
			name: "3 events-[2.a,1.b]a",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 0, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
			},
			trimOverlaps: true,
		},
		{
			// more desirable event :       [----)
			// less desirable event : [----)
			// least desirable event:    [-------)
			name: "3 events-[2.a,1.b]b",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 0, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			trimOverlaps: true,
		},
		{
			// more desirable event :       [----)
			// less desirable event : [----)
			// least desirable event:    [----)
			name: "3 events-[2.a,1.b]c",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 0, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			trimOverlaps: true,
		},
		{
			// more desirable event :   [----)
			// less desirable event : [----)
			// least desirable event:     [----)
			name: "3 events-[2.a,2.b]a",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
			},
			trimOverlaps: true,
		},
		{
			// more desirable event :   [------)
			// less desirable event : [----)
			// least desirable event:     [----)
			name: "3 events-[2.a,2.b]b",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			trimOverlaps: true,
		},
		{
			// more desirable event :   [--------)
			// less desirable event : [----)
			// least desirable event:     [----)
			name: "3 events-[2.a,2.b]c",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			trimOverlaps: true,
		},
		{
			// more desirable event : [----)
			// less desirable event : [----)
			// least desirable event:     [----)
			name: "3 events-[2.a,3.a]",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        3,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
			},
			trimOverlaps: true,
		},
		{
			// more desirable event : [--------)
			// less desirable event :   [----)
			// least desirable event:       [----)
			name: "3 events-[2.a,3.b]a",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 45, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 45, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        3,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 45, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
			},
			trimOverlaps: true,
		},
		{
			// more desirable event : [----------)
			// less desirable event :   [----)
			// least desirable event:       [----)
			name: "3 events-[2.a,3.b]b",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			trimOverlaps: true,
		},
		{
			// more desirable event : [------------)
			// less desirable event :   [----)
			// least desirable event:       [----)
			name: "3 events-[2.a,3.b]c",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			trimOverlaps: true,
		},
		{
			// more desirable event :  [----)
			// less desirable event :[---------)
			// least desirable event:        [----)
			name: "3 events-[2.a,3.c]a",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
			},
			trimOverlaps: true,
		},
		{
			// more desirable event :  [-------)
			// less desirable event :[-----------)
			// least desirable event:        [------)
			name: "3 events-[2.a,3.c]b",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 45, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 45, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 45, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
			},
			trimOverlaps: true,
		},
		{
			// more desirable event :[-----)
			// less desirable event :[----)
			// least desirable event:   [----)
			name: "3 events-[2.a,3.d]a",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 45, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 45, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 45, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
			},
			trimOverlaps: true,
		},
		{
			// more desirable event :[-------)
			// less desirable event :[----)
			// least desirable event:   [----)
			name: "3 events-[2.a,3.d]b",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			trimOverlaps: true,
		},
		{
			// more desirable event :[----------)
			// less desirable event :[----)
			// least desirable event:   [----)
			name: "3 events-[2.a,3.d]c",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			trimOverlaps: true,
		},
		{
			// more desirable event :[----)
			// less desirable event :   [----)
			// least desirable event:      [----)
			name: "3 events-[2.a,3.e]a",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
			},
			trimOverlaps: true,
		},
		{
			// more desirable event :[-------)
			// less desirable event :   [------)
			// least desirable event:      [-----)
			name: "3 events-[2.a,3.e]b",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 15, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 15, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 15, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
			},
			trimOverlaps: true,
		},
		{
			// more desirable event :[----)
			// less desirable event :          [----)
			// least desirable event:       [----)
			name: "3 events-[2.b,1.a]a",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
			},
			trimOverlaps: true,
		},
		{
			// more desirable event :[--------)
			// less desirable event :           [----)
			// least desirable event:       [-----)
			name: "3 events-[2.b,1.a]b",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
			},
			trimOverlaps: true,
		},
		{
			// more desirable event :         [----)
			// less desirable event :   [----)
			// least desirable event: [----)
			name: "3 events-[2.b,1.b]",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			trimOverlaps: true,
		},
		{
			// more desirable event : [-----)
			// less desirable event :    [----)
			// least desirable event:  [----)
			name: "3 events-[2.b,2.a]a",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 0, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 0, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
				{
					StartTime: time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
			},
			trimOverlaps: true,
		},
		{
			// more desirable event :   [---)
			// less desirable event :    [----)
			// least desirable event:  [----)
			name: "3 events-[2.b,2.a]b",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
				{
					StartTime: time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
			},
			trimOverlaps: true,
		},
		{
			// more desirable event :    [----)
			// less desirable event :    [----)
			// least desirable event:  [----)
			name: "3 events-[2.b,3.a]",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			trimOverlaps: true,
		},
		{
			// more desirable event : [---------)
			// less desirable event :    [----)
			// least desirable event:  [----)
			name: "3 events-[2.b,3.b]a",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			trimOverlaps: true,
		},
		{
			// more desirable event : [---------)
			// less desirable event :    [----)
			// least desirable event:  [----)
			name: "3 events-[2.b,3.b]a",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			trimOverlaps: true,
		},
		{
			// more desirable event :   [------)
			// less desirable event :    [----)
			// least desirable event:  [----)
			name: "3 events-[2.b,3.b]b",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			trimOverlaps: true,
		},
		{
			// more desirable event :     [----)
			// less desirable event :    [------)
			// least desirable event:  [----)
			name: "3 events-[2.b,3.c]",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
				{
					StartTime: time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
			},
			trimOverlaps: true,
		},
		{
			// more desirable event :    [-----)
			// less desirable event :    [----)
			// least desirable event:  [----)
			name: "3 events-[2.b,3.d]",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			trimOverlaps: true,
		},
		{
			// more desirable event : [-------)
			// less desirable event :    [----)
			// least desirable event:  [----)
			name: "3 events-[2.b,3.e]a",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			trimOverlaps: true,
		},
		{
			// more desirable event :   [-----)
			// less desirable event :    [----)
			// least desirable event:  [----)
			name: "3 events-[2.b,3.e]b",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			trimOverlaps: true,
		},
		{
			// more desirable event : [----)
			// less desirable event :      [----)
			// least desirable event:      [----)
			name: "3 events-[3.a,1.a]",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
			},
			trimOverlaps: true,
		},
		{
			// more desirable event :      [----)
			// less desirable event : [----)
			// least desirable event: [----)
			name: "3 events-[3.a,1.b]",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			trimOverlaps: true,
		},
		{
			// more desirable event : [----)
			// less desirable event :    [----)
			// least desirable event:    [----)
			name: "3 events-[3.a,2.a]",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
			},
			trimOverlaps: true,
		},
		{
			// more desirable event :   [----)
			// less desirable event : [----)
			// least desirable event: [----)
			name: "3 events-[3.a,2.b]",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			trimOverlaps: true,
		},
		{
			// more desirable event : [------)
			// less desirable event :  [----)
			// least desirable event:  [----)
			name: "3 events-[3.a,3.b]",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			trimOverlaps: true,
		},
		{
			// more desirable event :  [----)
			// less desirable event : [------)
			// least desirable event: [------)
			name: "3 events-[3.a,3.c]",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
			},
			trimOverlaps: true,
		},
		{
			// more desirable event : [------)
			// less desirable event : [----)
			// least desirable event: [----)
			name: "3 events-[3.a,3.d]",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			trimOverlaps: true,
		},
		{
			// more desirable event : [------)
			// less desirable event :   [----)
			// least desirable event:   [----)
			name: "3 events-[3.a,3.e]",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			trimOverlaps: true,
		},
		{
			// more desirable event : [----)
			// less desirable event :      [------)
			// least desirable event:       [----)
			name: "3 events-[3.b,1.a]",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
			},
			trimOverlaps: true,
		},
		{
			// more desirable event :        [----)
			// less desirable event : [------)
			// least desirable event:  [----)
			name: "3 events-[3.b,1.b]",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 3, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 4, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 3, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 4, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			trimOverlaps: true,
		},
		{
			// more desirable event : [----)
			// less desirable event :   [------)
			// least desirable event:    [----)
			name: "3 events-[3.b,2.a]",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
			},
			trimOverlaps: true,
		},
		{
			// more desirable event :      [----)
			// less desirable event : [------)
			// least desirable event:  [----)
			name: "3 events-[3.b,2.b]",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 4, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 4, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			trimOverlaps: true,
		},
		{
			// more desirable event : [------)
			// less desirable event : [------)
			// least desirable event:  [----)
			name: "3 events-[3.b,3.a]",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			trimOverlaps: true,
		},
		{
			// more desirable event :  [----)
			// less desirable event : [------)
			// least desirable event:  [----)
			name: "3 events-[3.b,3.c]",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
				{
					StartTime: time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
			},
			trimOverlaps: true,
		},
		{
			// more desirable event : [----)
			// less desirable event : [------)
			// least desirable event:  [----)
			name: "3 events-[3.b,3.d]",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
				{
					StartTime: time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
			},
			trimOverlaps: true,
		},
		{
			// more desirable event :   [----)
			// less desirable event : [------)
			// least desirable event:  [----)
			name: "3 events-[3.b,3.e]",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			trimOverlaps: true,
		},
		{
			// more desirable event : [----)
			// less desirable event :       [----)
			// least desirable event:      [------)
			name: "3 events-[3.c,1.a]a",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
			},
			trimOverlaps: true,
		},
		{
			// more desirable event : [----)
			// less desirable event :      [----)
			// least desirable event:     [------)
			name: "3 events-[3.c,1.a]b",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
			},
			trimOverlaps: true,
		},
		{
			// more desirable event :       [----)
			// less desirable event :  [----)
			// least desirable event: [------)
			name: "3 events-[3.c,1.b]a",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			trimOverlaps: true,
		},
		{
			// more desirable event :        [----)
			// less desirable event :  [----)
			// least desirable event: [------)
			name: "3 events-[3.c,1.b]b",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 3, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 4, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 3, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 4, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			trimOverlaps: true,
		},
		{
			// more desirable event : [----)
			// less desirable event :   [----)
			// least desirable event:  [------)
			name: "3 events-[3.c,2.a]a",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 45, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 45, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 45, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
			},
			trimOverlaps: true,
		},
		{
			// more desirable event :    [----)
			// less desirable event :      [----)
			// least desirable event:  [---------)
			name: "3 events-[3.c,2.a]a",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 45, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 45, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 45, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
			},
			trimOverlaps: true,
		},
		{
			// more desirable event :    [----)
			// less desirable event :  [----)
			// least desirable event: [--------)
			name: "3 events-[3.c,2.b]a",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 4, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 45, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 45, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 45, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
				{
					StartTime: time.Date(2020, 1, 1, 3, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 4, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
			},
			trimOverlaps: true,
		},
		{
			// more desirable event :    [-----)
			// less desirable event :  [----)
			// least desirable event: [------)
			name: "3 events-[3.c,2.b]b",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 45, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 45, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 4, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 45, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 45, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 4, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			trimOverlaps: true,
		},
		{
			// more desirable event :  [----)
			// less desirable event :  [----)
			// least desirable event: [------)
			name: "3 events-[3.c,3.a]",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 45, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
				{
					StartTime: time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 45, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
			},
			trimOverlaps: true,
		},
		{
			// more desirable event :  [------)
			// less desirable event :   [----)
			// least desirable event: [--------)
			name: "3 events-[3.c,3.b]a",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 4, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
				{
					StartTime: time.Date(2020, 1, 1, 3, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 4, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
			},
			trimOverlaps: true,
		},
		{
			// more desirable event : [-------)
			// less desirable event :   [----)
			// least desirable event: [--------)
			name: "3 events-[3.c,3.b]b",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 4, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
				{
					StartTime: time.Date(2020, 1, 1, 3, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 4, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
			},
			trimOverlaps: true,
		},
		{
			// more desirable event :  [-------)
			// less desirable event :   [----)
			// least desirable event: [--------)
			name: "3 events-[3.c,3.b]c",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 4, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 4, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 4, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			trimOverlaps: true,
		},
		{
			// more desirable event : [--------)
			// less desirable event :   [----)
			// least desirable event:  [------)
			name: "3 events-[3.c,3.b]d",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 4, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 4, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			trimOverlaps: true,
		},
		{
			// more desirable event :   [----)
			// less desirable event :  [------)
			// least desirable event: [--------)
			name: "3 events-[3.c,3.c]",
			testSchedule: schedule{
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 4, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
			},
			expectedSchedule: []event{
				{
					StartTime: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
				{
					StartTime: time.Date(2020, 1, 1, 1, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
					ID:        3,
				},
				{
					StartTime: time.Date(2020, 1, 1, 3, 0, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 3, 30, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
					ID:        2,
				},
				{
					StartTime: time.Date(2020, 1, 1, 3, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 1, 1, 4, 0, 0, 0, time.UTC),
					CreatedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					ID:        1,
				},
			},
			trimOverlaps: true,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			e := NewEngine(tc.testSchedule, tc.trimOverlaps)
			e.Merge()
			mergedSchedule := e.MergedSchedule

			evs := make([]event, len(mergedSchedule))
			for i := range mergedSchedule {
				evs[i] = *(mergedSchedule[i].(*event))
			}

			if len(evs) != len(tc.expectedSchedule) {
				t.Logf("mergedSchedule:\n\t%+v\n", evs)
				t.Logf("expectedSchedule:\n\t%+v\n", tc.expectedSchedule)
				t.Fatalf("expected %d events, got %d", len(tc.expectedSchedule), len(mergedSchedule))
			}

			var failed bool
			for i, gotEvent := range evs {
				expectedEvent := tc.expectedSchedule[i]

				if gotEvent != expectedEvent {
					t.Logf("event %d seems to be wrong:\n%s", i+1, cmp.Diff(gotEvent, expectedEvent))
					failed = true
				}
			}

			if failed {
				t.Fatalf("expected merged schedule to be:\n\t%+v\ngot:\n\t%+v", tc.expectedSchedule, evs)
			}
		})
	}
}
