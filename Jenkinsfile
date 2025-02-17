pipeline {
    agent any
    
    parameters {
        string(name: 'ENVIRONMENT_KEY', defaultValue: '', description: 'Ключ окружения для приложения')
    }
    
    environment {
        APP_PORT = '8888'
        CONTAINER_NAME = 'guess-game'
    }
    
    stages {
        stage('Checkout') {
            steps {
                checkout([$class: 'GitSCM',
                    branches: [[name: "${params.ENVIRONMENT_KEY}"]],
                    userRemoteConfigs: [[
                        url: 'https://github.com/NSWalpakhart/SRINIPOI-kursach.git',
                    ]]
                ])
            }
        }
        
        stage('Test') {
            steps {
                dir("SRINIPOI-kursach") {
                sh '''
                    go test ./... -v > test-report.txt || true
                '''
                archiveArtifacts artifacts: 'test-report.txt', allowEmptyArchive: true
                }
            }
        }
        
        stage('Build') {
            steps {
                dir('SRINIPOI-kursach'){
                sh '''
                    docker build -t guess-game-app .
                    docker save guess-game-app > guess-game-app.tar
                '''
                archiveArtifacts artifacts: 'guess-game-app.tar', allowEmptyArchive: false
                }
            }
        }
        
        stage('Deploy') {
            steps {
                dir("SRINIPOI-kursach") {
                sh '''
                    docker stop ${CONTAINER_NAME} || true
                    docker rm ${CONTAINER_NAME} || true
                    docker run -d \
                        --name ${CONTAINER_NAME} \
                        -p ${APP_PORT}:${APP_PORT} \
                        --restart unless-stopped \
                        guess-game-app
                '''
            }
            }
        }
        
        stage('Health Check') {
            steps {
                sh '''
                    sleep 10
                    curl -f http://localhost:${APP_PORT}
                '''
            }
        }
        
        stage('Error Handler') {
            steps {
                dir("SRINIPOI-kursach") {
                script {
                    if (currentBuild.result == 'FAILURE' || currentBuild.result == null) {
                        node('built-in') {
                            sh '''
                                docker start guess-game || true
                            '''
                        }
                    }
                }
            }
            }
        }
    }
}
