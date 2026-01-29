package utils

import (
	"errors"
	"fmt"
	"mime/multipart"
	"time"

	"github.com/google/uuid"
	"github.com/goutamkumar/golang_restapi_postgresql_test1/internal/models"
	"gorm.io/gorm"
)

type UserResponse struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Avatar   string    `json:"avatar,omitempty"` // optional
}

func ToUserResponse(user *models.User) UserResponse {
	return UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		//Avatar:   user.AvatarURL,
	}
}

func IsNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

// helper/upload.go

func UploadFile(file *multipart.FileHeader) (string, error) {
	// OPTION A: Save to Local Disk (Simple for dev)
	// filename := filepath.Base(file.Filename)
	// dst := "./uploads/" + filename
	// if err := c.SaveUploadedFile(file, dst); err != nil { return "", err }
	// return "/uploads/" + filename, nil

	// OPTION B: Upload to AWS S3 (Production)
	// This is where you would use the AWS SDK to put the file into a bucket.
	// For this example, I will simulate returning a fake S3 URL.

	// Simulate processing time
	// time.Sleep(100 * time.Millisecond)

	// Generate a unique filename to prevent overwrites
	uniqueName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename)

	return "https://my-bucket.s3.amazonaws.com/" + uniqueName, nil
}

func DeleteFile(fileURL string) error {
	// OPTION A: Delete from Local Disk
	// path := "." + fileURL // assuming fileURL is like "/uploads/filename.jpg"
	// return os.Remove(path)

	// OPTION B: Delete from AWS S3
	// This is where you would use the AWS SDK to delete the file from the bucket.
	// For this example, I will simulate a successful deletion.
	// Simulate processing time
	// time.Sleep(50 * time.Millisecond)
	return nil
}
