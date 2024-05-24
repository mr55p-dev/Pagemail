resource "aws_cloudwatch_log_group" "pagemail_dev" {
  name              = "/splat/pagemail-dev"
  retention_in_days = 1
  tags = {
    Environment = "dev"
    Application = "pagemail"
  }
}
resource "aws_cloudwatch_log_group" "pagemail_prd" {
  name              = "/splat/pagemail"
  retention_in_days = 5
  tags = {
    Environment = "prd"
    Application = "pagemail"
  }
}
resource "aws_cloudwatch_log_group" "pagemail_readability_dev" {
  name              = "/splat/pagemail-readability-dev"
  retention_in_days = 1
  tags = {
    Environment = "dev"
    Application = "pagemail-readability"
  }
}
resource "aws_cloudwatch_log_group" "pagemail_readability_prd" {
  name              = "/splat/pagemail-readability"
  retention_in_days = 5
  tags = {
    Environment = "prd"
    Application = "pagemail-readability"
  }
}
