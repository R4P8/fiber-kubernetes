pipeline {
    agent {
        docker {
            image 'rizqirafa8/go-with-docker-v2'
            args '-v /var/run/docker.sock:/var/run/docker.sock'
        }
    }

    environment {
        APP_NAME = 'go-fiber-example-v2'
        APP_VERSION = '1.25.0'
        DOCKER_USERNAME = 'rizqirafa8'
        DOCKER_IMAGE = "${DOCKER_USERNAME}/${APP_NAME}:${APP_VERSION}"

        // Harbor config
        HARBOR_URL = 'http://localhost:8083/'
        HARBOR_PROJECT = 'go-fiber'
        HARBOR_IMAGE = "${HARBOR_URL}/${HARBOR_PROJECT}/${APP_NAME}:${APP_VERSION}"

        // ArgoCD config
        ARGOCD_SERVER = 'argocd-server.argocd.svc.cluster.local:443'
        ARGOCD_APP = 'go-fiber-app'
    }

    stages {
        stage('Checkout') {
            steps {
                checkout scm
            }
        }

    stage('Prepare .env File') {
            steps {
                withCredentials([file(credentialsId: 'env-file', variable: 'ENV_FILE')]) {
                    sh '''
                        echo "🔐 Copying .env from Jenkins Credentials..."
                        cp $ENV_FILE .env
                        echo "✅ .env file ready in workspace"
                    '''
                }
            }
        }

        stage('Go Build & Unit Test') {
            steps {
                sh '''
                    go mod tidy
                    go build -o app
                    go test ./...
                '''
            }
        }

        stage('Build Docker Image') {
            steps {
                sh "docker build -t ${DOCKER_IMAGE} ."
            }
        }

        stage('Security Scan (Trivy)') {
            steps {
                sh '''
                    echo "🔍 Installing Trivy..."
                    apk add --no-cache curl
                    curl -sfL https://raw.githubusercontent.com/aquasecurity/trivy/main/contrib/install.sh | sh

                    echo "🔍 Running Trivy vulnerability scan..."
                    trivy image --exit-code 0 --severity HIGH,CRITICAL ${DOCKER_IMAGE}
                '''
            }
        }

        stage('Push to Harbor Registry') {
            steps {
                withCredentials([usernamePassword(credentialsId: 'harbor_cred',
                                                  usernameVariable: 'HARBOR_USER',
                                                  passwordVariable: 'HARBOR_PASS')]) {
                    sh '''
                        echo "$HARBOR_PASS" | docker login ${HARBOR_URL} -u "$HARBOR_USER" --password-stdin
                        docker tag ${DOCKER_IMAGE} ${HARBOR_IMAGE}
                        docker push ${HARBOR_IMAGE}
                        docker logout ${HARBOR_URL}
                    '''
                }
            }
        }

        stage('Trigger ArgoCD Sync') {
            steps {
                withCredentials([usernamePassword(credentialsId: 'argocd_cred',
                                                  usernameVariable: 'ARGO_USER',
                                                  passwordVariable: 'ARGO_PASS')]) {
                    sh '''
                        echo "🚀 Triggering ArgoCD Sync for ${ARGOCD_APP}..."
                        argocd login ${ARGOCD_SERVER} --username "$ARGO_USER" --password "$ARGO_PASS" --insecure
                        argocd app sync ${ARGOCD_APP}
                        argocd app wait ${ARGOCD_APP} --sync --health
                    '''
                }
            }
        }
    }

    post {
        success {
            echo "✅ Pipeline completed successfully!"
        }
        failure {
            echo "❌ Pipeline failed. Please check the logs."
        }
    }
}
