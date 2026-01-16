properties([
    parameters([
        choice(
            name: 'ACTION',
            defaultValue: 'test',
            choices: [
                    'test', 'release',
            ]
        ),
    ])
])
goVersion='1.24.1'

pipeline {
    agent any
    tools {
        go "${goVersion}"
        'org.jenkinsci.plugins.docker.commons.tools.DockerTool' 'myDocker'
    }
    stages {
        stage('Check') {
            steps {
                sh 'go version'
                sh 'docker version'
            }
        }
        stage('Unit tests') {
            when {
                expression { params.ACTION == 'test' }
            }
            steps {
                sh(
                    script: """
                        cd grpc_server
                        go test -v ./...""",
                )
            }
        }
    }
}