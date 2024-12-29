package models

type VoteRequest struct {
	ID       int    `json:"id" binding:"required"`       
	Type     string `json:"type" binding:"required"`    
	VoteType string `json:"vote_type" binding:"required"` 
}