package main

import (
	"NotaborEngine/internal/notatomic"
)

type AtomicFloat32 struct {
	handle notatomic.Float32
}

func (a *AtomicFloat32) Set(val float32) {
	a.handle.Set(val)
}

func (a *AtomicFloat32) Get() float32 {
	return a.handle.Get()
}

func (a *AtomicFloat32) Add(delta float32) float32 {
	return a.handle.Add(delta)
}

func (a *AtomicFloat32) Sub(delta float32) float32 {
	return a.handle.Sub(delta)
}

func (a *AtomicFloat32) CompareAndSwap(old, new float32) bool {
	return a.handle.CompareAndSwap(old, new)
}

func (a *AtomicFloat32) Inc() float32 {
	return a.handle.Inc()
}

func (a *AtomicFloat32) Dec() float32 {
	return a.handle.Dec()
}

func (a *AtomicFloat32) GetAndSet(val float32) float32 {
	return a.handle.GetAndSet(val)
}

func (a *AtomicFloat32) Reset() {
	a.handle.Reset()
}

func (a *AtomicFloat32) IsZero() bool {
	return a.handle.IsZero()
}

func (a *AtomicFloat32) TryAdd(delta float32, limit float32) bool {
	return a.handle.TryAdd(delta, limit)
}

func (a *AtomicFloat32) SetIfGreater(val float32) {
	a.handle.SetIfGreater(val)
}

func (a *AtomicFloat32) SetIfLess(val float32) {
	a.handle.SetIfLess(val)
}

func (a *AtomicFloat32) SetIfEqual(val float32) {
	a.handle.SetIfEqual(val)
}

func (a *AtomicFloat32) String() string {
	return a.handle.String()
}

type AtomicFloat64 struct {
	handle notatomic.Float64
}

func (a *AtomicFloat64) Set(val float64) {
	a.handle.Set(val)
}

func (a *AtomicFloat64) Get() float64 {
	return a.handle.Get()
}

func (a *AtomicFloat64) Add(delta float64) float64 {
	return a.handle.Add(delta)
}

func (a *AtomicFloat64) Sub(delta float64) float64 {
	return a.handle.Sub(delta)
}

func (a *AtomicFloat64) CompareAndSwap(old, new float64) bool {
	return a.handle.CompareAndSwap(old, new)
}

func (a *AtomicFloat64) Inc() float64 {
	return a.handle.Inc()
}

func (a *AtomicFloat64) Dec() float64 {
	return a.handle.Dec()
}

func (a *AtomicFloat64) GetAndSet(val float64) float64 {
	return a.handle.GetAndSet(val)
}

func (a *AtomicFloat64) Reset() {
	a.handle.Reset()
}

func (a *AtomicFloat64) IsZero() bool {
	return a.handle.IsZero()
}

func (a *AtomicFloat64) TryAdd(delta float64, limit float64) bool {
	return a.handle.TryAdd(delta, limit)
}

func (a *AtomicFloat64) SetIfGreater(val float64) {
	a.handle.SetIfGreater(val)
}

func (a *AtomicFloat64) SetIfLess(val float64) {
	a.handle.SetIfLess(val)
}

func (a *AtomicFloat64) SetIfEqual(val float64) {
	a.handle.SetIfEqual(val)
}

func (a *AtomicFloat64) String() string {
	return a.handle.String()
}

type AtomicInt32 struct {
	handle notatomic.Int32
}

func (a *AtomicInt32) Set(val int32) {
	a.handle.Set(val)
}

func (a *AtomicInt32) Get() int32 {
	return a.handle.Get()
}

func (a *AtomicInt32) Add(delta int32) int32 {
	return a.handle.Add(delta)
}

func (a *AtomicInt32) Sub(delta int32) int32 {
	return a.handle.Sub(delta)
}

func (a *AtomicInt32) CompareAndSwap(old, new int32) bool {
	return a.handle.CompareAndSwap(old, new)
}

func (a *AtomicInt32) Inc() int32 {
	return a.handle.Inc()
}

func (a *AtomicInt32) Dec() int32 {
	return a.handle.Dec()
}

func (a *AtomicInt32) GetAndSet(val int32) int32 {
	return a.handle.GetAndSet(val)
}

func (a *AtomicInt32) Reset() {
	a.handle.Reset()
}

func (a *AtomicInt32) IsZero() bool {
	return a.handle.IsZero()
}

func (a *AtomicInt32) TryAdd(delta int32, limit int32) bool {
	return a.handle.TryAdd(delta, limit)
}

func (a *AtomicInt32) SetIfGreater(val int32) {
	a.handle.SetIfGreater(val)
}

func (a *AtomicInt32) SetIfLess(val int32) {
	a.handle.SetIfLess(val)
}

func (a *AtomicInt32) SetIfEqual(val int32) {
	a.handle.SetIfEqual(val)
}

func (a *AtomicInt32) String() string {
	return a.handle.String()
}

type AtomicInt64 struct {
	handle notatomic.Int64
}

func (a *AtomicInt64) Set(val int64) {
	a.handle.Set(val)
}

func (a *AtomicInt64) Get() int64 {
	return a.handle.Get()
}

func (a *AtomicInt64) Add(delta int64) int64 {
	return a.handle.Add(delta)
}

func (a *AtomicInt64) Sub(delta int64) int64 {
	return a.handle.Sub(delta)
}

func (a *AtomicInt64) CompareAndSwap(old, new int64) bool {
	return a.handle.CompareAndSwap(old, new)
}

func (a *AtomicInt64) Inc() int64 {
	return a.handle.Inc()
}

func (a *AtomicInt64) Dec() int64 {
	return a.handle.Dec()
}

func (a *AtomicInt64) GetAndSet(val int64) int64 {
	return a.handle.GetAndSet(val)
}

func (a *AtomicInt64) Reset() {
	a.handle.Reset()
}

func (a *AtomicInt64) IsZero() bool {
	return a.handle.IsZero()
}

func (a *AtomicInt64) TryAdd(delta int64, limit int64) bool {
	return a.handle.TryAdd(delta, limit)
}

func (a *AtomicInt64) SetIfGreater(val int64) {
	a.handle.SetIfGreater(val)
}

func (a *AtomicInt64) SetIfLess(val int64) {
	a.handle.SetIfLess(val)
}

func (a *AtomicInt64) SetIfEqual(val int64) {
	a.handle.SetIfEqual(val)
}

func (a *AtomicInt64) String() string {
	return a.handle.String()
}

type AtomicUInt32 struct {
	handle notatomic.UInt32
}

func (a *AtomicUInt32) Set(val uint32) {
	a.handle.Set(val)
}

func (a *AtomicUInt32) Get() uint32 {
	return a.handle.Get()
}

func (a *AtomicUInt32) Add(delta uint32) uint32 {
	return a.handle.Add(delta)
}

func (a *AtomicUInt32) Sub(delta uint32) uint32 {
	return a.handle.Sub(delta)
}

func (a *AtomicUInt32) CompareAndSwap(old, new uint32) bool {
	return a.handle.CompareAndSwap(old, new)
}

func (a *AtomicUInt32) Inc() uint32 {
	return a.handle.Inc()
}

func (a *AtomicUInt32) Dec() uint32 {
	return a.handle.Dec()
}

func (a *AtomicUInt32) GetAndSet(val uint32) uint32 {
	return a.handle.GetAndSet(val)
}

func (a *AtomicUInt32) Reset() {
	a.handle.Reset()
}

func (a *AtomicUInt32) IsZero() bool {
	return a.handle.IsZero()
}

func (a *AtomicUInt32) TryAdd(delta uint32, limit uint32) bool {
	return a.handle.TryAdd(delta, limit)
}

func (a *AtomicUInt32) SetIfGreater(val uint32) {
	a.handle.SetIfGreater(val)
}

func (a *AtomicUInt32) SetIfLess(val uint32) {
	a.handle.SetIfLess(val)
}

func (a *AtomicUInt32) SetIfEqual(val uint32) {
	a.handle.SetIfEqual(val)
}

func (a *AtomicUInt32) Or(mask uint32) uint32 {
	return a.handle.Or(mask)
}

func (a *AtomicUInt32) And(mask uint32) uint32 {
	return a.handle.And(mask)
}

func (a *AtomicUInt32) Clear(mask uint32) uint32 {
	return a.handle.Clear(mask)
}

func (a *AtomicUInt32) Toggle(mask uint32) uint32 {
	return a.handle.Toggle(mask)
}

func (a *AtomicUInt32) String() string {
	return a.handle.String()
}

type AtomicUInt64 struct {
	handle notatomic.UInt64
}

func (a *AtomicUInt64) Set(val uint64) {
	a.handle.Set(val)
}

func (a *AtomicUInt64) Get() uint64 {
	return a.handle.Get()
}

func (a *AtomicUInt64) Add(delta uint64) uint64 {
	return a.handle.Add(delta)
}

func (a *AtomicUInt64) Sub(delta uint64) uint64 {
	return a.handle.Sub(delta)
}

func (a *AtomicUInt64) CompareAndSwap(old, new uint64) bool {
	return a.handle.CompareAndSwap(old, new)
}

func (a *AtomicUInt64) Inc() uint64 {
	return a.handle.Inc()
}

func (a *AtomicUInt64) Dec() uint64 {
	return a.handle.Dec()
}

func (a *AtomicUInt64) GetAndSet(val uint64) uint64 {
	return a.handle.GetAndSet(val)
}

func (a *AtomicUInt64) Reset() {
	a.handle.Reset()
}

func (a *AtomicUInt64) IsZero() bool {
	return a.handle.IsZero()
}

func (a *AtomicUInt64) TryAdd(delta uint64, limit uint64) bool {
	return a.handle.TryAdd(delta, limit)
}

func (a *AtomicUInt64) SetIfGreater(val uint64) {
	a.handle.SetIfGreater(val)
}

func (a *AtomicUInt64) SetIfLess(val uint64) {
	a.handle.SetIfLess(val)
}

func (a *AtomicUInt64) SetIfEqual(val uint64) {
	a.handle.SetIfEqual(val)
}

func (a *AtomicUInt64) Or(mask uint64) uint64 {
	return a.handle.Or(mask)
}

func (a *AtomicUInt64) And(mask uint64) uint64 {
	return a.handle.And(mask)
}

func (a *AtomicUInt64) Clear(mask uint64) uint64 {
	return a.handle.Clear(mask)
}

func (a *AtomicUInt64) Toggle(mask uint64) uint64 {
	return a.handle.Toggle(mask)
}

func (a *AtomicUInt64) String() string {
	return a.handle.String()
}

type AtomicBool struct {
	handle notatomic.Bool
}

func (a *AtomicBool) Set(val bool) {
	a.handle.Set(val)
}

func (a *AtomicBool) Get() bool {
	return a.handle.Get()
}

func (a *AtomicBool) CompareAndSwap(old, new bool) bool {
	return a.handle.CompareAndSwap(old, new)
}

func (a *AtomicBool) GetAndSet(val bool) bool {
	return a.handle.GetAndSet(val)
}

func (a *AtomicBool) Reset() {
	a.handle.Reset()
}

func (a *AtomicBool) IsZero() bool {
	return a.handle.IsZero()
}

func (a *AtomicBool) SetIfTrue(val bool) bool {
	return a.handle.SetIfTrue(val)
}

func (a *AtomicBool) SetIfFalse(val bool) bool {
	return a.handle.SetIfFalse(val)
}

type AtomicPointer[T any] struct {
	handle notatomic.Pointer[T]
}

func (a *AtomicPointer[T]) Set(val *T) {
	a.handle.Set(val)
}

func (a *AtomicPointer[T]) Get() *T {
	return a.handle.Get()
}

func (a *AtomicPointer[T]) CompareAndSwap(old, new *T) bool {
	return a.handle.CompareAndSwap(old, new)
}

func (a *AtomicPointer[T]) GetAndSet(new *T) *T {
	return a.handle.GetAndSet(new)
}

func (a *AtomicPointer[T]) IsNil() bool {
	return a.handle.IsNil()
}
