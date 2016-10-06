variable region {}
variable release_channel { description = "CoreOS Release Channel" }

provider aws {
    region = "${var.region}"
}

data aws_ami coreos {
    filter {
        name = "virtualization-type"
        values = [ "hvm" ]
    }
    filter {
        name = "state"
        values = [ "available" ]
    }
    filter {
        name = "architecture"
        values = [ "x86_64" ]
    }
    filter {
        name = "image-type"
        values = [ "machine" ]
    }
    filter {
        name = "name"
        values = [ "CoreOS-${var.release_channel}-*" ]
    }

    owners = [ "595879546273" ]
    most_recent = true
}

output id { value = "${data.aws_ami.coreos.id}" }
