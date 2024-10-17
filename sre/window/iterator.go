package window

import "fmt"

type Iterator struct {
	count         int     // 总共要迭代的桶的数量
	iteratedCount int     // 当前已迭代的桶的数量
	cur           *Bucket // 当前正在迭代的桶
}

func (i *Iterator) Next() bool {
	return i.count != i.iteratedCount
}
func (i *Iterator) Bucket() Bucket {
	if !(i.Next()) {
		panic(fmt.Errorf("stat/metric: iteration out of range iteratedCount: %d count: %d", i.iteratedCount, i.count))
	}
	bucket := *i.cur
	i.iteratedCount++
	i.cur = i.cur.Next()
	return bucket
}
