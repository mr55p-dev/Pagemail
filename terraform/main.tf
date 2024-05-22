terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "5.31.0"
    }
  }
  backend "s3" {
    region         = "eu-west-2"
    bucket         = "168938868801-terraform"
    dynamodb_table = "168938868801-terraform"
    key            = "pagemail.tfstate"
  }
}

provider "aws" {
  region = "eu-west-2"
  default_tags {
    tags = {
      terraform = true
    }
  }
}
