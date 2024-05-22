data "aws_iam_role" "ec2" {
  name = "pagemail-prod-ec2-role"
}

data "aws_iam_policy" "SESFullAccess" {
  arn = "arn:aws:iam::aws:policy/AmazonSESFullAccess"
}

data "aws_iam_policy_document" "EcrFullAccess" {
  statement {
    sid       = "VisualEditor0"
    effect    = "Allow"
    resources = ["*"]

    actions = [
      "ecr-public:DescribeRegistries",
      "ecr:DescribeImageReplicationStatus",
      "ecr:ListTagsForResource",
      "ecr:ListImages",
      "ecr:BatchGetRepositoryScanningConfiguration",
      "ecr:GetRegistryScanningConfiguration",
      "ecr:DescribeRepositories",
      "ecr:BatchCheckLayerAvailability",
      "ecr:GetLifecyclePolicy",
      "ecr-public:DescribeImageTags",
      "ecr:GetRegistryPolicy",
      "ecr-public:DescribeImages",
      "ecr-public:GetAuthorizationToken",
      "ecr:DescribeImageScanFindings",
      "ecr:GetLifecyclePolicyPreview",
      "ecr:GetDownloadUrlForLayer",
      "ecr:DescribeRegistry",
      "ecr:DescribePullThroughCacheRules",
      "ecr-public:GetRepositoryCatalogData",
      "ecr:GetAuthorizationToken",
      "ecr-public:GetRepositoryPolicy",
      "ecr-public:DescribeRepositories",
      "ecr:BatchGetImage",
      "ecr:DescribeImages",
      "ecr-public:GetRegistryCatalogData",
      "ecr-public:ListTagsForResource",
      "ecr-public:BatchCheckLayerAvailability",
      "ecr:GetRepositoryPolicy",
    ]
  }
}

resource "aws_iam_policy" "ECRFullAccess" {
  name        = "ECRFullAccess"
  description = "Full access to ECR"
  policy      = data.aws_iam_policy_document.EcrFullAccess.json
}

resource "aws_iam_role_policy_attachment" "ECRFullAccess" {
  policy_arn = aws_iam_policy.ECRFullAccess.arn
  role       = data.aws_iam_role.ec2.name
}

resource "aws_iam_role_policy_attachment" "SESFullAccess" {
  policy_arn = data.aws_iam_policy.SESFullAccess.arn
  role       = data.aws_iam_role.ec2.name
}
