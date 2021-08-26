
resource "random_password" "aurora_password" {
  length = 32
  special = false
  upper = true
  lower = true
  number = true
  min_lower = 1
  min_numeric = 1
  min_upper = 1
}

resource "aws_secretsmanager_secret" "aurora_password" {
  name = "aurora_password"
  description = "(Terraform) Database password for Aurora's Data API"
}
resource "aws_secretsmanager_secret_version" "aurora_password" {
  secret_id = aws_secretsmanager_secret.aurora_password.id
  secret_string = jsonencode({
    username: aws_rds_cluster.test_mysql.master_username,
    password: aws_rds_cluster.test_mysql.master_password,
  })
}

output "postgresql_resource_arn" {
  value = aws_rds_cluster.test_postgresql.arn
}
output "mysql_resource_arn" {
  value = aws_rds_cluster.test_mysql.arn
}
output "rds_secret_arn" {
  value = aws_secretsmanager_secret.aurora_password.arn
}
output "mysql_database_name" {
  value = aws_rds_cluster.test_mysql.database_name
}
output "postgresql_database_name" {
  value = aws_rds_cluster.test_postgresql.database_name
}
