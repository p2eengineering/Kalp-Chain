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
    
        stage('Transfer Files to Remote Server') {
            when {
                branch 'beta-compiler-kalp-296'
            }
            steps {
                sh 'sshpass -p "1Kamal-POC" ssh -o StrictHostKeyChecking=no ubuntu@161.35.226.109 "cp /home/ubuntu/peer1.example.com/peer /home/ubuntu/backup"'
                sh 'sshpass -p "1Kamal-POC" ssh -o StrictHostKeyChecking=no ubuntu@137.184.187.170 "cp /home/ubuntu/peer2.example.com/peer /home/ubuntu/backup"'
                sh 'sshpass -p "1Kamal-POC" scp -r /var/lib/jenkins/workspace/_compiler_beta-compiler-kalp-296/build/bin/peer ubuntu@161.35.226.109:/home/ubuntu/peer1.example.com'
                sh 'sshpass -p "1Kamal-POC" scp -r /var/lib/jenkins/workspace/_compiler_beta-compiler-kalp-296/build/bin/peer ubuntu@137.184.187.170:/home/ubuntu/peer2.example.com'
            }
        }

        

       
    }
}

 
