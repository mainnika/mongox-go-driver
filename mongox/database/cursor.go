package database

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/mainnika/mongox-go-driver/v2/mongox"
	"github.com/mainnika/mongox-go-driver/v2/mongox/query"
)

func (d *Database) createCursor(target interface{}, composed *query.Query) (cursor *mongox.Cursor, err error) {
	_, hasPreloader := composed.Preloader()
	if hasPreloader {
		return d.createAggregateCursor(target, composed)
	}

	return d.createSimpleCursor(target, composed)
}

func (d *Database) createSimpleCursor(target interface{}, composed *query.Query) (cursor *mongox.Cursor, err error) {
	collection, err := d.GetCollectionOf(target)
	if err != nil {
		return nil, err
	}

	opts := options.Find()
	opts.Sort = composed.Sorter()
	opts.Limit = composed.Limiter()
	opts.Skip = composed.Skipper()

	ctx := d.Context()
	m := composed.M()

	return collection.Find(ctx, m, opts)
}

func (d *Database) createAggregateCursor(target interface{}, composed *query.Query) (cursor *mongox.Cursor, err error) {
	collection, err := d.GetCollectionOf(target)
	if err != nil {
		return nil, err
	}

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
			return nil, fmt.Errorf("%w: private field is not preloadable", mongox.ErrMalformedBase)
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
			return nil, fmt.Errorf("%w: foreign field is not specified", mongox.ErrMalformedBase)
		}

		foreignField := strings.TrimSpace(preloadData[1])
		if len(foreignField) == 0 {
			return nil, fmt.Errorf("%w: foreign field is empty", mongox.ErrMalformedBase)
		}
		localField := strings.TrimSpace(preloadData[0])
		if len(localField) == 0 {
			localField = "_id"
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
			} else {
				return nil, fmt.Errorf("%w: preload limit should be an integer", mongox.ErrMalformedBase)
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
				return nil, fmt.Errorf("%w: preload field should have ptr type", mongox.ErrMalformedBase)
			}

			lookupCollection, err := d.GetCollectionOf(reflect.Zero(fieldType).Interface())
			if err != nil {
				return nil, err
			}

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

	ctx := d.Context()
	opts := options.Aggregate()

	return collection.Aggregate(ctx, pipeline, opts)
}
