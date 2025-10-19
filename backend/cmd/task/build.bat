@echo off
echo Starting Task Service Linux build...

REM Set environment variables for cross-compilation
set GOOS=linux
set GOARCH=amd64
set CGO_ENABLED=0

REM Create bin directory if it doesn't exist
if not exist ".\bin" mkdir ".\bin"

REM Build task service program
echo Building crawler program...
go build -o .\bin\crawler .

if %ERRORLEVEL% EQU 0 ( 
    echo Build completed successfully!
    echo Output: bin\crawler
) else (
    echo Build failed!
    exit /b 1
)

