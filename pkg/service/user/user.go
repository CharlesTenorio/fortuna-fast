package user

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strconv"

	"github.com/go-chi/jwtauth"
	"github.com/katana/fortuna/backend-go/pkg/adapter/mongodb"
	"github.com/katana/fortuna/backend-go/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type TokenResponse struct {
	Token string `json:"token"`
}

type UserServiceInterface interface {
	GetAll(ctx context.Context, filters model.FilterUsuario, limit, page int64) (*model.Paginate, error)
	GetByID(ctx context.Context, ID string) (*model.Usuario, error)
	GetByEmail(ctx context.Context, email string) (user *model.Usuario, err error)
	Create(ctx context.Context, user *model.Usuario) (*model.Usuario, error)
	Update(ctx context.Context, ID string, userToChange *model.Usuario) (bool, error)
	ChangePassword(ctx context.Context, currentPassword, newPassword string, userRequest *model.Usuario) error
}

type UserDataService struct {
	mdb         mongodb.MongoDBInterface
	Jwt         *jwtauth.JWTAuth
	JwtExpirado int
}

func NewUsuarioservice(mongo_connection mongodb.MongoDBInterface) *UserDataService {
	return &UserDataService{
		mdb: mongo_connection,
	}
}

func (uds *UserDataService) GetAll(ctx context.Context, filters model.FilterUsuario, limit, page int64) (*model.Paginate, error) {

	collection := uds.mdb.GetCollection("usuarios")

	query := bson.M{}

	if filters.Nome != "" || filters.Email != "" || filters.Enable != "" {

		if filters.Nome != "" {
			query["nome"] = bson.M{"$regex": filters.Nome, "$options": "i"}
		}

		if filters.Email != "" {
			query["email"] = bson.M{"$regex": filters.Email, "$options": "i"}
		}

		if filters.Enable != "" {
			enable, err := strconv.ParseBool(filters.Enable)
			if err != nil {
				log.Println(err)
				return nil, err
			}
			query["enable"] = enable
		}
	}

	count, err := collection.CountDocuments(ctx, query, &options.CountOptions{})

	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	paginate := model.NewPaginate(limit, page, count)

	curr, err := collection.Find(ctx, query, paginate.GetPaginatedOpts())
	if err != nil {
		return nil, err
	}

	result := make([]*model.Usuario, 0)
	for curr.Next(ctx) {
		user := &model.Usuario{}
		if err := curr.Decode(user); err != nil {
			log.Println(err.Error())
		}

		result = append(result, user)
	}

	paginate.Paginate(result)

	return paginate, nil
}

func (uds *UserDataService) GetByID(ctx context.Context, ID string) (*model.Usuario, error) {

	collection := uds.mdb.GetCollection("usuarios")

	user := &model.Usuario{}

	objectID, err := objectIDFromHex(ID)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	filter := bson.D{

		{Key: "_id", Value: objectID},
	}

	err = collection.FindOne(ctx, filter).Decode(user)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return user, nil
}

func (uds *UserDataService) GetByEmail(ctx context.Context, email string) (user *model.Usuario, err error) {

	collection := uds.mdb.GetCollection("usuarios")

	filter := bson.D{

		{Key: "email", Value: email},
	}

	err = collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return user, nil
}

func (uds *UserDataService) Create(ctx context.Context, user *model.Usuario) (*model.Usuario, error) {

	collection := uds.mdb.GetCollection("usuarios")
	usr, _ := model.NewUsuario(user.Nome, user.Senha, user.Email, user.Role)
	result, err := collection.InsertOne(ctx, usr)
	if err != nil {
		log.Println(err.Error())
		return usr, err
	}

	usr.ID = result.InsertedID.(primitive.ObjectID)

	return usr, nil
}

func (uds *UserDataService) Update(ctx context.Context, ID string, userToChange *model.Usuario) (bool, error) {

	collection := uds.mdb.GetCollection("usuarios")

	opts := options.Update().SetUpsert(true)

	objectID, err := objectIDFromHex(ID)
	if err != nil {
		return false, err
	}

	filter := bson.D{

		{Key: "_id", Value: objectID},
	}

	values := bson.D{
		{Key: "email", Value: userToChange.Email},
		{Key: "nome", Value: userToChange.Nome},
		{Key: "updated_at", Value: userToChange.UpdatedAt},
	}

	update := bson.D{{Key: "$set", Value: values}}

	_, err = collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		log.Println("Error while updating data")
		return false, err
	}

	return true, nil
}

func objectIDFromHex(hex string) (objectID primitive.ObjectID, err error) {
	objectID, err = primitive.ObjectIDFromHex(hex)
	if err != nil {
		log.Println(err.Error())
		return objectID, err
	}
	return objectID, nil
}

func (uds *UserDataService) ChangePassword(ctx context.Context, currentPassword, newPassword string, userRequest *model.Usuario) error {
	collection := uds.mdb.GetCollection("usuarios")
	opts := options.Update().SetUpsert(true)

	objectID, err := primitive.ObjectIDFromHex(userRequest.ID.Hex())
	if err != nil {
		log.Println(err.Error())
		return err
	}

	filter := bson.D{

		{Key: "_id", Value: objectID},
	}

	var user *model.Usuario

	if err := collection.FindOne(ctx, filter).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			log.Println("User not found")
			return err
		}
		log.Println(err.Error())
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Senha), []byte(currentPassword)); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			log.Println("The currentPassoword do not correspond")
			return err
		}
		log.Println("Error to compare the passwords")
		return err
	}

	if len(newPassword) < 8 {
		log.Println("New password length is too short")
		return fmt.Errorf("New password length is too short")
	}

	const (
		uppercasePattern = "^(.*[A-Z]).*$"
		numberPattern    = "^(.*[0-9]).*$"
		symbolPattern    = "^(.*[!@#$%^&*()\\-_+=]).*$"
	)

	matchCapital, errCapital := regexp.MatchString(uppercasePattern, newPassword)
	if errCapital != nil || !matchCapital {
		log.Println("New password must contain at least one uppercase letter")
		return errCapital
	}

	matchNumber, errNumber := regexp.MatchString(numberPattern, newPassword)
	if errNumber != nil || !matchNumber {
		log.Println("New password must contain at least one number")
		return errNumber
	}

	matchSymbol, errSymbol := regexp.MatchString(symbolPattern, newPassword)
	if errSymbol != nil || !matchSymbol {
		log.Println("New password must contain at least one symbol")
		return errSymbol
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error generating hashed password")
		return err
	}

	update := bson.D{{Key: "$set",
		Value: bson.D{
			{Key: "senha", Value: string(hashedPassword)},
		},
	}}

	result, err := collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		log.Println("Error while updating password")
		return err
	}

	if result.ModifiedCount == 0 {
		log.Println("Password not updated")
		return fmt.Errorf("password not updated")
	}

	return nil
}
