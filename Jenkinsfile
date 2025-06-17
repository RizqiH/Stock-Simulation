pipeline {
    agent {
        docker {
            image 'golang:1.24-alpine'
            args '-v /var/run/docker.sock:/var/run/docker.sock'
        }
    }
    
    environment {
        GO_VERSION = '1.24'
        DOCKER_REGISTRY = 'your-registry.com'
        IMAGE_NAME = 'stock-simulation-api'
        COMPOSE_PROJECT_NAME = 'stock-simulation-test'
    }
    
    stages {
        stage('Checkout') {
            steps {
                checkout scm
                echo 'Code checked out successfully'
            }
        }
        
        stage('Setup Environment') {
            steps {
                sh '''
                    apk add --no-cache docker docker-compose make curl
                    go version
                    docker --version
                    docker-compose --version
                '''
            }
        }
        
        stage('Install Dependencies') {
            steps {
                sh '''
                    go mod download
                    go mod verify
                    go mod tidy
                '''
            }
        }
        
        stage('Code Quality & Security') {
            parallel {
                stage('Lint') {
                    steps {
                        sh '''
                            go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
                            golangci-lint run --timeout=5m
                        '''
                    }
                }
                
                stage('Security Scan') {
                    steps {
                        sh '''
                            go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
                            gosec ./...
                        '''
                    }
                }
                
                stage('Vulnerability Check') {
                    steps {
                        sh '''
                            go install golang.org/x/vuln/cmd/govulncheck@latest
                            govulncheck ./...
                        '''
                    }
                }
            }
        }
        
        stage('Unit Tests') {
            steps {
                sh '''
                    echo "Running unit tests..."
                    go test -v -race -coverprofile=coverage.out ./...
                    go tool cover -html=coverage.out -o coverage.html
                    go tool cover -func=coverage.out
                '''
            }
            post {
                always {
                    publishHTML([
                        allowMissing: false,
                        alwaysLinkToLastBuild: true,
                        keepAll: true,
                        reportDir: '.',
                        reportFiles: 'coverage.html',
                        reportName: 'Coverage Report'
                    ])
                }
            }
        }
        
        stage('Integration Tests') {
            steps {
                script {
                    try {
                        sh '''
                            echo "Starting integration test environment..."
                            cp .env.example .env.test
                            
                            # Update test environment variables
                            sed -i 's/MYSQL_PORT=3306/MYSQL_PORT=3308/g' .env.test
                            sed -i 's/REDIS_PORT=6379/REDIS_PORT=6380/g' .env.test
                            sed -i 's/API_PORT=8080/API_PORT=8081/g' .env.test
                            
                            # Start test environment
                            docker-compose -f docker-compose.test.yml --env-file .env.test up -d --build
                            
                            # Wait for services to be ready
                            echo "Waiting for services to be ready..."
                            sleep 30
                            
                            # Run integration tests
                            docker-compose -f docker-compose.test.yml exec -T api go test -tags=integration ./tests/integration/...
                        '''
                    } catch (Exception e) {
                        echo "Integration tests failed: ${e.getMessage()}"
                        currentBuild.result = 'UNSTABLE'
                    } finally {
                        sh '''
                            echo "Cleaning up test environment..."
                            docker-compose -f docker-compose.test.yml down -v
                            docker system prune -f
                        '''
                    }
                }
            }
        }
        
        stage('Build Docker Image') {
            when {
                anyOf {
                    branch 'main'
                    branch 'develop'
                    changeRequest()
                }
            }
            steps {
                script {
                    def imageTag = env.BRANCH_NAME == 'main' ? 'latest' : "${env.BRANCH_NAME}-${env.BUILD_NUMBER}"
                    sh '''
                        echo "Building Docker image..."
                        docker build -t ${IMAGE_NAME}:${imageTag} .
                        docker tag ${IMAGE_NAME}:${imageTag} ${DOCKER_REGISTRY}/${IMAGE_NAME}:${imageTag}
                    '''
                    
                    if (env.BRANCH_NAME == 'main') {
                        sh '''
                            echo "Pushing to registry..."
                            docker push ${DOCKER_REGISTRY}/${IMAGE_NAME}:${imageTag}
                        '''
                    }
                }
            }
        }
        
        stage('Deploy to Staging') {
            when {
                branch 'develop'
            }
            steps {
                sh '''
                    echo "Deploying to staging environment..."
                    # Add your staging deployment commands here
                    # docker-compose -f docker-compose.staging.yml up -d
                '''
            }
        }
        
        stage('Deploy to Production') {
            when {
                branch 'main'
            }
            steps {
                input message: 'Deploy to production?', ok: 'Deploy'
                sh '''
                    echo "Deploying to production environment..."
                    # Add your production deployment commands here
                    # docker-compose -f docker-compose.prod.yml up -d
                '''
            }
        }
    }
    
    post {
        always {
            sh '''
                echo "Cleaning up workspace..."
                docker system prune -f
            '''
            
            // Archive artifacts
            archiveArtifacts artifacts: 'coverage.html,coverage.out', allowEmptyArchive: true
            
            // Publish test results if available
            publishTestResults testResultsPattern: 'test-results.xml', allowEmptyResults: true
        }
        
        success {
            echo 'Pipeline completed successfully!'
            // Add notification logic here (Slack, email, etc.)
        }
        
        failure {
            echo 'Pipeline failed!'
            // Add failure notification logic here
        }
        
        unstable {
            echo 'Pipeline completed with warnings!'
            // Add warning notification logic here
        }
    }
}