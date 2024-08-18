package main

import (
	"fmt"
	"github.com/exsql-io/platform/pkg/lib/sqlparser"
	"github.com/exsql-io/platform/pkg/lib/substrait"
	"google.golang.org/protobuf/encoding/protojson"
	"log"
)

func main() {
	_, err := SubstraitExtensions.ReadDir("third_party/substrait/extensions")
	if err != nil {
		log.Fatal(err)
	}

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

	substraitPlan := substrait.ConvertToSubstrait(plan)
	fmt.Printf("Generated Substrait Plan: %s\n", protojson.Format(substraitPlan))
}
