resource "firezone_rule" "allow_all" {
  action      = "accept"
  destination = "0.0.0.0/0"
  port_range  = "80 - 443"
  port_type   = "tcp"
}

resource "firezone_rule" "allow_https" {
  action      = "accept"
  destination = "0.0.0.0/0"
  port_range  = "1 - 443"
  port_type   = "tcp"
  user_id     = firezone_user.user.id
}

resource "firezone_rule" "allow_https" {
  action      = "accept"
  destination = "0.0.0.0/0"
  port_range  = "1 - 443"
  port_type   = "tcp"
  user_id     = data.firezone_user.user.id
}