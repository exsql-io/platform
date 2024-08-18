package main

import (
	"fmt"
	"github.com/exsql-io/platform/pkg/lib/sqlparser"
)

func main() {
	sql := `
		SELECT name, salary
		FROM employees
		WHERE department = 'Engineering' AND active = 'TRUE'
		ORDER BY name
		LIMIT 10;
	`

	sqlParser, err := sqlparser.NewSQLParser()
	if err != nil {
		fmt.Printf("Failed to initialize sqlParser: %v\n", err)
		return
	}

	plan, err := sqlParser.Parse(sql)
	if err != nil {
		fmt.Printf("Failed to parse SQL: %v\n", err)
		return
	}

	fmt.Printf("Generated Substrait Plan: %+v\n", plan)
}
