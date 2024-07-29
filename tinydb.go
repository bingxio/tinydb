package tinydb

import (
	"encoding/json"
	"os"
)

type Database struct {
	Table []Table `json:"table"`
}

type Table struct {
	Name string `json:"name"`
	Rows []any  `json:"rows"`
}

func OpenDB(path string) (*Database, error) {
	db := new(Database)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		_, err = os.Create(path)

		if err != nil {
			return nil, err
		}
		return db, nil
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if len(b) == 0 {
		return db, nil
	}
	if err = json.Unmarshal(b, db); err != nil {
		return nil, err
	}
	return db, nil
}

func (db *Database) table(name string) *Table {
	i := 0
	for i != len(db.Table) {
		if seek := db.Table[i]; seek.Name == name {
			break
		}
		i++
	}
	if len(db.Table) != 0 {
		return &db.Table[i]
	}
	return db.newTable(name)
}

func (db *Database) back() *Table {
	len := len(db.Table) - 1
	return &db.Table[len]
}

func (db *Database) newTable(name string) *Table {
	db.Table = append(db.Table, Table{Name: name, Rows: []any{}})
	return db.back()
}

func (db Database) WriteStorage(file *os.File) error {
	bytes, err := json.MarshalIndent(db, "", "  ")
	if err != nil {
		return err
	}
	bytes = append(bytes, '\n')
	_, err = file.Write(bytes)
	return err
}

type Cond func(any) bool

func (db *Database) iterator(name string, cond Cond, do func(any)) []int {
	index := []int{}
	tb := db.table(name)

	for i := 0; i < len(tb.Rows); i++ {
		v := true
		if cond != nil {
			v = cond(tb.Rows[i])
		}
		if v {
			index = append(index, i)
		}
		if v && do != nil {
			do(tb.Rows[i])
		}
	}
	return index
}

func (db *Database) Insert(name string, row any) {
	tb := db.table(name)
	tb.Rows = append(tb.Rows, row)
}

func (db *Database) Query(name string, cond Cond) []any {
	conform := []any{}
	db.iterator(
		name, cond, func(row any) { conform = append(conform, row) })
	return conform
}

func (db *Database) Delete(name string, cond Cond) uint {
	tb := db.table(name)
	index := db.iterator(name, cond, nil)

	for _, i := range index {
		copy(tb.Rows[i:], tb.Rows[i+1:])
		tb.Rows = tb.Rows[:len(tb.Rows)-1]
	}
	return uint(len(index))
}

func (db *Database) Update(name string, cond Cond, value any) uint {
	index := db.iterator(name, cond, nil)
	tb := db.table(name)

	for _, i := range index {
		tb.Rows[i] = value
	}
	return uint(len(index))
}
