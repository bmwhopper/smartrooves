<!-- : Begin batch script
@echo off
setlocal

set argC=0
for %%x in (%*) do Set /A argC+=1

if "%argC%" GEQ "2" (
  goto :usage
  exit /B %errorlevel%
)

set NPM_GLOBAL=0
set NPM_PACKAGE=fabric-client
set NPM_CA_PACKAGE=fabric-ca-client
set NPM_PACKAGE_VER=@1.3.0

if "%argC%" == "1" (
  set ARG=%1
  call :check_arg
  if errorlevel 1 (
    exit /B %errorlevel%
  )
)

for /f %%i in ('npm root -g') do set NPM_GLOBAL_ROOT=%%i
if "%NPM_GLOBAL%" == "1" (
  call npm install -g --ignore-scripts %NPM_PACKAGE%%NPM_PACKAGE_VER%
  call npm install -g --ignore-scripts %NPM_CA_PACKAGE%%NPM_PACKAGE_VER%
  if "%NPM_PACKAGE%" == "fabric-client" (
    set ARG=%NPM_GLOBAL_ROOT%\%NPM_PACKAGE%\node_modules\grpc\deps\grpc\src\core\lib\security\security_connector\security_connector.cc
    call :do_edit_code
  )
  cd %NPM_GLOBAL_ROOT%\%NPM_PACKAGE%
  call npm rebuild --unsafe-perm --build-from-source
) else (
  call npm install --ignore-scripts %NPM_PACKAGE%%NPM_PACKAGE_VER%
  call npm install --ignore-scripts %NPM_CA_PACKAGE%%NPM_PACKAGE_VER%
  if "%NPM_PACKAGE%" == "fabric-client" (
    set ARG=.\node_modules\grpc\deps\grpc\src\core\lib\security\security_connector\security_connector.cc
    call :do_edit_code 
  )
  call npm rebuild --unsafe-perm --build-from-source
)
call npm install
cmd
rem exit /B 0

:usage
  echo "Usage: %~nx0 [-g]"
  exit /B 1

:check_arg
  if  "%ARG%" == "-g" (
    set NPM_GLOBAL=1
  ) else (
    goto :usage
  )
  exit /B 0

:do_edit_code
  set filename="%ARG%"
  set toReplace="if (p == nullptr) {"
  set replaceWith="if (false) {"
  call :replace_in_file
  set toReplace="if (!grpc_chttp2_is_alpn_version_supported(p->value.data, p->value.length)) {"
  set replaceWith="if (p != nullptr && !grpc_chttp2_is_alpn_version_supported(p->value.data, p->value.length)) {"
  call :replace_in_file
  exit /B 0

:replace_in_file
cscript //nologo "%~f0?.wsf" %filename% %toReplace% %replaceWith%
exit /b

----- Begin wsf script --->
<job><script language="VBScript">
  Const ForReading = 1
  Const ForWriting = 2

  strFileName = Wscript.Arguments(0)
  strOldText = Wscript.Arguments(1)
  strNewText = Wscript.Arguments(2)

  Set objFSO = CreateObject("Scripting.FileSystemObject")
  Set objFile = objFSO.OpenTextFile(strFileName, ForReading)

  strText = objFile.ReadAll
  objFile.Close
  strNewText = Replace(strText, strOldText, strNewText)

  Set objFile = objFSO.OpenTextFile(strFileName, ForWriting)
  objFile.Write strNewText
  objFile.Close
</script></job>

