pipeline {
    agent any
    tools {
		go 'go v1.18'
	}
	environment { 
		GO111MODULE = 'on'
		GOPATH = "$HOME/go"
		PATH = "$PATH:$GOROOT/bin:$GOPATH/bin"
	}
    stages {
        stage('Download package from FTP server') {
            when {
                branch 'beta-compiler-kalp-296'
            }
            steps {
                sh 'scp -r ubuntu@13.233.229.232:/home/ubuntu/main/lib.linux-* /var/lib/jenkins/workspace/_compiler_beta-compiler-kalp-296/internal/peer/lifecycle/chaincode'
            }
        }
        stage('Run make peer') {
            when {
                branch 'beta-compiler-kalp-296'
            }
            steps {
                sh 'make clean peer'
            }
        }
    
        // stage('Transfer Files to Remote Server') {
        //     when {
        //         branch 'beta-compiler-kalp-296'
        //     }
        //     steps {
        //         sh 'scp -r ./build/lib.linux-* ubuntu@13.233.229.232:/home/ubuntu/main'
        //         sh 'scp -r ./dist/kpc ubuntu@13.233.229.232:/home/ubuntu/main'
        //     }
        // }
       
    }
}

 
