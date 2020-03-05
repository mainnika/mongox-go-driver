package database

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/mainnika/mongox-go-driver/v2/mongox"
	"github.com/mainnika/mongox-go-driver/v2/mongox/query"
)

// Database handler
type Database struct {
	client *mongo.Client
	dbname string
	ctx    context.Context
}

// NewDatabase function creates new database instance with mongo client and empty context
func NewDatabase(client *mongo.Client, dbname string) mongox.Database {

	db := &Database{}
	db.client = client
	db.dbname = dbname

	return db
}

// Client function returns a mongo client
func (d *Database) Client() mongox.MongoClient {
	return d.client
}

// Context function returns a context
func (d *Database) Context() context.Context {
	return d.ctx
}

// Name function returns a database name
func (d *Database) Name() string {
	return d.dbname
}

// New function creates new database context with same client
func (d *Database) New(ctx context.Context) mongox.Database {

	if ctx == nil {
		ctx = context.Background()
	}

	return &Database{
		client: d.client,
		dbname: d.dbname,
		ctx:    ctx,
	}
}

// GetCollectionOf returns the collection object by the «collection» tag of the given document;
// the «collection» tag should exists, e.g.:
// type Foobar struct {
//     base.ObjectID `bson:",inline" json:",inline" collection:"foobars"`
// 	   ...
// Will panic if there is no «collection» tag
func (d *Database) GetCollectionOf(document interface{}) mongox.MongoCollection {

	el := reflect.TypeOf(document).Elem()
	numField := el.NumField()

	for i := 0; i < numField; i++ {
		field := el.Field(i)
		tag := field.Tag
		found, ok := tag.Lookup("collection")
		if !ok {
			continue
		}

		return d.client.Database(d.dbname).Collection(found)
	}

	panic(fmt.Errorf("document %v does not have a collection tag", document))
}

func (d *Database) createSimpleLoad(target interface{}, composed *query.Query) (cursor *mongo.Cursor, err error) {

	collection := d.GetCollectionOf(target)
	opts := options.Find()

	opts.Sort = composed.Sorter()
	opts.Limit = composed.Limiter()
	opts.Skip = composed.Skipper()

	return collection.Find(d.Context(), composed.M(), opts)
}

func (d *Database) createAggregateLoad(target interface{}, composed *query.Query) (cursor *mongo.Cursor, err error) {

	collection := d.GetCollectionOf(target)
	opts := options.Aggregate()

	pipeline := primitive.A{}

	if !composed.Empty() {
		pipeline = append(pipeline, primitive.M{"$match": primitive.M{"$expr": composed.M()}})
	}
	if composed.Sorter() != nil {
		pipeline = append(pipeline, primitive.M{"$sort": composed.Sorter()})
	}
	if composed.Skipper() != nil {
		pipeline = append(pipeline, primitive.M{"$skip": *composed.Skipper()})
	}
	if composed.Limiter() != nil {
		pipeline = append(pipeline, primitive.M{"$limit": *composed.Limiter()})
	}

	el := reflect.ValueOf(target).Elem()
	elType := el.Type()
	numField := elType.NumField()
	_, preloads := composed.Preloader()

	for i := 0; i < numField; i++ {

		field := elType.Field(i)
		tag := field.Tag

		preloadTag, ok := tag.Lookup("preload")
		if !ok {
			continue
		}
		jsonTag, ok := tag.Lookup("json")
		if jsonTag == "-" {
			return nil, fmt.Errorf("preload private field is impossible")
		}

		jsonData := strings.SplitN(jsonTag, ",", 2)
		jsonName := field.Name
		if len(jsonData) > 0 {
			jsonName = strings.TrimSpace(jsonData[0])
		}

		preloadData := strings.Split(preloadTag, ",")
		if len(preloadData) == 0 {
			continue
		}
		if len(preloadData) == 1 {
			panic("there is no foreign field")
		}

		localField := strings.TrimSpace(preloadData[0])
		if len(localField) == 0 {
			localField = "_id"
		}

		foreignField := strings.TrimSpace(preloadData[1])
		if len(foreignField) == 0 {
			panic("there is no foreign field")
		}

		preloadLimiter := 100
		preloadReversed := false
		if len(preloadData) > 2 {
			stringLimit := strings.TrimSpace(preloadData[2])
			intLimit := preloadLimiter

			preloadReversed = strings.HasPrefix(stringLimit, "-")
			if preloadReversed {
				stringLimit = stringLimit[1:]
			}

			intLimit, err = strconv.Atoi(stringLimit)
			if err == nil {
				preloadLimiter = intLimit
			}
		}

		for _, preload := range preloads {
			if preload != jsonName {
				continue
			}

			isSlice := el.Field(i).Kind() == reflect.Slice

			typ := el.Field(i).Type()
			if typ.Kind() == reflect.Slice {
				typ = typ.Elem()
			}
			if typ.Kind() != reflect.Ptr {
				panic("preload field should have ptr type")
			}

			lookupCollection := d.GetCollectionOf(reflect.Zero(typ).Interface())
			lookupVars := primitive.M{"selector": "$" + localField}
			lookupPipeline := primitive.A{
				primitive.M{"$match": primitive.M{"$expr": primitive.M{"$eq": primitive.A{"$" + foreignField, "$$selector"}}}},
			}

			if preloadReversed {
				lookupPipeline = append(lookupPipeline, primitive.M{"$sort": primitive.M{"_id": -1}})
			}
			if isSlice && preloadLimiter > 0 {
				lookupPipeline = append(lookupPipeline, primitive.M{"$limit": preloadLimiter})
			} else if !isSlice {
				lookupPipeline = append(lookupPipeline, primitive.M{"$limit": 1})
			}

			pipeline = append(pipeline, primitive.M{
				"$lookup": primitive.M{
					"from":     lookupCollection.Name(),
					"let":      lookupVars,
					"pipeline": lookupPipeline,
					"as":       jsonName,
				},
			})

			if isSlice {
				continue
			}

			pipeline = append(pipeline, primitive.M{
				"$unwind": primitive.M{
					"preserveNullAndEmptyArrays": true,
					"path":                       "$" + jsonName,
				},
			})
		}
	}

	return collection.Aggregate(d.Context(), pipeline, opts)
}
