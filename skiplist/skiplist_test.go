package skiplist

// apache 2.0 antlabs
import (
	"fmt"
	"sync"
	"testing"

	"github.com/antlabs/gstl/cmp"
)

func Test_New(t *testing.T) {
	n := New[int, int]()
	if n == nil {
		t.Errorf("expected non-nil, got nil")
	}
}

func Test_SetGet(t *testing.T) {
	zset := New[float64, string]()
	max := 100.0
	for i := 0.0; i < max; i++ {
		zset.Set(i, fmt.Sprintf("%d", int(i)))
	}

	for i := 0.0; i < max; i++ {
		v := zset.Get(i)
		if v != fmt.Sprintf("%d", int(i)) {
			t.Errorf("expected %s, got %s", fmt.Sprintf("%d", int(i)), v)
		}
	}
}

// 测试插入重复
func Test_InsertRepeatingElement(t *testing.T) {
	sl := New[float64, string]()
	max := 100
	for i := 0; i < max; i++ {
		sl.Set(float64(i), fmt.Sprint(i))
	}

	for i := 0; i < max; i++ {
		sl.Set(float64(i), fmt.Sprint(i+1))
	}

	for i := 0; i < max; i++ {
		if sl.Get(float64(i)) != fmt.Sprint(i+1) {
			t.Errorf("expected %s, got %s", fmt.Sprint(i+1), sl.Get(float64(i)))
		}
	}
}

func Test_SetGetRemove(t *testing.T) {
	zset := New[float64, float64]()

	max := 100.0
	for i := 0.0; i < max; i++ {
		zset.Set(i, i)
	}

	for i := 0.0; i < max; i++ {
		zset.Remove(i)
		if float64(zset.Len()) != max-1 {
			t.Errorf("expected %f, got %f", max-1, float64(zset.Len()))
		}
		for j := 0.0; j < max; j++ {
			if j == i {
				continue
			}
			v, ok := zset.TryGet(j)
			if !ok {
				t.Errorf("expected true for score:%f, i:%f, j:%f", j, i, j)
				return
			}
			if v != j {
				t.Errorf("expected %f, got %f", j, v)
			}
		}
		zset.Set(i, i)
	}
}

// 测试TopMin, 它返回最小的几个值
func Test_Skiplist_TopMin(t *testing.T) {

	need := []int{}
	count10 := 10
	count100 := 100
	count1000 := 1000

	for i := 0; i < count1000; i++ {
		need = append(need, i)
	}

	needCount := []int{count10, count100, count100}
	for i, b := range []*SkipList[float64, int]{
		// btree里面元素 少于 TopMin 需要返回的值
		func() *SkipList[float64, int] {
			b := New[float64, int]()
			for i := 0; i < count10; i++ {
				b.Set(float64(i), i)
			}

			if b.Len() != count10 {
				t.Errorf("expected %d, got %d", count10, b.Len())
			}
			return b
		}(),
		// btree里面元素 等于 TopMin 需要返回的值
		func() *SkipList[float64, int] {

			b := New[float64, int]()
			for i := 0; i < count100; i++ {
				b.Set(float64(i), i)
			}
			if b.Len() != count100 {
				t.Errorf("expected %d, got %d", count100, b.Len())
			}
			return b
		}(),
		// btree里面元素 大于 TopMin 需要返回的值
		func() *SkipList[float64, int] {

			b := New[float64, int]()
			for i := 0; i < count1000; i++ {
				b.Set(float64(i), i)
			}
			if b.Len() != count1000 {
				t.Errorf("expected %d, got %d", count1000, b.Len())
			}
			return b
		}(),
	} {
		var key, val []int
		b.TopMin(count100, func(k float64, v int) bool {
			key = append(key, int(k))
			val = append(val, v)
			return true
		})
		if !equalSlices(key, need[:needCount[i]]) {
			t.Errorf("expected %v, got %v", need[:needCount[i]], key)
		}
		if !equalSlices(val, need[:needCount[i]]) {
			t.Errorf("expected %v, got %v", need[:needCount[i]], val)
		}
	}
}

// 测试下负数
func Test_Skiplist_TopMin2(t *testing.T) {
	start := -10
	max := 100
	limit := 10
	sl := New[float64, int]()

	need := make([]int, 0, limit)
	for i, l := start, limit; i < max && l > 0; i++ {
		sl.Set(float64(i), i)
		need = append(need, i)
		l--
	}

	got := make([]int, 0, limit)
	sl.TopMin(10, func(k float64, v int) bool {
		got = append(got, int(k))
		return true
	})

	if !equalSlices(need, got) {
		t.Errorf("expected %v, got %v", need, got)
	}
}

// debug, 指定层
func Test_SkipList_SetAndGet_Level(t *testing.T) {

	sl := New[float64, int]()

	keys := []int{5, 8, 10}
	level := []int{2, 3, 5}
	for i, key := range keys {
		sl.InsertInner(float64(key), key, level[i])
	}

	sl.Draw()
	for _, i := range keys {
		v, count, _ := sl.GetWithMeta(float64(i))
		fmt.Printf("get %v count = %v, nodes:%v, level:%v maxlevel:%v\n",
			float64(i),
			count.Total,
			count.Keys,
			count.Level,
			count.MaxLevel)
		if v != i {
			t.Errorf("expected %d, got %d", i, v)
		}
	}
}

// debug, 用的入口函数
func Test_SkipList_SetAndGet2(t *testing.T) {

	sl := New[float64, int]()

	max := 1000
	start := -1
	for i := max; i >= start; i-- {
		sl.Set(float64(i), i)
	}

	sl.Draw()
	for i := start; i < max; i++ {
		v, count, _ := sl.GetWithMeta(float64(i))
		fmt.Printf("get %v count = %v, nodes:%v, level:%v maxlevel:%v\n",
			float64(i),
			count.Total,
			count.Keys,
			count.Level,
			count.MaxLevel)
		if v != i {
			t.Errorf("expected %d, got %d", i, v)
		}
	}
}

// 测试TopMax, 返回最大的几个数据降序返回
func Test_Skiplist_TopMax(t *testing.T) {

	need := [3][]int{}
	count10 := 10
	count100 := 100
	count1000 := 1000
	count := []int{count10, count100, count1000}

	for i := 0; i < len(count); i++ {
		for j, k := count[i]-1, count100-1; j >= 0 && k >= 0; j-- {
			need[i] = append(need[i], j)
			k--
		}
	}

	for i, b := range []*SkipList[float64, int]{
		// btree里面元素 少于 TopMax 需要返回的值
		func() *SkipList[float64, int] {
			b := New[float64, int]()
			for i := 0; i < count10; i++ {
				b.Set(float64(i), i)
			}

			if b.Len() != count10 {
				t.Errorf("expected %d, got %d", count10, b.Len())
			}
			return b
		}(),
		// btree里面元素 等于 TopMax 需要返回的值
		func() *SkipList[float64, int] {

			b := New[float64, int]()
			for i := 0; i < count100; i++ {
				b.Set(float64(i), i)
			}
			if b.Len() != count100 {
				t.Errorf("expected %d, got %d", count100, b.Len())
			}
			return b
		}(),
		// btree里面元素 大于 TopMax 需要返回的值
		func() *SkipList[float64, int] {

			b := New[float64, int]()
			for i := 0; i < count1000; i++ {
				b.Set(float64(i), i)
			}
			if b.Len() != count1000 {
				t.Errorf("expected %d, got %d", count1000, b.Len())
			}
			return b
		}(),
	} {
		var key, val []int
		b.TopMax(count100, func(k float64, v int) bool {
			key = append(key, int(k))
			val = append(val, v)
			return true
		})
		length := cmp.Min(count[i], len(need[i]))
		if !equalSlices(key, need[i][:length]) {
			t.Errorf("expected %v, got %v", need[i][:length], key)
		}
		if !equalSlices(val, need[i][:length]) {
			t.Errorf("expected %v, got %v", need[i][:length], val)
		}
	}
}

func Test_ConcurrentSkipList_InsertGet(t *testing.T) {
	csl := NewConcurrent[int, string]()
	var wg sync.WaitGroup
	count := 1000

	// Concurrent inserts
	wg.Add(count)
	for i := 0; i < count; i++ {
		go func(i int) {
			defer wg.Done()
			csl.Insert(i, fmt.Sprintf("value%d", i))
		}(i)
	}

	wg.Wait()

	// Concurrent gets
	wg.Add(count)
	for i := 0; i < count; i++ {
		go func(i int) {
			defer wg.Done()
			if val, ok := csl.Get(i); !ok || val != fmt.Sprintf("value%d", i) {
				t.Errorf("expected value%d, got %v", i, val)
			}
		}(i)
	}

	wg.Wait()
}

func Test_ConcurrentSkipList_Delete(t *testing.T) {
	csl := NewConcurrent[int, string]()
	var wg sync.WaitGroup
	count := 1000

	// Insert elements
	for i := 0; i < count; i++ {
		csl.Insert(i, fmt.Sprintf("value%d", i))
	}

	// Concurrent deletes
	wg.Add(count)
	for i := 0; i < count; i++ {
		go func(i int) {
			defer wg.Done()
			csl.Delete(i)
		}(i)
	}

	wg.Wait()

	// Verify all elements are deleted
	for i := 0; i < count; i++ {
		if _, ok := csl.Get(i); ok {
			t.Errorf("expected element %d to be deleted", i)
		}
	}
}

func Test_ConcurrentSkipList_Get(t *testing.T) {
	csl := NewConcurrent[int, string]()
	var wg sync.WaitGroup
	count := 1000

	// Insert elements
	for i := 0; i < count; i++ {
		csl.Insert(i, fmt.Sprintf("value%d", i))
	}

	// Concurrent gets
	wg.Add(count)
	for i := 0; i < count; i++ {
		go func(i int) {
			defer wg.Done()
			if val, ok := csl.Get(i); !ok || val != fmt.Sprintf("value%d", i) {
				t.Errorf("expected value%d, got %v", i, val)
			}
		}(i)
	}

	wg.Wait()
}

func Test_ConcurrentSkipList_Range(t *testing.T) {
	csl := NewConcurrent[int, string]()
	count := 1000

	// Insert elements
	for i := 0; i < count; i++ {
		csl.Insert(i, fmt.Sprintf("value%d", i))
	}

	// Range over elements
	elements := make(map[int]string)
	csl.Range(func(score int, value string) bool {
		elements[score] = value
		return true
	})

	// Verify all elements are ranged
	for i := 0; i < count; i++ {
		if val, exists := elements[i]; !exists || val != fmt.Sprintf("value%d", i) {
			t.Errorf("expected value%d, got %v", i, val)
		}
	}
}

func Test_ConcurrentSkipList_Remove(t *testing.T) {
	csl := NewConcurrent[int, string]()
	var wg sync.WaitGroup
	count := 1000

	// Insert elements
	for i := 0; i < count; i++ {
		csl.Insert(i, fmt.Sprintf("value%d", i))
	}

	// Concurrent removes
	wg.Add(count)
	for i := 0; i < count; i++ {
		go func(i int) {
			defer wg.Done()
			csl.Remove(i)
		}(i)
	}

	wg.Wait()

	// Verify all elements are removed
	for i := 0; i < count; i++ {
		if _, ok := csl.Get(i); ok {
			t.Errorf("expected element %d to be removed", i)
		}
	}
}

// 辅助函数，用于比较两个切片是否相等
func equalSlices(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
