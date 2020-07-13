# Terraform Provider DND5e

This repository is a [Terraform](https://www.terraform.io) provider for Dungeons & Dragons 5th Edition. It is a terrible idea.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 0.13.x
- [Go](https://golang.org/doc/install) >= 1.12

## Installation

```hcl
terraform {
  required_providers {
    dnd5e = {
      source = "alisdair/dnd5e"
      version = "0.0.1"
    }
  }
}
```

## Using the provider

```hcl
resource "dnd5e_character" "julia" {
  name              = "Julia Axereaver"
  class             = "druid"
  alignment         = "lawful good"
  experience_points = 14900
  strength          = 13
  dexterity         = 9
  constitution      = 12
  intelligence      = 19
  wisdom            = 17
  charisma          = 10
}

resource "dnd5e_roll" "damage" {
  number = 3
  sides = 6
  modifier = dnd5e_character.julia.strength_modifier
}

output "level" {
  value = dnd5e_character.julia.level
}

output "damage" {
  value = "${join(" + ", dnd5e_roll.damage.values)} + ${dnd5e_roll.damage.modifier} = ${dnd5e_roll.damage.total}"
}
```

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go build`. You can use it locally by copying the provider to your global plugin cache directory:

```shellsession
$ mkdir -p ~/.terraform.d/plugins/registry.terraform.io/alisdair/dnd5e/0.0.1/darwin_amd64
$ cp terraform-provider-dnd5e ~/.terraform.d/plugins/registry.terraform.io/alisdair/dnd5e/0.0.1/darwin_amd64/terraform-provider-dnd5e_v0.0.1
```

Change `darwin_amd64` to the appropriate OS/arch for your system.
