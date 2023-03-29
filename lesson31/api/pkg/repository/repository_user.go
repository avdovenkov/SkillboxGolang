package repository

import (
	"context"
	"fmt"
	"http_service/pkg/user"
	"sync"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RepositoryUserMongo struct {
	sync.Mutex
	Client *mongo.Client `json:"storage"`
}
type idGenerator struct {
	key string `json:"key"`
	N   int64  `bson:"id"`
}
type FriendsList struct {
	Friends []*user.Friend `json:"friends"`
}

func NewRepositoryUser(client *mongo.Client) *RepositoryUserMongo {
	return &RepositoryUserMongo{Client: client}
}
func (r *RepositoryUserMongo) GetUser(ctx context.Context, id int64) (*user.User, error) {
	users := r.Client.Database("Users").Collection("users")
	user := user.User{}
	err := users.FindOne(ctx, bson.D{{"id", id}}).Decode(&user)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return &user, nil
}
func (r *RepositoryUserMongo) UpdateUser(ctx context.Context, user *user.User) error {
	r.Lock()
	users := r.Client.Database("Users").Collection("users")
	var updatedDocument bson.M
	filter := bson.D{{"id", user.Id}}
	update := bson.D{{"$set", bson.D{{"name", user.Name}, {"age", user.Age}, {"friends", user.Friends}}}}
	err := users.FindOneAndUpdate(ctx, filter, update).Decode(&updatedDocument)
	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in
		// the collection.
		if err == mongo.ErrNoDocuments {
			return err
		}
		logrus.Error(err)
	}
	fmt.Printf("updated document %v", updatedDocument)
	r.Unlock()
	return nil
}
func (r *RepositoryUserMongo) CreateUser(ctx context.Context, u *user.User) error {
	newid, err := r.getNextIdUser(ctx, r.Client, "counter")
	if err != nil {
		logrus.Fatal(err)
	}
	users := r.Client.Database("Users").Collection("users")
	res, err := users.InsertOne(ctx, bson.D{primitive.E{Key: "id", Value: newid}, primitive.E{Key: "name", Value: u.Name}, primitive.E{Key: "age", Value: u.Age}, primitive.E{Key: "friends", Value: u.Friends}})
	if err != nil {
		return err
	}
	fmt.Printf("inserted document with ID %v\n", res.InsertedID)
	return nil
}
func (r *RepositoryUserMongo) getNextIdUser(ctx context.Context, client *mongo.Client, key string) (int64, error) {
	counter := client.Database("Users").Collection("counter")
	//var result bson.M
	generator := &idGenerator{}
	filter := bson.D{{"counter", key}}
	update := bson.D{{"$inc", bson.D{{"id", 1}}}}
	err := counter.FindOneAndUpdate(ctx, filter, update).Decode(&generator)
	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in
		// the collection.
		if err == mongo.ErrNoDocuments {
			_, err := counter.InsertOne(ctx, bson.M{"counter": "counter", "id": 1})
			if err != nil {
				return -1, err
			}
			return 1, nil
		} else if err != nil {
			return -1, err

		}
	}
	fmt.Printf("updated document %v", r)
	return generator.N, nil
}
func (r *RepositoryUserMongo) GetUserFriend(ctx context.Context, id int64) ([]*user.Friend, error) {
	users := r.Client.Database("Users").Collection("users")
	u := user.User{}
	friends := FriendsList{}
	err := users.FindOne(ctx, bson.D{{"id", id}}).Decode(&u)
	for _, friendid := range u.Friends {
		one_friend := user.Friend{}
		err = users.FindOne(ctx, bson.D{{"id", friendid}}).Decode(&one_friend)
		if err != nil {
			return nil, err
		}
		friends.Friends = append(friends.Friends, &one_friend)
	}
	if err != nil {
		return nil, err
	}
	return friends.Friends, nil
}
func (r *RepositoryUserMongo) DeleteUser(ctx context.Context, id int64) error {
	r.Lock()
	users := r.Client.Database("Users").Collection("users")
	cur, err := users.Find(ctx, bson.D{})
	if err != nil {
		logrus.Fatal(err)
	}
	defer cur.Close(context.Background())
	for cur.Next(context.Background()) {
		user := user.User{}
		err := cur.Decode(&user)
		if err != nil {
			logrus.Fatal(err)
		}
		for i, one_user_id := range user.Friends {
			if one_user_id == id {
				user.Friends = append(user.Friends[:i], user.Friends[i+1:]...)
				var updatedDocument bson.M
				filter := bson.D{{"id", user.Id}}
				update := bson.D{{"$set", bson.D{{"friends", user.Friends}}}}
				err := users.FindOneAndUpdate(ctx, filter, update).Decode(&updatedDocument)
				if err != nil {
					// ErrNoDocuments means that the filter did not match any documents in
					// the collection.
					if err == mongo.ErrNoDocuments {
						return err
					}
					logrus.Fatal(err)
				}
				fmt.Printf("updated document %v", updatedDocument)

			}
		}
	}
	res, err := users.DeleteOne(ctx, bson.D{{"id", id}})
	if err != nil {
		logrus.Fatal(err)
	}
	fmt.Printf("deleted %v documents\n", res.DeletedCount)
	r.Unlock()
	return nil
}
