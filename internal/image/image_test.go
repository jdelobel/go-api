package image_test

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jdelobel/go-api/internal/image"
	"github.com/jdelobel/go-api/internal/platform/db"
)

const (

	// Succeed is the Unicode codepoint for a check mark.
	Succeed = "\u2713"

	// Failed is the Unicode codepoint for an X mark.
	Failed = "\u2717"
)

// TestImages validates an image can be created, retrieved and
// then removed from the system.
func TestImages(t *testing.T) {
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

	m := image.CreateImage{
		ID:        "47c658e0-68d7-4d79-9f9f-25ece8a1fb03",
		Title:     "Image Elijah Baley",
		URL:       "/images/1280/720/test-2260-b1396d-1@1x.jpeg",
		Slug:      "/images/1280/720/test-2260-b1396d-1@1x",
		Publisher: "etf1",
	}

	t.Log("Given the need to add a new image, retrieve and remove that image from the system.")
	{
		t.Log("\tTest 0:\tWhen using a valid CreateImage value")
		{
			imageId, err := image.Create(ctx, masterDB, nil, &m)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to create an image in the system : %v", Failed, err)
			}
			t.Logf("\t%s\tShould be able to create an image in the system.", Succeed)

			rm, err := image.Retrieve(ctx, masterDB, *imageId.ID)
			if err != nil {
				t.Fatalf("\t%s\tShould be able to retrieve the image back from the system : %v", Failed, err)
			}
			t.Logf("\t%s\tShould be able to retrieve the image back from the system.", Succeed)

			if rm == nil || *imageId.ID != *rm.ID {
				t.Fatalf("\t%s\tShould have a match between the created image and the one retrieved : %v", Failed, err)
			}
			t.Logf("\t%s\tShould have a match between the created image and the one retrieved.", Succeed)

			if _, err := image.Retrieve(ctx, masterDB, *rm.ID); err == nil {
				t.Fatalf("\t%s\tShould NOT be able to retrieve the image back from the system : %v", Failed, err)
			}
			t.Logf("\t%s\tShould NOT be able to retrieve the image back from the system.", Succeed)
		}
	}
}
