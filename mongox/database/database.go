package database

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/mainnika/mongox-go-driver/v2/mongox"
	"github.com/mainnika/mongox-go-driver/v2/mongox/query"
)

// Database handler
type Database struct {
	client *mongox.Client
	dbname string
	ctx    context.Context
}

// NewDatabase function creates new database instance with mongo client and empty context
func NewDatabase(client *mongox.Client, dbname string) (db mongox.Database) {

	db = &Database{
		client: client,
		dbname: dbname,
	}

	return
}

// Client function returns a mongo client
func (d *Database) Client() (client *mongox.Client) {
	return d.client
}

// Context function returns a context
func (d *Database) Context() (ctx context.Context) {

	ctx = d.ctx
	if ctx == nil {
		ctx = context.Background()
	}

	return
}

// Name function returns a database name
func (d *Database) Name() (name string) {
	return d.dbname
}

// New function creates new database context with same client
func (d *Database) New(ctx context.Context) (db mongox.Database) {

	if ctx == nil {
		ctx = context.Background()
	}

	db = &Database{
		client: d.client,
		dbname: d.dbname,
		ctx:    ctx,
	}

	return
}

// GetCollectionOf returns the collection object by the «collection» tag of the given document;
// the «collection» tag should exists, e.g.:
// type Foobar struct {
//     base.ObjectID `bson:",inline" json:",inline" collection:"foobars"`
// 	   ...
// Will panic if there is no «collection» tag
func (d *Database) GetCollectionOf(document interface{}) (collection *mongox.Collection) {

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

func (d *Database) createSimpleLoad(target interface{}, composed *query.Query) (cursor *mongox.Cursor, err error) {

	collection := d.GetCollectionOf(target)
	opts := options.Find()

	opts.Sort = composed.Sorter()
	opts.Limit = composed.Limiter()
	opts.Skip = composed.Skipper()

	return collection.Find(d.Context(), composed.M(), opts)
}

func (d *Database) createAggregateLoad(target interface{}, composed *query.Query) (cursor *mongox.Cursor, err error) {

	collection := d.GetCollectionOf(target)
	opts := options.Aggregate()

	pipeline := primitive.A{}

	if !composed.Empty() {
		pipeline = append(pipeline, primitive.M{"$match": composed.M()})
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

	el := reflect.ValueOf(target)
	elType := el.Type()
	if elType.Kind() == reflect.Ptr {
		elType = elType.Elem()
	}

	numField := elType.NumField()
	preloads, _ := composed.Preloader()

	for i := 0; i < numField; i++ {

		field := elType.Field(i)
		tag := field.Tag

		preloadTag, ok := tag.Lookup("preload")
		if !ok {
			continue
		}
		jsonTag, _ := tag.Lookup("json")
		if jsonTag == "-" {
			panic(fmt.Errorf("preload private field is impossible"))
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
			panic(fmt.Errorf("there is no foreign field"))
		}

		localField := strings.TrimSpace(preloadData[0])
		if len(localField) == 0 {
			localField = "_id"
		}

		foreignField := strings.TrimSpace(preloadData[1])
		if len(foreignField) == 0 {
			panic(fmt.Errorf("there is no foreign field"))
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

			field := elType.Field(i)
			fieldType := field.Type

			isSlice := fieldType.Kind() == reflect.Slice
			if isSlice {
				fieldType = fieldType.Elem()
			}

			isPtr := fieldType.Kind() != reflect.Ptr
			if isPtr {
				panic(fmt.Errorf("preload field should have ptr type"))
			}

			lookupCollection := d.GetCollectionOf(reflect.Zero(fieldType).Interface())
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
