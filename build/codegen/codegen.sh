KUBE_CODE_GEN_VERSION="kubernetes-1.19.3"
GROUP_VERSIONS="crd.example.com:v1"
scriptdir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
echo $scriptdir
codegendir="${scriptdir}/../../vendor/k8s.io/code-generator"
echo $codegendir
# vendoring k8s.io/code-generator temporarily
echo "require k8s.io/code-generator ${KUBE_CODE_GEN_VERSION}" >> ${scriptdir}/../../go.mod
go mod vendor
# git checkout HEAD ${scriptdir}/../../go.mod ${scriptdir}/../../go.sum

bash $codegendir/generate-groups.sh all github.com/murali-bashyam/crdexample/pkg/client ../pkg/apis crd.example.com:v1 --output-base "${scriptdir}/../../vendor" --go-header-file "${scriptdir}/boilerplate.go.txt"
cp -r "${scriptdir}/../../vendor/github.com/crdexample/pkg" "${scriptdir}/../../"

