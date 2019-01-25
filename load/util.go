package load

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
)

func compact(s string) string {
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

func exec(tx *sql.Tx, query string, args ...interface{}) (sql.Result, error) {
	fmt.Fprintf(os.Stderr, "%s{", query)
	for x, a := range args {
		if x != 0 {
			fmt.Fprintf(os.Stderr, ", ")
		}
		fmt.Fprintf(os.Stderr, "\"%s\"", a)
	}
	fmt.Fprintf(os.Stderr, "}\n\n")
	return tx.Exec(query, args...)
}
