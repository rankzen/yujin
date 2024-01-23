// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: summoners.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const countSoloqRecordsByName = `-- name: CountSoloqRecordsByName :one
SELECT COUNT(*)
FROM soloq_records
WHERE summoner_name = $1
`

func (q *Queries) CountSoloqRecordsByName(ctx context.Context, summonerName pgtype.Text) (int64, error) {
	row := q.db.QueryRow(ctx, countSoloqRecordsByName, summonerName)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const countSummonerRecordsByName = `-- name: CountSummonerRecordsByName :one
SELECT COUNT(*)
FROM summoner_records
WHERE name = $1
`

func (q *Queries) CountSummonerRecordsByName(ctx context.Context, name pgtype.Text) (int64, error) {
	row := q.db.QueryRow(ctx, countSummonerRecordsByName, name)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const deleteSoloqRecord = `-- name: DeleteSoloqRecord :one
DELETE FROM soloq_records
WHERE record_id = $1
RETURNING record_date, summoner_name
`

type DeleteSoloqRecordRow struct {
	RecordDate   pgtype.Timestamp `json:"record_date"`
	SummonerName pgtype.Text      `json:"summoner_name"`
}

func (q *Queries) DeleteSoloqRecord(ctx context.Context, recordID pgtype.UUID) (DeleteSoloqRecordRow, error) {
	row := q.db.QueryRow(ctx, deleteSoloqRecord, recordID)
	var i DeleteSoloqRecordRow
	err := row.Scan(&i.RecordDate, &i.SummonerName)
	return i, err
}

const deleteSummonerRecord = `-- name: DeleteSummonerRecord :one
DELETE FROM summoner_records
WHERE record_id = $1
RETURNING record_date, name
`

type DeleteSummonerRecordRow struct {
	RecordDate pgtype.Timestamp `json:"record_date"`
	Name       pgtype.Text      `json:"name"`
}

func (q *Queries) DeleteSummonerRecord(ctx context.Context, recordID pgtype.UUID) (DeleteSummonerRecordRow, error) {
	row := q.db.QueryRow(ctx, deleteSummonerRecord, recordID)
	var i DeleteSummonerRecordRow
	err := row.Scan(&i.RecordDate, &i.Name)
	return i, err
}

const insertSoloqRecord = `-- name: InsertSoloqRecord :one
INSERT INTO soloq_records
(record_date, league_id, summoner_id, summoner_name, tier, rank, league_points, wins, losses)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING record_id
`

type InsertSoloqRecordParams struct {
	RecordDate   pgtype.Timestamp `json:"record_date"`
	LeagueID     pgtype.Text      `json:"league_id"`
	SummonerID   pgtype.Text      `json:"summoner_id"`
	SummonerName pgtype.Text      `json:"summoner_name"`
	Tier         pgtype.Text      `json:"tier"`
	Rank         pgtype.Text      `json:"rank"`
	LeaguePoints pgtype.Int4      `json:"league_points"`
	Wins         pgtype.Int4      `json:"wins"`
	Losses       pgtype.Int4      `json:"losses"`
}

func (q *Queries) InsertSoloqRecord(ctx context.Context, arg InsertSoloqRecordParams) (pgtype.UUID, error) {
	row := q.db.QueryRow(ctx, insertSoloqRecord,
		arg.RecordDate,
		arg.LeagueID,
		arg.SummonerID,
		arg.SummonerName,
		arg.Tier,
		arg.Rank,
		arg.LeaguePoints,
		arg.Wins,
		arg.Losses,
	)
	var record_id pgtype.UUID
	err := row.Scan(&record_id)
	return record_id, err
}

const insertSummonerRecord = `-- name: InsertSummonerRecord :one
INSERT INTO summoner_records
(record_date, account_id, profile_icon_id, revision_date, name, id, puuid, summoner_level)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING record_id
`

type InsertSummonerRecordParams struct {
	RecordDate    pgtype.Timestamp `json:"record_date"`
	AccountID     pgtype.Text      `json:"account_id"`
	ProfileIconID pgtype.Int4      `json:"profile_icon_id"`
	RevisionDate  pgtype.Int8      `json:"revision_date"`
	Name          pgtype.Text      `json:"name"`
	ID            pgtype.Text      `json:"id"`
	Puuid         pgtype.Text      `json:"puuid"`
	SummonerLevel pgtype.Int8      `json:"summoner_level"`
}

func (q *Queries) InsertSummonerRecord(ctx context.Context, arg InsertSummonerRecordParams) (pgtype.UUID, error) {
	row := q.db.QueryRow(ctx, insertSummonerRecord,
		arg.RecordDate,
		arg.AccountID,
		arg.ProfileIconID,
		arg.RevisionDate,
		arg.Name,
		arg.ID,
		arg.Puuid,
		arg.SummonerLevel,
	)
	var record_id pgtype.UUID
	err := row.Scan(&record_id)
	return record_id, err
}

const selectSoloqRecordsByName = `-- name: SelectSoloqRecordsByName :many
SELECT record_id, record_date, league_id, summoner_id, summoner_name, tier, rank, league_points, wins, losses
FROM soloq_records
WHERE summoner_name = $1
ORDER BY record_date DESC
LIMIT $2 OFFSET $3
`

type SelectSoloqRecordsByNameParams struct {
	SummonerName pgtype.Text `json:"summoner_name"`
	Limit        int32       `json:"limit"`
	Offset       int32       `json:"offset"`
}

func (q *Queries) SelectSoloqRecordsByName(ctx context.Context, arg SelectSoloqRecordsByNameParams) ([]SoloqRecord, error) {
	rows, err := q.db.Query(ctx, selectSoloqRecordsByName, arg.SummonerName, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []SoloqRecord
	for rows.Next() {
		var i SoloqRecord
		if err := rows.Scan(
			&i.RecordID,
			&i.RecordDate,
			&i.LeagueID,
			&i.SummonerID,
			&i.SummonerName,
			&i.Tier,
			&i.Rank,
			&i.LeaguePoints,
			&i.Wins,
			&i.Losses,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const selectSoloqRecordsBySummonerId = `-- name: SelectSoloqRecordsBySummonerId :many
SELECT record_id, record_date, league_id, summoner_id, summoner_name, tier, rank, league_points, wins, losses
FROM soloq_records
WHERE summoner_id = $1
ORDER BY record_date DESC
LIMIT $2 OFFSET $3
`

type SelectSoloqRecordsBySummonerIdParams struct {
	SummonerID pgtype.Text `json:"summoner_id"`
	Limit      int32       `json:"limit"`
	Offset     int32       `json:"offset"`
}

func (q *Queries) SelectSoloqRecordsBySummonerId(ctx context.Context, arg SelectSoloqRecordsBySummonerIdParams) ([]SoloqRecord, error) {
	rows, err := q.db.Query(ctx, selectSoloqRecordsBySummonerId, arg.SummonerID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []SoloqRecord
	for rows.Next() {
		var i SoloqRecord
		if err := rows.Scan(
			&i.RecordID,
			&i.RecordDate,
			&i.LeagueID,
			&i.SummonerID,
			&i.SummonerName,
			&i.Tier,
			&i.Rank,
			&i.LeaguePoints,
			&i.Wins,
			&i.Losses,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const selectSummonerRecordsByName = `-- name: SelectSummonerRecordsByName :many
SELECT record_id, record_date, account_id, profile_icon_id, revision_date, name, id, puuid, summoner_level
FROM summoner_records
WHERE name = $1
ORDER BY record_date DESC
LIMIT $2 OFFSET $3
`

type SelectSummonerRecordsByNameParams struct {
	Name   pgtype.Text `json:"name"`
	Limit  int32       `json:"limit"`
	Offset int32       `json:"offset"`
}

func (q *Queries) SelectSummonerRecordsByName(ctx context.Context, arg SelectSummonerRecordsByNameParams) ([]SummonerRecord, error) {
	rows, err := q.db.Query(ctx, selectSummonerRecordsByName, arg.Name, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []SummonerRecord
	for rows.Next() {
		var i SummonerRecord
		if err := rows.Scan(
			&i.RecordID,
			&i.RecordDate,
			&i.AccountID,
			&i.ProfileIconID,
			&i.RevisionDate,
			&i.Name,
			&i.ID,
			&i.Puuid,
			&i.SummonerLevel,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const selectSummonerRecordsByPuuid = `-- name: SelectSummonerRecordsByPuuid :many
SELECT record_id, record_date, account_id, profile_icon_id, revision_date, name, id, puuid, summoner_level
FROM summoner_records
WHERE puuid = $1
ORDER BY record_date DESC
LIMIT $2 OFFSET $3
`

type SelectSummonerRecordsByPuuidParams struct {
	Puuid  pgtype.Text `json:"puuid"`
	Limit  int32       `json:"limit"`
	Offset int32       `json:"offset"`
}

func (q *Queries) SelectSummonerRecordsByPuuid(ctx context.Context, arg SelectSummonerRecordsByPuuidParams) ([]SummonerRecord, error) {
	rows, err := q.db.Query(ctx, selectSummonerRecordsByPuuid, arg.Puuid, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []SummonerRecord
	for rows.Next() {
		var i SummonerRecord
		if err := rows.Scan(
			&i.RecordID,
			&i.RecordDate,
			&i.AccountID,
			&i.ProfileIconID,
			&i.RevisionDate,
			&i.Name,
			&i.ID,
			&i.Puuid,
			&i.SummonerLevel,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}