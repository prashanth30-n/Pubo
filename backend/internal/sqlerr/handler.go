package sqlerr

import (
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/prashanth30-n/Pubo/internal/errs"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// ErrCode reports the error code for a given error.
// If the error is nil or is not of type *Error it reports sqlerr.Other.
func ErrCode(err error) Code {
	var pgerr *Error
	if errors.As(err, &pgerr) {
		return pgerr.Code
	}
	return Other
}
func ConvertPgError(src *pgconn.PgError) *Error {
	return &Error{
		Code:           MapCode(src.Code),
		Severity:       MapSeverity(src.Severity),
		DatabaseCode:   src.Code,
		Message:        src.Message,
		SchemaName:     src.SchemaName,
		TableName:      src.TableName,
		ColoumnName:    src.ColumnName,
		DataTypeName:   src.DataTypeName,
		ConstraintName: src.ConstraintName,
		driverErr:      src,
	}
}

func generateErrorCode(tableName string, errType Code) string {
	if tableName == "" {
		tableName = "RECORD"
	}
	domain := strings.ToUpper(tableName)
	// Singularize the table name
	if strings.HasSuffix(domain, "S") && len(domain) > 1 {
		domain = domain[:len(domain)-1]
	}
	action := "ERROR"
	switch errType {
	case ForeignKeyViolation:
		action = "NOT_FOUND"
	case UniqueViolation:
		action = "ALREADY_EXISTS"
	case NotNullViolation:
		action = "REQUIRED"
	case CheckViolation:
		action = "INVALID"

	}
	return fmt.Sprintf("%s_%s", domain, action)
}

// formatUserFriendlyMessage generates a user-friendly error message
func formatUserFriendlyMessage(sqlErr *Error) string {
	entityName := getEntityName(sqlErr.TableName, sqlErr.ColoumnName)

	switch sqlErr.Code {
	case ForeignKeyViolation:
		return fmt.Sprintf("The referenced %s does not exist", entityName)
	case UniqueViolation:
		return fmt.Sprintf("A %s with this identifier already exists", entityName)
	case NotNullViolation:
		fieldName := humanizeText(sqlErr.ColoumnName)
		if fieldName == "" {
			fieldName = "field"
		}
		return fmt.Sprintf("The %s is required", fieldName)
	case CheckViolation:
		fieldName := humanizeText(sqlErr.ColoumnName)
		if fieldName != "" {
			return fmt.Sprintf("The %s value does not meet required conditions", fieldName)
		}
		return "One or more values do not meet required conditions"
	default:
		return "An error occurred while processing your request"
	}
}

func getEntityName(tableName, coloumnName string) string {
	if coloumnName != "" && strings.HasSuffix(strings.ToLower(coloumnName), "_id") {
		entity := strings.TrimSuffix(strings.ToLower(coloumnName), "_id")
		return humanizeText(entity)
	}
	if tableName != "" {
		entity := tableName
		if strings.HasSuffix(entity, "s") && len(entity) > 1 {
			entity = entity[:len(entity)-1]
		}
		return humanizeText(entity)
	}
	return "record"
}

func humanizeText(text string) string {
	if text == "" {
		return ""
	}
	return cases.Title(language.English).String(strings.ReplaceAll(text, "_", " "))
}

// extractColumnForUniqueViolation gets field name from unique constraint
func extractColumnForUniqueViolation(constraintName string) string {
	if constraintName == "" {
		return ""
	}

	// Try standard naming convention first (unique_table_column)
	if strings.HasPrefix(constraintName, "unique_") {
		parts := strings.Split(constraintName, "_")
		if len(parts) >= 3 {
			return parts[len(parts)-1]
		}
	}

	// Try alternate convention (table_column_key)
	re := regexp.MustCompile(`_([^_]+)_(?:key|ukey)$`)
	matches := re.FindStringSubmatch(constraintName)
	if len(matches) > 1 {
		return matches[1]
	}

	return ""
}

// HandleError processes a database error into an appropriate application error

func HandleError(err error) error {
	var httpErr *errs.HTTPError
	if errors.As(err, &httpErr) {
		return err
	}
	var pgerr *pgconn.PgError
	if errors.As(err, &pgerr) {
		sqlErr := ConvertPgError(pgerr)

		errorCode := generateErrorCode(sqlErr.TableName, sqlErr.Code)
		userMessage := formatUserFriendlyMessage(sqlErr)

		switch sqlErr.Code {
		case ForeignKeyViolation:
			return errs.NewBadRequestError(userMessage, false, &errorCode, nil, nil)
		case UniqueViolation:
			columnName := extractColumnForUniqueViolation(sqlErr.ConstraintName)
			if columnName != "" {
				userMessage = strings.ReplaceAll(userMessage, "identifier", humanizeText(columnName))
			}
			return errs.NewBadRequestError(userMessage, true, &errorCode, nil, nil)
		case NotNullViolation:
			fieldErrors := []errs.FieldError{
				{
					Field: strings.ToLower(sqlErr.ColoumnName),
					Error: "is required",
				},
			}
			return errs.NewBadRequestError(userMessage, true, &errorCode, fieldErrors, nil)

		case CheckViolation:
			return errs.NewBadRequestError(userMessage, true, &errorCode, nil, nil)

		default:
			return errs.NewInternalServerError()
		}
	}
	//hanlde common pgx errors
	switch {
	case errors.Is(err, pgx.ErrNoRows), errors.Is(err, sql.ErrNoRows):
		errMsg := err.Error()
		tablePrefix := "table"
		if strings.Contains(errMsg, tablePrefix) {
			table := strings.Split(strings.Split(errMsg, tablePrefix)[1], ":")[0]
			entityName := getEntityName(table, "")
			return errs.NewNotFoundError(fmt.Sprintf("%s not found", entityName), true, nil)
		}
		return errs.NewNotFoundError("resource not found", false, nil)

	}
	return errs.NewInternalServerError()
}
