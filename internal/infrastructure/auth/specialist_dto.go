package authinfrastructure

import (
	"pethelp-backend/internal/domain/auth/account"
)

type SpecialistDTO struct {
	Name                 string `json:"name" validate:"required,min=2"`
    FamilyName           string `json:"family_name" validate:"required,min=2"`
	Phone 				 string `json:"phone" validate:"required,phone"`
    Email                string `json:"email" validate:"required,email"`
    Password             string `json:"password" validate:"required,min=12"`
    PasswordConfirmation string `json:"password_confirmation" validate:"required,eqfield=Password"`
}

func (d *SpecialistDTO) GetEmail() string {
	return d.Email
}

func (d *SpecialistDTO) GetPhone() string {
	return d.Phone
}

func (d *SpecialistDTO) GetPassword() string {
	return d.Password
}

func (d *SpecialistDTO) GetPasswordConfirmation() string {
	return d.PasswordConfirmation
}

func (d *SpecialistDTO) GetName() string {
	return d.Name
}

func (d *SpecialistDTO) GetFamilyName() string {
	return d.FamilyName
}

func NewSpecialistDTO() *SpecialistDTO {
	return &SpecialistDTO{}
}

func ToSpecialist(dto *SpecialistDTO) *account.Specialist {
	return &account.Specialist{
		Name:       dto.Name,
		FamilyName: dto.FamilyName,
		Phone:      dto.Phone,
		Email:      dto.Email,
	}
} 