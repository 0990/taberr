-- Generated by github.com/0990/taberr
-- Version: 0.0.1

local tab = {
	packageName = {
		{ ErrorID = 0, ErrorType = "Invalid", ErrorMsg = "" 	},
		{ ErrorID = 2, ErrorType = "systemfailed", ErrorMsg = "" 	},
		{ ErrorID = 1, ErrorType = "loginfailed", ErrorMsg = "登录失败" 	},
		{ ErrorID = 100, ErrorType = "systemfailed45", ErrorMsg = "系统失败78" 	},
		{ ErrorID = 101, ErrorType = "SYstemHello", ErrorMsg = "许家龙" 	}
	}

}


-- ErrorID
tab.packageNameByErrorID = {}
for _, rec in pairs(tab.packageName) do
	tab.packageNameByErrorID[rec.ErrorID] = rec
end

return tab