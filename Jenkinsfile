node {
  def gitVersion, binaryxmlBuildEnvImage

  try {
    slackSend(message:"${currentBuild.fullDisplayName} started (<${env.BUILD_URL}|Open>)")

    stage('Checkout') {
      deleteDir()
      checkout([
          $class: 'GitSCM',
          branches: scm.branches,
          doGenerateSubmoduleConfigurations: scm.doGenerateSubmoduleConfigurations,
          extensions: [[$class:'CloneOption', noTags:false]],
          userRemoteConfigs: scm.userRemoteConfigs
      ])
      gitVersion = sh(returnStdout:true, script:'./gitVersion.sh -f sem').trim()
      echo "Version tag derived from Git is ${gitVersion}"
    }

    stage('Build') {
      binaryxmlBuildEnvImage = docker.build("binaryxml-buildenv:${gitVersion}")
    }

    stage('Test') {
      sh 'mkdir -p target'
      binaryxmlBuildEnvImage.withRun("-v ${pwd()}/target:/target", "make check TARGET=/target") { c ->
        sh "docker wait ${c.id}"
        sh "docker logs ${c.id}"
        junit 'target/test-report.xml'
      }
    }

    stage('Report') {
      slackSend(color:'good', message:"${currentBuild.fullDisplayName} succeeded (<${env.BUILD_URL}|Open>)")
    }
  } catch (e) {
    slackSend(color:'#ff0000', message:"${currentBuild.fullDisplayName} failed (<${env.BUILD_URL}|Open>)")
    throw e
  }
}
