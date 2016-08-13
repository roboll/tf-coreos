# terraform-provider-coreos

[![CircleCI](https://circleci.com/gh/roboll/terraform-provider-coreos.svg?style=svg&circle-token=b7660d6420a79eaae9caf4645c59fb6bfe5b3342)](https://circleci.com/gh/roboll/terraform-provider-coreos)

A terraform provider for coreos cloud config templating and validation. Inspired by the builtin `template_file` resource, the `coreos_cloudconfig` resource accepts a [hil](https://github.com/hashicorp/hil) template and optionally applies validation via [coreos-cloudinit](https://github.com/coreos/coreos-cloudinit/) and gzip+base64 encoding.

## usage

```
resource coreos_cloudconfig userdata {
  gzip = true     # default true
  validate = true # default true

  template = <<TMPL
#cloud-config
hostname: something.example.com
coreos:
  units:

  ...
TMPL
}

output rendered { value = "${coreos_cloudconfig.userdata.rendered}" }
```

## get it

`go get github.com/roboll/terraform-provider-coreos`

## development

[govendor](https://github.com/kardianos/govendor) for vendoring
