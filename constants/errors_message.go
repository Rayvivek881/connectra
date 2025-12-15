package constants

import "errors"

var (
	InvalidCredentialsError     = errors.New("invalid credentials")
	UserAccountDeletedError     = errors.New("user account is deleted")
	UserAccountNotActiveError   = errors.New("user account is not active")
	UserEmailAlreadyExistsError = errors.New("user with this email already exists")
	FailedToHashPasswordError   = errors.New("failed to hash password")
	FailedToCreateUserError     = errors.New("failed to create user")

	FailedToGenerateTokenError = errors.New("failed to generate token")
	InvalidSigningMethodError  = errors.New("invalid signing method")
	InvalidTokenError          = errors.New("invalid token")
	InvalidOrExpiredTokenError = errors.New("invalid or expired token")

	AuthorizationHeaderRequiredError = errors.New("authorization header is required")
	InvalidAuthorizationFormatError  = errors.New("invalid authorization header format")

	InvalidRequestBodyError = errors.New("invalid request body")

	PageSizeExceededError   = errors.New("page size exceeds maximum limit")
	PageNumberExceededError = errors.New("page number exceeds maximum limit")
	FailedToFetchDataError  = errors.New("failed to fetch data")
)
