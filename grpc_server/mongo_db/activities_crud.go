package mongodb

import (
	"context"
	"grpcserver/internals/models"
	"grpcserver/internals/utils"
	pb "grpcserver/proto/generated_files"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

func AddActivitiesToDB(ctx context.Context, request_activities []*pb.Activity) ([]*pb.Activity, error) {
	client := MongoClient

	newActivities := make([]*models.Activities, len(request_activities))
	for i, pbActivities := range request_activities {
		newActivities[i] = MapPbActivitiesToModelActivities(pbActivities)
		newActivities[i].CreatedAt = time.Now().Format("2006-01-02")
	}
	utils.Logger.Info("Activities to add", zap.Any("newActivities", newActivities))

	var addedActivities []*pb.Activity
	for _, activity := range newActivities {
		result, err := client.Database("data").Collection("tracked_activities").InsertOne(ctx, activity)
		if err != nil {
			return nil, utils.ErrorHandler(err, "Error adding value to database")
		}
		objectID, ok := result.InsertedID.(primitive.ObjectID)
		if ok {
			activity.Id = objectID.Hex()
		}
		pbActivity := MapModelActivitiesToPb(activity)
		addedActivities = append(addedActivities, pbActivity)
	}
	return addedActivities, nil
}

func GetActivitiessFromDb(ctx context.Context, sortOptions primitive.D, filter primitive.M) ([]*pb.Activity, error) {
	client := MongoClient
	coll := client.Database("data").Collection("tracked_activities")

	var cursor *mongo.Cursor
	var err error
	if len(sortOptions) < 1 {
		cursor, err = coll.Find(ctx, filter)
	} else {
		cursor, err = coll.Find(ctx, filter, options.Find().SetSort(sortOptions))
	}
	if err != nil {
		return nil, utils.ErrorHandler(err, "Internal Error")
	}

	defer func() {
		if err := cursor.Close(ctx); err != nil {
			utils.Logger.Error("Failed to close the cursor", zap.Error(err))
		}
	}()

	activities, err := decodeEntities(ctx, cursor, func() *pb.Activity { return &pb.Activity{} }, newModel)
	if err != nil {
		return nil, err
	}
	return activities, nil
}
