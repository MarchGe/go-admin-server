package common

const UserCreateEmailTemplate = `
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <title>账号开通成功通知</title>
</head>
<body>
<div style="width: 100%; height: 100%; display: flex; flex-direction: column; align-items: center;">
    <div style="display: flex; flex-direction: column; align-items: center; font-size: 15px; color: #333333; padding: 0 50px; background-color: #FFFAF0;">
            <div style="height: 40px;"></div>
            <div>欢迎使用{{.SystemName}}，您的账号已开通成功：</div>
            <div style="background-color: #FFDD99; padding: 10px 40px; margin-top: 30px; border-radius: 4px; white-space: nowrap;">
                <div style="height: 25px; line-height: 25px; font-size: 14px; font-weight: 600;">邮箱：{{.Email}}</div>
                <div style="height: 25px; line-height: 25px; font-size: 14px; font-weight: 600;">密码：{{.Password}}</div>
            </div>
            <div style="margin-top: 20px;">登录 <a href="{{.AccessUrl}}" target="_blank">{{.AccessUrl}}</a> 使用，请及时修改密码！</div>
            <div style="height: 40px;"></div>
    </div>
</div>
</body>
</html>
`

const PasswordResetEmailTemplate = `
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <title>密码重置成功通知</title>
</head>
<body>
<div style="width: 100%; height: 100%; display: flex; flex-direction: column; align-items: center;">
    <div style="display: flex; flex-direction: column; align-items: center; font-size: 15px; color: #333333; padding: 0 50px; background-color: #FFFAF0;">
            <div style="height: 40px;"></div>
            <div>欢迎使用{{.SystemName}}，您的密码已重置成功：</div>
            <div style="background-color: #FFDD99; padding: 10px 40px; margin-top: 30px; border-radius: 4px; white-space: nowrap;">
                <div style="height: 25px; line-height: 25px; font-size: 14px; font-weight: 600;">邮箱：&emsp;{{.Email}}</div>
                <div style="height: 25px; line-height: 25px; font-size: 14px; font-weight: 600;">新密码：{{.Password}}</div>
            </div>
            <div style="margin-top: 20px;">登录 <a href="{{.AccessUrl}}" target="_blank">{{.AccessUrl}}</a> 使用，请及时修改密码！</div>
            <div style="height: 40px;"></div>
    </div>
</div>
</body>
</html>
`
