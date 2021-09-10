node('PlatformSoftware') {
    isMain = isMainBranch()

    stage('Clean') {
        cleanWs()
        cleanupDockerImages()
    }

    stage('Checkout') {
        checkoutRepo('cbdr', 'hound')
    }
    stage('Build') {
        try {
            withCredentials([
                usernamePassword(credentialsId: 'DockerHub', passwordVariable: 'DHPASSWORD', usernameVariable: 'DHUSERNAME'), 
                sshUserPrivateKey(credentialsId: 'GITHUB_PRIVATE_KEY', keyFileVariable: 'SSH_PRIVATE_KEY')
            ]) {

                sh '''#!/bin/bash -el

                docker login --username $DHUSERNAME --password $DHPASSWORD

                    # We double tag all images with the Jenkins Build ID ($BUILD_DISPLAY_NAME) and also "latest". The
                    # Jenkins Build ID tag will be used across this Jenkins build and will be removed after the build, to
                    # save disk space. The "latest" tag will be kept, so that docker build cache could be used. Untagged
                    # images will also be removed after the Jenkins build.
                    #
                    # We manually trigger the build-env stage and tag it to prevent it from always being pruned. Otherwise
                    # there would always be no build cache and would always need to build from scratch, which would be slow.
                
                
                docker build --label "GIT_COMMIT=$GIT_COMMIT" -t "cbdr/ps-hound:$BUILD_DISPLAY_NAME" -t cbdr/ps-hound:latest --build-arg SSH_PRIVATE_KEY="$(cat $SSH_PRIVATE_KEY)" --force-rm --quiet .

                '''
            }
        } catch (ex) {
            cleanupDockerImages()
            throw ex
        }

    }
    stage('Test') {

    }
    stage('Publish') {
        if (isMain) {
            try {
                withCredentials([usernamePassword(credentialsId: 'DockerHub', passwordVariable: 'DHPASSWORD', usernameVariable: 'DHUSERNAME'), string(credentialsId:'SSH_PRIVATE_KEY', variable:'SSH_PRIVATE_KEY')]) {
                    sh '''#!/bin/bash -el
                    echo "Running docker login"
                    docker login --username $DHUSERNAME --password $DHPASSWORD
                    echo "Pushing images ?"
                    docker push cbdr/ps-hound:$BUILD_DISPLAY_NAME
                    docker push cbdr/ps-hound:latest
                    '''
                }
            } catch (ex) {
                cleanupDockerImages()
                throw ex
            }
        } else {
            println "Not on main, on branch: ${env.BRANCH_NAME}, nothing to publish"
        }
    }

    stage('Post Cleanup') {
        cleanWs()
        cleanupDockerImages()
    }
}

def isMainBranch() {
    return 'main' == env.BRANCH_NAME
}
    
def cleanupDockerImages() {
    try {
        sh '''#!/bin/bash -l
        echo "Cleaning up docker images"
        TAGGED_IMAGES=$(docker images -f "reference=cbdr/ps-hound*:$BUILD_DISPLAY_NAME" | tail -n+2 | awk '{ print $1":"$2 }')
        if [ -n "$TAGGED_IMAGES" ]; then
            docker rmi $TAGGED_IMAGES
        fi
        if (ps ax | grep -q '[d]ocker build'); then
            echo "\\`docker system prune\\` is skipped because a docker build is in progress"
        else
            docker system prune -f
        fi
        '''
    } catch (ex) {
    // error is ignored for cleanup
    }
}