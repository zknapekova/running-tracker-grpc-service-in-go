db = db.getSiblingDB('data');

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

db.createCollection('activities', {
  validator: {
    $jsonSchema: {
      bsonType: 'object',
      required: ['type', 'distance', 'time', 'trainers_model_id',],
      properties: {
        type: {
          bsonType: 'string',
        },
        distance: {
          bsonType: 'double',
        },
        time: {
          bsonType: 'double'
        },
        trainers_model_id: {
          bsonType: 'objectId',
        },
      }
    }
  }
});


print('All collections created successfully');