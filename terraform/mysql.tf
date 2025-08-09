
resource "aws_rds_cluster" "test_mysql" {
  cluster_identifier = "mysql"
  engine = "aurora-mysql"
  engine_version = "8.0.mysql_aurora.3.08.2"
  engine_mode        = "provisioned"
  database_name = "go_rds_driver_mysql"
  master_username = "root"
  master_password = random_password.aurora_password.result
  enable_http_endpoint = true
  enabled_cloudwatch_logs_exports = []
  iam_roles = []
  skip_final_snapshot = true
  tags = {}
  enable_global_write_forwarding = false

  serverlessv2_scaling_configuration {
    max_capacity = 64
    min_capacity = 1
    seconds_until_auto_pause = 600
  }
}


resource "aws_rds_cluster_instance" "test_mysql" {
  cluster_identifier = aws_rds_cluster.test_mysql.id
  instance_class     = "db.serverless"
  engine             = aws_rds_cluster.test_mysql.engine
  engine_version     = aws_rds_cluster.test_mysql.engine_version
}