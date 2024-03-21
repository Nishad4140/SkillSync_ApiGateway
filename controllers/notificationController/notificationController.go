package notificationController

import (
	"context"

	"github.com/Nishad4140/SkillSync_ProtoFiles/pb"
)

func (n *NotificationController) SendOTP(email string) error {
	_, err := n.Conn.SendOTP(context.Background(), &pb.SendOTPRequest{
		Email: email,
	})
	return err
}
func (n *NotificationController) VerifyOTP(email, otp string) bool {
	res, err := n.Conn.VerifyOTP(context.Background(), &pb.VerifyOTPRequest{
		Otp:   otp,
		Email: email,
	})
	if err != nil {
		return false
	}
	return res.IsVerified
}
