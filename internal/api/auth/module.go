package apiauth

import (
	"pethelp-backend/internal/domain/auth/account"
	authinfra "pethelp-backend/internal/infrastructure/auth"
)

// Module sets up the /api/v1/specialists/register endpoint
var Module = AuthModule[*account.Specialist, *authinfra.SpecialistDTO](
    "specialist-auth",               // fx module name
    "specialists",                   // DB table
    "/api/v1/specialists",          // HTTP route prefix
    func() *authinfra.SpecialistDTO { return &authinfra.SpecialistDTO{} },
    func(d *authinfra.SpecialistDTO) *account.Specialist {
        return &account.Specialist{
            Name:       d.Name,
            FamilyName: d.FamilyName,
            Phone:      d.Phone,
            Email:      d.Email,
            // PasswordHash is set inside AuthService.Register
            IsBanned:   false,
            IsDeleted:  false,
            IsActive:   true,
            IsVerified: false,
        }
    },
)