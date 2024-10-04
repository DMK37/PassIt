package api

type LikeRequest struct {
	OwnerId string `json:"owner_id"`
	PostId string `json:"post_id"`
}