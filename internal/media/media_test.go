package media_test

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/go-api/internal/media"
	"github.com/go-api/internal/platform/db"
)

const (

	// Succeed is the Unicode codepoint for a check mark.
	Succeed = "\u2713"

	// Failed is the Unicode codepoint for an X mark.
	Failed = "\u2717"
)

// TestMedias validates a media can be created, retrieved and
// then removed from the system.
func TestMedias(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Check the environment for a configured port value.
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "got:got2015@ds039441.mongolab.com:39441/gotraining"
	}

	// Register the Master Session for the database.
	log.Println("main : Started : Capturing Master DB...")
	masterDB, err := db.NewPSQL("postgres", dbHost)
	if err != nil {
		t.Fatal(err)
	}

	m := media.CreateMedia{
		ID:        "47c658e0-68d7-4d79-9f9f-25ece8a1fb03",
		Title:     "Media Elijah Baley",
		URL:       "/images/1280/720/test-2260-b1396d-1@1x.jpeg",
		Slug:      "/images/1280/720/test-2260-b1396d-1@1x",
		Publisher: "etf1",
	}

	t.Log("Given the need to add a new media, retrieve and remove that media from the system.")
	{
		t.Log("\tTest 0:\tWhen using a valid CreateMedia value")
		{
			mediaId, err := media.Create(ctx, masterDB, &m)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to create a media in the system : %v", Failed, err)
			}
			t.Logf("\t%s\tShould be able to create a media in the system.", Succeed)

			rm, err := media.Retrieve(ctx, masterDB, mediaId)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to retrieve the media back from the system : %v", Failed, err)
			}
			t.Logf("\t%s\tShould be able to retrieve the media back from the system.", Succeed)

			if rm == nil || mediaId != *rm.ID {
				t.Fatalf("\t%s\tShould have a match between the created media and the one retrieved : %v", Failed, err)
			}
			t.Logf("\t%s\tShould have a match between the created media and the one retrieved.", Succeed)

			if _, err := media.Retrieve(ctx, masterDB, *rm.ID); err == nil {
				t.Fatalf("\t%s\tShould NOT be able to retrieve the media back from the system : %v", Failed, err)
			}
			t.Logf("\t%s\tShould NOT be able to retrieve the media back from the system.", Succeed)
		}
	}
}
