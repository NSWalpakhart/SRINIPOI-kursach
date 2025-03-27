pipeline {
    agent any
    
    parameters {
        string(name: 'ENVIRONMENT_KEY', defaultValue: 'main', description: 'Environment key for the application')
        string(name: 'LOCAL_PC_IP', defaultValue: '192.168.0.203', description: 'Local PC IP address')
        string(name: 'SSH_USER', defaultValue: 'walpakhart', description: 'Local PC username')
    }
    
    environment {
        APP_PORT = '8888'
        CONTAINER_NAME = 'guess-game'
        DOCKER_HUB_CREDENTIALS = credentials('docker-hub-credentials')
        DOCKER_IMAGE_NAME = 'nswalpakhart/guess-game-app'
    }
    
    stages {
        stage('Deploy to Local PC') {
            steps {
                withCredentials([sshUserPrivateKey(credentialsId: 'local-pc-ssh-key', keyFileVariable: 'SSH_KEY_FILE')]) {
                    sh """#!/bin/bash
                        # Check SSH connection
                        ssh -o StrictHostKeyChecking=no -i "\${SSH_KEY_FILE}" ${params.SSH_USER}@${params.LOCAL_PC_IP} 'echo "SSH connection successful"' || echo "SSH connection failed but continuing"
                        
                        # Install dependencies on local PC
                        ssh -o StrictHostKeyChecking=no -i "\${SSH_KEY_FILE}" ${params.SSH_USER}@${params.LOCAL_PC_IP} '
                            # Check and install all required tools
                            echo "Installing required tools..."
                            sudo apt update
                            sudo apt install -y docker.io git wget curl
                            
                            # Configure Docker
                            sudo systemctl start docker || true
                            sudo systemctl enable docker || true
                            sudo usermod -aG docker ${params.SSH_USER}
                            sudo chmod 666 /var/run/docker.sock || true
                            
                            # Clean previous Docker configurations
                            rm -f ~/.docker/config.json || true
                            mkdir -p ~/.docker
                            
                            # Docker Hub authentication
                            docker logout
                            echo "${DOCKER_HUB_CREDENTIALS_PSW}" | docker login -u "${DOCKER_HUB_CREDENTIALS_USR}" --password-stdin
                            docker info
                            
                            # Install Go with pre-download check
                            echo "Installing Go..."
                            if ! command -v go &> /dev/null || ! go version | grep -q "go1.2"; then
                                wget -q --timeout=30 --tries=3 https://go.dev/dl/go1.22.0.linux-amd64.tar.gz
                                if [ -f go1.22.0.linux-amd64.tar.gz ]; then
                                    sudo rm -rf /usr/local/go
                                    sudo tar -C /usr/local -xzf go1.22.0.linux-amd64.tar.gz
                                    ls -la /usr/local/go/bin || echo "Go directory not created"
                                    rm -f go1.22.0.linux-amd64.tar.gz
                                    
                                    # Add Go to PATH for current session
                                    export PATH=$PATH:/usr/local/go/bin
                                    echo "export PATH=\$PATH:/usr/local/go/bin" >> ~/.bashrc
                                    source ~/.bashrc
                                else
                                    echo "Error downloading Go"
                                fi
                            else
                                echo "Go is already installed"
                                go version
                            fi
                            
                            # Check tools availability
                            echo "Checking tools availability..."
                            docker --version || echo "Docker not available"
                            git --version || echo "Git not available"
                            wget --version || echo "Wget not available"
                            curl --version || echo "Curl not available"
                            /usr/local/go/bin/go version || echo "Go not available"
                            
                            # Create and switch to working directory
                            rm -rf ~/guess-game-app
                            mkdir -p ~/guess-game-app
                            cd ~/guess-game-app
                            
                            # Clone repository
                            git clone https://github.com/NSWalpakhart/SRINIPOI-kursach.git .
                            git checkout ${params.ENVIRONMENT_KEY}
                            
                            # Check files before starting
                            echo "Checking files..."
                            ls -la
                            
                            # Create go.mod file with correct Go version
                            if [ ! -f "go.mod" ]; then
                                echo "module github.com/NSWalpakhart/SRINIPOI-kursach" > go.mod
                                echo "" >> go.mod
                                echo "go 1.20" >> go.mod
                            fi
                            
                            # Check and create main.go if missing
                            if [ ! -f "main.go" ] && [ -f "app.go" ]; then
                                # Use app.go directly instead of copying
                                echo "Using app.go as main file"
                            fi
                            
                            # Remove main.go if both files exist to avoid conflicts
                            if [ -f "main.go" ] && [ -f "app.go" ]; then
                                echo "Found both main.go and app.go. Removing main.go to avoid conflicts"
                                rm main.go
                            fi
                            
                            # Check source code before building
                            echo "Checking server settings..."
                            grep -E "http.ListenAndServe|ListenAndServe" *.go || echo "Could not find server settings"
                            
                            # Run tests
                            echo "Running tests..."
                            /usr/local/go/bin/go test . -v || true
                            
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
                            docker push ${DOCKER_HUB_CREDENTIALS_USR}/guess-game-app:${BUILD_NUMBER} || echo "Error pushing version ${BUILD_NUMBER}"
                            docker push ${DOCKER_HUB_CREDENTIALS_USR}/guess-game-app:latest || echo "Error pushing latest version"

                            # Check and stop existing container
                            echo "Stopping existing container..."
                            sudo -S docker stop ${CONTAINER_NAME} || true
                            sudo -S docker rm ${CONTAINER_NAME} || true
                            
                            # Run new container with proper name and network settings
                            echo "Starting new container..."
                            sudo -S docker run -d \
                                --name ${CONTAINER_NAME} \
                                -p ${APP_PORT}:${APP_PORT} \
                                --restart unless-stopped \
                                --network=host \
                                ${DOCKER_IMAGE_NAME}:${BUILD_NUMBER}

                            # Increased wait time for application startup
                            echo "Waiting for application to start..."
                            sleep 10
                            
                            MAX_ATTEMPTS=5
                            ATTEMPT=1
                            
                            while [ \$ATTEMPT -le \$MAX_ATTEMPTS ]; do
                                echo "Checking application availability (attempt \$ATTEMPT)..."
                                if curl -s -f http://localhost:${APP_PORT} > /dev/null 2>&1; then
                                    echo "Application started successfully!"
                                    break
                                else
                                    if [ \$ATTEMPT -eq \$MAX_ATTEMPTS ]; then
                                        echo "Application not responding"
                                        docker logs ${CONTAINER_NAME}
                                    else
                                        sleep 10
                                        ATTEMPT=\$((ATTEMPT+1))
                                    fi
                                fi
                            done
                        ' || echo "Deployment failed but continuing"
                    """
                }
            }
        }
    }
    
    post {
        failure {
            withCredentials([sshUserPrivateKey(credentialsId: 'local-pc-ssh-key', keyFileVariable: 'SSH_KEY_FILE')]) {
                sh """#!/bin/bash
                    # Start container if stopped
                    ssh -o StrictHostKeyChecking=no -i "\${SSH_KEY_FILE}" ${params.SSH_USER}@${params.LOCAL_PC_IP} '
                        docker start ${CONTAINER_NAME} || true
                    ' || echo "Failed to restart container, but continuing"
                """
            }
        }
    }
}
