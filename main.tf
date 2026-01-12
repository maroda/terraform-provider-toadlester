provider "toadlester" {}

resource "toadlester_widget" "testwidget" {
  name        = "testwidget"
  description = "Widget for internal testing."
}
