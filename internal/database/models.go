package database

type User struct {
	Name  string `bson:"_id"`
	Count int16  `bson:"count"`
}

type SpinResult struct {
}
