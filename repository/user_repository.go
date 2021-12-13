package repository

import (
	"context"

	"github.com/project-jano/user-service.go/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	DBFieldUserId       = "userId"
	DBFieldCertificates = "certificates"
)

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(collection *mongo.Collection) UserRepository {
	return UserRepository{collection: collection}
}

func (repository *UserRepository) FindUser(ctx context.Context, userId string) *model.User {
	var user *model.User
	filter := bson.D{{Key: DBFieldUserId, Value: userId}}
	findError := repository.collection.FindOne(ctx, filter).Decode(&user)

	if findError != nil {
		if findError != mongo.ErrNoDocuments {
			return nil
		}

		user = &model.User{
			UserId:       userId,
			Certificates: []model.UserCertificate{},
		}
	}

	return user
}

func (repository *UserRepository) UpdateCertificates(ctx context.Context, user *model.User) error {
	filter := bson.D{{Key: DBFieldUserId, Value: user.UserId}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: DBFieldCertificates, Value: user.Certificates}}}}
	opts := options.Update().SetUpsert(true)

	_, upsertErr := repository.collection.UpdateOne(ctx, filter, update, opts)

	return upsertErr
}
