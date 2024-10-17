package window

type Options struct {
	Size int
}

type Window struct {
	buckets []Bucket
	size    int
}

func NewWindow(opt *Options) *Window {
	buckets := make([]Bucket, opt.Size)
	for i := range buckets {
		buckets[i].Points = make([]float64, 0)
		next := i + 1
		if next == opt.Size {
			next = 0
		}
		buckets[i].Points = append(buckets[i].Points, float64(next))
	}
	return &Window{buckets, opt.Size}
}
func (w *Window) ResetWindow() {
	for offset := range w.buckets {
		w.ResetBucket(offset)
	}
}
func (w *Window) ResetBucket(offset int) {
	w.buckets[offset%w.size].Reset()
}
func (w *Window) ResetBuckets(offset int, count int) {
	for i := 0; i < count; i++ {
		w.ResetBucket(offset + i)
	}
}
func (w *Window) Append(offset int, val float64) {
	w.buckets[offset%w.size].Append(val)
}
func (w *Window) Add(offset int, val float64) {
	offset %= w.size
	if w.buckets[offset].Count == 0 {
		w.buckets[offset].Append(val)
		return
	}
	w.buckets[offset].Add(0, val)
}
func (w *Window) Bucket(offset int) Bucket {
	return w.buckets[offset%w.size]
}
func (w *Window) Size() int {
	return w.size
}
func (w *Window) Iterator(offset int, count int) Iterator {
	return Iterator{
		count: count,
		cur:   &w.buckets[offset%w.size],
	}
}
