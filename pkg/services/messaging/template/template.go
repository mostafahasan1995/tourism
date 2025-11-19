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

	enums.CONTACTUS: &ContactUsTpl{
		Subject: "New Contact Us Submission",
		Message: `<!DOCTYPE html>
<html lang="en" xmlns="http://www.w3.org/1999/xhtml">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width,initial-scale=1">
	<meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
	<title>New Contact Us Submission</title>
</head>
<body style="margin:0;padding:0;font-family:Arial,sans-serif;background-color:#f4f4f4;">
	<table role="presentation" style="width:100%;border-collapse:collapse;margin:0;padding:0;background-color:#f4f4f4;">
		<tr>
			<td style="padding:20px 0;">
				<table role="presentation" style="width:600px;border-collapse:collapse;margin:0 auto;background-color:#ffffff;border-radius:8px;box-shadow:0 2px 4px rgba(0,0,0,0.1);">
					<tr>
						<td style="padding:30px;background-color:#ffffff;border-radius:8px 8px 0 0;">
							<h1 style="margin:0 0 20px 0;color:#333333;font-size:24px;font-weight:bold;">New Contact Us Submission</h1>
						</td>
					</tr>
					<tr>
						<td style="padding:0 30px 30px 30px;">
							<table role="presentation" style="width:100%;border-collapse:collapse;">
								<tr>
									<td style="padding:10px 0;border-bottom:1px solid #eeeeee;">
										<strong style="color:#333333;font-size:14px;">Full Name:</strong>
										<span style="color:#666666;font-size:14px;margin-left:10px;">{{.FullName}}</span>
									</td>
								</tr>
								<tr>
									<td style="padding:10px 0;border-bottom:1px solid #eeeeee;">
										<strong style="color:#333333;font-size:14px;">Email Address:</strong>
										<a href="mailto:{{.EmailAddress}}" style="color:#0066cc;font-size:14px;margin-left:10px;text-decoration:none;">{{.EmailAddress}}</a>
									</td>
								</tr>
								<tr>
									<td style="padding:10px 0;border-bottom:1px solid #eeeeee;">
										<strong style="color:#333333;font-size:14px;">Phone Number:</strong>
										<span style="color:#666666;font-size:14px;margin-left:10px;">{{.PhoneNumber}}</span>
									</td>
								</tr>
								<tr>
									<td style="padding:10px 0;border-bottom:1px solid #eeeeee;">
										<strong style="color:#333333;font-size:14px;">How Did You Find Us:</strong>
										<span style="color:#666666;font-size:14px;margin-left:10px;">{{.HowDidYouFindUs}}</span>
									</td>
								</tr>
								<tr>
									<td style="padding:15px 0;">
										<strong style="color:#333333;font-size:14px;display:block;margin-bottom:10px;">Message:</strong>
										<div style="color:#333333;font-size:14px;line-height:1.6;background-color:#f9f9f9;padding:15px;border-radius:4px;border-left:4px solid #0066cc;">
											{{.Message}}
										</div>
									</td>
								</tr>
								{{if .Additional}}
								<tr>
									<td style="padding:15px 0;border-top:1px solid #eeeeee;">
										<strong style="color:#333333;font-size:14px;display:block;margin-bottom:10px;">Additional Fields:</strong>
										<ul style="margin:0;padding-left:20px;color:#666666;font-size:14px;">
											{{range $k, $v := .Additional}}<li style="margin-bottom:5px;"><strong>{{$k}}:</strong> {{$v}}</li>{{end}}
										</ul>
									</td>
								</tr>
								{{end}}
							</table>
						</td>
					</tr>
					<tr>
						<td style="padding:20px 30px;background-color:#f9f9f9;border-radius:0 0 8px 8px;border-top:1px solid #eeeeee;">
							<p style="margin:0;color:#999999;font-size:12px;text-align:center;">This is an automated message from your contact form.</p>
						</td>
					</tr>
				</table>
			</td>
		</tr>
	</table>
</body>
</html>`,
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
