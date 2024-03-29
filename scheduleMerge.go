package scheduleMerge

import (
	"time"
)

// Event is a time-bounded entity. It has a start time and an end time.
//
// It is assumed that the start time is always before the end time.
//
// An Event is considered to be bounded by its start and end times as follows:
//
//	[StartTime, EndTime)
type Event interface {
	// GetStartTime returns the start time of the Event.
	GetStartTime() time.Time
	// GetEndTime returns the end time of the Event.
	GetEndTime() time.Time
	// SetStartTime sets the start time of the Event.
	SetStartTime(time.Time)
	// SetEndTime sets the end time of the Event.
	SetEndTime(time.Time)
	// Clone returns a deep copy of the Event.
	Clone() Event
}

// Schedule is a slice of Event(s).
type Schedule interface {
	SortByDesirability() // Sorts (in place) the schedule by desirability in ascending order.
	GetEvents() []Event  // Returns the events in the schedule.
}

func NewEngine(rawSchedule Schedule, trimOverlaps bool) *Engine {
	rawSchedule.SortByDesirability()
	return &Engine{
		RawSchedule:    rawSchedule.GetEvents(),
		MergedSchedule: []Event{},
		TrimOverlaps:   trimOverlaps,
	}
}

// Engine is the main struct that is used to merge the potentially conflicting raw events into a single conflict-free
// schedule.
type Engine struct {
	// The raw schedule passed to the engine via the NewEngine constructor.
	RawSchedule []Event
	// The merged schedule that is created by the engine.
	MergedSchedule []Event
	// Indicates whether the engine should trim the overlaps between the events. If true, the engine will trim the
	// overlaps between the events. If false, the engine will discard the less desirable conflicting event.
	TrimOverlaps bool

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
	// most desirable. Events in `e.MergedSchedule` are sorted by StartTime/EndTime from
	// oldest to newest and never overlap with each other.
	for _, rawEvent := range e.RawSchedule {
		if len(e.MergedSchedule) == 0 {
			e.MergedSchedule = append(e.MergedSchedule, rawEvent)
			continue
		}

		// At least one event has already been inserted into the `e.MergedSchedule`.
		// Find all events in `e.MergedSchedule` that are completely before the rawEvent. We can safely insert the
		// rawEvent after the last event that is completely before the rawEvent.
		//
		// rawEvent (more desirable):       [----)
		// PCME(s) (less desirable) : [----)
		lastSafeMergedEventIndex := findLastSafeMergedEventIndex(rawEvent, e.MergedSchedule)

		// We will isolate all the events in `e.MergedSchedule` that are potentially conflicting with the rawEvent and
		// check in detail.
		safeMergedEvents, potentialConflictMergedEvents := splitMergedEventsOnSafeInsert(lastSafeMergedEventIndex, e.MergedSchedule)

		if len(potentialConflictMergedEvents) == 0 {
			// There are no events in `e.MergedSchedule` that are potentially conflicting with the rawEvent.
			// Therefore, we can safely insert the rawEvent after the last event that is completely before the rawEvent.
			e.MergedSchedule = append(safeMergedEvents, rawEvent)
			continue
		}

		// There are events in `e.MergedSchedule` that are potentially conflicting with the rawEvent. We will check
		// each of them in detail.
		mergedSchedule := e.merge(rawEvent, potentialConflictMergedEvents)
		e.MergedSchedule = append(safeMergedEvents, mergedSchedule...)
	}

	e.mergingFinished = true
}

func (e *Engine) merge(rawEvent Event, PCMEs []Event) (mergedSchedule []Event) {
	var (
		rawStart         = rawEvent.GetStartTime()
		rawEnd           = rawEvent.GetEndTime()
		rawInserted      bool // Indicates whether the rawEvent has been inserted into the mergedSchedule.
		rawInsertedIndex int  // Indicates the index of the rawEvent in the mergedSchedule.
	)

	// PCME(s) (Potentially Conflicting Merged Event(s)) are sorted by StartTime from oldest to newest and never
	// overlap with each other.
	for PCMEIndex, PCME := range PCMEs {
		var (
			pcmeStart = PCME.GetStartTime()
			pcmeEnd   = PCME.GetEndTime()
		)

		// Check for all types of (non)overlaps between the rawEvent and the current PCME.
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

		// "1.b": rawEvent (more desirable):       [----)
		//        PCME (less desirable)    : [----)
		// if rawStart.After(pcmeEnd) || rawStart.Equal(pcmeEnd) {
		//	 // This case is already handled by `findLastSafeMergedEventIndex`. Therefore, we can safely
		//   // ignore it.
		//	 continue
		// }

		// "2.b": rawEvent (more desirable):    [----)
		//        PCME (less desirable)    : [----)
		if rawStart.After(pcmeStart) && rawStart.Before(pcmeEnd) && rawEnd.After(pcmeEnd) {
			if !e.TrimOverlaps {
				// If we are not trimming overlaps, we can safely ignore the current PCME and move on.
				if !rawInserted {
					mergedSchedule = append(mergedSchedule, rawEvent)
					rawInsertedIndex = PCMEIndex
					rawInserted = true
				}
				continue
			}

			// If we are trimming overlaps, we can trim the current PCME and insert the rawEvent after it.
			pcmePart := PCME.Clone()
			pcmePart.SetEndTime(rawStart)

			// if rawInserted {
			//     // Because we are processing the PCMEs in time order, we can safely assume that the rawEvent
			//     // has not yet been inserted.
			// }

			// e.TrimOverlaps == true && rawInserted == false
			mergedSchedule = append(mergedSchedule, pcmePart, rawEvent)
			rawInsertedIndex = PCMEIndex + 1
			rawInserted = true
			continue
		}

		// Check if the rawEvent fully contains the current mergedEvent.
		// Potential cases:
		// "3.a": rawEvent (more desirable): [----------)
		//        PCME (less desirable)    : [----------)
		//
		// "3.b": rawEvent (more desirable): [----------)
		//        PCME (less desirable)    :    [----)
		//
		// "3.d": rawEvent (more desirable): [----------)
		//        PCME (less desirable)    : [----)
		//
		// "3.e": rawEvent (more desirable): [----------)
		//        PCME (less desirable)    :       [----)
		if (pcmeStart.After(rawStart) || pcmeStart.Equal(rawStart)) &&
			(pcmeEnd.Before(rawEnd) || pcmeEnd.Equal(rawEnd)) {
			// If so, insert the rawEvent and ignore the current PCME.
			if !rawInserted {
				mergedSchedule = append(mergedSchedule, rawEvent)
				rawInsertedIndex = PCMEIndex
				rawInserted = true
			}
			continue
		}

		// Check if the rawEvent is fully contained by the current mergedEvent.
		// "3.c": rawEvent (more desirable):    [----)
		//        PCME (less desirable)    : [----------)
		if (rawStart.After(pcmeStart) || rawStart.Equal(pcmeStart)) &&
			(rawEnd.Before(pcmeEnd) || rawEnd.Equal(pcmeEnd)) {
			if !e.TrimOverlaps {
				// If we are not trimming overlaps, we can safely ignore the current PCME and move on.
				if !rawInserted {
					mergedSchedule = append(mergedSchedule, rawEvent)
				}
				break
			}

			// If we are trimming overlaps, we can split the current PCME into two parts and insert the rawEvent
			// between.
			var (
				pcmePart1 Event
				pcmePart2 Event
			)
			if !rawStart.Equal(pcmeStart) {
				pcmePart1 = PCME.Clone()
				pcmePart1.SetEndTime(rawStart)
			}
			if !rawEnd.Equal(pcmeEnd) {
				pcmePart2 = PCME.Clone()
				pcmePart2.SetStartTime(rawEnd)
			}

			if !rawInserted {
				if pcmePart1 != nil {
					mergedSchedule = append(mergedSchedule, pcmePart1)
				}
				mergedSchedule = append(mergedSchedule, rawEvent)
				if pcmePart2 != nil {
					mergedSchedule = append(mergedSchedule, pcmePart2)
				}

				mergedSchedule = append(mergedSchedule, PCMEs[PCMEIndex+1:]...)
				break
			}

			// e.TrimOverlaps == true && rawInserted == true
			var (
				wip       []Event
				beforeRaw = mergedSchedule[:rawInsertedIndex]
				afterRaw  = mergedSchedule[rawInsertedIndex+1:]
			)

			if pcmePart1 != nil {
				wip = append(wip, pcmePart1)
			}
			wip = append(wip, rawEvent)
			if pcmePart2 != nil {
				wip = append(wip, pcmePart2)
			}

			mergedSchedule = append(
				append(
					append(
						beforeRaw,
						wip...,
					),
					afterRaw...,
				),
				PCMEs[PCMEIndex+1:]...,
			)
			break
		}

		// "2.a": rawEvent (more desirable): [----)
		//        PCME (less desirable)    :    [----)
		if rawStart.Before(pcmeStart) && rawEnd.After(pcmeStart) && rawEnd.Before(pcmeEnd) {
			if !e.TrimOverlaps {
				// If we are not trimming overlaps, we can safely ignore the current PCME and move on.
				if !rawInserted {
					mergedSchedule = append(mergedSchedule, rawEvent)
					rawInsertedIndex = PCMEIndex
					rawInserted = true
				}
				continue
			}

			// If we are trimming overlaps, we can trim the current PCME and insert the rawEvent before it.
			pcmePart := PCME.Clone()
			pcmePart.SetStartTime(rawEnd)

			if !rawInserted {
				mergedSchedule = append(mergedSchedule, rawEvent, pcmePart)
				rawInsertedIndex = PCMEIndex + 1
				rawInserted = true
				continue
			}

			mergedSchedule = append(mergedSchedule, pcmePart)
			continue
		}

		// "1.a": rawEvent (more desirable): [----)
		//	      PCME (less desirable)    :       [----)
		if rawEnd.Before(pcmeStart) || rawEnd.Equal(pcmeStart) {
			// The rawEvent does not overlap with the current PCME. Therefore, we can safely insert the
			// rawEvent before the current PCME and safely move on. There will not be any other PCMEs
			// that overlap with the rawEvent.
			if !rawInserted {
				mergedSchedule = append(
					append(
						mergedSchedule,
						rawEvent,
					),
					PCMEs[PCMEIndex:]...,
				)
				break
			}

			mergedSchedule = append(mergedSchedule, PCMEs[PCMEIndex:]...)
			break
		}
	}

	return mergedSchedule
}

// findLastSafeMergedEventIndex returns the index of the last event in mergedEvents that is
// completely before the rawEvent. If there are no events in mergedEvents that are completely
// before the rawEvent, -1 is returned. It is assumed that mergedEvents is sorted by StartTime/EndTime
// from oldest to newest.
func findLastSafeMergedEventIndex(rawEvent Event, mergedEvents []Event) int {
	var (
		// Because this variable is used to store an index of a slice, we initialize it with -1 to indicate that no
		// safe index has been found yet.
		lastSafeMergedEventIndex = -1
	)

	if len(mergedEvents) == 0 {
		return lastSafeMergedEventIndex
	}

	rawStart := rawEvent.GetStartTime()
	for mergedEventIndex, mergedEvent := range mergedEvents {
		mergedEnd := mergedEvent.GetEndTime()
		if mergedEnd.Before(rawStart) || mergedEnd.Equal(rawStart) {
			lastSafeMergedEventIndex = mergedEventIndex
			continue
		}
		// Because the events in `e.MergedSchedule` are sorted by StartTime/EndTime from oldest to newest, we can
		// safely break the loop as soon as we find the first event that is not completely before the rawEvent.
		break
	}

	return lastSafeMergedEventIndex
}

func splitMergedEventsOnSafeInsert(lastSafeMergedEventIndex int, mergedEvents []Event) (safe, potentialConflict []Event) {
	if lastSafeMergedEventIndex == -1 {
		return []Event{}, mergedEvents
	}

	if len(mergedEvents) == lastSafeMergedEventIndex+1 {
		return mergedEvents, []Event{}
	}

	return mergedEvents[:lastSafeMergedEventIndex+1], mergedEvents[lastSafeMergedEventIndex+1:]
}
