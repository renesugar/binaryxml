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
            slackSend color:'good', message:"${currentBuild.fullDisplayName} completed successfully (<${env.BUILD_URL}|Open>)"
        }
        failure {
            slackSend color:'bad', message:"${currentBuild.fullDisplayName} failed (<${env.BUILD_URL}|Open>)"
        }
    }
}
