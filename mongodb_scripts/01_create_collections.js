db = db.getSiblingDB('data');

db.createCollection('tracked_activities', {
  validator: {
    $jsonSchema: {
      bsonType: 'object',
      required: ['name', 'duration', 'date', 'created_at'],
      properties: {
        name: {
          bsonType: 'string',
          description: "activity name"
        },
        duration: {
          bsonType: 'int',
          description: "measured in minutes",
        },
        date: {
          bsonType: "string",
          description: "date of the activity",
        },
        created_at: {
          bsonType: "string",
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