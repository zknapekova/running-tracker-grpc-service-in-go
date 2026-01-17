db = db.getSiblingDB('data');

db.tracked_activities.createIndex({name: 1})

print('Index added successfully');