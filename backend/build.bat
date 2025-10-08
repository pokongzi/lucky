@echo off
echo Starting Linux build...

REM Set environment variables for cross-compilation
set GOOS=linux
set GOARCH=amd64
set CGO_ENABLED=0

REM Create bin directory if it doesn't exist
if not exist "..\bin" mkdir "..\bin"

REM Build backend program
echo Building backend program...
go build -o lucky .

echo Build completed!