package main

import (
	"database/sql"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math/rand"
	"testing"
	"time"
)

import (
	_ "github.com/stretchr/testify/require"
)

var (
	randSource = rand.NewSource(time.Now().UnixNano())
	randRange  = rand.New(randSource)
)

func connectToDb() *sql.DB {
	db, err := sql.Open("sqlite", "tracker.db")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return db
}

func getTestParcel() Parcel {
	return Parcel{
		Client:    1000,
		Status:    ParcelStatusRegistered,
		Address:   "test",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

func TestAddGetDelete(t *testing.T) {
	store := NewParcelStore(connectToDb())
	parcel := getTestParcel()

	p, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotEmpty(t, p)

	gp, err := store.Get(p)
	require.NoError(t, err)
	require.Equal(t, parcel.Client, gp.Client)
	require.Equal(t, parcel.Status, gp.Status)
	require.Equal(t, parcel.Address, gp.Address)
	require.Equal(t, parcel.CreatedAt, gp.CreatedAt)

	err = store.Delete(p)
	require.NoError(t, err)

	_, err = store.Get(p)
	require.Equal(t, sql.ErrNoRows, err)
}

func TestSetAddress(t *testing.T) {
	store := NewParcelStore(connectToDb())
	parcel := getTestParcel()

	p, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotEmpty(t, p)

	newAddress := "new test address"
	err = store.SetAddress(p, newAddress)
	require.NoError(t, err)

	gp, err := store.Get(p)
	require.NoError(t, err)
	require.Equal(t, newAddress, gp.Address)
}

func TestSetStatus(t *testing.T) {
	store := NewParcelStore(connectToDb())
	parcel := getTestParcel()

	p, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotEmpty(t, p)

	err = store.SetStatus(p, ParcelStatusSent)
	require.NoError(t, err)

	gp, err := store.Get(p)
	require.NoError(t, err)
	require.Equal(t, ParcelStatusSent, gp.Status)
}

func TestGetByClient(t *testing.T) {
	store := NewParcelStore(connectToDb())

	parcels := []Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}
	parcelMap := map[int]Parcel{}

	client := randRange.Intn(10_000_000)
	parcels[0].Client = client
	parcels[1].Client = client
	parcels[2].Client = client

	// add
	for i := 0; i < len(parcels); i++ {
		id, err := store.Add(parcels[i])
		require.NoError(t, err)
		require.NotEmpty(t, id)
		parcels[i].Number = id
		parcelMap[id] = parcels[i]
	}

	var storedParcels []Parcel
	storedParcels, err := store.GetByClient(client)
	require.NoError(t, err)
	assert.Equal(t, len(parcels), len(storedParcels))

	for _, parcel := range storedParcels {
		assert.Equal(t, parcelMap[parcel.Number], parcel)
		assert.Equal(t, parcelMap[parcel.Number].Client, parcel.Client)
		assert.Equal(t, parcelMap[parcel.Number].Address, parcel.Address)
		assert.Equal(t, parcelMap[parcel.Number].Status, parcel.Status)
		assert.Equal(t, parcelMap[parcel.Number].CreatedAt, parcel.CreatedAt)
	}
}
