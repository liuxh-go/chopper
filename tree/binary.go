package tree

// BinaryNode 二叉树节点
type BinaryNode[T any] struct {
	node        *Node[T]
	left, right *BinaryNode[T]
}

// NewBinaryTree 新建二叉树
func NewBinaryTree[T any](list ...T) *BinaryNode[T] {
	if len(list) == 0 {
		return new(BinaryNode[T])
	}

	nodeList := make([]*BinaryNode[T], len(list))
	for i, t := range list {
		nodeList[i] = &BinaryNode[T]{
			node: &Node[T]{
				data: t,
			},
		}
	}

	for i := 0; i <= len(list)/2-1; i++ {
		nodeList[i].left = nodeList[2*i+1]

		if 2*i+2 > len(list)-1 {
			break
		}
		nodeList[i].right = nodeList[2*i+2]
	}

	return nodeList[0]
}

// DLR 先序遍历
func (bn *BinaryNode[T]) DLR(fn func(T)) {
	if bn.node == nil {
		return
	}

	fn(bn.node.data)

	if bn.left != nil {
		bn.left.DLR(fn)
	}

	if bn.right != nil {
		bn.right.DLR(fn)
	}
}

// LDR 中序遍历
func (bn *BinaryNode[T]) LDR(fn func(T)) {
	if bn.node == nil {
		return
	}

	if bn.left != nil {
		bn.left.LDR(fn)
	}

	fn(bn.node.data)

	if bn.right != nil {
		bn.right.LDR(fn)
	}
}

// LRD 后序遍历
func (bn *BinaryNode[T]) LRD(fn func(T)) {
	if bn.node == nil {
		return
	}

	if bn.left != nil {
		bn.left.LRD(fn)
	}

	if bn.right != nil {
		bn.right.LRD(fn)
	}

	fn(bn.node.data)
}

// AddNode 添加节点
func (bn *BinaryNode[T]) AddNode(t T) *BinaryNode[T] {
	//	TODO
	return bn
}
