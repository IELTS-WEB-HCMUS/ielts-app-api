package services

import (
	"context"
	"encoding/json"
	"errors"
	"ielts-web-api/common"
	"ielts-web-api/internal/models"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/dlclark/regexp2"
)

var JWTSecret = []byte("your_secret_key")

func (s *Service) CheckDuplicatedEmail(ctx context.Context, email string) (bool, error) {
	_, err := s.userRepo.GetDetailByConditions(ctx, func(tx *gorm.DB) {
		tx.Where("email = ?", email)
	})

	// If no user is found, the email is not duplicated
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil // Email is NOT duplicated
	}

	// If there's an error other than "record not found", return the error
	if err != nil {
		return false, err
	}

	// If no error and a user is found, the email is duplicated
	return true, nil // Email IS duplicated
}

func (s *Service) SignupUser(ctx context.Context, req models.SignupRequest) error {
	storedOTP, err := s.otpRepo.GetDetailByConditions(ctx, func(tx *gorm.DB) {
		tx.Where("target = ? AND type = ?", req.Email, common.VERIFY_EMAIL_TYPE)
		tx.Order("created_at desc")
	})

	if err != nil {
		return err
	}

	if storedOTP.VerifyToken != req.VerifyToken || !storedOTP.IsVerified {
		return common.ErrInvalidVerifyToken
	}

	mailPattern := regexp2.MustCompile(`^((?!\.)[\w\-_.]*[^.])(@\w+)(\.\w+(\.\w+)?[^.\W])$`, regexp2.None)
	isValidMail, _ := mailPattern.MatchString(req.Email)
	if !isValidMail {
		return common.ErrInvalidEmailFormat
	}

	passwordPattern := regexp2.MustCompile(`^(?=.*[0-9])(?=.*[a-z])(?=.*[A-Z])(?=.*\W)(?!.* ).{8,40}$`, regexp2.None)
	isStrongPassword, _ := passwordPattern.MatchString(req.Password)
	if !isStrongPassword {
		return common.ErrWeakPassword
	}

	_, err = s.userRepo.GetDetailByConditions(ctx, func(tx *gorm.DB) {
		tx.Where("email = ?", req.Email)
	})

	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	if req.Role == common.ROLE_END_USER {
		newUser := models.User{
			Email:              req.Email,
			Password:           string(hashedPassword),
			RoleID:             common.ROLE_END_USER_UUID,
			FirstName:          &req.FirstName,
			LastName:           &req.LastName,
			EmailNotifications: true,
			Avatar:             common.DEFAULT_AVATAR,
			VocabUsageCount:    common.DEFAULT_VOCAB_COUNT,
			IsBanned:           false,
		}
		user, err := s.userRepo.Create(ctx, &newUser)
		if err != nil {
			return err
		}

		defaultDate := "1900-01-01" // default date
		parsedTime, err := time.Parse(time.DateOnly, defaultDate)
		if err != nil {
			return err
		}
		newUserTarget := models.Target{
			ID:              user.ID,
			TargetReading:   -1,
			TargetListening: -1,
			TargetSpeaking:  -1,
			TargetWriting:   -1,
			NextExamDate:    parsedTime,
		}
		_, err = s.targetRepo.Create(ctx, &newUserTarget)
		if err != nil {
			return err
		}

		categories := []string{"Topic 1", "Topic 2", "Topic 3", "Topic 4"}
		for _, category := range categories {
			newUserVocabCategory := models.UserVocabCategory{
				Name:   category,
				UserID: user.ID,
			}
			_, err = s.vocabCategoriesRepo.Create(ctx, &newUserVocabCategory)
			if err != nil {
				return err
			}
		}
	} else {
		return common.ErrRoleNotFound
	}

	return nil
}

func (s *Service) LoginUser(ctx context.Context, req models.LoginRequest) (*string, error) {
	var user *models.User
	var err error

	if req.IdToken != nil {
		googleUser, err := verifyGoogleOAuthToken(*req.IdToken)
		if err != nil {
			return nil, err
		}

		user, err = s.userRepo.GetDetailByConditions(ctx, func(tx *gorm.DB) {
			tx.Where("email = ? AND provider= ?", googleUser.Email, common.USER_PROVIDER_GOOGLE)
		})
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				newUser := models.User{
					FirstName:          &googleUser.GivenName,
					LastName:           &googleUser.FamilyName,
					Email:              googleUser.Email,
					RoleID:             common.ROLE_END_USER_UUID,
					Provider:           common.USER_PROVIDER_GOOGLE,
					EmailNotifications: true,
					Avatar:             common.DEFAULT_AVATAR,
					VocabUsageCount:    common.DEFAULT_VOCAB_COUNT,
					IsBanned:           false,
				}
				user, err = s.userRepo.Create(ctx, &newUser)
				if err != nil {
					return nil, err
				}
				defaultDate := "1900-01-01" // default date
				parsedTime, err := time.Parse(time.DateOnly, defaultDate)
				if err != nil {
					return nil, err
				}
				newUserTarget := models.Target{
					ID:              user.ID,
					TargetReading:   -1,
					TargetListening: -1,
					TargetSpeaking:  -1,
					TargetWriting:   -1,
					NextExamDate:    parsedTime,
				}
				_, err = s.targetRepo.Create(ctx, &newUserTarget)
				if err != nil {
					return nil, err
				}
			} else {
				return nil, err
			}
		}
	} else {
		user, err = s.userRepo.GetDetailByConditions(ctx, func(tx *gorm.DB) {
			tx.Where("email = ?", req.Email)
		})

		if user.Provider == common.USER_PROVIDER_GOOGLE {
			return nil, common.ErrGoogleAccount
		}

		if err != nil {
			return nil, err
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(*req.Password)); err != nil {
			return nil, common.ErrInvalidEmailOrPassWord
		}
	}

	return generateJWTToken(user)
}

func generateJWTToken(user *models.User) (*string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		"role":  user.RoleID,
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(JWTSecret)
	if err != nil {
		return nil, err
	}

	return &tokenString, nil
}

func verifyGoogleOAuthToken(idToken string) (*models.GoogleUser, error) {
	resp, err := http.Get("https://www.googleapis.com/oauth2/v3/tokeninfo?id_token=" + idToken)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, common.ErrInvalidGoogleAuthenToken
	}
	var googleUser models.GoogleUser
	if err := json.Unmarshal(bodyBytes, &googleUser); err != nil {
		return nil, err
	}
	return &googleUser, nil
}

func (s *Service) GenerateOTP(ctx context.Context, email string, typeToSend string) (string, error) {
	otp := common.GenerateRandomOTP()

	expiry := time.Now().UTC().Add(1 * time.Minute)

	existingOTP, err := s.otpRepo.GetDetailByConditions(ctx, func(tx *gorm.DB) {
		tx.Where("target = ? AND type = ?", email, typeToSend)
	})

	if err == nil {
		existingOTP.IsVerified = true
		_, err = s.otpRepo.Update(ctx, existingOTP.ID, existingOTP)
		if err != nil {
			return "", common.ErrFailedToInValidateExistingOTP
		}
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", err
	}

	newOTP := models.OTP{
		Target:     email,
		Type:       typeToSend,
		OTPCode:    otp,
		ExpiredAt:  expiry,
		IsVerified: false,
	}

	_, err = s.otpRepo.Create(ctx, &newOTP)
	if err != nil {
		return "", err
	}

	return otp, nil
}

func (s *Service) ValidateOTP(ctx context.Context, email, otp string, typeToValidate string) (string, error) {
	storedOTP, err := s.otpRepo.GetDetailByConditions(ctx, func(tx *gorm.DB) {
		tx.Where("target = ? AND type = ?", email, typeToValidate)
		tx.Order("created_at desc")
	})
	if err != nil {
		return "", err
	}
	if storedOTP.IsVerified {
		return "", common.ErrOTPAlreadyVerified
	}

	expiryTime, err := common.NormalizeToBangkokTimezone(storedOTP.ExpiredAt)
	if err != nil {
		return "", err
	}
	currentTime, err := common.NormalizeToBangkokTimezone(time.Now())
	if err != nil {
		return "", err
	}

	newAttempt := models.OTPAttempt{
		OTPID:     storedOTP.ID,
		IsSuccess: false,
		CreatedAt: currentTime,
	}

	if expiryTime.Before(currentTime) {
		newAttempt.IsSuccess = false
		_, _ = s.otpAttemptRepo.Create(ctx, &newAttempt)
		return "", common.ErrOTPExpired
	}

	if storedOTP.OTPCode != otp {
		newAttempt.IsSuccess = false
		_, _ = s.otpAttemptRepo.Create(ctx, &newAttempt)
		return "", common.ErrInvalidOTP
	}

	storedOTP.IsVerified = true
	verifyToken, err := common.GenerateBase64Token(32)
	storedOTP.VerifyToken = verifyToken
	if err != nil {
		return "", common.ErrOtpVerityTokenCreateFailed
	}

	_, err = s.otpRepo.Update(ctx, storedOTP.ID, storedOTP)
	if err != nil {
		return "", common.ErrFailedToUpdateOTPStatus
	}

	newAttempt.IsSuccess = true
	_, _ = s.otpAttemptRepo.Create(ctx, &newAttempt)

	return verifyToken, nil
}

func (s *Service) ResetPassword(ctx context.Context, req models.ResetPasswordRequest) error {
	user, err := s.userRepo.GetDetailByConditions(ctx, func(tx *gorm.DB) {
		tx.Where("email = ?", req.Email)
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return common.ErrEmailNotFound
		}
		return err
	}

	storedOTP, err := s.otpRepo.GetDetailByConditions(ctx, func(tx *gorm.DB) {
		tx.Where("target = ? AND type = ?", req.Email, common.RESET_PASSSWORD_TYPE)
		tx.Order("created_at desc")
	})

	if err != nil {
		return err
	}

	if storedOTP.VerifyToken != req.VerifyToken || !storedOTP.IsVerified {
		return common.ErrInvalidVerifyToken
	}

	passwordPattern := regexp2.MustCompile(`^(?=.*[0-9])(?=.*[a-z])(?=.*[A-Z])(?=.*\W)(?!.* ).{8,40}$`, regexp2.None)
	isStrongPassword, _ := passwordPattern.MatchString(req.NewPassword)
	if !isStrongPassword {
		return common.ErrWeakPassword
	}

	// Check if the new password is the same as the old one
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.NewPassword))
	if err == nil {
		// If no error, it means the password matches
		return common.ErrPasswordDuplicated
	} else if !errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		// If the error is something other than mismatch, return it
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	updatedUser := models.User{
		Password: string(hashedPassword),
	}

	return s.userRepo.UpdatesByConditions(ctx, &updatedUser, func(tx *gorm.DB) {
		tx.Where("email = ?", req.Email)
	})
}
