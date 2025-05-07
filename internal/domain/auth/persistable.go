package auth

type Persistable interface {
    // the name of the DB table
    TableName() string

    // a slice of column names (in order)
    Columns() []string

    // the *pointer* to each column value in the same order
    // (so QueryRowContext can scan directly into a struct)
    Values() []interface{}

    // set the auto-generated ID back on a struct
    SetID(int64)

	GetID() int64  
}

