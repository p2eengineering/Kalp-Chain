pipeline {
    agent any
    stages {
        stage('Download package from FTP server') {
            when {
                branch 'beta-compiler-kalp-296'
            }
            steps {
                sh 'scp -r ubuntu@13.233.229.232:/home/ubuntu/main/lib.linux-* /internal/peer/lifecycle/chaincode'
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

 
