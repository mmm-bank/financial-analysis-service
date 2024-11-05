package storage

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/mmm-bank/financial-analysis-service/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"time"
)

var _ AnalysisStorage = &MongoAnalysis{}

type AnalysisStorage interface {
	GetExpensesAnalysis(ctx context.Context, userID uuid.UUID, monthYear string) (models.ExpensesAnalysis, error)
	GetIncomeAnalysis(ctx context.Context, userID uuid.UUID, monthYear string) (models.IncomeAnalysis, error)
}

type MongoAnalysis struct {
	expensesCollection *mongo.Collection
	incomeCollection   *mongo.Collection
}

func NewMongoAnalysis(mongoURL string) *MongoAnalysis {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURL))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	err = client.Ping(context.TODO(), readpref.Primary())
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	return &MongoAnalysis{
		expensesCollection: client.Database("financial_analysis").Collection("expenses"),
		incomeCollection:   client.Database("financial_analysis").Collection("income"),
	}
}

func (m *MongoAnalysis) GetExpensesAnalysis(ctx context.Context, userID uuid.UUID, monthYear string) (models.ExpensesAnalysis, error) {
	var analysis models.ExpensesAnalysis
	filter := bson.M{"user_id": userID, "month_year": monthYear}
	err := m.expensesCollection.FindOne(ctx, filter).Decode(&analysis)
	return analysis, err
}

func (m *MongoAnalysis) GetIncomeAnalysis(ctx context.Context, userID uuid.UUID, monthYear string) (models.IncomeAnalysis, error) {
	var analysis models.IncomeAnalysis
	filter := bson.M{"user_id": userID, "month_year": monthYear}
	err := m.incomeCollection.FindOne(ctx, filter).Decode(&analysis)
	return analysis, err
}

func (m *MongoAnalysis) AddTransfer(transfer *models.Transfer) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	now := time.Now()
	monthYear := fmt.Sprintf("%d-%d", now.Month(), now.Year())

	_, err := m.expensesCollection.UpdateOne(ctx,
		bson.M{"user_id": transfer.SenderID, "month_year": monthYear},
		bson.M{
			"$push": bson.M{"transactions": bson.M{
				"id":          transfer.ID,
				"receiver_id": transfer.ReceiverID,
				"category":    "transfer",
				"cost":        transfer.Amount,
				"created_at":  time.Now(),
			}},
			"$inc": bson.M{"total_cost": transfer.Amount, "transaction_count": 1},
		},
		options.Update().SetUpsert(true),
	)
	if err != nil && !mongo.IsDuplicateKeyError(err) {
		return err
	}

	_, err = m.incomeCollection.UpdateOne(ctx,
		bson.M{"user_id": transfer.ReceiverID, "month_year": monthYear},
		bson.M{
			"$push": bson.M{"transactions": bson.M{
				"id":         transfer.ID,
				"sender_id":  transfer.SenderID,
				"category":   "transfer",
				"cost":       transfer.Amount,
				"created_at": time.Now(),
			}},
			"$inc": bson.M{"total_cost": transfer.Amount, "transaction_count": 1},
		},
		options.Update().SetUpsert(true),
	)
	if err != nil && !mongo.IsDuplicateKeyError(err) {
		return err
	}
	return nil
}
