@Library('whiteblock-dev@go-get-then-go-vet')_

def DEFAULT_BRANCH = 'dev'
def GOLANG_IMAGE = 'golang:1.12'

pipeline {
  options {
    buildDiscarder(logRotator(numToKeepStr: '10'))
  }
  agent any

  stages {
    stage('Static-Analysis') {
      when {
        anyOf {
          changeRequest target: DEFAULT_BRANCH
          changeRequest target: 'master'
        }
      }
      steps {
        goFmt()
        goVet()
      }
    }
  }
}
