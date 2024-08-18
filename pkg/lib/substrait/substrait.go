package substrait

import (
	"fmt"
	substraitpb "github.com/exsql-io/platform/pkg/lib/proto/substrait" // Replace with the actual path
	"strings"
	"vitess.io/vitess/go/vt/sqlparser"
)

// ConvertToSubstrait converts a Vitess AST to a Substrait plan
func ConvertToSubstrait(stmt sqlparser.Statement) *substraitpb.Plan {
	switch stmt := stmt.(type) {
	case *sqlparser.Select:
		return convertSelect(stmt)
	default:
		fmt.Println("Unsupported SQL statement type")
		return nil
	}
}

// convertSelect handles the conversion of a SELECT statement
func convertSelect(sel *sqlparser.Select) *substraitpb.Plan {
	var columns []string
	for _, expr := range sel.SelectExprs {
		if col, ok := expr.(*sqlparser.AliasedExpr); ok {
			if colName, ok := col.Expr.(*sqlparser.ColName); ok {
				columns = append(columns, colName.Name.String())
			}
		}
	}

	tableName := ""
	if len(sel.From) > 0 {
		if tableExpr, ok := sel.From[0].(*sqlparser.AliasedTableExpr); ok {
			if tableNameExpr, ok := tableExpr.Expr.(sqlparser.TableName); ok {
				tableName = tableNameExpr.Name.String()
			}
		}
	}

	var conditions *substraitpb.Expression
	if sel.Where != nil {
		conditions = extractConditions(sel.Where.Expr)
	}

	// Construct the Substrait plan
	plan := &substraitpb.Plan{
		Relations: []*substraitpb.PlanRel{
			{
				RelType: &substraitpb.PlanRel_Root{
					Root: &substraitpb.RelRoot{
						Names: columns,
						Input: &substraitpb.Rel{
							RelType: &substraitpb.Rel_Read{
								Read: &substraitpb.ReadRel{
									BaseSchema: &substraitpb.NamedStruct{
										Names: columns,
									},
									ReadType: &substraitpb.ReadRel_NamedTable_{
										NamedTable: &substraitpb.ReadRel_NamedTable{
											Names: strings.Split(tableName, "."),
										},
									},
									Filter: conditions,
								},
							},
						},
					},
				},
			},
		},
	}

	return plan
}

// extractConditions recursively extracts conditions from the WHERE clause
func extractConditions(expr sqlparser.Expr) *substraitpb.Expression {
	switch condition := expr.(type) {
	case *sqlparser.ComparisonExpr:
		return &substraitpb.Expression{
			RexType: &substraitpb.Expression_ScalarFunction_{
				ScalarFunction: &substraitpb.Expression_ScalarFunction{
					FunctionReference: 0,
					Arguments: []*substraitpb.FunctionArgument{
						functionArgument(condition.Left),
						functionArgument(condition.Right),
					},
					OutputType: &substraitpb.Type{
						Kind: &substraitpb.Type_Bool{
							Bool: &substraitpb.Type_Boolean{
								Nullability: substraitpb.Type_NULLABILITY_REQUIRED,
							},
						},
					},
				},
			},
		}
	case *sqlparser.AndExpr:
		return &substraitpb.Expression{}
	default:
		panic(fmt.Sprintf("unknown expression type: %T", expr))
	}
}

func functionArgument(expr sqlparser.Expr) *substraitpb.FunctionArgument {
	switch expr.(type) {
	case *sqlparser.Literal:
	case *sqlparser.ColName:
		return &substraitpb.FunctionArgument{
			ArgType: &substraitpb.FunctionArgument_Value{
				Value: &substraitpb.Expression{
					RexType: &substraitpb.Expression_Selection{
						Selection: &substraitpb.Expression_FieldReference{
							ReferenceType: &substraitpb.Expression_FieldReference_DirectReference{
								DirectReference: &substraitpb.Expression_ReferenceSegment{
									ReferenceType: &substraitpb.Expression_ReferenceSegment_StructField_{
										StructField: &substraitpb.Expression_ReferenceSegment_StructField{
											Field: 0,
										},
									},
								},
							},
							RootType: &substraitpb.Expression_FieldReference_RootReference_{
								RootReference: &substraitpb.Expression_FieldReference_RootReference{},
							},
						},
					},
				},
			},
		}
	}

	return &substraitpb.FunctionArgument{}
}
