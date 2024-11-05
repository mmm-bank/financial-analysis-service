db = db.getSiblingDB("financial_analysis");

db.createCollection("expenses");
db.createCollection("income");

db.expenses.createIndex({ "user_id": 1, "month_year": 1 }, { unique: true });
db.expenses.createIndex({ "transactions.id": 1 }, { unique: true });

db.income.createIndex({ "user_id": 1, "month_year": 1 }, { unique: true });
db.income.createIndex({ "transactions.id": 1 }, { unique: true });
