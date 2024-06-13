package models

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Artist struct {
	ID        int64             `json:"id"`
	Response  int               `json:"response"`
	Name      string            `json:"name"`
	SourceUrl string            `json:"source"`
	Url       string            `json:"url"`
	Path      string            `json:"path"`
	Tags      []string          `json:"tag"`
	Similar   []string          `json:"similar"`
	Videos    map[string]string `json:"videos"`
}

func ArtistExists(path string) bool {
	var exists bool
	db, err := sql.Open("sqlite3", "music.db")
	if err != nil {
		log.Println("Something went wrong opening database:", err)
	}
	query := `SELECT 1 FROM Artists WHERE Path = ?`
	err = db.QueryRow(query, path).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		log.Fatalf("querying for url existence: %v", err)
	}
	return exists
}

func (a *Artist) Save() (*Artist, error) {
	var id int64
	fail := func(err error) (*Artist, error) {
		log.Println("Something went wrong in a transaction:", err)
		return nil, fmt.Errorf("Saving Artist: %v", err)
	}
	db, err := sql.Open("sqlite3", "music.db")
	if err != nil {
		log.Println("Something went wrong opening database:", err)
		return fail(err)
	}
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Println("Error beginning the transaction", err)
		return fail(err)
	}
	err = tx.QueryRowContext(ctx, "INSERT INTO Artists (Response, Name, Url, Path, LastCrawled) VALUES (?, ?, ?, ?, datetime('now')) ON CONFLICT(Path) DO UPDATE SET  Path = ? RETURNING ID", a.Response, a.Name, a.Url, a.Path, a.Path).Scan(&id)
	if err != nil {
		log.Println("Something went wrong inserting:", err)
		err := tx.Rollback()
		if err != nil {
			log.Println("Error rolling back transaction on the insertion of artist", err)
		}
		return fail(err)
	}
	fmt.Printf("NEW ID: %+v\n", id)
	// Insert Tags

	for _, tag := range a.Tags {
		var tagID int
		err := tx.QueryRowContext(ctx, "INSERT INTO Tags (Name) VALUES (?) ON CONFLICT(Name) DO UPDATE SET Name = ? RETURNING ID", tag, tag).Scan(&tagID)
		if err != nil {
			tx.Rollback()
			return fail(err)
		}
		_, err = tx.Exec("INSERT OR IGNORE INTO Artist_Tags (ArtistUrl, TagName) VALUES (?, ?)", a.Url, tag)
		if err != nil {
			tx.Rollback()
			return fail(err)
		}
	}

	// Insert Videos

	for videoUrl, videoName := range a.Videos {
		var videoID int
		err := tx.QueryRowContext(ctx, "INSERT INTO Videos (Name, Url) VALUES (?, ?) ON CONFLICT(Url) DO UPDATE SET Name = ? RETURNING ID", videoName, videoUrl, videoName).Scan(&videoID)
		if err != nil {
			tx.Rollback()
			return fail(err)
		}
		_, err = tx.Exec("INSERT OR IGNORE INTO Artist_Videos (ArtistUrl, VideoUrl) VALUES (?, ?)", a.Url, videoUrl)
		if err != nil {
			tx.Rollback()
			return fail(err)
		}
	}
	// Insert Seeds
	for _, seed := range a.Similar {
		url := "https://www.last.fm" + seed
		if !ArtistExists(seed) {
			_, err = tx.Exec("INSERT OR IGNORE INTO Seeds (SourceUrl, Url) VALUES (?, ?)", a.Url, url)
			if err != nil {
				tx.Rollback()
				return fail(err)
			}
		} else { //Artist exists, so we need to insert the relationship
			_, err = tx.Exec("INSERT OR IGNORE INTO Similar_Artists (ArtistUrl1, ArtistUrl2) VALUES (?, ?)", a.Url, url)
			if err != nil {
				tx.Rollback()
				return fail(err)
			}
		}
	}
	if a.SourceUrl != "" {
		_, err = tx.Exec("INSERT OR IGNORE INTO Similar_Artists (ArtistUrl1, ArtistUrl2) VALUES (?, ?)", a.Url, a.SourceUrl)
		if err != nil {
			tx.Rollback()
			return fail(err)
		}
		_, err = tx.Exec("DELETE FROM Seeds WHERE Url = ?", a.Url)
		if err != nil {
			tx.Rollback()
			return fail(err)
		}
	}

	fmt.Printf("%+v\n", a)
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
	return a, err
}
