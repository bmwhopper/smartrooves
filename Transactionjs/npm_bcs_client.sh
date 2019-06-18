#!/bin/sh
WORK_DIR="${PWD}"
NPM_GLOBAL=0
NPM_PACKAGE=fabric-client
NPM_CA_PACKAGE=fabric-ca-client
NPM_PACKAGE_VER=@1.3.0
NPM_ROOT_DIR=""

usage()
{
  echo "Usage: $0 [-g]" >&2
  echo """
Commands:
      -g                Install a global package
       """
  exit 1
}

check_arg()
{
  if [ "$1" = "-g" ]; then
    NPM_GLOBAL=1
    NPM_ROOT_DIR=$(npm root -g)
  else
    usage
  fi
}

disable_grpc_alpn()
{
  declare -r file="${BASE_DIR}/binding.gyp"
  sed -i".orig" "s/'grpc_alpn%'\s*:.*,/'grpc_alpn%': 'false',/" "${file}"
}

rebuild_code() {
  if [ "$NPM_GLOBAL" -eq 1 ]; then
    cd "${NPM_ROOT_DIR}""/""${NPM_PACKAGE}"
  fi

  npm rebuild --unsafe-perm --build-from-source
}

upgrade_grpc_package()
{
  sed -i".orig" "s/\"grpc\"\s*:.*,/\"grpc\": \"1.17.0\",/" "${PACK_DEFINE_FILE}"
  rm -rf ${BASE_DIR}
  rm -f ${PACK_LOCK_FILE}
  if [ "$NPM_GLOBAL" -eq 1 ]; then
    cd ${NPM_ROOT_DIR}/$NPM_PACKAGE
  fi
  npm install
  disable_grpc_alpn
  rebuild_code
}

main() {
  declare -r npm_global_option=$1
  declare -r package_name=$2

  if [ "$#" -gt 1 ]; then
    usage
  fi

  if [ "$#" -eq 1 ]; then
    check_arg "${npm_global_option}"
  fi

  if [ "$NPM_GLOBAL" -eq 1 ]; then
    BASE_DIR="${NPM_ROOT_DIR}""/$NPM_PACKAGE/node_modules/grpc"
    PACK_DEFINE_FILE="${NPM_ROOT_DIR}""/$NPM_PACKAGE/package.json"
    PACK_LOCK_FILE="${NPM_ROOT_DIR}""/$NPM_PACKAGE/package-lock.json"
    declare -r npm_option="--global --ignore-scripts"
  else
    BASE_DIR="./node_modules/grpc"
    PACK_DEFINE_FILE="./node_modules/$NPM_PACKAGE/package.json"
    PACK_LOCK_FILE="./package-lock.json"
    declare -r npm_option="--ignore-scripts"
  fi

  npm install ${npm_option} log4js@0.6.38
  npm install ${npm_option} fs-extra@6.0.1
  npm install ${npm_option} $NPM_PACKAGE$NPM_PACKAGE_VER
  npm install ${npm_option} $NPM_CA_PACKAGE$NPM_PACKAGE_VER
  upgrade_grpc_package

  cd "${WORK_DIR}"
}

main "$@"
