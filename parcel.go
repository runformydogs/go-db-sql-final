package main

import (
	"database/sql"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	//query := `INSERT INTO parcel (client, status, address, created_at) VALUES (:client, :status, :address, :created_at)`
	res, err := s.db.Exec("INSERT INTO parcel (client, status, address, created_at) VALUES (?, ?, ?, ?)", p.Client, p.Status, p.Address, p.CreatedAt)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	query := `SELECT number, status, address, created_at FROM parcel WHERE number = :number`
	row := s.db.QueryRow(query, sql.Named("number", number))

	p := Parcel{}
	err := row.Scan(&p.Number, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {
		return Parcel{}, err
	}
	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	query := `SELECT number, status, address, created_at FROM parcel WHERE client = :client`
	rows, err := s.db.Query(query, sql.Named("client", client))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []Parcel

	for rows.Next() {
		p := Parcel{}
		if err := rows.Scan(&p.Number, &p.Status, &p.Address, &p.CreatedAt); err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	query := `UPDATE parcel SET status = :status WHERE number = :number`
	_, err := s.db.Exec(query, sql.Named("status", status), sql.Named("number", number))
	return err
}

func (s ParcelStore) SetAddress(number int, address string) error {
	query := `UPDATE parcel SET address = :address WHERE number = :number AND status = :status`
	_, err := s.db.Exec(query, sql.Named("address", address), sql.Named("number", number), sql.Named("status", ParcelStatusRegistered))
	return err
}

func (s ParcelStore) Delete(number int) error {
	query := `DELETE FROM parcel WHERE number = :number AND status = :status`
	_, err := s.db.Exec(query, sql.Named("number", number), sql.Named("status", ParcelStatusRegistered))
	return err
}
