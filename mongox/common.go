package mongox

type Saver interface {
	Save(db *Database) error
}

type Deleter interface {
	Delete(db *Database) error
}

type Loader interface {
	Load(db *Database, filters ...interface{}) error
}
