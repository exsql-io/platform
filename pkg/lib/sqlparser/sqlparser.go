package sqlparser

import "vitess.io/vitess/go/vt/sqlparser"

type SQLParser struct {
	parser *sqlparser.Parser
}

func NewSQLParser() (*SQLParser, error) {
	parser, err := sqlparser.New(sqlparser.Options{})
	if err != nil {
		return nil, err
	}

	return &SQLParser{
		parser: parser,
	}, nil
}

func (p *SQLParser) Parse(sql string) (sqlparser.Statement, error) {
	return p.parser.Parse(sql)
}
