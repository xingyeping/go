/*****************************************************************************/
/**
* \file       ptree.go
* \author     xingyeping
* \date       2016/11/01
* \version    V1
* \brief      Ptree,used for route information database
******************************************************************************/
package ptree

import (
	"fmt"
)

//Ptree inner node
type PtreeNode struct {
	left, right, parent *PtreeNode
	tree                *Ptree
	keyLen              uint32
	key                 []byte
	Info                interface{} //User Information
}

type Ptree struct {
	root      *PtreeNode
	maxKeyLen uint32
	pType     uint32 //Ptree Type
}

const (
	PtreeIPv4 = 1
	PtreeIPv6 = 2
)

const (
	PtreeMinKeyLen = 1
	PtreeMaxKeyLen = 16
)

var maskbit = []byte{0x00, 0x80, 0xc0, 0xe0, 0xf0, 0xf8, 0xfc, 0xfe, 0xff}

func Max(v1 uint32, v2 uint32) uint32 {
	if v1 >= v2 {
		return v1
	} else {
		return v2
	}
}

func Min(v1 uint32, v2 uint32) uint32 {
	if v1 <= v2 {
		return v1
	} else {
		return v2
	}

}

func (t *Ptree) Init(maxKeyLen uint32) {
	t.maxKeyLen = maxKeyLen
}

func PtreeNew(maxKeyLen uint32) *Ptree {
	var tree *Ptree
	tree = new(Ptree)
	tree.maxKeyLen = maxKeyLen
	return tree
}

func (t *Ptree) Free(f func(interface{})) {
	if t == nil {
		return
	}
	var tmp *PtreeNode
	node := t.root
	for node != nil {
		if node.left != nil {
			node = node.left
			continue
		}
		if node.right != nil {
			node = node.right
			continue
		}
		tmp = node
		node = node.parent
		if node != nil {
			if node.left == tmp {
				node.left = nil
			} else {
				node.right = nil
			}
			if f != nil {
				f(tmp.Info)
				tmp.Info = nil
			}
		} else {
			if f != nil {
				f(node.Info)
				node.Info = nil
			}
		}
	}
	t.root = nil
}

func (t *Ptree) GetNode(key []byte, keyLen uint32) *PtreeNode {
	var newNode, node, match *PtreeNode

	if keyLen > t.maxKeyLen {
		return nil
	}

	node = t.root
	for node != nil && node.keyLen <= keyLen && ptreeKeyMatch(node.key, node.keyLen, key, keyLen) {
		if node.keyLen == keyLen {
			return node
		}
		match = node
		node = node.ptreeGetChild(key, node.keyLen)
	}

	if node == nil {
		newNode = ptreeNodeSet(t, key, keyLen)
		if match != nil {
			ptreeSetLink(match, newNode)
		} else {
			t.root = newNode
		}
	} else {
		newNode = ptreeNodeCommon(node, key, keyLen)
		if newNode == nil {
			return nil
		}
		newNode.tree = t
		ptreeSetLink(newNode, node)

		if match != nil {
			ptreeSetLink(match, newNode)
		} else {
			t.root = newNode
		}

		if newNode.keyLen != keyLen {
			match = newNode
			newNode = ptreeNodeSet(t, key, keyLen)
			ptreeSetLink(match, newNode)
		}
	}
	return newNode
}

func (t *Ptree) LookUpNode(key []byte, keyLen uint32) *PtreeNode {
	if keyLen > t.maxKeyLen {
		return nil
	}
	node := t.root
	for (node != nil) && (node.keyLen <= keyLen) && (ptreeKeyMatch(node.key, node.keyLen, key, keyLen)) {
		if node.keyLen == keyLen && node.Info != nil {
			return node
		}
		node = node.ptreeGetChild(key, node.keyLen)
	}
	return nil
}

func (t *Ptree) LookUpNodeEvenNoInfo(key []byte, keyLen uint32) *PtreeNode {
	if keyLen > t.maxKeyLen {
		return nil
	}
	node := t.root
	for node != nil && node.keyLen <= keyLen && ptreeKeyMatch(node.key, node.keyLen, key, keyLen) {
		if node.keyLen == keyLen {
			return node
		}
		node = node.ptreeGetChild(key, node.keyLen)
	}
	return nil
}
func (t *Ptree) MatchNode(key []byte, keyLen uint32) *PtreeNode {
	var node, matched *PtreeNode
	if keyLen > t.maxKeyLen {
		return nil
	}
	node = t.root
	for node != nil && node.keyLen <= keyLen {
		if node.Info == nil {
			node = node.ptreeGetChild(key, node.keyLen)
			continue
		}
		if !ptreeKeyMatch(node.key, node.keyLen, key, keyLen) {
			break
		}
		matched = node
		node = node.ptreeGetChild(key, node.keyLen)
	}
	return matched
}

func (t *Ptree) GetRoot() *PtreeNode {
	return t.root
}

func (t *Ptree) PrintIPv4(f func(interface{})) (totalNode, totalInfo uint32) {
	var node *PtreeNode
	var nodeNum, infoNum uint32
	for node = t.GetRoot(); node != nil; node = node.PtreeNodeNext() {
		infoNum += node.ptreeNodePrintIPv4(f)
		nodeNum++
	}
	fmt.Printf("TotalNode:%d,TotalInfo:%d\n", nodeNum, infoNum)
	return nodeNum, infoNum
}

func ptreeNodePrintIPv4Key(key []byte, keyLen uint32) {
	if keyLen == 0 {
		fmt.Printf("0.0.0.0/0\n")
	} else if keyLen <= 8 {
		fmt.Printf("%d.0.0.0/%d\n", key[0], keyLen)
	} else if keyLen <= 16 {
		fmt.Printf("%d.%d.0.0/%d\n", key[0], key[1], keyLen)
	} else if keyLen <= 24 {
		fmt.Printf("%d.%d.%d.0/%d\n", key[0], key[1], key[2], keyLen)
	} else {
		fmt.Printf("%d.%d.%d.%d/%d\n", key[0], key[1], key[2], key[3], keyLen)
	}
}

func (node *PtreeNode) ptreeNodePrintIPv4(f func(interface{})) uint32 {
	var infoNum uint32
	fmt.Printf("node:%p,left:%p,right:%p,parent:%p,",
		node, node.left, node.right, node.parent)
	ptreeNodePrintIPv4Key(node.key, node.keyLen)
	if node.Info != nil {
		infoNum = 1
		if f != nil {
			f(node.Info)
		}
	}
	return infoNum
}

func (node *PtreeNode) PtreeNodeDelete() {
	var child, parent *PtreeNode
	if node.left != nil && node.right != nil {
		return
	}
	if node.left != nil {
		child = node.left
	} else {
		child = node.right
	}

	parent = node.parent

	if child != nil {
		child.parent = parent
	}

	if parent != nil {
		if parent.left == node {
			parent.left = child
		} else {
			parent.right = child
		}
	} else {
		node.tree.root = child
	}

	//if parent is not used, then we also delete it
	if (parent.Info == nil) && (parent.left == nil) && (parent.right == nil) {
		parent.PtreeNodeDelete()
	}
}

func (node *PtreeNode) PtreeNodeNext() *PtreeNode {
	if node.left != nil {
		return node.left
	}
	if node.right != nil {
		return node.right
	}

	for node.parent != nil {
		if node.parent.left == node && node.parent.right != nil {
			return node.parent.right
		}
		node = node.parent
	}
	return nil
}

func (node *PtreeNode) PtreeNodeNextUntil(limit *PtreeNode) *PtreeNode {
	if node.left != nil {
		return node.left
	}
	if node.right != nil {
		return node.right
	}

	for node.parent != nil && node != limit {
		if node.parent.left == node && node.parent.right != nil {
			return node.parent.right
		}
		node = node.parent
	}
	return nil
}

func PtreeBitToOctets(keyLen uint32) uint32 {
	return Max((keyLen+7)/8, PtreeMinKeyLen)
}

func ptreeKeyCopy(node *PtreeNode, key []byte, keyLen uint32) {
	var octets uint32
	if keyLen == 0 {
		return
	}
	octets = PtreeBitToOctets(keyLen)
	node.key = make([]byte, octets+1)
	copy(node.key, key)
	node.keyLen = keyLen
}

func ptreeKeyMatch(np []byte, nLen uint32, pp []byte, pLen uint32) bool {
	var shift, offset, i uint32
	if nLen > pLen {
		return false
	}
	offset = Min(nLen, pLen) / 8
	shift = Min(nLen, pLen) % 8

	if shift > 0 {
		if (maskbit[shift] & (np[offset] ^ pp[offset])) > 0 {
			return false
		}
	}
	for i = 0; i < offset; i++ {
		if np[i] != pp[i] {
			return false
		}
	}
	return true
}

func ptreeCheckBit(key []byte, keyLen uint32) bool {
	var offset, shift uint32
	offset = keyLen / 8
	shift = 7 - (keyLen % 8)
	if (key[offset] >> shift & 1) > 0 {
		return true
	}
	return false
}

func (node *PtreeNode) ptreeGetChild(key []byte, keyLen uint32) *PtreeNode {
	if ptreeCheckBit(key, keyLen) {
		return node.right
	}
	return node.left
}

func (node *PtreeNode) PtreeNodeGetParent() *PtreeNode {
	return node.parent
}

func (node *PtreeNode) PtreeNodeGetLeftChild() *PtreeNode {
	return node.left
}

func (node *PtreeNode) PtreeNodeGetRightChild() *PtreeNode {
	return node.right
}

func ptreeNodeCreate(keyLen uint32) *PtreeNode {
	node := new(PtreeNode)
	octets := PtreeBitToOctets(keyLen)
	node.keyLen = keyLen
	node.key = make([]byte, octets+1)
	return node
}

func ptreeNodeSet(tree *Ptree, key []byte, keyLen uint32) *PtreeNode {
	node := ptreeNodeCreate(keyLen)
	ptreeKeyCopy(node, key, keyLen)
	node.tree = tree
	return node
}

func ptreeSetLink(node *PtreeNode, newNode *PtreeNode) {
	if ptreeCheckBit(newNode.key, node.keyLen) {
		node.right = newNode
	} else {
		node.left = newNode
	}
	newNode.parent = node
}

func ptreeNodeCommon(node *PtreeNode, pp []byte, pLen uint32) *PtreeNode {
	var i, j, keyLen uint32
	var diff, mask, boundary byte
	var newNode *PtreeNode

	np := node.key
	for i = 0; i < pLen/8; i++ {
		if np[i] != pp[i] {
			break
		}
	}

	keyLen = i * 8

	if keyLen != pLen {
		diff = np[i] ^ pp[i]
		mask = 0x80
		for (keyLen < pLen) && ((mask & diff) == 0) {
			if boundary == 0 {
				boundary = 1
			}
			mask >>= 1
			keyLen++
		}
	}

	newNode = ptreeNodeCreate(keyLen)
	newp := newNode.key

	for j = 0; j < i; j++ {
		newp[j] = np[j]
	}
	if boundary > 0 {
		newp[j] = np[j] & maskbit[newNode.keyLen%8]
	}
	return newNode
}
