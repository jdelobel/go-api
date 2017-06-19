package media

import "time"

// Media contains information about a media.
type CreateMedia struct {
	ID          string    `json:"media_id"`
	Title       string    `json:"title" validate:"required,min=3"`
	URL         string    `json:"url" validate:"required,min=3"`
	Slug        string    `json:"slug" validate:"required,min=3"`
	Publisher   string    `json:"publisher" validate:"required,min=3"`
	PublishedAt time.Time `json:"published_at"`
	ExpiredAt   time.Time `json:"expired_at"`
	Specific    struct{}  `json:"specific"`
}

// Media contains information about a media.
type Media struct {
	ID          *string    `db:"media_id" json:"media_id"`
	Title       *string    `db:"title" json:"title" validate:"required,min=3"`
	URL         *string    `db:"url" json:"url" validate:"required,min=3"`
	Slug        *string    `db:"slug" json:"slug" validate:"required,min=3"`
	Publisher   *string    `db:"publisher" json:"publisher" validate:"required,min=3"`
	PublishedAt *time.Time `db:"published_at" json:"published_at"`
	ExpiredAt   *time.Time `db:"expired_at" json:"expired_at"`
	Specific    *struct{}  `db:"specific" json:"specific"`
	CreatedAt   *time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   *time.Time `db:"updated_at" json:"updated_at"`
	RestoredAt  *time.Time `db:"restored_at" json:"restored_at"`
	DeletedAt   *time.Time `db:"deleted_at" json:"deleted_at"`
}
