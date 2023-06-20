package auth

import "context"

// ICredentialStore Interface for credential stores which are used to store and lookup
// credentials to authenticate against external services.
type ICredentialStore interface {
	// Store stores credential parameters into the store.
	Store(ctx context.Context, key string, credential *CredentialParams) error

	// Lookup lookups credential parameters by its key.
	Lookup(ctx context.Context, key string) (*CredentialParams, error)
}
