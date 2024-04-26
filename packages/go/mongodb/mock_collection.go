package mongodb

import (
	"context"
	"errors"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	UnexpectedType = errors.New("unexpected type")
)

type MockCollection struct {
	mock.Mock
}

func (c *MockCollection) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (response *mongo.InsertOneResult, e error) {
	// InsertOneResult is the result type returned by an InsertOne operation.
	// type InsertOneResult struct {
	// The _id of the inserted document. A value generated by the driver will be of type primitive.ObjectID.
	// InsertedID interface{}
	// }
	args := c.Called(ctx, document, opts)
	e = args.Error(1)
	firstResponseArg := args.Get(0)

	if firstResponseArg != nil {
		if v, ok := firstResponseArg.(*mongo.InsertOneResult); !ok {
			return nil, UnexpectedType
		} else {
			response = v
		}
	}
	return
}

func (c *MockCollection) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) (response *mongo.SingleResult) {
	// use mongo.NewSingleResultFromDocument tyo create the response for this mocked function
	args := c.Called(ctx, filter, opts)
	firstResponseArg := args.Get(0)

	if firstResponseArg != nil {
		if v, ok := firstResponseArg.(*mongo.SingleResult); !ok {
			panic(UnexpectedType)
		} else {
			response = v
		}
	}

	return
}

func (c *MockCollection) DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (response *mongo.DeleteResult, e error) {
	// DeleteResult is the result type returned by DeleteOne and DeleteMany operations.
	// type DeleteResult struct {
	//	 DeletedCount int64 `bson:"n"` // The number of documents deleted.
	// }
	args := c.Called(ctx, filter, opts)
	e = args.Error(1)
	firstResponseArg := args.Get(0)

	if firstResponseArg != nil {
		if v, ok := firstResponseArg.(*mongo.DeleteResult); !ok {
			return nil, UnexpectedType
		} else {
			response = v
		}
	}
	return
}

func (c *MockCollection) UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (response *mongo.UpdateResult, e error) {
	// UpdateResult is the result type returned from UpdateOne, UpdateMany, and ReplaceOne operations.
	// type UpdateResult struct {
	//   MatchedCount  int64 // The number of documents matched by the filter.
	//	 ModifiedCount int64 // The number of documents modified by the operation.
	//	 UpsertedCount int64 // The number of documents upserted by the operation.
	//	 UpsertedID    interface{} // The _id field of the upserted document, or nil if no upsert was done.
	// }
	args := c.Called(ctx, filter, opts)
	e = args.Error(1)
	firstResponseArg := args.Get(0)

	if firstResponseArg != nil {
		if v, ok := firstResponseArg.(*mongo.UpdateResult); !ok {
			return nil, UnexpectedType
		} else {
			response = v
		}
	}
	return
}

func (c *MockCollection) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (response *mongo.Cursor, e error) {
	// use mongo.NewCursorFromDocuments to provide the result
	args := c.Called(ctx, filter, opts)
	e = args.Error(1)
	firstResponseArg := args.Get(0)

	if firstResponseArg != nil {
		if v, ok := firstResponseArg.(*mongo.Cursor); !ok {
			return nil, UnexpectedType
		} else {
			response = v
		}
	}
	return
}

func (c *MockCollection) CountDocuments(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (n int64, e error) {
	args := c.Called(ctx, filter, opts)
	e = args.Error(1)
	firstResponseArg := args.Get(0)

	if firstResponseArg != nil {
		if v, ok := firstResponseArg.(int64); !ok {
			return 0, UnexpectedType
		} else {
			n = v
		}
	}
	return
}
