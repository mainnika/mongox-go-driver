package mongox

import (
	"errors"

	"go.mongodb.org/mongo-driver/mongo"
)

// Reexported mongo errors
var (
	ErrMissingResumeToken  = mongo.ErrMissingResumeToken
	ErrNilCursor           = mongo.ErrNilCursor
	ErrUnacknowledgedWrite = mongo.ErrUnacknowledgedWrite
	ErrClientDisconnected  = mongo.ErrClientDisconnected
	ErrNilDocument         = mongo.ErrNilDocument
	ErrEmptySlice          = mongo.ErrEmptySlice
	ErrInvalidIndexValue   = mongo.ErrInvalidIndexValue
	ErrNonStringIndexName  = mongo.ErrNonStringIndexName
	ErrMultipleIndexDrop   = mongo.ErrMultipleIndexDrop
	ErrWrongClient         = mongo.ErrWrongClient
	ErrNoDocuments         = mongo.ErrNoDocuments
)

var (
	ErrMalformedBase     = errors.New("source contains malformed document base")
	ErrUninitializedBase = errors.New("uninitialized document")
	ErrNoCollection      = errors.New("no collection found")
)
