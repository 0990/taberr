-- Generated by github.com/davyxu/tabtoy
-- Version: 2.9.1
-- Modify by nuyan

local tab = {
	emsg = {
		{ ErrorCode = "LoginFailure", Msg = "登录失败" 	},
		{ ErrorCode = "DiamondNotEnough", Msg = "元宝不足" 	},
		{ ErrorCode = "InGame", Msg = "正在游戏中" 	},
		{ ErrorCode = "ConfigNotExist", Msg = "无此配置" 	},
		{ ErrorCode = "HeroConfigNotExist", Msg = "无此英雄配置" 	},
		{ ErrorCode = "PveHeroIDRepeat", Msg = "主副将不能一样" 	},
		{ ErrorCode = "PveHeroIDLimit", Msg = "英雄限制" 	},
		{ ErrorCode = "PveFriendHaveNoStone", Msg = "友将没有这个符石" 	},
		{ ErrorCode = "PveMasterHaveNoStone", Msg = "主将没有这个符石" 	},
		{ ErrorCode = "PveChapterConfigNotExist", Msg = "无章节配置" 	},
		{ ErrorCode = "PveSectionConfigNotExist", Msg = "无关卡配置" 	},
		{ ErrorCode = "PveMasterHeroIDLimit", Msg = "主将英雄限制" 	},
		{ ErrorCode = "PveMinisterHeroIDLimit", Msg = "副将英雄限制" 	},
		{ ErrorCode = "PveMasterHaveNoHero", Msg = "主将无此英雄" 	},
		{ ErrorCode = "PveMinisterHaveNoHero", Msg = "副将无此英雄" 	},
		{ ErrorCode = "PveMinisterHaveNoStone", Msg = "副将没有这个符石" 	},
		{ ErrorCode = "PveStarNoEnough", Msg = "星数不足" 	},
		{ ErrorCode = "PveNoChapterID", Msg = "无此章节" 	},
		{ ErrorCode = "PveNoBoxIndex", Msg = "无此宝箱" 	},
		{ ErrorCode = "PveBoxHaveReceived", Msg = "已领取过此宝箱" 	},
		{ ErrorCode = "PveSectionLocked", Msg = "关卡锁住状态" 	},
		{ ErrorCode = "PveEnergyNotEnough", Msg = "体力不足" 	},
		{ ErrorCode = "PveNoSectionID", Msg = "无此关卡" 	},
		{ ErrorCode = "PveSweepCountLimit", Msg = "扫荡次数限制" 	},
		{ ErrorCode = "PveEnergyBuyCountLimit", Msg = "体力购买次数限制" 	}
	}

}


-- ErrorCode
tab.emsgByErrorCode = {}
for _, rec in pairs(tab.emsg) do
	tab.emsgByErrorCode[rec.ErrorCode] = rec
end

tab.Enum = {
	Err = {
		[50006] = "PveMasterHeroIDLimit",
		[50007] = "PveMinisterHeroIDLimit",
		[50008] = "PveMasterHaveNoHero",
		LoginFailure = 10000,
		DiamondNotEnough = 20000,
		InGame = 30000,
		ConfigNotExist = 40000,
		HeroConfigNotExist = 40001,
		PveHeroIDLimit = 50000,
		PveHeroIDRepeat = 50001,
		PveFriendHaveNoStone = 50002,
		PveMasterHaveNoStone = 50003,
		PveChapterConfigNotExist = 50004,
		PveSectionConfigNotExist = 50005,
	},
}

return tab