resource "aws_dynamodb_table" "birthday" {
  name           = "birthday"
  billing_mode   = "PROVISIONED"
  read_capacity  = 25
  write_capacity = 25
  hash_key       = "user"

  attribute {
    name = "user"
    type = "S"
  }

  tags = {
    Name        = "birthday"
  }

  point_in_time_recovery {
    enabled = true
  }
}
