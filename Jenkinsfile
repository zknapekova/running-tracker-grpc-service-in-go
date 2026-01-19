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
        stage('Build Image') {
            steps {
                sh 'docker compose build grpcserver'
            }
        }
        stage('Run Tests') {
            parallel {
                stage('Unit Tests') {
                    steps {
                        sh '''
                        cd grpc_server
                        go test -v ./...'''
                    }
                }
                stage('Integration Tests') {
                    steps {
                        sh '''
                            docker compose up -d mongodb grpcserver
                            cd grpc_server
                            go test -v ./tests/...
                            docker compose down'''
                    }
                }
            }
        }
    }
}