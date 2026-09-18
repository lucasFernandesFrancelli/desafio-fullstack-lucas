package auth_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"ekaizen-backend/internal/auth"
	"ekaizen-backend/internal/models"
)

func TestJWTIssuer_IssueAndParse(t *testing.T) {
	issuer := auth.NewJWTIssuer("segredo-de-teste")
	user := models.User{ID: uuid.New(), Name: "Ana", Role: models.RoleAnalista}

	token, err := issuer.Issue(user)
	require.NoError(t, err)

	claims, err := issuer.Parse(token)
	require.NoError(t, err)
	assert.Equal(t, user.ID, claims.UserID)
	assert.Equal(t, user.Role, claims.Role)
	assert.Equal(t, user.Name, claims.Name)
}

func TestJWTIssuer_Parse_MalformedToken(t *testing.T) {
	issuer := auth.NewJWTIssuer("segredo-de-teste")

	_, err := issuer.Parse("isto-não-é-um-jwt")

	assert.ErrorIs(t, err, auth.ErrInvalidToken)
}

func TestJWTIssuer_Parse_WrongSecret(t *testing.T) {
	issued := auth.NewJWTIssuer("segredo-a")
	other := auth.NewJWTIssuer("segredo-b")

	token, err := issued.Issue(models.User{ID: uuid.New()})
	require.NoError(t, err)

	_, err = other.Parse(token)
	assert.ErrorIs(t, err, auth.ErrInvalidToken)
}
