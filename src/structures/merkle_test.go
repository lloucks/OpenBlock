package structures

/*
func Test_merkle_tree(t *testing.T) {
	fmt.Printf("Starting merkle tree test\n")
	var transactions []Transaction

	t1 := CreateTransaction("hi hi hi", 1)
	transactions = append(transactions, *t1)
	t2 := CreateTransaction("how are you", 2)
	transactions = append(transactions, *t2)
	t3 := CreateTransaction("I am good", 1)
	transactions = append(transactions, *t3)
	t4 := CreateTransaction("how are you?", 2)
	transactions = append(transactions, *t4)
	t5 := CreateTransaction("Terrible", 1)
	transactions = append(transactions, *t5)
	t6 := CreateTransaction("Do you own a PS7?", 3)
	transactions = append(transactions, *t6)
	t7 := CreateTransaction("It was just released yesterday", 3)
	transactions = append(transactions, *t7)
	t8 := CreateTransaction("Where can you purchase it?", 4)
	transactions = append(transactions, *t8)
	t9 := CreateTransaction("Cant find it anywhere", 4)
	transactions = append(transactions, *t9)

	mTree1 := CreateMerkleTree(1, transactions)
	fmt.Printf("Root Hash: %d\n", mTree1.Root.HashedData)
	treeDiff, err := mTree1.TreeDifference()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Difference in tree: %t\n", treeDiff)

	faultyTransation := CreateTransaction("Leaderboard scores::", 3)
	result, err := mTree1.VerifyTransaction(faultyTransation)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Trusted transaction: %t\n", result)

	mTree1 = mTree1.AddTransaction(faultyTransation)
	result2, err := mTree1.VerifyTransaction(faultyTransation)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Trusted transaction: %t\n", result2)
	result3, err := mTree1.VerifyTransaction(t5)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Trusted transaction: %t\n", result3)
	b := mTree1.Serialize()
	fmt.Println(b)


	attackNumber := (len(mTree1.Leafs) / 2)
	oldNode := mTree1.Leafs[attackNumber]
	faultyTransaction := CreateTransaction("Faulty Transaction", 4)
	faultyNode := CreateMerkleNode(oldNode.Index, nil, nil,
		faultyTransaction.ID, true)
	faultyNode.Parent = oldNode.Parent
	mTree1.Leafs[attackNumber] = faultyNode
	result4, err := mTree1.VerifyTransaction(faultyTransaction)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Trusted transaction: %t\n", result4)

}
*/
