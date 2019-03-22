@echo off
set GOOS=linux
cscript /nologo change_package_name.vbs
gopherjs build -m eclib.go -o "../static/eclib.js" > NUL 2>&1
cscript /nologo revert_package_name.vbs
if %ERRORLEVEL% neq 0 GOTO ERROR
echo Success 
exit /b 0
:ERROR
    echo Failure
    exit /b 1