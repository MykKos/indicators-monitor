package influx

import "sync"

type Queue struct {
	sync.Mutex
	Items []interface{}
}

func (q *Queue) Pop() interface{} {
	q.Lock()
	defer q.Unlock()
	if len(q.Items) == 0 {
		return nil
	}
	item := q.Items[0]
	q.Items = q.Items[1:]
	return item
}

func (q *Queue) Push(item interface{}) {
	q.Lock()
	defer q.Unlock()
	q.Items = append(q.Items, item)
}

func (q *Queue) PushLocked(item interface{}) {
	q.Lock()
	defer q.Unlock()
	q.Items = append(q.Items, item)
}

func (q *Queue) PopMultiple(amount int) (items []interface{}) {
	q.Lock()
	defer q.Unlock()
	if len(q.Items) == 0 {
		return nil
	}
	if len(q.Items) < amount {
		items = q.Items
		q.Items = nil
	} else {
		items = q.Items[:amount]
		q.Items = q.Items[amount:]
	}
	return items
}
