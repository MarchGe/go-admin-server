package devops

import (
	devops2 "github.com/MarchGe/go-admin-server/app/admin/apis/devops"
	"testing"
)

func TestFilenameCheck(t *testing.T) {
	tests := []struct {
		FileName     string
		Description  string
		ExpectedPass bool
	}{
		{FileName: "-xxx.doc", Description: "测试文件名不能以“-”开头", ExpectedPass: false},
		{FileName: "xxx-.doc", Description: "测试文件名非开头可以包含“-”", ExpectedPass: true},
		{FileName: "xx!#x.doc", Description: "测试文件名不能包含特殊字符", ExpectedPass: false},
	}
	for _, test := range tests {
		err := devops2.FilenameCheck(test.FileName)
		if (test.ExpectedPass && err == nil) || (!test.ExpectedPass && err != nil) {
			t.Logf("\n-----\n用例描述：%s\n文件名：  %s		预期：%t		结果：%s	原因：%v", test.Description, test.FileName, test.ExpectedPass, "一致", err)
		} else {
			t.Errorf("\n-----\n用例描述：%s\n文件名：  %s		预期：%t		结果：%s	原因：%v", test.Description, test.FileName, test.ExpectedPass, "不一致", err)
		}
	}
}

func TestCleanFilename(t *testing.T) {
	tests := []struct {
		FileName       string
		Description    string
		ExpectedOutput string
	}{
		{FileName: "-xxx.doc", Description: "测试文件名不能以“-”开头", ExpectedOutput: "_xxx.doc"},
		{FileName: "xxx-.doc", Description: "测试文件名非开头可以包含“-”", ExpectedOutput: "xxx-.doc"},
		{FileName: "%xx!#x.d@%oc&", Description: "测试文件名不能包含特殊字符", ExpectedOutput: "_xx__x.d@_oc_"},
	}
	for _, test := range tests {
		fileName := devops2.CleanFilename(test.FileName)
		if test.ExpectedOutput == fileName {
			t.Logf("\n-----\n用例描述：%s\n文件名：  %s		预期：%s		输出：%s	结果：%s", test.Description, test.FileName, test.ExpectedOutput, fileName, "一致")
		} else {
			t.Errorf("\n-----\n用例描述：%s\n文件名：  %s		预期：%s		输出：%s	结果：%s", test.Description, test.FileName, test.ExpectedOutput, fileName, "不一致")
		}
	}
}
