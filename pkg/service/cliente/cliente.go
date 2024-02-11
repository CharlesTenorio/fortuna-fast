package cliente

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/katana/fortuna/backend-go/internal/config/logger"
	"github.com/katana/fortuna/backend-go/pkg/adapter/mongodb"
	"github.com/katana/fortuna/backend-go/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ClienteServiceInterface interface {
	Create(ctx context.Context, Cliente model.Cliente) (*model.Cliente, error)
	Update(ctx context.Context, ID string, clienteToChange *model.Cliente) (bool, error)
	GetByID(ctx context.Context, ID string) (*model.Cliente, error)
	GetAll(ctx context.Context, filters model.FilterCliente, limit, page int64) (*model.Paginate, error)
	GetByDocumento(ctx context.Context, Documento string) bool
}

type ClienteDataService struct {
	mdb mongodb.MongoDBInterface
}

func NewClienteervice(mongo_connection mongodb.MongoDBInterface) *ClienteDataService {
	return &ClienteDataService{
		mdb: mongo_connection,
	}
}

func (cat *ClienteDataService) Create(ctx context.Context, Cliente model.Cliente) (*model.Cliente, error) {
	collection := cat.mdb.GetCollection("cfStore")
	cli := model.NewCliente(Cliente)
	result, err := collection.InsertOne(ctx, cli)
	if err != nil {
		logger.Error("erro salvar  Cliente", err)
		return &Cliente, err
	}

	cli.ID = result.InsertedID.(primitive.ObjectID)

	return cli, nil
}

func (cat *ClienteDataService) Update(ctx context.Context, ID string, Cliente *model.Cliente) (bool, error) {
	collection := cat.mdb.GetCollection("cfStore")

	opts := options.Update().SetUpsert(true)

	objectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {

		logger.Error("Error to parse ObjectIDFromHex", err)
		return false, err
	}

	filter := bson.D{

		{Key: "_id", Value: objectID},
		{Key: "data_type", Value: "cliente"},
	}

	update := bson.D{{Key: "$set",
		Value: bson.D{
			{Key: "nome", Value: Cliente.Nome},
			{Key: "enabled", Value: Cliente.Enabled},
			{Key: "updated_at", Value: time.Now().Format(time.RFC3339)},
		},
	}}

	_, err = collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		logger.Error("Error while updating data", err)

		return false, err
	}

	return true, nil
}

func (cat *ClienteDataService) GetByID(ctx context.Context, ID string) (*model.Cliente, error) {

	collection := cat.mdb.GetCollection("cfStore")

	Cliente := &model.Cliente{}

	objectID, err := primitive.ObjectIDFromHex(ID)
	if err != nil {

		logger.Error("Error to parse ObjectIDFromHex", err)
		return nil, err
	}

	filter := bson.D{
		{Key: "data_type", Value: "cliente"},
		{Key: "_id", Value: objectID},
	}

	err = collection.FindOne(ctx, filter).Decode(Cliente)
	if err != nil {
		logger.Error("erro ao consultar Cliente", err)
		return nil, err
	}

	return Cliente, nil
}

func (cat *ClienteDataService) GetAll(ctx context.Context, filters model.FilterCliente, limit, page int64) (*model.Paginate, error) {
	collection := cat.mdb.GetCollection("cfStore")

	query := bson.M{"data_type": "cliente"}

	if filters.Nome != "" || filters.Enabled != "" {
		if filters.Nome != "" {
			query["nome"] = bson.M{"$regex": fmt.Sprintf(".*%s.*", filters.Nome), "$options": "i"}
		}
		if filters.Enabled != "" {
			enable, err := strconv.ParseBool(filters.Enabled)
			if err != nil {
				logger.Error("erro converter campo enabled", err)
				return nil, err
			}
			query["enabled"] = enable
		}
	}
	count, err := collection.CountDocuments(ctx, query, &options.CountOptions{})

	if err != nil {
		logger.Error("erro ao consultar todas as Clientes", err)
		return nil, err
	}

	pagination := model.NewPaginate(limit, page, count)

	curr, err := collection.Find(ctx, query, pagination.GetPaginatedOpts())
	if err != nil {
		return nil, err
	}

	result := make([]*model.Cliente, 0)
	for curr.Next(ctx) {
		cat := &model.Cliente{}
		if err := curr.Decode(cat); err != nil {
			logger.Error("erro ao consulta todas as Clientes", err)
		}
		result = append(result, cat)
	}

	pagination.Paginate(result)

	return pagination, nil
}

func (cat *ClienteDataService) GetByDocumento(ctx context.Context, Doc string) bool {

	collection := cat.mdb.GetCollection("cfStore")

	// Utilizando o método CountDocuments para verificar a existência
	filter := bson.D{
		{Key: "cpf_cnpj", Value: Doc},
		{Key: "data_type", Value: "cliente"},
	}
	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		logger.Error("erro ao consultar Cliente pelo doc", err)
		return false
	}

	// Se count for maior que zero, o fornecedor existe
	return count > 0
}
