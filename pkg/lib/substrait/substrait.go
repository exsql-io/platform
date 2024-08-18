package substrait

import (
	"fmt"
	substraitpb "github.com/exsql-io/platform/pkg/lib/proto/substrait" // Replace with the actual path
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
	columns := []string{}
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

	conditions := []*substraitpb.FilterRel{}
	if sel.Where != nil {
		conditions = extractConditions(sel.Where.Expr)
	}

	// Construct the Substrait plan
	plan := &substraitpb.Plan{
		Roots: []*substraitpb.Plan_RelRoot{
			{
				Input: &substraitpb.Rel{
					RelType: &substraitpb.Rel_Read{
						Read: &substraitpb.ReadRel{
							BaseSchema: &substraitpb.NamedStruct{
								Names: columns,
							},
							TableName: tableName,
						},
					},
				},
			},
		},
	}

	if len(conditions) > 0 {
		plan.Roots[0].Input = &substraitpb.Rel{
			RelType: &substraitpb.Rel_Filter{
				Filter: &substraitpb.FilterRel{
					Input: plan.Roots[0].Input,
					Condition: &substraitpb.Expression{
						ExprType: &substraitpb.Expression_Condition{
							Condition: &substraitpb.Condition{
								Clauses: conditions,
							},
						},
					},
				},
			},
		}
	}

	return plan
}

// extractConditions recursively extracts conditions from the WHERE clause
func extractConditions(expr sqlparser.Expr) *substraitpb.Expression {
	switch node := expr.(type) {
	case *sqlparser.ComparisonExpr:
		return &substraitpb.Expression{
			RexType: &substraitpb.Expression_ScalarFunction_{
				ScalarFunction: &substraitpb.Expression_ScalarFunction{
					FunctionReference: 0,
					Arguments:         nil,
					Options:           nil,
					OutputType:        nil,
				},
			},
		}
	case *sqlparser.AndExpr:
		return &substraitpb.Expression{}
	default:
		panic(fmt.Sprintf("unknown expression type: %T", expr))
	}
}
