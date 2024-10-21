package models

import (
    "context"
    "github.com/jackc/pgx/v4"
)

// Trail struct represents a trail entry in the database
type Trail struct {
    FID         int    `json:"fid"` // Updated to int
    Name        string `json:"name"`
    Restrooms   string `json:"restrooms"`
    Picnic      string `json:"picnic"`
    Fishing     string `json:"fishing"`
    Type        string `json:"type"`
    Difficulty  string `json:"difficulty"`
    AccessType  string `json:"access_type"`
    THLeash     string `json:"th_leash"`
    BikeTrail   string `json:"bike_trail"`
    HorseTrail  string `json:"horse_trail"`
    Fee         string `json:"fee"`
    RecycleBin  string `json:"recycle_bin"`
    Grills      string `json:"grills"`
    BikeRack    string `json:"bike_rack"`
    DogTube     string `json:"dog_tube"`
}

// CreateTable creates the trails table in PostgreSQL
func CreateTable(conn *pgx.Conn) error {
    _, err := conn.Exec(context.Background(), `
        CREATE TABLE IF NOT EXISTS trails (
            fid INTEGER PRIMARY KEY,
            name TEXT,
            restrooms TEXT,
            picnic TEXT,
            fishing TEXT,
            type TEXT,
            difficulty TEXT,
            access_type TEXT,
            th_leash TEXT,
            bike_trail TEXT,
            horse_trail TEXT,
            fee TEXT,
            recycle_bin TEXT,
            grills TEXT,
            bike_rack TEXT,
            dog_tube TEXT
        )
    `)
    return err
}

