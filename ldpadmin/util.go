package ldpadmin

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

func trimSql(s string) string {
	sp := strings.Split(s, "\n")
	var b strings.Builder
	for _, line := range sp {
		if len(line) == 0 {
			continue
		}
		b.WriteString(
			strings.TrimRight(strings.TrimPrefix(line, "  "), " ") +
				"\n")
	}
	return b.String()
}

func exec(tx *sql.Tx, opts *LoadOptions, query string,
	args ...interface{}) (sql.Result, error) {
	if opts.Debug {
		var q string = query
		var a string
		for x := len(args) - 1; x >= 0; x-- {
			switch t := args[x].(type) {
			case int64, float64:
				a = fmt.Sprintf("%v", args[x])
			case string, time.Time, bool, []byte:
				a = fmt.Sprintf("'%v'", args[x])
			default:
				return nil, fmt.Errorf("unknown type %T", t)
			}
			q = strings.Replace(q, fmt.Sprintf("$%v", x+1), a, -1)
		}
		fmt.Printf("%s", q)
	}
	return tx.Exec(query, args...)
}
