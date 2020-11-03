KUBE_CODE_GEN_VERSION="kubernetes-1.19.3"
GROUP_VERSIONS="crd.example.com:v1"
PKGDIR="github.com/murali-bashyam/crdexample/pkg"
scriptdir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
codegendir="${scriptdir}/../../vendor/k8s.io/code-generator"
echo "require k8s.io/code-generator ${KUBE_CODE_GEN_VERSION}" >> ${scriptdir}/../../go.mod
go mod vendor
cd $codegendir
go mod edit -replace=github.com/murali-bashyam/crdexample=../../../../crdexample
chmod +x ./generate-groups.sh

./generate-groups.sh all ${PKGDIR}/client ${PKGDIR}/apis ${GROUP_VERSIONS} --output-base "${scriptdir}/../../vendor" --go-header-file "${scriptdir}/boilerplate.go.txt"
cd ../../..
cp -r "vendor/${PKGDIR}" .
rm -rf vendor

