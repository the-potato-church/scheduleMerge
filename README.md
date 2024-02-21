![ScheduleMerge Logo](ScheduleMerge.png)

# ScheduleMerge

This is a simple package to produce conflict-free lists of time-bound objects from a list of potentially conflicting
objects.

ScheduleMerge exposes two interfaces:
- `Event` - a time-bound object with a start and end time
- `Schedule` - a slice of `Event`s

The `Event` has to implement the following functions:
- `GetStartTime() time.Time` - returns the start time of the event.
- `GetEndTime() time.Time` - returns the end time of the event.
- `SetStartTime(time.Time)` - sets the start time of the event (this is used when trimming an event).
- `SetEndTime(time.Time)` - sets the end time of the event (this is used when trimming an event).
- `Clone() Event` - returns a deep copy of the event (this is used when trimming an event).

The `Schedule` has to implement the following functions:
- `GetEvents() []Event` - returns a slice of the `Event`s in the `Schedule`.
- `SortByDesirability()` - an in-place sort function for the `Event`s in the `Schedule` (the concept of *desirability* 
  is further explained below).

This package further provides an `Engine` (and it's constructor function `NewEngine(rawSchedule Schedule, trimOverlaps 
bool) *Engine`) which is used to do the actual sorting and merging of the `Event`s. The `Engine` is used via the
`Merge()` function. The `Engine` exposes the following fields:
- `RawSchedule` - the original `Schedule` passed to the `Engine` constructor.
- `MergedSchedule` - the `Schedule` produced by the `Merge()` function.
- `TrimOverlaps` - a boolean flag which determines whether the `Engine` should trim overlapping `Event`s or not.

If the `TrimOverlaps` (or `trimOverlaps` for its constructor) flag is set to `true`, the `Engine` will trim conflicting
`Event`s to produce a conflict-free `Schedule`. If the flag is set to `false`, the `Engine` will discard conflicting
`Event` with lower desirability to produce a conflict-free `Schedule`.

## Event

The `Event` interface is used to represent a time-bound object with a start and end time as follows: **[start, end)**.

## Desirability

The concept of *desirability* is used to decide which `Event` to determine which `Event` to prioritise when a conflict
occurs. The `Event` with the higher desirability is kept in the `Schedule` while the `Event` with the lower desirability
is either trimmed or discarded (depending on the `TrimOverlaps` flag). The `Event` with the higher desirability is 
determined by the `SortByDesirability()` function of the `Schedule` interface.

What *desirability* means is up to the user to decide. The `SortByDesirability()` function can be implemented in any way
to sort the `Event`s in the `Schedule` by any criteria. This can be anything from the length of the `Event` to the
number of coffee breaks you head on that day.

## Usage

This package is shared under the Apache License, Version 2.0. See the [LICENSE.md](LICENSE.md) file for details.
