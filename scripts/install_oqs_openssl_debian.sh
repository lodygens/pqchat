#!/usr/bin/env bash

set -x

ARCH=$(uname -m)

if [ "$(id -u)" -eq 0 ]; then
    SUDO=""
else
    if command -v sudo >/dev/null 2>&1; then
        SUDO="sudo"
    else
        echo "❌ ERROR: Must be run as root or have sudo installed."
        exit 1
    fi
fi

# --- CONFIG ---
LIBOQS_PREFIX=/opt/liboqs
OPENSSL_PREFIX=/usr
BUILD_DIR=$HOME/tmp/oqs_build


# Run update + install + cleanup
$SUDO apt-get update && \
$SUDO apt-get install -y \
    build-essential cmake ninja-build git pkg-config \
    libssl-dev wget perl \
    && \
$SUDO rm -rf /var/lib/apt/lists/*


rm -Rf liboqs
git clone --depth 1 https://github.com/open-quantum-safe/liboqs.git
cd liboqs
mkdir -p build && cd build
cmake -GNinja     -DBUILD_SHARED_LIBS=ON     -DOQS_ENABLE_KEM_HYBRID=ON     -DCMAKE_INSTALL_PREFIX="$LIBOQS_PREFIX"     -DCMAKE_BUILD_TYPE=Release ..
ninja
$SUDO ninja install

#[0/1] Install the project...
#-- Install configuration: "Release"
#-- Installing: /lib/cmake/liboqs/liboqsConfig.cmake
#-- Installing: /lib/cmake/liboqs/liboqsConfigVersion.cmake
#-- Installing: /lib/pkgconfig/liboqs.pc
#-- Installing: /lib/liboqs.so.0.15.0
#-- Installing: /lib/liboqs.so.9
#-- Installing: /lib/liboqs.so
#-- Installing: /lib/cmake/liboqs/liboqsTargets.cmake
#-- Installing: /lib/cmake/liboqs/liboqsTargets-release.cmake
#-- Installing: /include/oqs/oqs.h
#-- Installing: /include/oqs/aes_ops.h
#-- Installing: /include/oqs/common.h
#-- Installing: /include/oqs/rand.h
#-- Installing: /include/oqs/sha2_ops.h
#-- Installing: /include/oqs/sha3_ops.h
#-- Installing: /include/oqs/sha3x4_ops.h
#-- Installing: /include/oqs/kem.h
#-- Installing: /include/oqs/sig.h
#-- Installing: /include/oqs/sig_stfl.h
#-- Installing: /include/oqs/kem_bike.h
#-- Installing: /include/oqs/kem_frodokem.h
#-- Installing: /include/oqs/kem_ntruprime.h
#-- Installing: /include/oqs/kem_ntru.h
#-- Installing: /include/oqs/kem_classic_mceliece.h
#-- Installing: /include/oqs/kem_kyber.h
#-- Installing: /include/oqs/kem_ml_kem.h
#-- Installing: /include/oqs/sig_ml_dsa.h
#-- Installing: /include/oqs/sig_falcon.h
#-- Installing: /include/oqs/sig_sphincs.h
#-- Installing: /include/oqs/sig_mayo.h
#-- Installing: /include/oqs/sig_cross.h
#-- Installing: /include/oqs/sig_uov.h
#-- Installing: /include/oqs/sig_snova.h
#-- Installing: /include/oqs/sig_slh_dsa.h
#-- Installing: /include/oqs/oqsconfig.h

cd ../..

rm -Rf  oqs-provider
git clone --depth 1 https://github.com/open-quantum-safe/oqs-provider.git
cd oqs-provider/
mkdir -p build && cd build
cmake -GNinja   -DOQS_INSTALL_PATH="$LIBOQS_PREFIX"   -DOPENSSL_ROOT_DIR="$OPENSSL_PREFIX"   -DCMAKE_BUILD_TYPE=Release ..
ninja
$SUDO ninja install
#[0/1] Install the project...
#-- Install configuration: "Release"
#-- Installing: /usr/lib/aarch64-linux-gnu/ossl-modules/oqsprovider.so
#-- Installing: /usr/local/include/oqs-provider/oqs_prov.h


# Create the profile.d script
cat <<EOF | $SUDO tee /etc/profile.d/oqs-provider.sh >/dev/null
# OQS provider environment for OpenSSL 3
export OPENSSL_MODULES=$OPENSSL_PREFIX/lib/$ARCH-linux-gnu/ossl-modules
export OPENSSL_CONF=/etc/ssl/openssl.cnf
export LD_LIBRARY_PATH=$LIBOQS_PREFIX/lib:\$LD_LIBRARY_PATH
EOF

# Apply to current shell
export OPENSSL_MODULES=$OPENSSL_PREFIX/lib/$ARCH-linux-gnu/ossl-modules
export OPENSSL_CONF=/etc/ssl/openssl.cnf
export LD_LIBRARY_PATH=$LIBOQS_PREFIX/lib:$LD_LIBRARY_PATH


cat <<EOF | $SUDO tee -a /etc/ssl/openssl.cnf >/dev/null

# ======================================================
# REQUIRED FOR LOADING EXTERNAL PROVIDERS IN DEBIAN
# ======================================================

openssl_conf = openssl_init

[openssl_init]
providers = provider_sect
alg_section = algorithm_sect

[algorithm_sect]


############### OQS PROVIDER CONFIGURATION ###############
# Ajout minimal sans conflit : étend provider_sect et active le provider OQS.

[provider_sect]
default = default_sect
oqsprovider = oqsprovider_sect

[default_sect]
activate = 1

[oqsprovider_sect]
activate = 1
module = /usr/lib/$ARCH-linux-gnu/ossl-modules/oqsprovider.so

EOF

echo "✅ Testing OQS provider..."
openssl list -providers | grep OQS && echo "✅ OQS provider active" || echo "⚠️ OQS provider not detected"

echo "✅ Testing KEM algorithms..."
openssl list -kem-algorithms | grep -i kem && echo "✅ KEM algorithms active" || echo "⚠️ KEM algorithms not detected"

# Install liboqs-go
git clone --depth=1 https://github.com/open-quantum-safe/liboqs-go
cd liboqs-go
cat  .config/liboqs-go.pc| sed "s/LIBOQS_INCLUDE_DIR=.*/LIBOQS_INCLUDE_DIR=\/include/"     > temp.pc
mv temp.pc .config/liboqs-go.pc
cat  .config/liboqs-go.pc| sed "s/LIBOQS_LIB_DIR=.*/LIBOQS_LIB_DIR=\/lib/"     > temp.pc
mv temp.pc .config/liboqs-go.pc


echo "export PKG_CONFIG_PATH=$PKG_CONFIG_PATH:$(pwd)/liboqs-go/.config" >>  /etc/profile.d/oqs-provider.sh
export PKG_CONFIG_PATH=$PKG_CONFIG_PATH:$(pwd)/liboqs-go/.config
go run examples/kem/kem.go && echo "✅ OQS Golang bindings installed" || echo "⚠️ OQS Golang bindings installation error"
