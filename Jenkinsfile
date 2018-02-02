pipeline {
    agent { dockerfile true }
    stages {
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
            slackSend color:'good', message:"Job ${currentBuild.fullDisplayName} completed successfully"
        }
        failure {
            slackSend color:'bad', message:"Job ${currentBuild.fullDisplayName} failed"
        }
    }
}
