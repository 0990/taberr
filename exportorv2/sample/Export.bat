..\..\..\..\..\..\bin\tabtoy.exe ^
--mode=exportorv2 ^
--csharp_out=.\Config.cs ^
--binary_out=.\Config.bin ^
--pbt_out=.\Config.pbt ^
--proto_out=.\Config.proto ^
--json_out=.\Config.json ^
--lua_out=.\Config.lua ^
--go_out=.\Config.go ^
--combinename=Config ^
--lan=zh_cn ^
--goimportpkg=github.com/davyxu/tabtoy/exportorv2/sample/gamedef ^
Globals.xlsx ^
Sample.xlsx ^
Info.xlsx

@IF %ERRORLEVEL% NEQ 0 pause

: proto转go
..\..\proto\protoc.exe --plugin=protoc-gen-go=..\..\proto\protoc-gen-go.exe --go_out .\gamedef --proto_path=. .\Config.proto
@IF %ERRORLEVEL% NEQ 0 pause


: 表索引
copy .\Config.go .\table\tableindex.go

