package post

type Vote struct {
	User      string `json:"user"`
	VoteScore int    `json:"vote" bson:"vote"`
}

func (v *Vote) WithVote(value int) *Vote {
	v.VoteScore = value

	return v
}
