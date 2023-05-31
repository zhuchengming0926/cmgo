package main

import "fmt"

func main() {
	s3 := []int{1, 2, 3, 4, 5, 6, 7, 8}
	s4 := s3[3:6]
	fmt.Printf("The length of s4: %d\n", len(s4))   // 6-3=3
	fmt.Printf("The capacity of s4: %d\n", cap(s4)) // 8-3=5
	fmt.Printf("The value of s4: %d\n", s4)         //
	//mysql_pool.InitGormPool()
	//fmt.Println("人兴财旺")
	//f := excelize.NewFile()
	//defer func() {
	//	if err := f.Close(); err != nil {
	//		fmt.Println(err)
	//	}
	//}()
	//// Create a new sheet.
	//index, err := f.NewSheet("Sheet2")
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//// Set value of a cell.
	//f.SetCellValue("Sheet2", "A2", "Hello world.")
	//f.SetCellValue("Sheet1", "B2", 100)
	//// Set active sheet of the workbook.
	//f.SetActiveSheet(index)
	//// Save spreadsheet by the given path.
	//if err := f.SaveAs("Book1.xlsx"); err != nil {
	//	fmt.Println(err)
	//}
}
