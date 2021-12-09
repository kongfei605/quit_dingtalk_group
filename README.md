# quit_dingtalk_group

批量退出钉钉群: 通过exclude指定关键字 排除不想退出的群,逗号分割

操作步骤
1. 运行后 会打开浏览器 提示扫码。
2. 扫码后，请手动点击收起提示下载的浮层
3. 手动点击 联系人，触发自动退群

示例

go run dingding.go (默认不退出带有快猫字样的群)  或者
go run dingding.go -exclude="快猫,Gopher,xxx" (逗号分割 ，不想退出的群名关键字)
