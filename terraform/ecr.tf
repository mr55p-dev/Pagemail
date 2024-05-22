resource "aws_ecr_repository" "pagemail" {
  name                 = "pagemail"
  image_tag_mutability = "MUTABLE"

  image_scanning_configuration {
    scan_on_push = true
  }
}

resource "aws_ecr_lifecycle_policy" "pagemail" {
  repository = aws_ecr_repository.pagemail.name
  policy     = <<EOF
{
  "rules": [
    {
      "rulePriority": 1,
      "description": "cleanup untagged images",
      "selection": {
        "tagStatus": "untagged",
        "countType": "sinceImagePushed",
        "countUnit": "days",
        "countNumber": 1
      },
      "action": {
        "type": "expire"
      }
    }
  ]
}
EOF
}

resource "aws_ecr_repository" "pagemail_readability" {
  name                 = "pagemail-readability"
  image_tag_mutability = "MUTABLE"

  image_scanning_configuration {
    scan_on_push = true
  }
}

resource "aws_ecr_lifecycle_policy" "pagemail_readability" {
  repository = aws_ecr_repository.pagemail_readability.name
  policy     = <<EOF
{
  "rules": [
    {
      "rulePriority": 1,
      "description": "cleanup untagged images",
      "selection": {
        "tagStatus": "untagged",
        "countType": "sinceImagePushed",
        "countUnit": "days",
        "countNumber": 1
      },
      "action": {
        "type": "expire"
      }
    }
  ]
}
EOF
}
