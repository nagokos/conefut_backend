// Code generated by entc, DO NOT EDIT.

package ent

import (
	"time"

	"github.com/nagokos/connefut_backend/ent/competition"
	"github.com/nagokos/connefut_backend/ent/prefecture"
	"github.com/nagokos/connefut_backend/ent/recruitment"
	"github.com/nagokos/connefut_backend/ent/schema"
	"github.com/nagokos/connefut_backend/ent/user"
)

// The init function reads all schema descriptors with runtime code
// (default values, validators, hooks and policies) and stitches it
// to their package variables.
func init() {
	competitionMixin := schema.Competition{}.Mixin()
	competitionMixinFields0 := competitionMixin[0].Fields()
	_ = competitionMixinFields0
	competitionMixinFields1 := competitionMixin[1].Fields()
	_ = competitionMixinFields1
	competitionFields := schema.Competition{}.Fields()
	_ = competitionFields
	// competitionDescCreatedAt is the schema descriptor for created_at field.
	competitionDescCreatedAt := competitionMixinFields0[0].Descriptor()
	// competition.DefaultCreatedAt holds the default value on creation for the created_at field.
	competition.DefaultCreatedAt = competitionDescCreatedAt.Default.(func() time.Time)
	// competitionDescUpdatedAt is the schema descriptor for updated_at field.
	competitionDescUpdatedAt := competitionMixinFields0[1].Descriptor()
	// competition.DefaultUpdatedAt holds the default value on creation for the updated_at field.
	competition.DefaultUpdatedAt = competitionDescUpdatedAt.Default.(func() time.Time)
	// competition.UpdateDefaultUpdatedAt holds the default value on update for the updated_at field.
	competition.UpdateDefaultUpdatedAt = competitionDescUpdatedAt.UpdateDefault.(func() time.Time)
	// competitionDescID is the schema descriptor for id field.
	competitionDescID := competitionMixinFields1[0].Descriptor()
	// competition.DefaultID holds the default value on creation for the id field.
	competition.DefaultID = competitionDescID.Default.(func() string)
	// competition.IDValidator is a validator for the "id" field. It is called by the builders before save.
	competition.IDValidator = competitionDescID.Validators[0].(func(string) error)
	prefectureMixin := schema.Prefecture{}.Mixin()
	prefectureMixinFields0 := prefectureMixin[0].Fields()
	_ = prefectureMixinFields0
	prefectureMixinFields1 := prefectureMixin[1].Fields()
	_ = prefectureMixinFields1
	prefectureFields := schema.Prefecture{}.Fields()
	_ = prefectureFields
	// prefectureDescCreatedAt is the schema descriptor for created_at field.
	prefectureDescCreatedAt := prefectureMixinFields0[0].Descriptor()
	// prefecture.DefaultCreatedAt holds the default value on creation for the created_at field.
	prefecture.DefaultCreatedAt = prefectureDescCreatedAt.Default.(func() time.Time)
	// prefectureDescUpdatedAt is the schema descriptor for updated_at field.
	prefectureDescUpdatedAt := prefectureMixinFields0[1].Descriptor()
	// prefecture.DefaultUpdatedAt holds the default value on creation for the updated_at field.
	prefecture.DefaultUpdatedAt = prefectureDescUpdatedAt.Default.(func() time.Time)
	// prefecture.UpdateDefaultUpdatedAt holds the default value on update for the updated_at field.
	prefecture.UpdateDefaultUpdatedAt = prefectureDescUpdatedAt.UpdateDefault.(func() time.Time)
	// prefectureDescID is the schema descriptor for id field.
	prefectureDescID := prefectureMixinFields1[0].Descriptor()
	// prefecture.DefaultID holds the default value on creation for the id field.
	prefecture.DefaultID = prefectureDescID.Default.(func() string)
	// prefecture.IDValidator is a validator for the "id" field. It is called by the builders before save.
	prefecture.IDValidator = prefectureDescID.Validators[0].(func(string) error)
	recruitmentMixin := schema.Recruitment{}.Mixin()
	recruitmentMixinFields0 := recruitmentMixin[0].Fields()
	_ = recruitmentMixinFields0
	recruitmentMixinFields1 := recruitmentMixin[1].Fields()
	_ = recruitmentMixinFields1
	recruitmentFields := schema.Recruitment{}.Fields()
	_ = recruitmentFields
	// recruitmentDescCreatedAt is the schema descriptor for created_at field.
	recruitmentDescCreatedAt := recruitmentMixinFields1[0].Descriptor()
	// recruitment.DefaultCreatedAt holds the default value on creation for the created_at field.
	recruitment.DefaultCreatedAt = recruitmentDescCreatedAt.Default.(func() time.Time)
	// recruitmentDescUpdatedAt is the schema descriptor for updated_at field.
	recruitmentDescUpdatedAt := recruitmentMixinFields1[1].Descriptor()
	// recruitment.DefaultUpdatedAt holds the default value on creation for the updated_at field.
	recruitment.DefaultUpdatedAt = recruitmentDescUpdatedAt.Default.(func() time.Time)
	// recruitment.UpdateDefaultUpdatedAt holds the default value on update for the updated_at field.
	recruitment.UpdateDefaultUpdatedAt = recruitmentDescUpdatedAt.UpdateDefault.(func() time.Time)
	// recruitmentDescTitle is the schema descriptor for title field.
	recruitmentDescTitle := recruitmentFields[0].Descriptor()
	// recruitment.TitleValidator is a validator for the "title" field. It is called by the builders before save.
	recruitment.TitleValidator = recruitmentDescTitle.Validators[0].(func(string) error)
	// recruitmentDescContent is the schema descriptor for content field.
	recruitmentDescContent := recruitmentFields[4].Descriptor()
	// recruitment.ContentValidator is a validator for the "content" field. It is called by the builders before save.
	recruitment.ContentValidator = recruitmentDescContent.Validators[0].(func(string) error)
	// recruitmentDescID is the schema descriptor for id field.
	recruitmentDescID := recruitmentMixinFields0[0].Descriptor()
	// recruitment.DefaultID holds the default value on creation for the id field.
	recruitment.DefaultID = recruitmentDescID.Default.(func() string)
	// recruitment.IDValidator is a validator for the "id" field. It is called by the builders before save.
	recruitment.IDValidator = recruitmentDescID.Validators[0].(func(string) error)
	userMixin := schema.User{}.Mixin()
	userMixinFields0 := userMixin[0].Fields()
	_ = userMixinFields0
	userMixinFields1 := userMixin[1].Fields()
	_ = userMixinFields1
	userFields := schema.User{}.Fields()
	_ = userFields
	// userDescCreatedAt is the schema descriptor for created_at field.
	userDescCreatedAt := userMixinFields0[0].Descriptor()
	// user.DefaultCreatedAt holds the default value on creation for the created_at field.
	user.DefaultCreatedAt = userDescCreatedAt.Default.(func() time.Time)
	// userDescUpdatedAt is the schema descriptor for updated_at field.
	userDescUpdatedAt := userMixinFields0[1].Descriptor()
	// user.DefaultUpdatedAt holds the default value on creation for the updated_at field.
	user.DefaultUpdatedAt = userDescUpdatedAt.Default.(func() time.Time)
	// user.UpdateDefaultUpdatedAt holds the default value on update for the updated_at field.
	user.UpdateDefaultUpdatedAt = userDescUpdatedAt.UpdateDefault.(func() time.Time)
	// userDescName is the schema descriptor for name field.
	userDescName := userFields[0].Descriptor()
	// user.NameValidator is a validator for the "name" field. It is called by the builders before save.
	user.NameValidator = userDescName.Validators[0].(func(string) error)
	// userDescEmail is the schema descriptor for email field.
	userDescEmail := userFields[1].Descriptor()
	// user.EmailValidator is a validator for the "email" field. It is called by the builders before save.
	user.EmailValidator = userDescEmail.Validators[0].(func(string) error)
	// userDescAvatar is the schema descriptor for avatar field.
	userDescAvatar := userFields[3].Descriptor()
	// user.DefaultAvatar holds the default value on creation for the avatar field.
	user.DefaultAvatar = userDescAvatar.Default.(string)
	// userDescIntroduction is the schema descriptor for introduction field.
	userDescIntroduction := userFields[4].Descriptor()
	// user.IntroductionValidator is a validator for the "introduction" field. It is called by the builders before save.
	user.IntroductionValidator = userDescIntroduction.Validators[0].(func(string) error)
	// userDescEmailVerificationStatus is the schema descriptor for email_verification_status field.
	userDescEmailVerificationStatus := userFields[5].Descriptor()
	// user.DefaultEmailVerificationStatus holds the default value on creation for the email_verification_status field.
	user.DefaultEmailVerificationStatus = userDescEmailVerificationStatus.Default.(bool)
	// userDescID is the schema descriptor for id field.
	userDescID := userMixinFields1[0].Descriptor()
	// user.DefaultID holds the default value on creation for the id field.
	user.DefaultID = userDescID.Default.(func() string)
	// user.IDValidator is a validator for the "id" field. It is called by the builders before save.
	user.IDValidator = userDescID.Validators[0].(func(string) error)
}
