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
	//...other templates
}
