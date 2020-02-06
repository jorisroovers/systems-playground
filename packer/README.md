# Packer

Packer allows you to easily build machine images (like qcow, img, vmdk, ovf, etc).

It does this by booting a machine, simulating putting an ISO image in it and then performing steps that you specify (so-called "provisioners") in shell, ansible, or others. At the end it takes a disk snapshot of that machine which you can use to boot new machines of this type. IOW, packer allows you to easily bake images.

Packer supports multiple machine backends to do this: AWS, VBox, GCP, Openstack and many more, as well as multiple provisioners (shell, ansible, salt, etc). The real value is that packer takes care of the machine provisioning and orchestrating your custom provisioners. This makes it easy to build the same images for multiple target platforms.

NOTE: packer does NOT allow you to build/customize ISO images. ISO images are optical disks and don't represent machine images. ISOs are used to run installers that copy files from the optical medium to physical hard disk.

## Hands-on
```sh
brew install pack
packer build example-centos.json
```

Installation heavily based of:
https://github.com/geerlingguy/packer-boxes/blob/master/centos7/box-config.json

The whole installation can take a while, especially the 'Performing post-installation setup tasks' step within the CentOS VM can take several minutes.

After the process is done, a vmdk file will be generated in `output-virtualbox-iso`. Tried it, it actually boots and works correctly in VBox :-)

Interesting call-out from `example-centos.json`:

```json
"boot_command": [
    "<tab> text ks=http://{{ .HTTPIP }}:{{ .HTTPPort }}/centos7-ks.cfg<enter><wait>"
],
```
Boot command that will be entered. Basically this tells the boot loader to load the kickstarter file from an http server. This http server is hosted by packer and will host all files in the directory specified by the `http_directory` configuration property in `example-centos.json`.
