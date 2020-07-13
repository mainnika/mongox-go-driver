package database

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"text/template"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// IndexEnsure function ensures index in mongo collection of document
//   `index:""` -- https://docs.mongodb.com/manual/indexes/#create-an-index
//   `index:"-"` -- (descending)
//   `index:"-,+foo,+-bar"` -- https://docs.mongodb.com/manual/core/index-compound
//   `index:"-,+foo,+-bar,unique"` -- https://docs.mongodb.com/manual/core/index-unique
//   `index:"-,+foo,+-bar,unique,allowNull"` -- https://docs.mongodb.com/manual/core/index-partial
//   `index:"-,unique,allowNull,expireAfter=86400"` -- https://docs.mongodb.com/manual/core/index-ttl
//   `index:"-,unique,allowNull,expireAfter={{.Expire}}"` -- evaluate index as a golang template with `cfg` arguments
func (d *Database) IndexEnsure(cfg interface{}, document interface{}) (err error) {

	el := reflect.ValueOf(document).Elem().Type()
	numField := el.NumField()
	documents := d.GetCollectionOf(document)

	for i := 0; i < numField; i++ {

		field := el.Field(i)
		tag := field.Tag

		indexTag, ok := tag.Lookup("index")
		if !ok {
			continue
		}
		bsonTag, ok := tag.Lookup("bson")
		if !ok {
			return fmt.Errorf("bson tag is not defined for field:%v document:%v", field, document)
		}

		var tmpBuffer = &bytes.Buffer{}
		var tpl *template.Template

		tpl, err = template.New("").Parse(indexTag)
		if err != nil {
			panic(fmt.Errorf("invalid prop template %v, err:%w", indexTag, err))
		}
		err = tpl.Execute(tmpBuffer, cfg)
		if err != nil {
			panic(fmt.Errorf("failed to evaluate prop template %v, err:%w", indexTag, err))
		}

		indexString := tmpBuffer.String()
		indexValues := strings.Split(indexString, ",")
		bsonValues := strings.Split(bsonTag, ",")

		var f = false
		var t = true
		var key = bsonValues[0]
		var name = fmt.Sprintf("%s_%s_", indexString, key)

		if len(key) == 0 {
			panic(fmt.Errorf("cannot evaluate index key"))
		}

		opts := &options.IndexOptions{
			Background: &f,
			Unique:     &f,
			Name:       &name,
		}

		index := primitive.D{{Key: key, Value: 1}}
		if indexValues[0] == "-" {
			index = primitive.D{{Key: key, Value: -1}}
		}

		for _, prop := range indexValues[1:] {
			var left string
			var right string

			pair := strings.SplitN(prop, "=", 2)
			left = pair[0]
			if len(pair) > 1 {
				right = pair[1]
			}

			switch {
			case left == "unique":
				opts.Unique = &t

			case left == "allowNull":
				expression, isMap := opts.PartialFilterExpression.(primitive.M)
				if !isMap || expression == nil {
					expression = primitive.M{}
				}

				expression[key] = primitive.M{"$exists": true}
				opts.PartialFilterExpression = expression

			case left == "expireAfter":
				expireAfter, err := strconv.Atoi(right)
				if err != nil || expireAfter < 1 {
					panic(fmt.Errorf("invalid expireAfter value, err: %w", err))
				}

				expireAfterInt32 := int32(expireAfter)
				opts.ExpireAfterSeconds = &expireAfterInt32

			case len(left) > 0 && left[0] == '+':
				compoundValue := left[1:]
				if len(compoundValue) == 0 {
					panic(fmt.Errorf("invalid compound value"))
				}

				if compoundValue[0] == '-' {
					index = append(index, primitive.E{compoundValue[1:], -1})
				} else {
					index = append(index, primitive.E{compoundValue, 1})
				}

			default:
				panic(fmt.Errorf("unsupported flag:%q in tag:%q of type:%s", prop, tag, el))
			}
		}

		_, err = documents.Indexes().CreateOne(d.Context(), mongo.IndexModel{Keys: index, Options: opts})
		if err != nil {
			return
		}
	}

	return
}
