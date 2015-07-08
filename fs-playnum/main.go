package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const fs_playnum = "select mediaid,playnum, daynum,seven_daysnum,weeknum,monthnum,modifydate from fs_media_playnum"

func main() {
	db, err := sql.Open("mysql", "dbs:R4XBfuptAH@tcp(192.168.8.121:3306)/corsair_0")
	panic_error(err)
	defer db.Close()
	rows, err := db.Query(fs_playnum)
	panic_error(err)
	for rows.Next() {
		var id, playnum, daynum, seven_daysnum, weeknum, monthnum int64
		var modifydate []byte
		if err = rows.Scan(&id, &playnum, &daynum, &seven_daysnum, &weeknum, &monthnum, &modifydate); err == nil {
			fmt.Println(time_parse(string(modifydate)))
		} else {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}

func panic_error(err error) {
	if err != nil {
		panic(err)
	}
}

//Mon Jan 2 15:04:05 -0700 MST 2006
const tmlayout = "2006-01-02 15:04:05 -0700"

func time_parse(t string) time.Time {
	v, _ := time.Parse(tmlayout, t+" +0800")
	return v
}
