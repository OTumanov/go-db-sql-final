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
	// реализуйте добавление строки в таблицу parcel, используйте данные из переменной p
	ip, err := s.db.Exec("INSERT INTO parcel (client, status, address, created_at) VALUES (?, ?, ?, ?)", p.Client, p.Status, p.Address, p.CreatedAt)
	if err != nil {
		return 0, err
	}
	currentNumberParcel, _ := ip.LastInsertId()
	return int(currentNumberParcel), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	// реализуйте чтение строки по заданному number
	// здесь из таблицы должна вернуться только одна строка

	// заполните объект Parcel данными из таблицы
	p := Parcel{}

	err := s.db.QueryRow("SELECT client, status, address, created_at FROM parcel WHERE number = ?", number).Scan(&p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {
		return p, err
	}
	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	// реализуйте чтение строк из таблицы parcel по заданному client
	// здесь из таблицы может вернуться несколько строк
	rows, err := s.db.Query("SELECT client, status, address, created_at FROM parcel WHERE client = ?", client)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var parcels []Parcel
	for rows.Next() {
		var p Parcel
		err := rows.Scan(&p.Client, &p.Status, &p.Address, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		parcels = append(parcels, p)
	}
	return parcels, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	// реализуйте обновление статуса в таблице parcel
	s.db.QueryRow("UPDATE parcel SET status = ? WHERE number = ?", status, number)
	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	// реализуйте обновление адреса в таблице parcel
	// менять адрес можно только если значение статуса registered
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
	s.db.QueryRow("UPDATE parcel SET address = ? WHERE number = ?", address, number)
	return nil
}

func (s ParcelStore) Delete(number int) error {
	// реализуйте удаление строки из таблицы parcel
	// удалять строку можно только если значение статуса registered

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
	s.db.QueryRow("DELETE FROM parcel WHERE number = ?", number)
	return nil
}
