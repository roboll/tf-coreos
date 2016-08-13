resource coreos_cloudconfig userdata {
    template = "${file("${path.root}/cloud-config.yaml")}"
    gzip = false
}

output rendered { value = "${coreos_cloudconfig.userdata.rendered}" }
