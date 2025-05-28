const mongoHost = process.env.AMBULANCE_API_MONGODB_HOST
const mongoPort = process.env.AMBULANCE_API_MONGODB_PORT

const mongoUser = process.env.AMBULANCE_API_MONGODB_USERNAME
const mongoPassword = process.env.AMBULANCE_API_MONGODB_PASSWORD

const database = process.env.AMBULANCE_API_MONGODB_DATABASE
const collection = process.env.AMBULANCE_API_MONGODB_COLLECTION

const retrySeconds = parseInt(process.env.RETRY_CONNECTION_SECONDS || "5") || 5;

// try to connect to mongoDB until it is not available
let connection;
while(true) {
    try {
        connection = Mongo(`mongodb://${mongoUser}:${mongoPassword}@${mongoHost}:${mongoPort}`);
        break;
    } catch (exception) {
        print(`Cannot connect to mongoDB: ${exception}`);
        print(`Will retry after ${retrySeconds} seconds`)
        sleep(retrySeconds * 1000);
    }
}

// if database and collection exists, exit with success - already initialized
const databases = connection.getDBNames()
if (databases.includes(database)) {
    const dbInstance = connection.getDB(database)
    collections = dbInstance.getCollectionNames()
    if (collections.includes(collection)) {
            print(`Collection '${collection}' already exists in database '${database}'`)
        process.exit(0);
    }
}

// initialize
// create database and collection
const db = connection.getDB(database)
db.createCollection(collection)

// create indexes
db[collection].createIndex({ "id": 1 })
db[collection].createIndex({ "serialNumber": 1 }, { unique: true })
db[collection].createIndex({ "status": 1 })
db[collection].createIndex({ "location": 1 })
db[collection].createIndex({ "manufacturer": 1 })

// insert sample equipment data
let result = db[collection].insertMany([
    {
        "id": "123e4567-e89b-12d3-a456-426614174000",
        "name": "MRI Scanner",
        "serialNumber": "MRI-2023-001",
        "manufacturer": "GE Healthcare",
        "model": "SIGNA Architect",
        "installationDate": "2023-01-15",
        "location": "Radiology Department, Room 3",
        "serviceInterval": 90,
        "lastService": "2023-03-15",
        "nextService": "2023-06-13",
        "lifeExpectancy": 10,
        "status": "operational",
        "notes": "3T MRI scanner, installed by technician John Doe"
    },
    {
        "id": "223e4567-e89b-12d3-a456-426614174001",
        "name": "X-Ray Machine",
        "serialNumber": "XR-2022-105",
        "manufacturer": "Siemens Healthineers",
        "model": "Ysio Max",
        "installationDate": "2022-06-20",
        "location": "Emergency Department, Room 5",
        "serviceInterval": 60,
        "lastService": "2023-04-20",
        "nextService": "2023-06-19",
        "lifeExpectancy": 8,
        "status": "operational",
        "notes": "Digital X-ray system"
    },
    {
        "id": "323e4567-e89b-12d3-a456-426614174002",
        "name": "CT Scanner",
        "serialNumber": "CT-2021-087",
        "manufacturer": "Philips Healthcare",
        "model": "Ingenuity CT",
        "installationDate": "2021-09-10",
        "location": "Radiology Department, Room 1",
        "serviceInterval": 120,
        "lastService": "2023-02-10",
        "nextService": "2023-06-10",
        "lifeExpectancy": 12,
        "status": "in_repair",
        "notes": "128-slice CT scanner, currently undergoing maintenance"
    },
    {
        "id": "423e4567-e89b-12d3-a456-426614174003",
        "name": "Ultrasound Machine",
        "serialNumber": "US-2023-042",
        "manufacturer": "Canon Medical Systems",
        "model": "Aplio i800",
        "installationDate": "2023-02-28",
        "location": "Cardiology Department, Room 2",
        "serviceInterval": 30,
        "lastService": "2023-04-28",
        "nextService": "2023-05-28",
        "lifeExpectancy": 7,
        "status": "operational",
        "notes": "High-end ultrasound system for cardiac imaging"
    },
    {
        "id": "523e4567-e89b-12d3-a456-426614174004",
        "name": "Ventilator",
        "serialNumber": "VENT-2020-156",
        "manufacturer": "Medtronic",
        "model": "Puritan Bennett 980",
        "installationDate": "2020-04-15",
        "location": "ICU, Bed 3",
        "serviceInterval": 14,
        "lastService": "2023-05-01",
        "nextService": "2023-05-15",
        "lifeExpectancy": 10,
        "status": "operational",
        "notes": "Critical care ventilator for ICU patients"
    },
    {
        "id": "623e4567-e89b-12d3-a456-426614174005",
        "name": "Defibrillator",
        "serialNumber": "DEF-2019-234",
        "manufacturer": "Zoll Medical",
        "model": "R Series Plus",
        "installationDate": "2019-08-22",
        "location": "Emergency Department, Trauma Bay 1",
        "serviceInterval": 30,
        "lastService": "2023-04-22",
        "nextService": "2023-05-22",
        "lifeExpectancy": 12,
        "status": "faulty",
        "notes": "Battery replacement needed, scheduled for repair"
    },
    {
        "id": "723e4567-e89b-12d3-a456-426614174006",
        "name": "Anesthesia Machine",
        "serialNumber": "ANES-2018-078",
        "manufacturer": "Dräger",
        "model": "Perseus A500",
        "installationDate": "2018-11-12",
        "location": "Operating Room 2",
        "serviceInterval": 45,
        "lastService": "2023-03-12",
        "nextService": "2023-04-26",
        "lifeExpectancy": 15,
        "status": "operational",
        "notes": "Multi-purpose anesthesia workstation"
    },
    {
        "id": "823e4567-e89b-12d3-a456-426614174007",
        "name": "Patient Monitor",
        "serialNumber": "MON-2017-191",
        "manufacturer": "Philips Healthcare",
        "model": "IntelliVue MX800",
        "installationDate": "2017-05-30",
        "location": "ICU, Bed 1",
        "serviceInterval": 60,
        "lastService": "2023-01-30",
        "nextService": "2023-03-31",
        "lifeExpectancy": 10,
        "status": "decommissioned",
        "notes": "End of life, replaced with newer model"
    }
]);

if (result.writeError) {
    console.error(result)
    print(`Error when writing the data: ${result.errmsg}`)
} else {
    print(`Successfully inserted ${result.insertedIds.length} equipment records`)
}

// create additional indexes for better query performance
db[collection].createIndex({ "nextService": 1 }) // for service scheduling
db[collection].createIndex({ "installationDate": 1 }) // for age analysis
db[collection].createIndex({ "name": "text", "manufacturer": "text", "model": "text", "notes": "text" }) // for text search

print(`Database '${database}' and collection '${collection}' initialized successfully`)
print(`Created indexes on: id, serialNumber (unique), status, location, manufacturer, nextService, installationDate, and text search`)

// exit with success
process.exit(0);