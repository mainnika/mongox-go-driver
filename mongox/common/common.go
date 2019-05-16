package common

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/mainnika/mongox-go-driver/mongox"
	"github.com/mainnika/mongox-go-driver/mongox/errors"
	"github.com/mainnika/mongox-go-driver/mongox/query"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func createSimpleLoad(db *mongox.Database, target interface{}, composed *query.Query) (cursor *mongo.Cursor, err error) {

	collection := db.GetCollectionOf(target)
	opts := options.Find()

	opts.Sort = composed.Sorter()
	opts.Limit = composed.Limiter()
	opts.Skip = composed.Skipper()

	return collection.Find(db.Context(), composed.M(), opts)
}

func createAggregateLoad(db *mongox.Database, target interface{}, composed *query.Query) (cursor *mongo.Cursor, err error) {

	collection := db.GetCollectionOf(target)
	opts := options.Aggregate()

	pipelineHead := primitive.A{primitive.M{"$match": composed.M()}}
	pipelineTail := primitive.A{}

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
			return nil, errors.Malformedf("preload private field is impossible")
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

		preloadName := strings.TrimSpace(preloadData[0])
		if len(preloadName) == 0 {
			preloadName = jsonName
		}

		foreignField := strings.TrimSpace(preloadData[1])
		if len(foreignField) == 0 {
			panic("there is no foreign field")
		}

		preloadLimiter := 100
		if len(preloadData) > 2 {

			stringLimit := strings.TrimSpace(preloadData[2])
			intLimit := preloadLimiter

			intLimit, err = strconv.Atoi(stringLimit)
			if err == nil {
				preloadLimiter = intLimit
			}
		}

		for _, preload := range preloads {
			if preload != preloadName {
				continue
			}

			isPtr := el.Field(i).Kind() == reflect.Ptr
			isSlice := el.Field(i).Kind() == reflect.Slice
			isIface := el.Field(i).CanInterface()
			if (!isPtr && !isSlice) || !isIface {
				continue
			}

			typ := el.Field(i).Type()
			lookupCollection := db.GetCollectionOf(reflect.Zero(typ).Interface())
			lookupVars := primitive.M{"selector": "$_id"}
			lookupPipeline := primitive.A{
				// todo: make match from composed query
				primitive.M{"$match": primitive.M{"$expr": primitive.M{"$eq": primitive.A{"$" + foreignField, "$$selector"}}}},
			}

			if isSlice && preloadLimiter > 0 {
				lookupPipeline = append(lookupPipeline, primitive.M{"$limit": preloadLimiter})
			} else if !isSlice {
				lookupPipeline = append(lookupPipeline, primitive.M{"$limit": 1})
			}

			pipelineTail = append(pipelineTail, primitive.M{
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

			pipelineTail = append(pipelineTail, primitive.M{
				"$unwind": primitive.M{
					"preserveNullAndEmptyArrays": true,
					"path":                       "$" + jsonName,
				},
			})
		}
	}

	return collection.Aggregate(db.Context(), append(pipelineHead, pipelineTail...), opts)
}
