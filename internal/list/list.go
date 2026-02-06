package list

type Node[K comparable, V any] struct {
	Key  K
	Val  V
	next *Node[K, V]
	prev *Node[K, V]
}

type LinkedList[K comparable, V any] struct {
	head *Node[K, V]
	tail *Node[K, V]
}

func (ll *LinkedList[K, V]) AddToHead(key K, val V) *Node[K, V] {
	newNode := &Node[K, V]{Key: key, Val: val}
	if ll.head == nil {
		ll.head = newNode
		ll.tail = ll.head
	} else {
		newNode.next = ll.head
		ll.head.prev = newNode
		ll.head = newNode
	}
	return newNode
}

func (ll *LinkedList[K, V]) MoveToHead(node *Node[K, V]) {
	if node == ll.head {
		return
	}

	// Remove from current position
	if node == ll.tail {
		ll.tail = node.prev
		ll.tail.next = nil
	} else {
		node.prev.next = node.next
		node.next.prev = node.prev
	}

	// Add to head
	node.prev = nil
	node.next = ll.head
	ll.head.prev = node
	ll.head = node
}

func (ll *LinkedList[K, V]) RemoveTail() *Node[K, V] {
	if ll.tail == nil {
		return nil
	}

	removed := ll.tail

	if ll.tail.prev == nil {
		// Only one node in the list
		ll.head = nil
		ll.tail = nil
	} else {
		ll.tail = ll.tail.prev
		ll.tail.next = nil
	}

	return removed
}
