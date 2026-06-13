# frozen_string_literal: true

Vagrant.configure("2") do |config|
  config.vm.define :k3s do |node|
    node.vm.box       = "generic/ubuntu2204"
    node.vm.hostname  = "k3s-server"

    node.vm.network "private_network",
                    ip: "192.168.56.10"

    node.vm.network "forwarded_port",
                    guest: 6443, host: 6443,
                    hostip: "0.0.0.0",
                    protocol: "tcp"

    node.vm.provider :libvirt do |lv|
      lv.cpus     = 6
      lv.memory   = 16384
      lv.cpu_mode = "host-passthrough"
    end

    node.vm.provision "shell",
                      name: "base-setup",
                      inline: <<~SHELL
                        set -e
                        apt-get update
                        apt-get install -y curl apt-transport-https gnupg jq vim
                        sysctl -w net.ipv4.ip_forward=1
                        sysctl -w net.ipv6.conf.all.forwarding=1
                      SHELL
  end
end