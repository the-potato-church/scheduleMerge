package scheduleMerge

import (
	"time"
)

type Event interface {
	GetStartTime() time.Time
	GetEndTime() time.Time
	SetStartTime(time.Time)
	SetEndTime(time.Time)
	Clone() Event
}

// Schedule is a slice of Event(s).
type Schedule interface {
	SortByDesirability() // Sorts (in place) the schedule by desirability in ascending order.
	GetEvents() []Event
}

func NewEngine(rawSchedule Schedule, trimOverlaps bool) *Engine {
	rawSchedule.SortByDesirability()
	return &Engine{
		RawSchedule:    rawSchedule.GetEvents(),
		MergedSchedule: []Event{},
		TrimOverlaps:   trimOverlaps,
	}
}

type Engine struct {
	RawSchedule    []Event
	MergedSchedule []Event
	TrimOverlaps   bool

	mergingFinished bool
}

func (e *Engine) Merge() {
	if e.mergingFinished {
		return
	}
	if len(e.RawSchedule) == 0 {
		return
	}

	// Incoming rawEvents are sorted by Desirability from the least desirable to the
	// most desirable.
	for _, rawEvent := range e.RawSchedule {
		if len(e.MergedSchedule) == 0 {
			e.MergedSchedule = append(e.MergedSchedule, rawEvent)
			continue
		}

		rawStart := rawEvent.GetStartTime()
		rawEnd := rawEvent.GetEndTime()

		var (
			mergedSchedule []Event
			rawInserted    bool
		)
		// mergedEvents are sorted by StartTime from oldest to newest.
		for mergedEventIndex, mergedEvent := range e.MergedSchedule {
			// fact: none of the previous elements (before the current `mergedEvent`)
			// in the mergedSchedule overlap with the current rawEvent.

			mergedStart := mergedEvent.GetStartTime()
			mergedEnd := mergedEvent.GetEndTime()

			// Check if the rawEvent does not overlap with the current mergedEvent.
			if rawEnd.Before(mergedStart) || rawEnd.Equal(mergedStart) {
				// The rawEvent does not overlap with the current mergedEvent. Therefore, we can safely insert the
				// rawEvent before the current mergedEvent and safely move on. There will not be any other mergedEvents.
				// rawEvent (more desirable)   : [--------)
				// mergedEvent (less desirable):       	   [----------)
				if !rawInserted {
					mergedSchedule = append(
						append(
							mergedSchedule,
							rawEvent,
						),
						e.MergedSchedule[mergedEventIndex:]...,
					)
					break
				}
				mergedSchedule = append(mergedSchedule, e.MergedSchedule[mergedEventIndex:]...)
				break
			}
			if rawStart.After(mergedEnd) || rawStart.Equal(mergedEnd) {
				// The rawEvent does not overlap with the current mergedEvent. Therefore, we can safely insert the
				// rawEvent after the current mergedEvent. The rawEvent could not have been inserted before the current
				// mergedEvent.
				// rawEvent (more desirable)   :             [--------)
				// mergedEvent (less desirable): [----------)
				mergedSchedule = append(mergedSchedule, mergedEvent, rawEvent)
				rawInserted = true
				continue
			}

			// Check if the rawEvent collides with the current mergedEvent in all possible ways:

			// Check if the rawEvent fully contains the current mergedEvent.
			// Potential cases:
			// 1. rawEvent (more desirable)   : [------------------)
			//    mergedEvent (less desirable):      [--------)
			//
			// 2. rawEvent (more desirable)   : [------------------)
			//    mergedEvent (less desirable): [------------------)
			//
			// 3. rawEvent (more desirable)   : [------------------)
			//    mergedEvent (less desirable): [-------------)
			//
			// 4. rawEvent (more desirable)   : [------------------)
			//    mergedEvent (less desirable):      [-------------)
			if (mergedStart.After(rawStart) || mergedStart.Equal(rawStart)) &&
				(mergedEnd.Before(rawEnd) || mergedEnd.Equal(rawEnd)) {
				// If so, insert the rawEvent and ignore the current mergedEvent.
				if rawInserted {
					continue
				}
				mergedSchedule = append(mergedSchedule, rawEvent)
				rawInserted = true
				continue
			}

			// Check if the rawEvent is fully contained by the current mergedEvent.
			//   rawEvent (more desirable)   :      [--------)
			//   mergedEvent (less desirable): [------------------)
			if (rawStart.After(mergedStart) || rawStart.Equal(mergedStart)) &&
				(rawEnd.Before(mergedEnd) || rawEnd.Equal(mergedEnd)) {
				// If so, we split the current mergedEvent into two parts and insert the rawEvent between them or
				// ignore the mergedEvent all together. Because the rawEvent is fully contained by the current
				// mergedEvent, we know that  none of the following mergedEvents will conflict with the rawEvent.
				// Therefore, we can safely move on.
				//   rawEvent (more desirable)    :      [--------)
				//   mergedEvent1 (less desirable): [----)
				//   mergedEvent2 (less desirable):               [----)
				if e.TrimOverlaps {
					if !rawStart.Equal(mergedStart) {
						// rawEvent (more desirable)    :   [--------)
						// mergedEvent1 (less desirable): [----?
						event1 := mergedEvent.Clone()
						event1.SetEndTime(rawStart)
						if rawInserted {
							if len(mergedSchedule) > 1 {
								mergedSchedule = append(
									mergedSchedule[:len(mergedSchedule)-1],
									event1,
									rawEvent,
								)
							} else {
								mergedSchedule = []Event{event1, rawEvent}
							}
						} else {
							mergedSchedule = append(mergedSchedule, event1, rawEvent)
						}
					} else {
						if !rawInserted {
							mergedSchedule = append(mergedSchedule, rawEvent)
						}
					}

					if !rawEnd.Equal(mergedEnd) {
						// rawEvent (more desirable)    :  [--------)
						// mergedEvent2 (less desirable):         ?----)
						event2 := mergedEvent.Clone()
						event2.SetStartTime(rawEnd)
						mergedSchedule = append(mergedSchedule, event2)
					}
				} else {
					if rawInserted {
						continue
					}
					if len(mergedSchedule) > 1 {
						// If the mergedEvent is not the first element, we can safely append the rest of
						// the mergedSchedule.
						mergedSchedule = append(mergedSchedule[:len(mergedSchedule)-1], rawEvent)
					} else {
						mergedSchedule = []Event{rawEvent}
					}
				}

				if len(e.MergedSchedule) > mergedEventIndex+1 {
					// If the mergedEvent is not the last element, we can safely append the rest of the mergedSchedule.
					mergedSchedule = append(mergedSchedule, e.MergedSchedule[mergedEventIndex+1:]...)
				}
				break
			}

			// Check if the rawEvent overlaps with the current mergedEvent in other ways.
			// Cases:
			// 1. rawEvent (more desirable)   : [--------)
			//    mergedEvent (less desirable):      [--------)
			if rawStart.Before(mergedStart) && rawEnd.After(mergedStart) && rawEnd.Before(mergedEnd) {
				// If so, we can insert the rawEvent before the mergedEvent and trim the current mergedEvent/ignore it
				// and safely move on. There will not be any other mergedEvents that overlap with the rawEvent.
				if !rawInserted {
					mergedSchedule = append(mergedSchedule, rawEvent)
				}

				if e.TrimOverlaps {
					// rawEvent (more desirable)   : [--------)
					// mergedEvent (less desirable):          [----)
					event := mergedEvent.Clone()
					event.SetStartTime(rawEnd)
					mergedSchedule = append(mergedSchedule, event)
				}

				if len(e.MergedSchedule) > mergedEventIndex+1 {
					// If the mergedEvent is not the last element, we can safely append the rest of the mergedSchedule.
					mergedSchedule = append(mergedSchedule, e.MergedSchedule[mergedEventIndex+1:]...)
				}
				break
			}

			// 2. rawEvent (more desirable)   :      [--------)
			//    mergedEvent (less desirable): [--------)
			if rawStart.After(mergedStart) && rawStart.Before(mergedEnd) && rawEnd.After(mergedEnd) {
				// If so, we can trim the current mergedEvent/remove it from the mergeSchedule and insert the rawEvent
				// after the mergedEvent. The rawEvent could not have been inserted before the current mergedEvent.
				if e.TrimOverlaps {
					// rawEvent (more desirable)   :      [--------)
					// mergedEvent (less desirable): [----)
					event := mergedEvent.Clone()
					event.SetEndTime(rawStart)
					if len(mergedSchedule) > 1 {
						// If the mergedEvent is not the first element, we can safely append the rest of
						// the mergedSchedule.
						mergedSchedule = append(mergedSchedule[:len(mergedSchedule)-1], event)
					} else {
						mergedSchedule = []Event{event}
					}
				} else {
					if len(mergedSchedule) > 1 {
						// If the mergedEvent is not the first element, we can safely append the rest of
						// the mergedSchedule.
						mergedSchedule = mergedSchedule[:len(mergedSchedule)-1]
					} else {
						mergedSchedule = []Event{}
					}
				}

				mergedSchedule = append(mergedSchedule, rawEvent)
				rawInserted = true
				continue
			}
		}
		e.MergedSchedule = mergedSchedule
	}

	e.mergingFinished = true
}
