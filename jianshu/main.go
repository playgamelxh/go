package main

import "fmt"
import "database/sql"
import _ "./mysql"

func main() {

	fmt.Printf("Hello, jianshu")

	db, err := sql.Open("mysql", "root:123456@/jianshu")

	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	//insert
	stmtIns, err := db.Prepare("INSERT INTO squareNum VALUES( ?, ? )")
	if err != nil {
		panic(err.Error())
	}
	defer stmtIns.Close()

	//select
	stmtOut, err := db.Prepare("SELECT squareNumber FROM squarenum WHERE number = ?")
	if err != nil {
		panic(err.Error())
	}
	defer stmtOut.Close()

	// Insert square numbers for 0-24 in the database
	for i := 0; i < 25; i++ {
		_, err = stmtIns.Exec(i, (i * i)) // Insert tuples (i, i^2)
		if err != nil {
			panic(err.Error())
		}
	}

	var squareNum int // we "scan" the result in here

	// Query the square-number of 13
	err = stmtOut.QueryRow(13).Scan(&squareNum) // WHERE number = 13
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	fmt.Printf("The square number of 13 is: %d", squareNum)

	// Query another number.. 1 maybe?
	err = stmtOut.QueryRow(1).Scan(&squareNum) // WHERE number = 1
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	fmt.Printf("The square number of 1 is: %d", squareNum)

}
