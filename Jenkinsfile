pipeline {
    agent any
    
    parameters {
        string(name: 'ENVIRONMENT_KEY', defaultValue: 'main', description: 'Ключ окружения для приложения')
    }
    
    environment {
        APP_PORT = '8888'
        CONTAINER_NAME = 'guess-game'
        DOCKER_HUB_CREDENTIALS = credentials('docker-hub-credentials')
        DOCKER_IMAGE_NAME = 'nswalpakhart/guess-game-app'
    }
    
    stages {
        stage('Install Dependencies') {
            steps {
                sh '''
                    sudo apt-get update
                    sudo DEBIAN_FRONTEND=noninteractive apt-get install -y git wget curl docker.io golang-go
                '''
            }
        }
        
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
                sh '''
                    wget https://go.dev/dl/go1.20.linux-amd64.tar.gz
                    sudo rm -rf /usr/local/go && sudo -S tar -C /usr/local -xzf go1.20.linux-amd64.tar.gz
                    export PATH=$PATH:/usr/local/go/bin
                    
                    [ ! -f go.mod ] && go mod init guess-game
                    go mod tidy
                    
                    go test -v > test-report.txt || true
                '''
                archiveArtifacts artifacts: 'test-report.txt', allowEmptyArchive: true
            }
        }
        
        stage('Build') {
            steps {
                sh '''
                    sudo chmod 666 /var/run/docker.sock
                    echo ${DOCKER_HUB_CREDENTIALS_PSW} | docker login -u ${DOCKER_HUB_CREDENTIALS_USR} --password-stdin
                    docker build -t ${DOCKER_IMAGE_NAME}:${BUILD_NUMBER} -t ${DOCKER_IMAGE_NAME}:latest .
                    docker save ${DOCKER_IMAGE_NAME}:${BUILD_NUMBER} > guess-game-app.tar
                    docker push ${DOCKER_IMAGE_NAME}:${BUILD_NUMBER}
                '''
                archiveArtifacts artifacts: 'guess-game-app.tar', allowEmptyArchive: false
                }
        }
        
        stage('Deploy') {
            steps {
                sh '''
                    docker stop ${CONTAINER_NAME} || true
                    docker rm ${CONTAINER_NAME} || true
                    docker run -d \
                        --name ${CONTAINER_NAME} \
                        -p ${APP_PORT}:${APP_PORT} \
                        --restart unless-stopped \
                        ${DOCKER_IMAGE_NAME}:${BUILD_NUMBER}
                '''
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
