IMAGE_NAME = "ubuntu/trusty64"
N = 2

Vagrant.configure("2") do |config|
	config.vm.provider "virtualbox" do |v|
		v.memory = 1024
		v.cpus = 2
	end

	(1..N).each do |i|
		config.vm.define "node-#{i}" do |node|
			node.vm.box = IMAGE_NAME
			node.vm.network "private_network", ip: "192.168.56.#{i + 10}"
			node.vm.hostname = "node-#{i}"
            node.vm.synced_folder "~/Downloads", "/vagrant_data"
		end
	end
end
