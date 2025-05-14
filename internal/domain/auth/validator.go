package auth


type FieldsValidator interface {
    Validate(data interface{}) error
}

