package main

import (
	"container/heap"
	"math"
)

type State interface {
	Expand() []State
	Cost() float64
	Heuristic() float64
	IsGoal() bool
}

type Node struct {
	parent *Node
	state  State
	value  float64
}

type strategy interface {
	Peek() *Node
	Take() *Node
	Add(node *Node)
}

// A PriorityQueue implements heap.Interface and holds Nodes
type PriorityQueue []*Node

func (pq PriorityQueue) Len() int {
	return len(pq)
}

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].value < pq[j].value
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x interface{}) {
	item := x.(*Node)
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n - 1]
	*pq = old[0 : n - 1]
	return item
}

// Piority queue as A* Strategy
func (pq *PriorityQueue) Peek() *Node {
	n := len(*pq)
	if n > 0 {
		return (*pq)[n - 1]
	}
	return nil
}

func (pq *PriorityQueue) Take() *Node {
	return heap.Pop(pq).(*Node)
}

func (pq *PriorityQueue) Add(node *Node) {
	heap.Push(pq, node)
}

// A possibly mutable constraint, returns true if a node is constraint, so it should not be expanded further.
type constraint interface {
	onVisit(node *Node) bool
	onExpand(node *Node) bool
}

func aStar() strategy {
	pq := make(PriorityQueue, 0, 64)
	heap.Init(&pq)
	return &pq
}

type noConstraint bool

func (c noConstraint) onVisit(node *Node) bool {
	return false
}

func (c noConstraint) onExpand(node *Node) bool {
	return false
}

/*
(defn- general-search [state expand-fn h-fn constraint goal-fn the-limit]
  (let [limit ^double the-limit]
    (loop [queue          (:strategy state)
           contour        (Double/POSITIVE_INFINITY)
           visited        (long (get state :visited 0))
           expanded       (long (get state :expanded 0))]
      (when (Thread/interrupted) (throw (InterruptedException.)))
      (if-let [^Node node (s-peek queue)]
        (if (on-visit constraint node)
          (recur (s-pop! queue) contour (inc visited) expanded)
          (let [f-cost (.value node)
                queue  (s-pop! queue)]
              (if (goal-fn (.state node))
                {:node node :contour contour :visited visited :expanded expanded
                 :next-solver #(general-search { :strategy queue :visited visited :expanded expanded} expand-fn h-fn constraint goal-fn limit)}
                (let [moves (expand-fn (.state node))
                      [queue expanded contour]
                        (loop [queue    queue
                               contour  contour
                               expanded (long expanded)
                               moves    moves]
                          (if-let [move (first moves)]
                            (let [childnode ^Node (expand-node node h-fn move)]
                              (if (on-expand constraint childnode)
                                (recur queue contour expanded (next moves))
                                (if (> (.value childnode) limit)
                                  (recur queue (double (min contour (.value childnode))) expanded (next moves))
                                  (recur (s-conj! queue childnode) contour (inc expanded) (next moves)))))
                            [queue expanded contour]))]
                  (recur queue (double contour) (inc visited) (long expanded))))))
        {:node nil :contour contour :visited visited :expanded expanded}))))
 */

type result struct {
	node     *Node
	contour  float64
	visited  int
	expanded int
}

func generalSearch(queue strategy, visited int, expanded int, constr constraint, limit float64) result {
	contour := math.MaxFloat64

	for {
		node := queue.Take()
		if node == nil {
			return result{nil, contour, visited, expanded}
		}
		visited++
		if constr.onVisit(node) {
			continue
		}
		if node.state.IsGoal() {
			return result{node, contour, visited, expanded}
		}
		for _, child := range node.state.Expand() {
			childNode := &Node{node, child, math.Max(node.value, child.Cost() + child.Heuristic())}
			if constr.onExpand(childNode) {
				continue
			}
			if childNode.value > limit {
				contour = math.Min(contour, childNode.value)
				continue
			}
			queue.Add(childNode)
			expanded++

		}

	}
	return result{nil, contour, visited, expanded}
}

type Result struct {
	Solution *Node
	Visited  int
	Expanded int
}

func Solve(rootState State) Result {
	var s strategy
	s = aStar()
	s.Add(&Node{nil, rootState, rootState.Cost() + rootState.Heuristic()})
	var constraint noConstraint

	result := generalSearch(s, 0, 0, constraint, math.MaxFloat64)
	return Result{result.node, result.visited, result.expanded}
}