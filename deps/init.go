package deps

import (
	"fmt"
)

type State struct {
	err     error
	depName string
}

func NewState() *State {
	return &State{} //nolint:exhaustruct
}

func (s *State) Err() error {
	return s.err
}

func (s *State) HasError() bool {
	return s.err != nil
}

func (s *State) Name(depName string) *State {
	if s.err == nil {
		s.depName = depName
	}

	return s
}

func New[T any](state *State, initFn func() T) T {
	if !valid(state, initFn) {
		var zero T
		return zero
	}

	return initFn()
}

func Init[T any](state *State, initFn func() (T, error)) T {
	if !valid(state, initFn) {
		var zero T
		return zero
	}

	val, err := initFn()
	if err != nil {
		state.err = handleErr(state, err)
	}

	return val
}

func Init2[T1, T2 any](state *State, initFn func() (T1, T2, error)) (T1, T2) {
	if !valid(state, initFn) {
		var z1 T1
		var z2 T2
		return z1, z2
	}

	val1, val2, err := initFn()
	if err != nil {
		state.err = handleErr(state, err)
	}

	return val1, val2
}

func Init3[T1, T2, T3 any](state *State, initFn func() (T1, T2, T3, error)) (T1, T2, T3) {
	if !valid(state, initFn) {
		var z1 T1
		var z2 T2
		var z3 T3
		return z1, z2, z3
	}

	val1, val2, val3, err := initFn()
	if err != nil {
		state.err = handleErr(state, err)
	}

	return val1, val2, val3
}

func valid(state *State, initFn any) bool {
	return state != nil && state.err == nil && initFn != nil
}

func handleErr(state *State, err error) error {
	if state.depName == "" {
		return fmt.Errorf("failed to init dependency: %w", err)
	}

	return fmt.Errorf("failed to init %q dependency: %w", state.depName, err)
}
