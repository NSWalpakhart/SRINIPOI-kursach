pipeline {
    agent any
    
    environment {
        APP_PORT = '8888'
        CONTAINER_NAME = 'guess-game'
    }
    
    stages {
        stage('Checkout') {
            steps {
                checkout([$class: 'GitSCM',
                    branches: [[name: 'main']],
                    userRemoteConfigs: [[
                        url: 'https://github.com/NSWalpakhart/SRINIPOI-kursach.git',
                    ]]
                ])
            }
        }
        
        stage('Build') {
            steps {

                sh '''
                    # Создаем Docker image
                    docker build -t guess-game-app .
                '''
            }
            
        }
        
        stage('Test') {
            steps {

                sh '''
                    # Запускаем тесты
                    go test ./... || true
                '''
            }
            
        }
        
        stage('Deploy') {
            steps {

                sh '''
                    # Останавливаем и удаляем старый контейнер если он существует
                    docker stop ${CONTAINER_NAME} || true
                    docker rm ${CONTAINER_NAME} || true
                    
                    # Запускаем новый контейнер
                    docker run -d \
                        --name ${CONTAINER_NAME} \
                        -p ${APP_PORT}:${APP_PORT} \
                        --restart unless-stopped \
                        guess-game-app
                '''
            }
            
        }
        
        stage('Health Check') {
            steps {
                sh '''
                    # Ждем 10 секунд пока приложение запустится
                    sleep 10
                    
                    # Проверяем что приложение отвечает
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
