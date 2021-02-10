package main

import (
	"database/sql"
	"fmt"
	"github.com/fatih/color"
	_ "github.com/lib/pq"
	"strings"
)

const (
	database = "" // Database name
	user     = "" // User name for the database
	password = "" // password for the database
	host     = "" // host name for the database
	port     = 5432 // port number used for the database operations
	tableName = "" // enter table name for the database
)

func main() {
	fmt.Println("Hello, world")
	postgresURL := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, database)
	db, err := sql.Open("postgres", postgresURL)
	if err != nil {
		panic(err)
	}

	var choice int = 9 // dummy assignment
	fmt.Println("Successfully connected!")
	if tableExist(db){
		fmt.Printf("Table is present: %T", db)
	} else {
		fmt.Println("Table is created")
	}
	for choice != 0{
		fmt.Println("Choices are as follow")
		fmt.Println("1. For register employee\t2. For get employee details")
		fmt.Println("3. For update employee data\t4. For print all employees")
		fmt.Println("5. For delete a employee details\t0. For exit()")

		fmt.Println("Enter your choice: ")
		fmt.Scan(&choice)
		switch choice {
		case 1:
			go registerEmployee(db)
		case 2:
			go getEmployeeDetails(db)
		case 3:
			go updateEmployeeDetails(db)
		case 4:
			go getAllEmployeeDetails(db)
		case 5:
			go deleteEmployeeFromDB(db)
		case 0:

		default:
			color.Set(color.FgHiRed)
			fmt.Println("Enter correct choice")
			color.Unset()
		}
	}


	defer db.Close()
	color.Set(color.FgGreen)
	fmt.Println("Program exited successfully...")
	fmt.Println("Connection closed")
	color.Unset()

}

// Schema of table is dependent to the need if you need any additional columns than add them as need
// this program can't handle the case of multiple employees with same id but it give tht appropriate error and
// close the program due to the exception occur
func registerEmployee(db *sql.DB) {
	var id, name, password, salary, age string
	fmt.Print("Enter employee id: ")
	fmt.Scanln(&id)
	fmt.Print("Enter employee name: ")
	fmt.Scanln(&name)
	fmt.Print("Enter employee password: ")
	fmt.Scanln(&password)
	fmt.Print("Enter employee salary: ")
	fmt.Scanln(&salary)
	fmt.Print("Enter employee age: ")
	fmt.Scanln(&age)

	query := fmt.Sprintf("INSERT INTO %s VALUES('%s', '%s', '%s', '%s', '%s');",
		tableName, id, name, password, salary, age)

	_, err := db.Exec(query)
	if err!=nil{
		panic(err)
	}
	color.Set(color.FgHiGreen)
	fmt.Println("Employee is inserted")
	color.Unset()
}
func getEmployeeDetails(db *sql.DB) {
	var id string
	fmt.Print("Enter employee id: ")
	fmt.Scanln(&id)
	query := fmt.Sprintf("SELECT * FROM %s WHERE ID = '%s';", tableName, id)
	//fmt.Println(query)
	rows, err := db.Query(query)
	if err!=nil{
		panic(err)
	}
	var ans string = "["
	//rows.Close()
	for rows.Next() {
		var id, name, password, salary, age string

		err = rows.Scan(&id, &name, &password, &salary, &age)
		if err!=nil{
			panic(err)
		}
		ans += fmt.Sprintf("{id: \"%s\", name: \"%s\", password: \"%s\", salary: \"%s\", age: \"%s\"}",
			id, name, password, salary, age)
	}
	ans += "]"
	color.Set(color.FgCyan)
	fmt.Println(ans)
	color.Unset()
}
func getAllEmployeeDetails(db *sql.DB) {
	query := fmt.Sprintf("SELECT * FROM %s;", tableName)
	rows, err := db.Query(query)
	if err!=nil{
		panic(err)
	}
	var ans string = "["
	//rows.Close()
	for rows.Next() {
		var id, name, pass, salary, age string

		err = rows.Scan(&id, &name, &pass, &salary, &age)
		if err!=nil{
			panic(err)
		}
		ans += fmt.Sprintf("{id: \"%s\", name: \"%s\", password: \"%s\", salary: \"%s\", age: \"%s\"}, ",
			id, name, pass, salary, age)
	}
	ans = ans[:len(ans)-2] + "]"
	color.Set(color.FgCyan)
	fmt.Println(ans)
	color.Unset()
}
func updateEmployeeDetails(db *sql.DB) {
	query := fmt.Sprintf("UPDATE %s SET ", tableName)
	var id, name, pass, salary, age string
	fmt.Scanln()
	fmt.Print("Enter employee id: ")
	fmt.Scanln(&id)
	fmt.Print("Enter Name(if need change) :")
	fmt.Scanln(&name)
	if len(name) != 0{
		query += fmt.Sprintf("NAME = '%s' ", name)
	}
	fmt.Print("Enter Password(if need change) :")
	fmt.Scanln(&pass)
	if len(pass) != 0{
		if query[len(query)-2:] == "' "{
			query = query[:len(query)-1] +", "
		}
		query += fmt.Sprintf("PASSWORD = '%s' ", pass)
	}
	fmt.Print("Enter Salary(if need change) :")
	fmt.Scanln(&salary)
	if len(salary) != 0{
		if query[len(query)-2:] == "' "{
			query = query[:len(query)-1] +", "
		}
		query += fmt.Sprintf("SALARY = '%s' ", salary)
	}
	fmt.Print("Enter Age(if need change) :")
	fmt.Scanln(&age)
	if len(age) != 0{
		if query[len(query)-2:] == "' "{
			query = query[:len(query)-1] +", "
		}
		query += fmt.Sprintf("AGE = '%s' ", age)
	}
	query += fmt.Sprintf("WHERE ID = '%s';", id)
	fmt.Println(query)
	_, err := db.Exec(query)
	if err!=nil{
		panic(err)
	}
	color.Set(color.FgHiGreen)
	fmt.Printf("Employee with ID: %s data updated\n", id)
	color.Unset()
}
func deleteEmployeeFromDB(db *sql.DB) {
	var id string
	fmt.Print("Enter employee id: ")
	fmt.Scanln(&id)

	query := fmt.Sprintf("DELETE FROM %s WHERE ID = '%s';", tableName, id)

	_, err := db.Exec(query)
	if err!=nil{
		panic(err)
	}
	color.Set(color.FgHiGreen)
	fmt.Printf("Employee with ID: %s deleted\n", id)
	color.Unset()
}

func tableExist(db *sql.DB) bool{
	query := fmt.Sprintf("SELECT * FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_NAME = '%s'", strings.ToLower(tableName))
	rows, err := db.Query(query)
	if err!=nil{
		panic(err)
	}

	for rows.Next(){
		return true
	}
	query = fmt.Sprintf("CREATE TABLE %s(ID VARCHAR(10) PRIMARY KEY, NAME VARCHAR(30), PASSWORD VARCHAR(30), SALARY VARCHAR(15), AGE VARCHAR(3) )", tableName)
	_, err = db.Exec(query)

	if err!=nil{
		panic(err)
	}

	return false
}
