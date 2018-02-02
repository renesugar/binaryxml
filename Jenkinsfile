pipeline {
    agent { dockerfile true }
    stages {
        stage('Start') {
            steps {
                slackSend message:"${currentBuild.fullDisplayName} started (<${env.BUILD_URL}|Open>)"
            }
        }
        stage('Test') {
            steps {
                sh 'make check'
            }
        }
    }
    post {
        always {
            junit 'target/test-report.xml'
        }
        success {
            slackSend color:'good', message:"${currentBuild.fullDisplayName} completed successfully (<${env.BUILD_URL}|Open>)"
        }
        unstable {
            slackSend color:'warning', message:"${currentBuild.fullDisplayName} has become unstable (<${env.BUILD_URL}|Open>)"
        }
        failure {
            slackSend color:'bad', message:"${currentBuild.fullDisplayName} failed (<${env.BUILD_URL}|Open>)"
        }
    }
}
