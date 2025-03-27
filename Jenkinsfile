pipeline {
    agent any
    
    parameters {
        string(name: 'ENVIRONMENT_KEY', defaultValue: 'main', description: '\u041a\u043b\u044e\u0447 \u043e\u043a\u0440\u0443\u0436\u0435\u043d\u0438\u044f \u0434\u043b\u044f \u043f\u0440\u0438\u043b\u043e\u0436\u0435\u043d\u0438\u044f')
    }
    
    environment {
        APP_PORT = '8888'
        CONTAINER_NAME = 'guess-game'
        DOCKER_HUB_CREDENTIALS = credentials('docker-hub-credentials')
        DOCKER_IMAGE_NAME = 'nswalpakhart/guess-game-app'
        GOROOT = "${WORKSPACE}/go"
        PATH = "${WORKSPACE}/go/bin:${env.PATH}"
    }
    
    stages {
        stage('Install Dependencies') {
            steps {
                sh '''
                    sudo -S apt-get update || true
                    sudo -S DEBIAN_FRONTEND=noninteractive apt-get install -y git wget curl docker.io || true
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
                    wget https://go.dev/dl/go1.22.0.linux-amd64.tar.gz
                    sudo rm -rf /usr/local/go
                    sudo tar -C /usr/local -xzf go1.22.0.linux-amd64.tar.gz
                    export PATH=$PATH:/usr/local/go/bin
                    
                    if [ ! -f go.mod ]; then
                        go mod tidy || true
                    fi

                    go test . -v || true
                '''
                archiveArtifacts artifacts: 'test-report.txt', allowEmptyArchive: true
            }
        }
        
        stage('Build') {
            steps {
                sh '''
                    chmod 666 /var/run/docker.sock || true

                    rm -f ~/.docker/config.json || true
                    mkdir -p ~/.docker

                    docker logout
                    echo "${DOCKER_HUB_CREDENTIALS_PSW}" | docker login -u "${DOCKER_HUB_CREDENTIALS_USR}" --password-stdin

                    docker info

                    mkdir -p build_temp
                    cp app.go app_test.go go.mod go.sum* build_temp/ || true
                    cp -r templates build_temp/ || true
                    cp Dockerfile build_temp/

                    cd build_temp
                    docker build -t ${DOCKER_HUB_CREDENTIALS_USR}/guess-game-app:${BUILD_NUMBER} -t ${DOCKER_HUB_CREDENTIALS_USR}/guess-game-app:latest .
                    cd ..

                    docker save ${DOCKER_HUB_CREDENTIALS_USR}/guess-game-app:${BUILD_NUMBER} > guess-game-app.tar

                    ls -la guess-game-app.tar

                    echo "${DOCKER_HUB_CREDENTIALS_PSW}" | docker login -u "${DOCKER_HUB_CREDENTIALS_USR}" --password-stdin
                    docker push ${DOCKER_HUB_CREDENTIALS_USR}/guess-game-app:${BUILD_NUMBER} || echo "Ошибка при публикации версии ${BUILD_NUMBER}"
                    docker push ${DOCKER_HUB_CREDENTIALS_USR}/guess-game-app:latest || echo "Ошибка при публикации версии latest"
                '''
                archiveArtifacts artifacts: 'guess-game-app.tar', allowEmptyArchive: false
            }
        }
        
        stage('Deploy') {
            steps {
                sh '''
                    sudo -S docker stop ${CONTAINER_NAME} || true
                    sudo -S docker rm ${CONTAINER_NAME} || true
                    sudo -S docker run -d \
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
                                sudo -S docker start guess-game || true
                            '''
                        }
                    }
                }
            }
        }
    }
}
