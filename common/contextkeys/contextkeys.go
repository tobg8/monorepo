package contextkeys

// Using a specific type as context keys to avoid possible collisions (golint)
type ctxKey uint8

// These constants are used as arguments for context.Context.Value() and context.WithValue()
const (
	// UserID is an authenticated account user_id
	UserID ctxKey = iota
	// StoreID is an authenticated account store_id
	StoreID
	// UniqueID is a request unique ID
	UniqueID
	// AdminID is an authenticated admin user id
	AdminID
	// Privileges is privileges as set by the introspect middleware
	Privileges
	// RefusedScopes is a list of blacklisted scopes for the user
	RefusedScopes
	// AuthToken is an authorization token
	AuthToken
	// UserMetadata is all metadata attached to the http request
	UserMetadata
	// Sub is authorization subject
	Sub
	// SessionID is the session ID associated with the authorization token as set by the introspect middleware
	SessionID
	// ClientID is the client ID associated with the authorization token as set by the introspect middleware
	ClientID
	// InstallID is the install ID associated with the authorization token as set by the introspect middleware
	InstallID
	// MemberID is the member ID of the authenticated member as set by the introspect middleware
	MemberID
	// Brand is brand identifier, value is a string
	Brand
	// BrandOriginHost is brand hostname value for which the brand is deduced
	BrandOriginHost
	// BrandConfig is brand configuration properties
	BrandConfig
	// Language identifier
	Language
)
