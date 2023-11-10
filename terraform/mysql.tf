
resource "aws_rds_cluster" "test_mysql" {
  cluster_identifier = "mysql"
  engine = "aurora-mysql"
  engine_version = "5.7.mysql_aurora.2.11.3"
  engine_mode = "serverless"
  database_name = "go_rds_driver_mysql"
  master_username = "root"
  master_password = random_password.aurora_password.result
  enable_http_endpoint = true
  enabled_cloudwatch_logs_exports = []
  iam_roles = []
  skip_final_snapshot = true
  tags = {}
  enable_global_write_forwarding = false

  scaling_configuration {
    auto_pause = true
    max_capacity = 64
    min_capacity = 1
    seconds_until_auto_pause = 300
    timeout_action = "RollbackCapacityChange"
  }
}
