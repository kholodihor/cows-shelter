terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.0"
    }
  }
  required_version = ">= 1.0.0"
}

provider "aws" {
  region = var.aws_region
}

# Create VPC
module "vpc" {
  source  = "terraform-aws-modules/vpc/aws"
  version = "~> 3.0"

  name = "${var.project_name}-vpc"
  cidr = var.vpc_cidr

  azs             = ["${var.aws_region}a", "${var.aws_region}b", "${var.aws_region}c"]
  private_subnets = var.private_subnet_cidrs
  public_subnets  = var.public_subnet_cidrs

  enable_nat_gateway   = true
  enable_dns_hostnames = true
  enable_dns_support   = true

  tags = {
    Environment = var.environment
    Project     = var.project_name
  }
}

# Security Group for MinIO
resource "aws_security_group" "minio" {
  name        = "${var.project_name}-minio-sg"
  description = "Security group for MinIO service"
  vpc_id      = module.vpc.vpc_id

  ingress {
    from_port   = 9000
    to_port     = 9000
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }


  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }


  tags = {
    Name        = "${var.project_name}-minio-sg"
    Environment = var.environment
  }
}

# ECS Cluster
resource "aws_ecs_cluster" "main" {
  name = "${var.project_name}-cluster"

  setting {
    name  = "containerInsights"
    value = "enabled"
  }

  tags = {
    Environment = var.environment
    Project     = var.project_name
  }
}

# ECS Task Definition
resource "aws_ecs_task_definition" "minio" {
  family                   = "${var.project_name}-minio"
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = 1024
  memory                   = 2048
  execution_role_arn       = aws_iam_role.ecs_task_execution_role.arn
  task_role_arn            = aws_iam_role.ecs_task_role.arn

  container_definitions = jsonencode([{
    name      = "minio"
    image     = "minio/minio:latest"
    essential = true
    portMappings = [
      {
        containerPort = 9000
        hostPort      = 9000
        protocol      = "tcp"
      },
      {
        containerPort = 9001
        hostPort      = 9001
        protocol      = "tcp"
      }
    ]
    environment = [
      {
        name  = "MINIO_ROOT_USER"
        value = var.minio_root_user
      },
      {
        name  = "MINIO_ROOT_PASSWORD"
        value = var.minio_root_password
      }
    ]
    command = ["server", "/data", "--console-address", ":9001"]
    mountPoints = [
      {
        containerPath = "/data"
        sourceVolume  = "minio_data"
      }
    ]
    logConfiguration = {
      logDriver = "awslogs"
      options = {
        "awslogs-group"         = "/ecs/${var.project_name}-minio"
        "awslogs-region"        = var.aws_region
        "awslogs-stream-prefix" = "ecs"
      }
    }
  }])

  volume {
    name = "minio_data"

    efs_volume_configuration {
      file_system_id     = aws_efs_file_system.minio_data.id
      transit_encryption = "ENABLED"

      authorization_config {
        access_point_id = aws_efs_access_point.minio_data.id
        iam             = "ENABLED"
      }
    }
  }

  tags = {
    Environment = var.environment
    Project     = var.project_name
  }
}

# ECS Service
resource "aws_ecs_service" "minio" {
  name            = "${var.project_name}-minio-service"
  cluster         = aws_ecs_cluster.main.id
  task_definition = aws_ecs_task_definition.minio.arn
  desired_count   = 1
  launch_type     = "FARGATE"

  network_configuration {
    security_groups  = [aws_security_group.minio.id]
    subnets          = module.vpc.private_subnets
    assign_public_ip = false
  }

  load_balancer {
    target_group_arn = aws_lb_target_group.minio.arn
    container_name   = "minio"
    container_port   = 9000
  }

  depends_on = [
    aws_iam_role_policy_attachment.ecs_task_execution_role,
    aws_lb_listener.minio
  ]

  tags = {
    Environment = var.environment
    Project     = var.project_name
  }
}

# EFS for persistent storage
resource "aws_efs_file_system" "minio_data" {
  creation_token = "${var.project_name}-minio-data"
  encrypted      = true

  tags = {
    Name        = "${var.project_name}-minio-data"
    Environment = var.environment
    Project     = var.project_name
  }
}

resource "aws_efs_mount_target" "minio_data" {
  count           = length(module.vpc.private_subnets)
  file_system_id  = aws_efs_file_system.minio_data.id
  subnet_id       = module.vpc.private_subnets[count.index]
  security_groups = [aws_security_group.efs.id]
}

resource "aws_efs_access_point" "minio_data" {
  file_system_id = aws_efs_file_system.minio_data.id
  posix_user {
    gid = 1000
    uid = 1000
  }
  root_directory {
    path = "/minio"
    creation_info {
      owner_gid   = 1000
      owner_uid   = 1000
      permissions = "755"
    }
  }
}

# Security Group for EFS
resource "aws_security_group" "efs" {
  name        = "${var.project_name}-efs-sg"
  description = "Security group for EFS"
  vpc_id      = module.vpc.vpc_id

  # NFS
  ingress {
    from_port       = 2049
    to_port         = 2049
    protocol        = "tcp"
    security_groups = [aws_security_group.minio.id]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name        = "${var.project_name}-efs-sg"
    Environment = var.environment
  }
}

# Application Load Balancer
resource "aws_lb" "minio" {
  name               = "${var.project_name}-minio-alb"
  internal           = false
  load_balancer_type = "application"
  security_groups    = [aws_security_group.alb.id]
  subnets            = module.vpc.public_subnets

  enable_deletion_protection = false

  tags = {
    Environment = var.environment
    Project     = var.project_name
  }
}

resource "aws_lb_target_group" "minio" {
  name        = "${var.project_name}-minio-tg"
  port        = 9000
  protocol    = "HTTP"
  vpc_id      = module.vpc.vpc_id
  target_type = "ip"

  health_check {
    enabled             = true
    interval            = 30
    path                = "/minio/health/live"
    port                = "traffic-port"
    healthy_threshold   = 3
    unhealthy_threshold = 3
    timeout             = 5
    protocol            = "HTTP"
    matcher             = "200"
  }


  tags = {
    Environment = var.environment
    Project     = var.project_name
  }
}

resource "aws_lb_listener" "minio" {
  load_balancer_arn = aws_lb.minio.arn
  port              = 80
  protocol          = "HTTP"

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.minio.arn
  }
}

# Security Group for ALB
resource "aws_security_group" "alb" {
  name        = "${var.project_name}-alb-sg"
  description = "Security group for ALB"
  vpc_id      = module.vpc.vpc_id

  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name        = "${var.project_name}-alb-sg"
    Environment = var.environment
  }
}

# IAM Role for ECS Task Execution
resource "aws_iam_role" "ecs_task_execution_role" {
  name = "${var.project_name}-ecs-task-execution-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "ecs-tasks.amazonaws.com"
        }
      },
    ]
  })
}

resource "aws_iam_role_policy_attachment" "ecs_task_execution_role" {
  role       = aws_iam_role.ecs_task_execution_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}

# IAM Role for ECS Task
resource "aws_iam_role" "ecs_task_role" {
  name = "${var.project_name}-ecs-task-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "ecs-tasks.amazonaws.com"
        }
      },
    ]
  })
}

# CloudWatch Log Group
resource "aws_cloudwatch_log_group" "minio" {
  name              = "/ecs/${var.project_name}-minio"
  retention_in_days = 30
  tags = {
    Environment = var.environment
    Project     = var.project_name
  }
}

# Outputs
output "minio_alb_dns_name" {
  value = aws_lb.minio.dns_name
}

output "minio_console_url" {
  value = "http://${aws_lb.minio.dns_name}:9001"
}
