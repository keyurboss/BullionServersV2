package services

import (
	"fmt"
	"math/rand"

	"github.com/go-faker/faker/v4"
	"github.com/mitchellh/mapstructure"
	"github.com/rpsoftech/bullion-server/src/interfaces"
	"github.com/rpsoftech/bullion-server/src/mongodb/repos"
	"github.com/rpsoftech/bullion-server/src/validator"
)

type generalUserService struct {
	generalUserReqRepo  *repos.GeneralUserReqRepoStruct
	GeneralUserRepo     *repos.GeneralUserRepoStruct
	BullionSiteInfoRepo *repos.BullionSiteInfoRepoStruct
}

var GeneralUserService *generalUserService

func init() {
	GeneralUserService = &generalUserService{
		GeneralUserRepo:     repos.GeneralUserRepo,
		BullionSiteInfoRepo: repos.BullionSiteInfoRepo,
		generalUserReqRepo:  repos.GeneralUserReqRepo,
	}
}

func (service *generalUserService) RegisterNew(bullionId string, user interface{}) (*interfaces.GeneralUserEntity, error) {
	var baseGeneralUser interfaces.GeneralUser
	var entity interfaces.GeneralUserEntity

	Bullion, err := service.BullionSiteInfoRepo.FindOne(bullionId)
	if err != nil {
		return &entity, err
	}
	if Bullion.GeneralUserInfo.AutoLogin {
		baseGeneralUser = interfaces.GeneralUser{
			FirstName:     faker.FirstName(),
			LastName:      faker.LastName(),
			FirmName:      faker.Username(),
			ContactNumber: faker.Phonenumber(),
			GstNumber:     fmt.Sprintf("%dAAAAA%dA1ZA", rand.Intn(99-10)+10, rand.Intn(9999-1000)+1000),
			OS:            "AUTO",
			IsAuto:        true,
			DeviceId:      faker.UUIDDigit(),
			DeviceType:    interfaces.DEVICE_TYPE_IOS,
		}
	} else {
		baseGeneralUser = interfaces.GeneralUser{
			IsAuto: false,
		}
	}
	baseGeneralUser.RandomPass = faker.Password()
	err = mapstructure.Decode(user, &baseGeneralUser)
	if err != nil {
		return &entity, err
	}
	errs := validator.Validator.Validate(&baseGeneralUser)
	if len(errs) > 0 {
		reqErr := &interfaces.RequestError{
			StatusCode: 400,
			Code:       interfaces.ERROR_INVALID_INPUT,
			Message:    "",
			Name:       "INVALID_INPUT",
			Extra:      errs,
		}
		return &entity, reqErr.AppendValidationErrors(errs)
	}
	entity = interfaces.GeneralUserEntity{
		BaseEntity:  interfaces.BaseEntity{},
		GeneralUser: baseGeneralUser,
		UserRolesInterface: interfaces.UserRolesInterface{
			Role: interfaces.ROLE_GENERAL_USER,
		},
	}
	entity.CreateNewId()

	errs = validator.Validator.Validate(&entity)
	if len(errs) > 0 {
		reqErr := &interfaces.RequestError{
			StatusCode: 400,
			Code:       interfaces.ERROR_INVALID_INPUT,
			Message:    "",
			Name:       "INVALID_INPUT",
			Extra:      errs,
		}
		return &entity, reqErr.AppendValidationErrors(errs)
	}
	service.GeneralUserRepo.Save(&entity)
	_, err = service.sendApprovalRequest(&entity, Bullion)
	if err != nil {
		return &entity, err
	}
	return &entity, err
}

func (service *generalUserService) CreateApprovalRequest(userId string, password string, bullionId string) (reqEntity *interfaces.GeneralUserReqEntity, err error) {
	var userEntity *interfaces.GeneralUserEntity
	var bullionEntity *interfaces.BullionSiteInfoEntity
	if userEntity, err = service.GetGeneralUserDetailsByIdPassword(userId, password); err == nil {
		if bullionEntity, err = service.BullionSiteInfoRepo.FindOne(bullionId); err == nil {
			reqEntity, err = service.sendApprovalRequest(userEntity, bullionEntity)
		}
	}
	return
}
func (service *generalUserService) sendApprovalRequest(user *interfaces.GeneralUserEntity, bullion *interfaces.BullionSiteInfoEntity) (reqEntity *interfaces.GeneralUserReqEntity, err error) {
	reqEntity, err = service.generalUserReqRepo.FindOneByGeneralUserIdAndBullionId(user.ID, bullion.ID)
	if err == nil {
		err = &interfaces.RequestError{
			StatusCode: 400,
			Code:       interfaces.ERROR_GENERAL_USER_REQ_EXISTS,
			Message:    "REQUEST ALREADY EXISTS",
			Name:       "ERROR_GENERAL_USER_REQ_EXISTS",
		}
		return
	} else {
		err = nil
	}
	reqEntity = &interfaces.GeneralUserReqEntity{
		GeneralUserId: user.ID,
		BullionId:     bullion.ID,
		Status:        interfaces.GENERAL_USER_AUTH_STATUS_REQUESTED,
	}
	if bullion.GeneralUserInfo.AutoApprove {
		reqEntity.Status = interfaces.GENERAL_USER_AUTH_STATUS_AUTHORIZED
	}
	reqEntity.CreateNewId()
	reqEntity, err = service.generalUserReqRepo.Save(reqEntity)
	return
}
func (service *generalUserService) GetGeneralUserDetailsByIdPassword(id string, password string) (*interfaces.GeneralUserEntity, error) {
	entity, err := service.GeneralUserRepo.FindOne(id)
	if err != nil {
		return entity, err
	}
	if entity.RandomPass != password {
		err = &interfaces.RequestError{
			StatusCode: 400,
			Code:       interfaces.ERROR_GENERAL_USER_INVALID_PASSWORD,
			Message:    fmt.Sprintf("GeneralUser Entity invalid password %s ", password),
			Name:       "ERROR_GENERAL_USER_INVALID_PASSWORD",
		}
		return entity, err
	}
	return entity, err
}
