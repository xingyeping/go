/*****************************************************************************/
/**
* \file       ptree_test.go
* \author     xingyeping
* \date       2016/11/02
* \brief      Ptree，用于路由的存储管理
******************************************************************************/
package ptree

import "testing"

func TestIPv4(t *testing.T) {
	var table *Ptree

	//table.Init(32)
	table = PtreeNew(32)

	prefix1 := []byte{5, 1, 1, 1}
	prefix2 := []byte{6, 1, 1, 1}
	prefix3 := []byte{7, 1, 1, 1}
	prefix4 := []byte{30, 30, 30, 30}
	prefix5 := []byte{40, 40, 40, 1}
	prefix6 := []byte{88, 8, 8, 8}
	prefix7 := []byte{101, 0, 0, 2}
	node1 := table.GetNode(prefix1, 32)
	node1.Info = prefix1
	//table.PrintIPv4(nil)

	node2 := table.GetNode(prefix2, 32)
	node2.Info = prefix2
	//table.PrintIPv4(nil)

	node3 := table.GetNode(prefix2, 24)
	node3.Info = prefix2
	//table.PrintIPv4(nil)

	node4 := table.GetNode(prefix3, 32)
	node4.Info = prefix3

	node5 := table.GetNode(prefix4, 32)
	node5.Info = prefix4

	node6 := table.GetNode(prefix4, 24)
	node6.Info = prefix4

	node7 := table.GetNode(prefix5, 32)
	node7.Info = prefix5

	node8 := table.GetNode(prefix6, 32)
	node8.Info = prefix6

	node9 := table.GetNode(prefix7, 32)
	node9.Info = prefix7

	node10 := table.GetNode(prefix7, 24)
	node10.Info = prefix7

	nodeNum, infoNum := table.PrintIPv4(nil)
	t.Logf("TotalNode:%d,TotalInfo:%d", nodeNum, infoNum)
	if infoNum != 10 {
		t.Fatalf("Error")
	}

	//Get same node
	node11 := table.GetNode(prefix1, 32)
	if node11 != node1 {
		t.Fatalf("Get same node failed")
	} else {
		t.Logf("Get same node success")
	}
	//LookUpNode found
	node12 := table.LookUpNode(prefix2, 32)
	if node12 != node2 {
		t.Fatalf("LookUpNode found failed")
	} else {
		t.Logf("LookUpNode found success")
	}
	//LookUpNode found, but no info
	prefix13 := []byte{0, 0, 0, 0}
	node13 := table.LookUpNode(prefix13, 2)
	if node13 != nil {
		t.Fatalf("LookUpNode found, but no info failed")
	} else {
		t.Logf("LookUpNode found, but no info success")
	}
	//LookUpNodeEvenNoInfo found, no info
	node14 := table.LookUpNodeEvenNoInfo(prefix13, 2)
	if node14 == nil {
		t.Fatalf("LookUpNodeEvenNoInfo found, no info, failed")
	} else {
		t.Logf("LookUpNodeEvenNoInfo found, no info, success")
	}

	//MatchNode found
	prefix15 := []byte{101, 0, 0, 1}
	node15 := table.MatchNode(prefix15, 32)
	if node15 != node10 {
		t.Fatalf("MatchNode found failed")
	} else {
		t.Logf("MatchNode found success")
	}
	//Delete node with only one child
	node6.PtreeNodeDelete()
	node16 := table.LookUpNode(prefix4, 24)
	nodeNum, infoNum = table.PrintIPv4(nil)
	if node16 != nil && infoNum != 9 {
		t.Fatalf("Delete node with only one child failed")
	} else {
		t.Logf("Delete node with only one child success")
	}

	//Add a mask 27
	prefix17 := []byte{7, 1, 1, 0}
	node17 := table.GetNode(prefix17, 27)
	if node17.PtreeNodeGetParent().PtreeNodeGetLeftChild() != node3 {
		t.Fatalf("Add a mask 27 failed")
	} else {
		t.Logf("Add a mask 27 success")
	}

}
