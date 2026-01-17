package models

type Activities struct {
	Id            string `protobuf:"id,omitempty" bson:"_id,omitempty"`
	Name          string `protobuf:"name,omitempty" bson:"name,omitempty"`
	Duration      int32  `protobuf:"duration,omitempty" bson:"duration,omitempty"`
	Distance      int32  `protobuf:"distance,omitempty" bson:"distance,omitempty"`
	Pace          int32  `protobuf:"pace,omitempty" bson:"pace,omitempty"`
	Date          string `protobuf:"date,omitempty" bson:"date,omitempty"`
	TrainersBrand string `protobuf:"trainers_brand,omitempty" bson:"trainers_brand,omitempty"`
	TrainersModel string `protobuf:"trainers_model,omitempty" bson:"trainers_model,omitempty"`
	CreatedAt     string `protobuf:"created_at,omitempty" bson:"created_at,omitempty"`
}
