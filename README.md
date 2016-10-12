# tf-coreos [![CircleCI](https://circleci.com/gh/roboll/tf-coreos.svg?style=svg)](https://circleci.com/gh/roboll/tf-coreos)

## Module

A terraform module that uses the AWS AMI data source to determine the correct AMI for a given CoreOS release channel.

## Plugin

A terraform provider for coreos cloud config templating and validation. Inspired by the builtin `template_file` resource, the `coreos_cloudconfig` resource accepts a [hil](https://github.com/hashicorp/hil) template and optionally applies validation via [coreos-cloudinit](https://github.com/coreos/coreos-cloudinit/) and gzip+base64 encoding.

### usage

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

### get it

`go get github.com/roboll/terraform-coreos/plugins/...`

_or_

`curl -L -o /usr/local/bin/terraform-provider-coreos https://github.com/roboll/terraform-coreos/releases/download/{VERSION}/terraform-provider-coreos_{OS}_{ARCH}`

## development

[govendor](https://github.com/kardianos/govendor) for vendoring
