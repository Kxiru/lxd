package auth

import (
	"context"
	"net/http"

	"github.com/canonical/lxd/shared/api"
	"github.com/canonical/lxd/shared/entity"
)

const (
	// AuthenticationMethodCluster is set in the request context as request.CtxProtocol when the request is authenticated
	// via mTLS and the peer certificate is present in the trust store as type certificate.TypeServer.
	AuthenticationMethodCluster string = "cluster"

	// AuthenticationMethodUnix is set in the request context as request.CtxProtocol when the request is made over the
	// unix socket.
	AuthenticationMethodUnix string = "unix"

	// AuthenticationMethodPKI is set in the request context as request.CtxProtocol when a `server.ca` file exists in
	// LXD_DIR, the peer certificate of the request was signed by the CA file, and core.trust_ca_certificates is true.
	//
	// Note: If core.trust_ca_certificates is false, the peer certificate is additionally verified via mTLS and the
	// value of request.CtxProtocol is set to api.AuthenticationMethodTLS.
	//
	// Note: Regardless of whether `core.trust_ca_certificates` is enabled, we still check if the client certificate
	// fingerprint is in the identity cache. If they are found, standard TLS restrictions will apply.
	AuthenticationMethodPKI string = "pki"

	// AuthenticationMethodDevLXD is the authentication method for interacting with the devlxd API.
	AuthenticationMethodDevLXD = "devlxd"
)

// PermissionChecker is a type alias for a function that returns whether a user has required permissions on an object.
// It is returned by Authorizer.GetPermissionChecker.
type PermissionChecker func(entityURL *api.URL) bool

// EntitlementReporter is an interface for adding entitlements to an entity.
type EntitlementReporter interface {
	// ReportEntitlements adds entitlements to the entity.
	// Note: this needs to be a list of string because the implementations of this method will be for the API types.
	ReportEntitlements([]string)
}

// Authorizer is the primary external API for this package.
type Authorizer interface {
	// Driver returns the driver name.
	Driver() string

	// CheckPermission checks if the caller has the given entitlement on the entity found at the given URL.
	//
	// Note: When a project does not have a feature enabled, the given URL should contain the request project, and the
	// effective project for the entity should be set on the request.Info in the given context.
	CheckPermission(ctx context.Context, entityURL *api.URL, entitlement Entitlement) error

	// GetPermissionChecker returns a PermissionChecker for a particular entity.Type.
	//
	// Note: As with CheckPermission, arguments to the returned PermissionChecker should contain the request project for
	// the entity. The effective project for the entity must be set on the request.Info in the given context before
	// calling the PermissionChecker.
	GetPermissionChecker(ctx context.Context, entitlement Entitlement, entityType entity.Type) (PermissionChecker, error)

	// CheckPermissionWithoutEffectiveProject checks a permission, but does not replace the project in the entity URL
	// with the effective project stored in the context.
	//
	// Warn: You almost never need this function. You should use CheckPermission instead.
	CheckPermissionWithoutEffectiveProject(ctx context.Context, entityURL *api.URL, entitlement Entitlement) error

	// GetPermissionCheckerWithoutEffectiveProject returns a PermissionChecker does not replace the project in the entity URL
	// with the effective project stored in the context.
	//
	// Warn: You almost never need this function. You should use GetPermissionChecker instead.
	GetPermissionCheckerWithoutEffectiveProject(ctx context.Context, entitlement Entitlement, entityType entity.Type) (PermissionChecker, error)
}

// IsDeniedError returns true if the error is not found or forbidden. This is because the CheckPermission method on
// Authorizer will return a not found error if the requestor does not have access to view the resource. If a requestor
// has view access, but not edit access a forbidden error is returned.
func IsDeniedError(err error) bool {
	return api.StatusErrorCheck(err, http.StatusNotFound, http.StatusForbidden)
}
