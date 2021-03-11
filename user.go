package docbase

type (
	UserID int
	User   struct {
		ID              UserID `json:"id"`
		Name            string `json:"name"`
		ProfileImageURL string `json:"profile_image_url"`
	}
)
