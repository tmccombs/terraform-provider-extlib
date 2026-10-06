terraform {
  required_providers {
    ext = {
      source = "registry.opentofu.org/tmccombs/extlib"
    }
  }
}

provider "ext" {}

output "boolnum" {
  value = {
    // 0
    f = provider::ext::boolnum(false)
    // 1
    t = provider::ext::boolnum(true)
  }
}

output "boolset" {
  value = {
    // 0
    f = provider::ext::boolset(false)
    // 1
    t = provider::ext::boolset(true)
  }
}
