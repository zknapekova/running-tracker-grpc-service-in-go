db = db.getSiblingDB('data');

db.createCollection('tracked_activities', {
  validator: {
    $jsonSchema: {
      bsonType: 'object',
      required: ['name', 'duration', 'date', 'trainers_model_id', 'created_at'],
      properties: {
        name: {
          bsonType: 'string',
          descrption: "activity name"
        },
        duration: {
          bsonType: 'int',
          description: "measured in minutes",
        },
        date: {
          bsonType: "date",
          description: "date of the activity",
        },
        trainers_model_id: {
          bsonType: 'objectId',
          description: "reference to the trainers collection"
        },
        created_at: {
          bsonType: "date",
          description: "date of creation",
        }
      }
    }
  }
});

db.createCollection('trainers', {
  validator: {
    $jsonSchema: {
      bsonType: 'object',
      required: ['brand', 'model'],
      properties: {
        brand: {
          bsonType: 'string',
        },
        model: {
          bsonType: 'string',
        }
      }
    }
  }
});


print('All collections created successfully');