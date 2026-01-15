pipeline {
    agent any
    tools {
        go 'go1.24'
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
            steps {
                sh 'cd grpc_server'
                sh 'go test -v ./...'
            }
        }
    }
}