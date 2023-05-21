package cycache

// 简易int64队列，先进先出，用于缓存过期检查
type stringQueue struct {
	Data []string
}

func (q *stringQueue) push(v string) { // key不可以为空
	if v == "" {
		return
	}
	q.Data = append(q.Data, v)
}

// 只取值
func (q *stringQueue) popNoDelete() string {
	if len(q.Data) == 0 {
		return ""
	}
	return (q.Data)[0]
}

// 只删除
func (q *stringQueue) popOnlyDelete() {
	if len(q.Data) == 0 {
		return
	}
	q.Data = q.Data[1:]
}

// 清理重复数据
func (q *stringQueue) clearDuplicateData() {
	if len(q.Data) == 0 {
		return
	}

	var newQueue []string
	var keySet = map[string]bool{}
	for i := len(q.Data) - 1; i >= 0; i-- {
		key := (q.Data)[i]
		v := keySet[key]
		if !v {
			newQueue = append(newQueue, key)
			keySet[key] = true
		}
	}
	lenQueue := len(newQueue)
	for i := 0; i < lenQueue/2; i++ {
		newQueue[i], newQueue[lenQueue-i-1] = newQueue[lenQueue-i-1], newQueue[i]
	}
	q.Data = newQueue
}
