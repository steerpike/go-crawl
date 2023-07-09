package models

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Artist struct {
	ID       int64             `json:"id"`
	FromID   int64             `json:"fromId"`
	Response int               `json:"response"`
	Name     string            `json:"name"`
	Url      string            `json:"url"`
	Path     string            `json:"path"`
	Tags     []string          `json:"tag"`
	Similar  []string          `json:"similar"`
	Videos   map[string]string `json:"videos"`
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
	fmt.Println("ID is:", id)
	a.ID = id
	// Check if there is a similar artist id to relate this to:
	if a.ID != 0 {
		if a.FromID != 0 && a.ID != 0 {
			_, err = tx.Exec("INSERT OR IGNORE INTO Similar_Artists (ArtistID1, ArtistID2) VALUES (?, ?)", a.FromID, a.ID)
			if err != nil {
				tx.Rollback()
				return fail(err)
			}
		}
		// Insert tags
		for _, tag := range a.Tags {
			res, err := tx.Exec("INSERT INTO Tags (TagName) VALUES (?) ON CONFLICT(TagName) DO UPDATE SET TagName = ? RETURNING ID", tag, tag)
			if err != nil {
				tx.Rollback()
				return fail(err)
			}

			tagId, err := res.LastInsertId()
			if err != nil {
				tx.Rollback()
				return fail(err)
			}

			_, err = tx.Exec("INSERT OR IGNORE INTO Artist_Tags (ArtistID, TagID) VALUES (?, ?)", a.ID, tagId)
			if err != nil {
				tx.Rollback()
				return fail(err)
			}
		}
		for videoName, videoUrl := range a.Videos {
			res, err := tx.Exec("INSERT INTO Videos (VideoName, VideoUrl) VALUES (?, ?) ON CONFLICT(VideoUrl) DO UPDATE SET (VideoName) = ? RETURNING ID", videoName, videoUrl, videoName)
			if err != nil {
				tx.Rollback()
				return fail(err)
			}

			videoId, err := res.LastInsertId()
			if err != nil {
				tx.Rollback()
				return fail(err)
			}

			_, err = tx.Exec("INSERT OR IGNORE INTO Artist_Videos (ArtistID, VideoID) VALUES (?, ?)", a.ID, videoId)
			if err != nil {
				tx.Rollback()
				return fail(err)
			}
		}
		fmt.Printf("%+v\n", a)
	} else {
		fmt.Println("Artist id was 0 for artist: ", a.Path)
		fmt.Printf("%+v\n", a)
	}
	tx.Commit()
	return a, err
}
