resource "firezone_device" "dev1" {
  user_id     = firezone_user.user.id
  name        = "dev1"
  description = "dev1"
  public_key  = random_string.public_key.result
}