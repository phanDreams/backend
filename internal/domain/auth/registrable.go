package auth

type Registrable interface {
    Persistable
    Credentialed
}