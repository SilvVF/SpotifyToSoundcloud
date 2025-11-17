package soundcloud

type TracksPage []Track

type Track struct {
	Kind                  string `json:"kind"`
	ID                    int    `json:"id"`
	Urn                   string `json:"urn"`
	CreatedAt             string `json:"created_at"`
	Duration              int    `json:"duration"`
	Commentable           bool   `json:"commentable"`
	CommentCount          int    `json:"comment_count"`
	Sharing               string `json:"sharing"`
	TagList               string `json:"tag_list"`
	Streamable            bool   `json:"streamable"`
	EmbeddableBy          string `json:"embeddable_by"`
	PurchaseURL           any    `json:"purchase_url"`
	PurchaseTitle         any    `json:"purchase_title"`
	Genre                 string `json:"genre"`
	Title                 string `json:"title"`
	Description           any    `json:"description"`
	LabelName             string `json:"label_name"`
	Release               any    `json:"release"`
	KeySignature          any    `json:"key_signature"`
	Isrc                  string `json:"isrc"`
	Bpm                   any    `json:"bpm"`
	ReleaseYear           int    `json:"release_year"`
	ReleaseMonth          int    `json:"release_month"`
	ReleaseDay            int    `json:"release_day"`
	License               string `json:"license"`
	URI                   string `json:"uri"`
	User                  User   `json:"user"`
	PermalinkURL          string `json:"permalink_url"`
	ArtworkURL            string `json:"artwork_url"`
	StreamURL             string `json:"stream_url"`
	DownloadURL           any    `json:"download_url"`
	WaveformURL           string `json:"waveform_url"`
	AvailableCountryCodes any    `json:"available_country_codes"`
	SecretURI             any    `json:"secret_uri"`
	UserFavorite          bool   `json:"user_favorite"`
	UserPlaybackCount     int    `json:"user_playback_count"`
	PlaybackCount         int    `json:"playback_count"`
	DownloadCount         int    `json:"download_count"`
	FavoritingsCount      int    `json:"favoritings_count"`
	RepostsCount          int    `json:"reposts_count"`
	Downloadable          bool   `json:"downloadable"`
	Access                string `json:"access"`
	Policy                any    `json:"policy"`
	MonetizationModel     any    `json:"monetization_model"`
	MetadataArtist        string `json:"metadata_artist"`
}

type Product struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
type Subscriptions struct {
	Product Product `json:"product"`
}
type User struct {
	AvatarURL            string          `json:"avatar_url"`
	ID                   int             `json:"id"`
	Urn                  string          `json:"urn"`
	Kind                 string          `json:"kind"`
	PermalinkURL         string          `json:"permalink_url"`
	URI                  string          `json:"uri"`
	Username             string          `json:"username"`
	Permalink            string          `json:"permalink"`
	CreatedAt            string          `json:"created_at"`
	LastModified         string          `json:"last_modified"`
	FirstName            any             `json:"first_name"`
	LastName             any             `json:"last_name"`
	FullName             string          `json:"full_name"`
	City                 any             `json:"city"`
	Description          any             `json:"description"`
	Country              any             `json:"country"`
	TrackCount           int             `json:"track_count"`
	PublicFavoritesCount int             `json:"public_favorites_count"`
	RepostsCount         int             `json:"reposts_count"`
	FollowersCount       int             `json:"followers_count"`
	FollowingsCount      int             `json:"followings_count"`
	Plan                 string          `json:"plan"`
	MyspaceName          any             `json:"myspace_name"`
	DiscogsName          any             `json:"discogs_name"`
	WebsiteTitle         any             `json:"website_title"`
	Website              any             `json:"website"`
	CommentsCount        int             `json:"comments_count"`
	Online               bool            `json:"online"`
	LikesCount           int             `json:"likes_count"`
	PlaylistCount        int             `json:"playlist_count"`
	Subscriptions        []Subscriptions `json:"subscriptions"`
}

type CreatedPlaylist struct {
	ArtworkURL    string `json:"artwork_url"`
	CreatedAt     string `json:"created_at"`
	Description   any    `json:"description"`
	Downloadable  bool   `json:"downloadable"`
	Duration      int    `json:"duration"`
	Ean           any    `json:"ean"`
	EmbeddableBy  string `json:"embeddable_by"`
	Genre         string `json:"genre"`
	Urn           string `json:"urn"`
	Kind          string `json:"kind"`
	Label         any    `json:"label"`
	LabelID       any    `json:"label_id"`
	LabelName     any    `json:"label_name"`
	LastModified  string `json:"last_modified"`
	License       string `json:"license"`
	LikesCount    int    `json:"likes_count"`
	Permalink     string `json:"permalink"`
	PermalinkURL  string `json:"permalink_url"`
	PlaylistType  string `json:"playlist_type"`
	PurchaseTitle any    `json:"purchase_title"`
	PurchaseURL   any    `json:"purchase_url"`
	Release       any    `json:"release"`
	ReleaseDay    any    `json:"release_day"`
	ReleaseMonth  any    `json:"release_month"`
	ReleaseYear   any    `json:"release_year"`
	Sharing       string `json:"sharing"`
	Streamable    bool   `json:"streamable"`
	TagList       string `json:"tag_list"`
	Tags          string `json:"tags"`
	Title         string `json:"title"`
	TrackCount    int    `json:"track_count"`
	Tracks        struct {
		Ref string `json:"$ref"`
	} `json:"tracks"`
	TracksURI string `json:"tracks_uri"`
	Type      string `json:"type"`
	URI       string `json:"uri"`
	User      struct {
		Ref string `json:"$ref"`
	} `json:"user"`
	UserID int `json:"user_id"`
}
