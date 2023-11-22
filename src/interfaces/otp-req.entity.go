package interfaces

import "time"

type (
	OTPReqBase struct {
		BullionId string    `bson:"bullionId" json:"bullionId" validate:"required,uuid"`
		Name      string    `bson:"name" json:"name" validate:"required"`
		Number    string    `bson:"number" json:"number" validate:"required,min=10,max=12"`
		Attempt   int16     `bson:"attempt" json:"attempt" validate:"required"`
		ExpiresAt time.Time `bson:"expiresAt" json:"expiresAt" validate:"required"`
	}
	OTPReqEntity struct {
		*BaseEntity `bson:"inline"`
		*OTPReqBase `bson:"inline"`
		OTP         string `bson:"otp" json:"otp" validate:"required"`
	}
)

func (otp *OTPReqEntity) NewAttempt() {
	otp.Attempt = otp.Attempt + 1
	otp.Updated()
}

func CreateOTPEntity(otpBase *OTPReqBase, OTP string) *OTPReqEntity {
	entity := &OTPReqEntity{
		BaseEntity: &BaseEntity{},
		OTPReqBase: otpBase,
		OTP:        OTP,
	}
	entity.createNewId()
	return entity
}
