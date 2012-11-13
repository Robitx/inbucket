# Compile and package inbucket dist for windows

param([Parameter(Mandatory=$true)]$versionLabel)

set DESKTOP ([Environment]::GetFolderPath("Desktop"))
set GOOS $(go env GOOS)
set GOARCH $(go env GOARCH)

set distname "inbucket-${versionLabel}-${GOOS}_${GOARCH}"
set distdir "$DESKTOP\$distname"

if (Test-Path $distdir) {
    Remove-Item -Force -Recurse $distdir
}

echo "Building $distname..."
md $distdir | Out-Null
go build -o "$distdir/inbucket.exe" -a -v "github.com/jhillyerd/inbucket"

echo "Copying resources..."
Copy-Item LICENSE -Destination "$distdir\LICENSE.txt"
Copy-Item README.md -Destination "$distdir\README.txt"
Copy-Item bin\inbucket.bat -Destination $distdir
Copy-Item etc -Destination $distdir -Recurse
Copy-Item themes -Destination $distdir -Recurse
