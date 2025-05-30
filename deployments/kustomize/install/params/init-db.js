const mongoHost = process.env.AMBULANCE_API_MONGODB_HOST
const mongoPort = process.env.AMBULANCE_API_MONGODB_PORT
const mongoUser = process.env.AMBULANCE_API_MONGODB_USERNAME
const mongoPassword = process.env.AMBULANCE_API_MONGODB_PASSWORD
const database = process.env.AMBULANCE_API_MONGODB_DATABASE

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

// Check if database exists and get collections
const databases = connection.getDBNames()
const db = connection.getDB(database)
let existingCollections = [];

if (databases.includes(database)) {
    existingCollections = db.getCollectionNames();
    print(`Database '${database}' exists with collections: ${existingCollections}`)
    
    // Check if both collections already exist and are properly initialized
    if (existingCollections.includes("equipment") && existingCollections.includes("orders")) {
        // Check if collections have data (basic validation)
        const equipmentCount = db.equipment.countDocuments();
        const ordersCount = db.orders.countDocuments();
        
        if (equipmentCount > 0 && ordersCount > 0) {
            print(`Collections already exist and contain data (equipment: ${equipmentCount}, orders: ${ordersCount})`)
            print(`Database '${database}' is already initialized`)
            process.exit(0);
        } else {
            print(`Collections exist but are empty, proceeding with initialization...`)
        }
    }
} else {
    print(`Database '${database}' does not exist, creating...`)
}

// Initialize equipment collection
if (!existingCollections.includes("equipment")) {
    print("Creating equipment collection...")
    db.createCollection("equipment")
    
    // Create indexes for equipment
    db.equipment.createIndex({ "id": 1 })
    db.equipment.createIndex({ "serialNumber": 1 }, { unique: true })
    db.equipment.createIndex({ "status": 1 })
    db.equipment.createIndex({ "location": 1 })
    db.equipment.createIndex({ "manufacturer": 1 })
} else {
    print("Equipment collection already exists, skipping creation...")
}

const equipmentCount = db.equipment.countDocuments();
if (equipmentCount === 0) {
    print("Inserting equipment sample data...")
    let result = db.equipment.insertMany([
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
        print(`Error when writing equipment data: ${result.errmsg}`)
    } else {
        print(`Successfully inserted ${result.insertedIds.length} equipment records`)
    }
} else {
    print(`Equipment collection already contains ${equipmentCount} records, skipping data insertion...`)
}

try {
    db.equipment.createIndex({ "nextService": 1 });
    db.equipment.createIndex({ "installationDate": 1 }); 
    db.equipment.createIndex({ "name": "text", "manufacturer": "text", "model": "text", "notes": "text" }); // ✅ Fixed: use "equipment" directly
    print("Equipment indexes created successfully")
} catch (e) {
    print(`Equipment indexes may already exist: ${e}`)
}

// Initialize orders collection
if (!existingCollections.includes("orders")) {
    print("Creating orders collection...")
    db.createCollection("orders")
    
    // Create basic orders indexes
    db.orders.createIndex({ "id": 1 })
    db.orders.createIndex({ "status": 1 })
    db.orders.createIndex({ "requestedBy": 1 })
    db.orders.createIndex({ "requestorDepartment": 1 })
    db.orders.createIndex({ "createdAt": 1 })
} else {
    print("Orders collection already exists, skipping creation...")
}

const ordersCount = db.orders.countDocuments();
if (ordersCount === 0) {
    print("Inserting orders sample data...")
    let ordersResult = db.orders.insertMany([
        {
            "id": "ord-123e4567-e89b-12d3-a456-426614174000",
            "items": [
                {
                    "equipmentName": "Surgical Gloves (Box of 100)",
                    "quantity": 50,
                    "unitPrice": 25.99,
                    "totalPrice": 1299.50
                },
                {
                    "equipmentName": "Face Masks N95 (Box of 20)",
                    "quantity": 30,
                    "unitPrice": 45.00,
                    "totalPrice": 1350.00
                }
            ],
            "requestedBy": "Dr. Sarah Johnson",
            "requestorDepartment": "Surgery",
            "status": "pending",
            "notes": "Urgent order for upcoming surgery schedule",
            "createdAt": new Date("2025-05-20T09:30:00Z"),
            "updatedAt": new Date("2025-05-20T09:30:00Z")
        },
        {
            "id": "ord-223e4567-e89b-12d3-a456-426614174001",
            "items": [
                {
                    "equipmentName": "Digital Thermometer",
                    "quantity": 25,
                    "unitPrice": 89.99,
                    "totalPrice": 2249.75
                }
            ],
            "requestedBy": "Nurse Manager Lisa Chen",
            "requestorDepartment": "Emergency Department",
            "status": "sent",
            "notes": "Replacement for faulty thermometers",
            "createdAt": new Date("2025-05-18T14:15:00Z"),
            "updatedAt": new Date("2025-05-19T10:00:00Z")
        },
        {
            "id": "ord-323e4567-e89b-12d3-a456-426614174002",
            "items": [
                {
                    "equipmentName": "Blood Pressure Cuff",
                    "quantity": 15,
                    "unitPrice": 125.50,
                    "totalPrice": 1882.50
                },
                {
                    "equipmentName": "Stethoscope",
                    "quantity": 10,
                    "unitPrice": 180.00,
                    "totalPrice": 1800.00
                },
                {
                    "equipmentName": "Pulse Oximeter",
                    "quantity": 20,
                    "unitPrice": 95.75,
                    "totalPrice": 1915.00
                }
            ],
            "requestedBy": "Dr. Michael Rodriguez",
            "requestorDepartment": "Cardiology",
            "status": "delivered",
            "notes": "Annual equipment refresh for cardiology department",
            "createdAt": new Date("2025-05-15T11:20:00Z"),
            "updatedAt": new Date("2025-05-20T16:45:00Z")
        },
        {
            "id": "ord-423e4567-e89b-12d3-a456-426614174003",
            "items": [
                {
                    "equipmentName": "Wheelchair",
                    "quantity": 5,
                    "unitPrice": 450.00,
                    "totalPrice": 2250.00
                }
            ],
            "requestedBy": "Physical Therapy Coordinator",
            "requestorDepartment": "Rehabilitation",
            "status": "pending",
            "notes": "New wheelchairs for patient mobility program",
            "createdAt": new Date("2025-05-21T08:00:00Z"),
            "updatedAt": new Date("2025-05-21T08:00:00Z")
        },
        {
            "id": "ord-523e4567-e89b-12d3-a456-426614174004",
            "items": [
                {
                    "equipmentName": "Disposable Syringes (Pack of 100)",
                    "quantity": 100,
                    "unitPrice": 35.99,
                    "totalPrice": 3599.00
                },
                {
                    "equipmentName": "IV Bags 0.9% Saline",
                    "quantity": 200,
                    "unitPrice": 12.50,
                    "totalPrice": 2500.00
                }
            ],
            "requestedBy": "Pharmacy Director",
            "requestorDepartment": "Pharmacy",
            "status": "cancelled",
            "notes": "Order cancelled due to supplier issues",
            "createdAt": new Date("2025-05-10T13:30:00Z"),
            "updatedAt": new Date("2025-05-12T09:15:00Z")
        }
    ]);

    if (ordersResult.writeError) {
        console.error(ordersResult)
        print(`Error when writing orders data: ${ordersResult.errmsg}`)
    } else {
        print(`Successfully inserted ${ordersResult.insertedIds.length} order records`)
    }
} else {
    print(`Orders collection already contains ${ordersCount} records, skipping data insertion...`)
}

try {
    db.orders.createIndex({ "updatedAt": 1 })
    db.orders.createIndex({ "items.equipmentName": 1 })
    print("Orders indexes created successfully")
} catch (e) {
    print(`Orders indexes may already exist: ${e}`)
}

print(`Database '${database}' initialized successfully`)
print(`Created collections: 'equipment' and 'orders'`)

// exit with success
process.exit(0);