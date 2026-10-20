package storage

type Exercise struct {
	ID     int     `json:"id"`
	Name   string  `json:"name"`
	Repeat int     `json:"reps"`
	Set    int     `json:"sets"`
	Weight float32 `json:"weight"`
}
