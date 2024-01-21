package main

import (
	"database/sql"
	"fmt"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	ip, err := s.db.Exec("INSERT INTO parcel (client, status, address, created_at) VALUES (?, ?, ?, ?)", p.Client, p.Status, p.Address, p.CreatedAt)
	if err != nil {
		return 0, err
	}
	currentNumberParcel, _ := ip.LastInsertId()
	return int(currentNumberParcel), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	p := Parcel{}
	err := s.db.QueryRow("SELECT client, status, address, created_at FROM parcel WHERE number = ?", number).Scan(&p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {
		return p, err
	}
	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	rows, err := s.db.Query("SELECT number, client, status, address, created_at FROM parcel WHERE client = ?", client)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var parcels []Parcel
	for rows.Next() {
		var p Parcel
		err := rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		parcels = append(parcels, p)
	}

	return parcels, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	s.db.QueryRow("UPDATE parcel SET status = ? WHERE number = ?", status, number)
	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	row, err := s.db.Query("SELECT status FROM parcel WHERE number = ?", number)
	if err != nil {
		return err
	}
	if row.Next() {
		var status string
		err = row.Scan(&status)
		if err != nil {
			return err
		}
		if status != "registered" {
			return fmt.Errorf("менять адрес можно только если значение статуса registered")
		}
	}
	row.Close()

	s.db.QueryRow("UPDATE parcel SET address = ? WHERE number = ?", address, number)
	return nil
}

func (s ParcelStore) Delete(number int) error {
	row, err := s.db.Query("SELECT status FROM parcel WHERE number = ?", number)
	if err != nil {
		return err
	}
	if row.Next() {
		var status string
		err = row.Scan(&status)
		if err != nil {
			return err
		}
		if status != "registered" {
			return fmt.Errorf("удалять строку можно только если значение статуса registered")
		}
	}

	row.Close()

	s.db.Exec("DELETE FROM parcel WHERE number = ?", number)
	return nil
}
