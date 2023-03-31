package tree

import (
	"strings"

	"github.com/liuxh-go/chopper/bytesconv"
	"github.com/liuxh-go/chopper/math"
)

/*
	路径树
	ex:
	1->(1)2->(12)3
		  |->(12)4
	2->(2)1->(21)3
*/

type nodeType uint8

const (
	static nodeType = iota
	root
	catchAll
)

// PathNode 路径节点
type PathNode[T any] struct {
	path      string
	indices   string
	wildChild bool
	nType     nodeType
	priority  uint32
	childList []*PathNode[T]
	fullPath  string
	data      *T
}

// AddNode 添加节点
func (pn *PathNode[T]) AddNode(path string, t *T) {
	fullPath := path
	pn.priority++

	// empty tree
	if len(pn.path) == 0 && len(pn.childList) == 0 {
		pn.insertChild(path, fullPath, t)
		pn.nType = root
		return
	}

	parentFullPathIndex := 0

walk:
	for {
		// 找到新节点和老节点的最长前缀
		i := longestCommonPrefix(path, pn.path)

		// 前缀小于老节点的路径，可以附加到老节点后
		if i < len(pn.path) {
			child := &PathNode[T]{
				path:      pn.path[i:],
				wildChild: pn.wildChild,
				nType:     static,
				indices:   pn.indices,
				priority:  pn.priority - 1,
				childList: pn.childList,
				fullPath:  pn.fullPath,
				data:      pn.data,
			}

			pn.childList = []*PathNode[T]{child}
			pn.indices = bytesconv.BytesToString([]byte{pn.path[i]})
			pn.path = path[:i]
			pn.wildChild = false
			pn.data = nil
			pn.fullPath = fullPath[:parentFullPathIndex+i]
		}

		// 新建子节点
		if i < len(path) {
			path = path[i:]
			c := path[0]

			// Check if a child with the next path byte exists
			for i, max := 0, len(pn.indices); i < max; i++ {
				if c == pn.indices[i] {
					parentFullPathIndex += len(pn.path)
					i = pn.incrementChildPrio(i)
					pn = pn.childList[i]
					continue walk
				}
			}

			// insert node
			if c != '*' && pn.nType != catchAll {
				pn.indices += bytesconv.BytesToString([]byte{c})
				child := &PathNode[T]{
					fullPath: fullPath,
				}
				pn.addNode(child)
				pn.incrementChildPrio(len(pn.indices) - 1)
				pn = child
			} else if pn.wildChild {
				pn = pn.childList[len(pn.childList)-1]
				pn.priority++

				// Check if the wildcard matches
				if len(path) >= len(pn.path) && pn.path == path[:len(pn.path)] &&
					// Adding a child to a catchAll is not possible
					pn.nType != catchAll &&
					// Check for longer wildcard, e.g. :name and :names
					(len(pn.path) >= len(path) || path[len(pn.path)] == '/') {
					continue walk
				}

				// Wildcard conflict
				pathSeg := path
				if pn.nType != catchAll {
					pathSeg = strings.SplitN(pathSeg, "/", 2)[0]
				}
				prefix := fullPath[:strings.Index(fullPath, pathSeg)] + pn.path
				panic("'" + pathSeg +
					"' in new path '" + fullPath +
					"' conflicts with existing wildcard '" + pn.path +
					"' in existing prefix '" + prefix +
					"'")
			}

			pn.insertChild(path, fullPath, t)
			return
		}

		// Otherwise add handle to current node
		if pn.data != nil {
			panic("data are already registered for path '" + fullPath + "'")
		}
		pn.data = t
		pn.fullPath = fullPath
		return
	}
}

func findWildcard(path string) (wildcard string, i int, valid bool) {
	// Find start
	for start, c := range []byte(path) {
		// A wildcard starts with '*' (catch-all)
		if c != '*' {
			continue
		}

		// Find end and check for invalid characters
		valid = true
		for end, c := range []byte(path[start+1:]) {
			switch c {
			case '/':
				return path[start : start+1+end], start, valid
			case '*':
				valid = false
			}
		}
		return path[start:], start, valid
	}
	return "", -1, false
}

func (pn *PathNode[T]) insertChild(path, fullPath string, t *T) {
	for {
		wildcard, i, valid := findWildcard(path)
		if i < 0 {
			break
		}

		// The wildcard name must only contain one '*' character
		if !valid {
			panic("only one wildcard per path segment is allowed, has: '" +
				wildcard + "' in path '" + fullPath + "'")
		}

		// check if the wildcard has a name
		if len(wildcard) < 2 {
			panic("wildcards must be named with a non-empty name in path '" + fullPath + "'")
		}

		if i+len(wildcard) != len(path) {
			panic("catch-all routes are only allowed at the end of the path in path '" + fullPath + "'")
		}

		if len(pn.path) > 0 && pn.path[len(pn.path)-1] == '/' {
			pathSeg := strings.SplitN(pn.childList[0].path, "/", 2)[0]
			panic("catch-all wildcard '" + path +
				"' in new path '" + fullPath +
				"' conflicts with existing path segment '" + pathSeg +
				"' in existing prefix '" + pn.path + pathSeg +
				"'")
		}

		// currently fixed width 1 for '/'
		i--
		if path[i] != '/' {
			panic("no / before catch-all in path '" + fullPath + "'")
		}

		pn.path = path[:i]

		child := &PathNode[T]{
			wildChild: true,
			nType:     catchAll,
			fullPath:  fullPath,
		}

		pn.addNode(child)
		pn.indices = string('/')
		pn = child
		pn.priority++

		child = &PathNode[T]{
			path:     path[i:],
			nType:    catchAll,
			data:     t,
			priority: 1,
			fullPath: fullPath,
		}
		pn.childList = []*PathNode[T]{child}

		return
	}

	pn.path = path
	pn.data = t
	pn.fullPath = fullPath
}

func (pn *PathNode[T]) incrementChildPrio(pos int) int {
	cs := pn.childList
	cs[pos].priority++
	prio := cs[pos].priority

	newPos := pos
	for ; newPos > 0 && cs[newPos-1].priority < prio; newPos-- {
		cs[newPos-1], cs[newPos] = cs[newPos], cs[newPos-1]
	}

	if newPos != pos {
		pn.indices = pn.indices[:newPos] +
			pn.indices[pos:pos+1] +
			pn.indices[newPos:pos] + pn.indices[pos+1:]
	}

	return newPos
}

func (pn *PathNode[T]) addNode(child *PathNode[T]) {
	if pn.wildChild && len(pn.childList) > 0 {
		wildcardChild := pn.childList[len(pn.childList)-1]
		pn.childList = append(pn.childList[:len(pn.childList)-1], child, wildcardChild)
	} else {
		pn.childList = append(pn.childList, child)
	}
}

func longestCommonPrefix(s1, s2 string) int {
	i := 0
	max := math.Min(len(s1), len(s2))
	for i < max && s1[i] == s2[i] {
		i++
	}

	return i
}

// NodeValue 节点值对象
type NodeValue[T any] struct {
	Data     *T
	Tsr      bool
	FullPath string
}

type skippedNode[T any] struct {
	path        string
	node        *PathNode[T]
	paramsCount int16
}

// GetValue 获取节点值
func (pn *PathNode[T]) GetValue(path string, skippedNodes *[]skippedNode[T]) (value NodeValue[T]) {
walk:
	for {
		prefix := pn.path
		if len(path) > len(prefix) {
			if path[:len(prefix)] == prefix {
				path = path[len(prefix):]

				idxc := path[0]
				for i, c := range []byte(pn.indices) {
					if c == idxc {
						if pn.wildChild {
							index := len(*skippedNodes)
							*skippedNodes = (*skippedNodes)[:index+1]
							(*skippedNodes)[index] = skippedNode[T]{
								path: prefix + path,
								node: &PathNode[T]{
									path:      pn.path,
									wildChild: pn.wildChild,
									nType:     pn.nType,
									priority:  pn.priority,
									childList: pn.childList,
									fullPath:  pn.fullPath,
									data:      pn.data,
								},
							}
						}

						pn = pn.childList[i]
						continue walk
					}
				}

				if !pn.wildChild {
					if path != "/" {
						for length := len(*skippedNodes); length > 0; length-- {
							skippedNode := (*skippedNodes)[length-1]
							*skippedNodes = (*skippedNodes)[:length-1]
							if strings.HasSuffix(skippedNode.path, path) {
								path = skippedNode.path
								pn = skippedNode.node
								continue walk
							}
						}
					}

					value.Tsr = path == "/" && pn.data != nil
					return
				}

				pn = pn.childList[len(pn.childList)-1]

				switch pn.nType {
				case catchAll:
					value.Data = pn.data
					value.FullPath = pn.fullPath

				default:
					panic("invalid node type")
				}
			}
		}

		if path == prefix {
			if pn.data == nil && path != "/" {
				for length := len(*skippedNodes); length > 0; length-- {
					skippedNode := (*skippedNodes)[length-1]
					*skippedNodes = (*skippedNodes)[:length-1]
					if strings.HasSuffix(skippedNode.path, path) {
						path = skippedNode.path
						pn = skippedNode.node
						continue walk
					}
				}
			}

			if value.Data = pn.data; value.Data != nil {
				value.FullPath = pn.fullPath
				return
			}

			if path == "/" && pn.wildChild && pn.nType != root {
				value.Tsr = true
				return
			}

			if path == "/" && pn.nType == static {
				value.Tsr = true
				return
			}

			for i, c := range []byte(pn.indices) {
				if c == '/' {
					pn = pn.childList[i]
					value.Tsr = (len(pn.path) == 1 && pn.data != nil) ||
						(pn.nType == catchAll && pn.childList[0].data != nil)
					return
				}
			}

			return
		}

		value.Tsr = path == "/" ||
			(len(prefix) == len(path)+1 && prefix[len(path)] == '/' &&
				path == prefix[:len(prefix)-1] && pn.data != nil)

		if !value.Tsr && path != "/" {
			for length := len(*skippedNodes); length > 0; length-- {
				skippedNode := (*skippedNodes)[length-1]
				*skippedNodes = (*skippedNodes)[:length-1]
				if strings.HasSuffix(skippedNode.path, path) {
					path = skippedNode.path
					pn = skippedNode.node
					continue walk
				}
			}
		}

		return
	}
}
