package ldpadmin

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/lib/pq"
)

func (l *Loader) sqlCopyStage(stagingTable string,
	columns ...string) (*sql.Stmt, error) {
	if l.opts.Debug {
		fmt.Printf("COPY %s.%s (stage.", stagingTable)
		for x, c := range columns {
			if x != 0 {
				fmt.Printf(", ")
			}
			fmt.Printf("%v", c)
		}
		fmt.Printf(") FROM stdin;\n")
	}
	stmt, err := l.tx.Prepare(pq.CopyInSchema("stage",
		stagingTable, columns...))
	if err != nil {
		return nil, err
	}
	return stmt, nil
}

func (l *Loader) sqlCopyExec(stmt *sql.Stmt,
	args ...interface{}) (sql.Result, error) {
	if l.opts.Debug {
		if len(args) == 0 {
			fmt.Printf("\\.\n")
		} else {
			for x, a := range args {
				if x != 0 {
					fmt.Printf("\t")
				}
				fmt.Printf("%v", a)
			}
			fmt.Printf("\n")
		}
	}
	r, err := stmt.Exec(args...)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (l *Loader) sqlExec(query string,
	args ...interface{}) (sql.Result, error) {
	if l.opts.Debug {
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
	r, err := l.tx.Exec(query, args...)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (l *Loader) sqlMergePlaceholders(
	targetTable, targetId, stagingTable, stagingId string) error {
	cmd := fmt.Sprintf(""+
		"INSERT INTO %s\n"+
		"    (%s)\n"+
		"    SELECT %s\n"+
		"        FROM stage.%s\n"+
		"    ON CONFLICT (%s) DO NOTHING;\n",
		targetTable, targetId, stagingId, stagingTable, targetId)
	_, err := l.sqlExec(cmd)
	if err != nil {
		return err
	}
	return nil
}

func (l *Loader) sqlTruncateStage(stagingTable string) error {
	cmd := fmt.Sprintf(""+
		"TRUNCATE stage.%s;\n",
		stagingTable)
	_, err := l.sqlExec(cmd)
	if err != nil {
		return err
	}
	return nil
}

// OLD
/*
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
*/
