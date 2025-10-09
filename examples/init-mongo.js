// MongoDB initialization script for examples
// This script creates a database and some sample data for testing migrations

// Switch to the examples database
db = db.getSiblingDB('migration_examples');

// Create sample users collection with some test data
db.users.insertMany([
  {
    _id: ObjectId(),
    email: "JOHN.DOE@EXAMPLE.COM", // Intentionally uppercase to test normalization
    first_name: "John",
    last_name: "Doe",
    status: "active",
    created_at: new Date("2024-01-01T10:00:00Z")
  },
  {
    _id: ObjectId(),
    email: "jane.smith@EXAMPLE.COM", // Mixed case
    first_name: "Jane",
    last_name: "Smith",
    status: "inactive",
    created_at: new Date("2024-01-02T14:30:00Z")
  },
  {
    _id: ObjectId(),
    email: "bob.johnson@example.com",
    first_name: "Bob",
    last_name: "Johnson",
    status: "active",
    created_at: new Date("2024-01-03T09:15:00Z"),
    updated_at: new Date("2024-01-03T09:15:00Z") // This one already has updated_at
  },
  {
    _id: ObjectId(),
    email: "ALICE.WILLIAMS@EXAMPLE.COM",
    first_name: "Alice",
    // last_name intentionally missing to test the migration logic
    status: "pending",
    created_at: new Date("2024-01-04T16:45:00Z")
  },
  {
    _id: ObjectId(),
    email: "charlie.brown@example.com",
    // first_name intentionally missing
    last_name: "Brown",
    status: "active",
    created_at: new Date("2024-01-05T11:20:00Z")
  }
]);

print("✅ Created migration_examples database with sample users");
print("📊 Inserted " + db.users.countDocuments({}) + " sample users");

// Show what we created
print("\n📋 Sample data overview:");
db.users.find({}, {email: 1, first_name: 1, last_name: 1, status: 1}).forEach(printjson);

print("\n🚀 Ready for migration examples!");