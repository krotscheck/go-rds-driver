resource "aws_rds_cluster" "test_postgresql" {
  cluster_identifier = "postgresql"
  engine = "aurora-postgresql"
  engine_version = "10.14"
  engine_mode = "serverless"
  database_name = "go_rds_driver_postgresql"
  master_username = "root"
  master_password = random_password.aurora_password.result
  enable_http_endpoint = true
  enabled_cloudwatch_logs_exports = []
  iam_roles = []
  skip_final_snapshot = true
  tags = {}

  scaling_configuration {
    auto_pause = true
    max_capacity = 64
    min_capacity = 2
    seconds_until_auto_pause = 300
    timeout_action = "RollbackCapacityChange"
  }
}
