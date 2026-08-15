# download_protoc.ps1
# Script to automate downloading and unzipping protoc compiler on Windows.

$ProgressPreference = 'SilentlyContinue'

$version = "25.3"
$url = "https://github.com/protocolbuffers/protobuf/releases/download/v$version/protoc-$version-win64.zip"
$toolsDir = "c:\Users\ranji\OneDrive\Desktop\MobileAutomation\tools"
$zipFile = Join-Path $toolsDir "protoc.zip"
$extractDir = Join-Path $toolsDir "protoc"

Write-Output "Creating tools directories..."
if (!(Test-Path $extractDir)) {
    New-Item -ItemType Directory -Force -Path $extractDir | Out-Null
}

Write-Output "Downloading protoc v$version from $url..."
Invoke-WebRequest -Uri $url -OutFile $zipFile

Write-Output "Extracting compiler zip..."
Expand-Archive -Path $zipFile -DestinationPath $extractDir -Force

Write-Output "Cleaning up zip file..."
Remove-Item -Path $zipFile -Force

Write-Output "Verification: Checking protoc.exe path..."
$protocPath = Join-Path $extractDir "bin\protoc.exe"
if (Test-Path $protocPath) {
    Write-Output "Successfully installed protoc to $protocPath"
    & $protocPath --version
} else {
    Write-Error "Failed to install protoc!"
    exit 1
}
