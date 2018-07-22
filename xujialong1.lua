-- Generated by github.com/0990/taberr
-- Version: 0.0.1

local tab = {
	packageName = {
		{ ErrorType = "Invalid", ErrorMsg = "" 	},
		{ ErrorType = "systemfailed", ErrorMsg = "" 	},
		{ ErrorType = "loginfailed", ErrorMsg = "登录失败" 	},
		{ ErrorType = "systemfailed45", ErrorMsg = "系统失败78" 	},
		{ ErrorType = "SYstemHello", ErrorMsg = "许家龙" 	}
	}

}


-- ErrorType
tab.packageNameByErrorType = {}
for _, rec in pairs(tab.packageName) do
	tab.packageNameByErrorType[rec.ErrorType] = rec
end

tab.Enum = {
	enumName = {
		Invalid = 0,
		systemfailed = 2,
		loginfailed = 1,
		systemfailed45 = 100,
		SYstemHello = 101,
	},
}

return tab