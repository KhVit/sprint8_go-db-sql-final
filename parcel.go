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
	// реализуйте добавление строки в таблицу parcel, используйте данные из переменной p
	ins, err := s.db.Exec("INSERT INTO parcel (client, status, address, created_at) VALUES (:client, :status, :address, :created_at)",
		sql.Named("client", p.Client),
		sql.Named("status", p.Status),
		sql.Named("address", p.Address),
		sql.Named("created_at", p.CreatedAt))
	if err != nil {
		return 0, err
	}

	lastId, err := ins.LastInsertId()
	if err != nil {
		return 0, err
	}

	// верните идентификатор последней добавленной записи
	return int(lastId), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	// реализуйте чтение строки по заданному number
	// здесь из таблицы должна вернуться только одна строка
	selRow := s.db.QueryRow("SELECT * FROM parcel WHERE number = :num",
		sql.Named("num", number))

	// заполните объект Parcel данными из таблицы
	p := Parcel{}

	err := selRow.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {
		return p, err
	}

	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	// реализуйте чтение строк из таблицы parcel по заданному client
	// здесь из таблицы может вернуться несколько строк
	selRows, err := s.db.Query("SELECT * FROM parcel WHERE client = :client_id",
		sql.Named("client_id", client))
	if err != nil {
		return nil, err
	}
	defer selRows.Close()

	// заполните срез Parcel данными из таблицы
	var res []Parcel

	for selRows.Next() {
		p := Parcel{}

		err := selRows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
		if err != nil {
			return nil, err
		}

		res = append(res, p)
	}

	if err := selRows.Err(); err != nil { // проверка курсора на наличие ошибок
		return nil, err
	}

	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	// реализуйте обновление статуса в таблице parcel
	_, err := s.db.Exec("UPDATE parcel SET status = :new_status WHERE number = :num",
		sql.Named("new_status", status),
		sql.Named("num", number))

	return err
}

func (s ParcelStore) SetAddress(number int, address string) error {
	// реализуйте обновление адреса в таблице parcel
	// менять адрес можно только если значение статуса registered
	_, err := s.db.Exec("UPDATE parcel SET address = :new_address WHERE number = :num AND status = :status_reg",
		sql.Named("new_address", address),
		sql.Named("num", number),
		sql.Named("status_reg", ParcelStatusRegistered))

	return err
}

func (s ParcelStore) Delete(number int) error {
	// реализуйте удаление строки из таблицы parcel
	// удалять строку можно только если значение статуса registered
	_, err := s.db.Exec("DELETE FROM parcel WHERE number = :num AND status = :status_reg",
		sql.Named("num", number),
		sql.Named("status_reg", ParcelStatusRegistered))

	return err
}
