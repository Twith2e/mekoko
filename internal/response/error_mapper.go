package response

import (
	appErr "mekoko/internal/errors"
	"net/http"
)

type ErrorMapping struct {
	Status int
	Error  APIError
}

func MapError(err error) ErrorMapping {
	switch err {
	case appErr.ErrInvalidRequestBody:
		return ErrorMapping{
			Status: http.StatusBadRequest,
			Error: APIError{
				Code:    "INVALID_REQUEST_BODY",
				Message: appErr.ErrInvalidRequestBody.Error(),
			},
		}

	case appErr.ErrInvalidPasswordLength:
		return ErrorMapping{
			Status: http.StatusUnprocessableEntity,
			Error: APIError{
				Code:    "INVALID_PASSWORD_LENGTH",
				Message: appErr.ErrInvalidPasswordLength.Error(),
			},
		}

	case appErr.ErrInvalidPassword:
		return ErrorMapping{
			Status: http.StatusUnprocessableEntity,
			Error: APIError{
				Code:    "INVALID_PASSWORD",
				Message: appErr.ErrInvalidPassword.Error(),
			},
		}

	case appErr.ErrRegisteringUser:
		return ErrorMapping{
			Status: http.StatusInternalServerError,
			Error: APIError{
				Code:    "FAILED_TO_REGISTER_USER",
				Message: appErr.ErrRegisteringUser.Error(),
			},
		}

	case appErr.ErrPasswordMismatch:
		return ErrorMapping{
			Status: http.StatusUnprocessableEntity,
			Error: APIError{
				Code:    "PASSWORDS_DO_NOT_MATCH",
				Message: appErr.ErrPasswordMismatch.Error(),
			},
		}

	case appErr.ErrFindingUser:
		return ErrorMapping{
			Status: http.StatusNotFound,
			Error: APIError{
				Code:    "USER_NOT_FOUND",
				Message: appErr.ErrFindingUser.Error(),
			},
		}

	case appErr.ErrInvalidCredentials:
		return ErrorMapping{
			Status: http.StatusUnauthorized,
			Error: APIError{
				Code:    "INVALID_CREDENTIALS",
				Message: appErr.ErrInvalidCredentials.Error(),
			},
		}

	case appErr.ErrUserExists:
		return ErrorMapping{
			Status: http.StatusConflict,
			Error: APIError{
				Code:    "USER_ALREADY_EXISTS",
				Message: appErr.ErrFindingUser.Error(),
			},
		}

	case appErr.ErrUnauthorized:
		return ErrorMapping{
			Status: http.StatusUnauthorized,
			Error: APIError{
				Code:    "UNAUTHORIZED",
				Message: appErr.ErrUnauthorized.Error(),
			},
		}

	case appErr.ErrInvalidSession:
		return ErrorMapping{
			Status: http.StatusUnauthorized,
			Error: APIError{
				Code:    "INVALID_SESSION",
				Message: appErr.ErrInvalidSession.Error(),
			},
		}

	case appErr.ErrRefreshingAccessToken:
		return ErrorMapping{
			Status: http.StatusInternalServerError,
			Error: APIError{
				Code:    "REFRESH_TOKEN_ERROR",
				Message: appErr.ErrRefreshingAccessToken.Error(),
			},
		}

	case appErr.ErrTooManyRequests:
		return ErrorMapping{
			Status: http.StatusTooManyRequests,
			Error: APIError{
				Code:    "TOO_MANY_REQUESTS",
				Message: appErr.ErrTooManyRequests.Error(),
			},
		}

	case appErr.ErrPasswordResetUserNotFound:
		return ErrorMapping{
			Status: http.StatusNotFound,
			Error: APIError{
				Code:    "EMAIL_DOES_NOT_EXIST",
				Message: appErr.ErrPasswordResetUserNotFound.Error(),
			},
		}

	default:
		return ErrorMapping{
			Status: http.StatusInternalServerError,
			Error: APIError{
				Code:    "INTERNAL_SERVER_ERROR",
				Message: "An error occured, please try again",
			},
		}
	}
}
