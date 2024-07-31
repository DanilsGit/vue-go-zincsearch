terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.60.0"
    }
  }

  required_version = ">= 1.2.0"
}

provider "aws" {
  region  = "us-west-1"
}

data "aws_ami" "ubuntu" {
  most_recent = true

  filter {
    name   = "name"
    values = ["ubuntu/images/hvm-ssd/ubuntu-focal-20.04-amd64-server-*"]
  }

  filter {
    name = "virtualization-type"
    values = ["hvm"]
  }

  owners = ["099720109477"] # Canonical
}

resource "aws_security_group" "app_sg" {
  name        = "app_security_group"
  description = "Security group for application server"

  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

    ingress {
    from_port   = 5173
    to_port     = 5173
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 8080
    to_port     = 8080
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 4080
    to_port     = 4080
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_instance" "app_server" {
  ami           = data.aws_ami.ubuntu.id
  instance_type = "t2.micro"

  vpc_security_group_ids = [aws_security_group.app_sg.id]

 user_data = <<-EOF
              #!/bin/bash
              # Actualizar e instalar dependencias
              sudo apt-get update
              sudo apt-get install -y \
                  apt-transport-https \
                  ca-certificates \
                  curl \
                  gnupg \
                  lsb-release \
                  git

              #Establecer las variables de entorno
              export ADMIN=admin
              export ADMIN_PASS=admin123

              # Instalar Docker
              curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg
              echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
              sudo apt-get update
              sudo apt-get install -y docker-ce docker-ce-cli containerd.io

              # Instalar Docker Compose
              sudo curl -L "https://github.com/docker/compose/releases/download/1.29.2/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
              sudo chmod +x /usr/local/bin/docker-compose

              # Instalar go
              sudo snap install go --classic

              # Añade el usuario ubuntu al grupo go
              sudo usermod -aG go ubuntu

              # Añade el usuario ubuntu al grupo docker
              sudo usermod -aG docker ubuntu

              # Clona el repositorio
              git clone https://github.com/DanilsGit/vue-go-zincsearch.git /home/ubuntu/app

              # Navega al directorio del proyecto
              #cd /home/ubuntu/app

              # Construye las imágenes de Docker
              #sudo docker-compose build

              # Inicia los contenedores en segundo plano
              #sudo docker-compose up -d

              # Da tiempo de espera de 10 segundos para que se inicie la base de datos
              #sleep 10

              # Crea el directorio si no existe y navega a la carpeta indexerDatabase/emails
              sudo mkdir -p /home/ubuntu/app/indexerDatabase/emails
              #cd /home/ubuntu/app/indexerDatabase/emails

              # Navega a la carpeta indexerDatabase y a la subcarpeta emails
              #cd /home/ubuntu/app/indexerDatabase/emails

              # Descarga el archivo extraíble
              #sudo wget http://www.cs.cmu.edu/~enron/enron_mail_20110402.tgz

              # Descomprime el archivo extraíble
              #sudo tar -xvzf enron_mail_20110402.tgz

              # Navega de regreso a la carpeta indexerDatabase
              #cd /home/ubuntu/app/indexerDatabase

              # Ejecuta el programa Go
              #go run main.go
              EOF

  tags = {
    Name = "test-app-server"
  }
}

