def DEFAULT_BRANCH = 'dev'
def GOLANG_IMAGE = 'golang:1.12'

pipeline {
  options {
    buildDiscarder(logRotator(numToKeepStr: '10'))
  }
  agent {
    kubernetes {
      cloud 'kubernetes-dev-gke'
      yaml """
  apiVersion: v1
  kind: Pod
  metadata:
  labels:
    cicd: true
  spec:
    containers:
    - name: golang
      image: ${GOLANG_IMAGE}
      command:
      - cat
      tty: true
  """
    }
  }

  stages {
    stage('Static-Analysis') {
      when {
        anyOf {
          changeRequest target: DEFAULT_BRANCH
          changeRequest target: 'master'
        }
      }
      steps {
        container("golang") {
          sh "find . -name '*.go' | xargs gofmt -w"
        }
        script {
          dirtyRepo = false
          try {
            sh 'git diff-index --quiet HEAD'
          } catch(Exception e) {
            dirtyRepo = true
          }

          if (dirtyRepo) {
            /* Need to use a personal-access-token and not an
               ssh deploy key with write access because
               the repo url uses https protocol not ssh
            */
            withCredentials([
              usernamePassword(
                credentialsId: 'github-repo-pac',
                passwordVariable: 'GIT_PASSWORD',
                usernameVariable: 'GIT_USERNAME'
              )
            ]) {
              sh """#!/bin/bash
              set -e
              set -x

              git config user.name ${env.CHANGE_AUTHOR}
              git config user.email ${env.CHANGE_AUTHOR_EMAIL}

              # Jenkins checks out workspace in detached HEAD state
              # so need to get on a branch to commit.
              # Also, if this is a PR, we can't use env.BRANCH_NAME
              # because it will be "PR-*"
              if [[ ${env.BRANCH_NAME} == PR-* ]]; then
                BRANCH=${env.CHANGE_BRANCH}
              else
                BRANCH=${env.BRANCH_NAME}
              fi

              git checkout -b \$BRANCH
              git commit -am 'gofmt'

              # used to read creds from environment during git push
              echo '#!/bin/bash' > ./credential-helper.sh
              echo 'echo username=$GIT_USERNAME' >> ./credential-helper.sh
              echo 'echo password=$GIT_PASSWORD' >> ./credential-helper.sh
              chmod 0755 ./credential-helper.sh

              git config credential.helper "/bin/bash ${env.WORKSPACE}/credential-helper.sh"
              git push ${env.GIT_URL} \$BRANCH
              """
              // return early if a gofmt commit was made
              // and let the build rerun itself
              error("Aborting the build early. Rerunning with gofmt'ed files")
            }
          }
        }
        container("golang") {
          sh "go vet ."
        }
      }
    }
  }
}
