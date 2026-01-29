@Library('jenkins-shared-libraries') _

// Read available options in jenkins-shared-libraries
// https://github.com/stayforlong/jenkins-shared-libraries/blob/master/README.md
testPipeline(
    projectName: "krakend-httpcache",
    language: "go",
    slackChannel: "#builds-platform",
    dockerImage: "golang:1.25-bookworm",
    dockerImagePrivateRegistry: false,
    sonarEnabled: false
)