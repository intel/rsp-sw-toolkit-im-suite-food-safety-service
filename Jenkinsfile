rrpBuildGoCode {
    projectKey = 'food-safety-sample'
    dockerBuildOptions = ['--squash', '--build-arg GIT_COMMIT=$GIT_COMMIT']
    testStepsInParallel = false
    buildImage = 'amr-registry.caas.intel.com/rrp/ci-go-build-image:1.12.0-alpine'
    dockerImageName = "rsp/${projectKey}"
    ecrRegistry = "280211473891.dkr.ecr.us-west-2.amazonaws.com"
    customBuildScript = "./build.sh"
    protexProjectName = 'bb-food-safety-service'


    infra = [
        stackName: 'RSP-Codepipeline-food-safety-sample'
    ]

    notify = [
        slack: [ success: '#ima-build-success', failure: '#ima-build-failed' ]
    ]
}
