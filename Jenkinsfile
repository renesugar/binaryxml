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
    }
}
