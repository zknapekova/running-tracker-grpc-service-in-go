package models

import pb "grpcserver/proto/generated_files"

type Trainers struct {
	Id				string					`protobuf:"id,omitempty" bson:"_id,omitempty"`
	Brand            string                 `protobuf:"brand,omitempty" bson:"brand,omitempty"`
	Model            string                 `protobuf:"model,omitempty" bson:"model,omitempty"`
	PurchaseDate     string                 `protobuf:"purchase_date,omitempty" bson:"purchase_date,omitempty"`
	ExpectedLifespan int64                  `protobuf:"expected_lifespan,omitempty" bson:"expected_lifespan,omitempty"`
	SurfaceType      pb.SurfaceType         `protobuf:"surface_type,omitempty" bson:"surface_type,omitempty"`
	Status           pb.TrainerStatus        `protobuf:"status,omitempty" bson:"status,omitempty"`
}