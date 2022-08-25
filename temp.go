// scl1 := []int{400, 600, 100, 300, 500, 200, 900}
// sort.Ints(scl1)
// fmt.Println("Sorted: " + fmt.Sprint(scl1))
// sort.Sort(sort.Reverse(sort.IntSlice(scl1)))
// fmt.Println("Reversed sort: " + fmt.Sprint(scl1))

//  // sort structs ascending
//  sort.Slice(fruits, less)aa
//  fmt.Println(fruits)

//  // sort structs descending
//  sort.Slice(fruits, reverse(less))
//  fmt.Println(fruits)

// func reverse(less func(i, j int) bool) func(i, j int) bool {
// 	return func(i, j int) bool {
// 		return !less(i, j)
// 	}
// }

// f, err := os.Create("data.txt")
// if err != nil {
// 	log.Fatal(err)
// }
// defer f.Close()

// _, err2 := f.WriteString("old falcon2\n")
// if err2 != nil {
// 	log.Fatal(err2)
// }
// fmt.Println("done")
// return