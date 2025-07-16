package template

import (
	"context"
	"larsa-tourism-microservices/pkg/services/messaging/enums"
)

type MsgTpl interface {
	Construct(ctx context.Context, data any) (string, string, error)
}

var Templates = map[enums.MsgTyps]MsgTpl{
	enums.INVITATION: &InvitationTpl{
		Subject: `Welcome to %s on %s`, //no need for html tpl
		Message: `<div dir="ltr">
					<p>Hi {{.MemberName}}</p>
					<p>We're excited to welcome you to {{.CompanyName}} on {{.PlatformName}}. Your account has been successfully created, and you can start using our platform immediately.
					</p>
					<p>Here are your login details:</p>
					<p>Email: {{.Email}}</p>
					<p>Temp password: {{.Password}}</p>
					<p>Please log in using the link below and update your password for security reasons: <a href="{{.Link}}">Click here to login</a></p>
					<p>Looking forward to collaborating with you!</p>

					<p><strong>Best Regards,</strong></p>
					<p><strong>{{.CompanyName}}</strong></p>
				</div>
				`,
	},
	enums.ACCOUNTUPDATED: &AccountUpdatedTpl{
		Subject: `Your Account Has Been Updated`,
		Message: `<div dir="ltr">
					<p>Hi {{.MemberName}},</p>
					<p>This email is to confirm that your account on {{.PlatformName}} has been successfully updated.</p>
					<p>Your account details have been modified as requested. If you did not request these changes, please contact our support team immediately.</p>
					<p>Here are your login details:</p>
					<p>Email: {{.Email}}</p>
					<p>Password: {{.Password}}</p>
					<p>You can access your account at: <a href="{{.Link}}">Login to your account</a></p>
					<p>If you have any questions or need assistance, please don't hesitate to contact us.</p>

					<p><strong>Best Regards,</strong></p>
					<p><strong>{{.CompanyName}}</strong></p>
				</div>
				`,
	},
	enums.WELCOME_NEWSLETTER: &NewsletterWelcomeTpl{
		Subject: "Welcome to our Newsletter!",
		Message: `<div dir="ltr">
		<p>Hi,</p>
		<p>Thank you for subscribing to our newsletter!</p>
		<p>We're excited to keep you updated with our latest news, offers, and updates.</p>
		<p>If you have any questions, feel free to reply to this email.</p>
		<p><strong>Best Regards,</strong></p>
		<p><strong>{{.CompanyName}}</strong></p>
	</div>`,
	},
	//...other templates
}
