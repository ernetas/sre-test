terraform {
  backend "s3" {
    bucket = "revolut-terraform"
    key    = "eu-central-1-revolut/terraform.tfstate"
    region = "eu-central-1"
  }
}

provider "aws" {
  region = "eu-central-1"
}
