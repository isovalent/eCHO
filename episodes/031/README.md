# Episode 31: Openshift Test Environment with Cilium
## Headlines

* [Cilium Service Mesh beta program](https://cilium.io/blog/2021/12/01/cilium-service-mesh-beta)
* [How Cloudflare uses eBPF to Build Programmable Packet Filtering in Magic Firewall](https://blog.cloudflare.com/programmable-packet-filtering-with-magic-firewall/)
* [Thomas Graf on Service Mesh and ebpf](https://thenewstack.io/how-ebpf-streamlines-the-service-mesh/)

## Cilium v1.11 released!
[release notes](https://github.com/cilium/cilium/releases)
[release blog](https://isovalent.com/blog/post/2021-12-release-111)

#### Nested Virt in GCP:
```
gcloud compute instances create $(whoami)-vagrant-dev   --enable-nested-virtualization   --zone=us-west1-a   --min-cpu-platform="Intel Haswell" --labels=owner=$(whoami) --machine-type=n1-standard-64 --image-project=ubuntu-os-cloud --image-family=ubuntu-2004-lts --boot-disk-size=1TB
gcloud compute config-ssh 
```

#### virsh network:
```
<network xmlns:dnsmasq='http://libvirt.org/schemas/network/dnsmasq/1.0' connections='3' ipv6='yes'>
  <name>ocp1</name>
  <uuid>b470337a-3fd9-4b8f-9eaa-c0d64c74cc3b</uuid>
  <forward mode='nat'>
    <nat>
      <port start='1024' end='65535'/>
    </nat>
  </forward>
  <bridge name='br-ocp1' stp='on' delay='0'/>
  <mac address='52:54:00:19:e0:d4'/>
  <domain name='ocp1.k8s.work'/>
  <dns>
    <host ip='192.168.200.1'>
      <hostname>api-int.ocp1.k8s.work</hostname>
      <hostname>api.ocp1.k8s.work</hostname>
    </host>
  </dns>
  <ip address='192.168.200.1' netmask='255.255.255.0' localPtr='yes'>
    <dhcp>
      <range start='192.168.200.2' end='192.168.200.9'/>
      <host mac='52:54:00:00:01:01' name='bootstrap01.ocp1.k8s.work' ip='192.168.200.11'/>
      <host mac='52:54:00:00:02:01' name='cp01.ocp1.k8s.work' ip='192.168.200.21'/>
      <host mac='52:54:00:00:02:02' name='cp02.ocp1.k8s.work' ip='192.168.200.22'/>
      <host mac='52:54:00:00:02:03' name='cp03.ocp1.k8s.work' ip='192.168.200.23'/>
      <bootp file='pxelinux.0' server='192.168.200.1'/>
    </dhcp>
  </ip>
  <dnsmasq:options>
    <dnsmasq:option value='cname=*.apps.ocp1.k8s.work,apps.ocp1.k8s.work,api.ocp1.k8s.work'/>
    <dnsmasq:option value='auth-zone=ocp1.k8s.work'/>
    <dnsmasq:option value='auth-server=ocp1.k8s.work,*'/>
  </dnsmasq:options>
</network>
```

#### Coredns!
```
.:53 {
    bind 127.0.0.1
    forward . /run/systemd/resolve/resolv.conf
    reload
    log
    errors
}
# BEGIN ocp1 dns managed by ansible 

ocp1.k8s.work {
    bind 127.0.0.1
    forward . 192.168.200.1
    log
    errors
}

# END ocp1 dns managed by ansible 
# BEGIN ocp2 dns managed by ansible 

ocp2.k8s.work {
    bind 127.0.0.1
    forward . 192.168.210.1
    log
    errors
}

# END ocp2 dns managed by ansible 
```
#### haproxy config:
```
#---------------------------------------------------------------------
# Global settings
#---------------------------------------------------------------------
global
    # to have these messages end up in /var/log/haproxy.log you will
    # need to:
    #
    # 1) configure syslog to accept network log events.  This is done
    #    by adding the '-r' option to the SYSLOGD_OPTIONS in
    #    /etc/sysconfig/syslog
    #
    # 2) configure local2 events to go to the /var/log/haproxy.log
    #   file. A line like the following can be added to
    #   /etc/sysconfig/syslog
    #
    #    local2.*                       /var/log/haproxy.log
    #
    log         127.0.0.1 local2

    chroot      /var/lib/haproxy
    pidfile     /var/run/haproxy.pid
    maxconn     4000
    user        haproxy
    group       haproxy
    daemon

    # turn on stats unix socket
    stats socket /var/lib/haproxy/stats

#---------------------------------------------------------------------
# common defaults that all the 'listen' and 'backend' sections will
# use if not designated in their block
#---------------------------------------------------------------------
defaults
    mode                    tcp
    log                     global
    option                  httplog
    option                  dontlognull
    option http-server-close
    option forwardfor       except 127.0.0.0/8
    option                  redispatch
    retries                 3
    timeout http-request    10s
    timeout queue           1m
    timeout connect         10s
    timeout client          4h
    timeout server          4h
    timeout http-keep-alive 10s
    timeout check           10s
    maxconn                 3000

#---------------------------------------------------------------------

listen stats
    bind :9000
    mode http
    stats enable
    stats uri /
    monitor-uri /healthz

frontend openshift_apiserver
    bind *:6443
    mode tcp
    tcp-request inspect-delay 5s
    tcp-request content accept if { req_ssl_hello_type 1 }
    acl ocp1_apiserver_tls req_ssl_sni -m end .ocp1.k8s.work
    use_backend ocp1-apiserver-backend if ocp1_apiserver_tls
    acl ocp2_apiserver_tls req_ssl_sni -m end .ocp2.k8s.work
    use_backend ocp2-apiserver-backend if ocp2_apiserver_tls
    option tcplog

frontend machine-config-server
    bind *:22623
    mode tcp
    tcp-request inspect-delay 5s
    tcp-request content accept if { req_ssl_hello_type 1 }
    acl ocp1_machine_config_tls req_ssl_sni -m end .ocp1.k8s.work
    use_backend ocp1-machine-config-server if ocp1_machine_config_tls
    acl ocp2_machine_config_tls req_ssl_sni -m end .ocp2.k8s.work
    use_backend ocp2-machine-config-server if ocp2_machine_config_tls
    option tcplog

frontend ingress-http
    bind *:80
    acl ocp1_ingress_http hdr_end(Host) -i .apps.ocp1.k8s.work
    use_backend ocp1-ingress-http
    acl ocp2_ingress_http hdr_end(Host) -i .apps.ocp2.k8s.work
    use_backend ocp2-ingress-http
    option tcplog

frontend ingress-https
    bind *:443
    mode tcp
    tcp-request inspect-delay 5s
    tcp-request content accept if { req_ssl_hello_type 1 }
    acl ocp1_ingress_tls req_ssl_sni -m end .apps.ocp1.k8s.work
    use_backend ocp1-ingress-https if ocp1_ingress_tls
    acl ocp2_ingress_tls req_ssl_sni -m end .apps.ocp2.k8s.work
    use_backend ocp2-ingress-https if ocp2_ingress_tls
    option tcplog

backend ocp1-apiserver-backend
    balance source
    mode tcp
    option ssl-hello-chk
    option httpclose
    server bootstrap01.ocp1.k8s.work 192.168.200.11:6443 check inter 2000 rise 2 fall 5
    server cp01.ocp1.k8s.work 192.168.200.21:6443 check inter 2000 rise 2 fall 5
    server cp02.ocp1.k8s.work 192.168.200.22:6443 check inter 2000 rise 2 fall 5
    server cp03.ocp1.k8s.work 192.168.200.23:6443 check inter 2000 rise 2 fall 5

backend ocp2-apiserver-backend
    balance source
    mode tcp
    option ssl-hello-chk
    option httpclose
    server bootstrap01.ocp2.k8s.work 192.168.210.11:6443 check inter 2000 rise 2 fall 5
    server cp01.ocp2.k8s.work 192.168.210.21:6443 check inter 2000 rise 2 fall 5
    server cp02.ocp2.k8s.work 192.168.210.22:6443 check inter 2000 rise 2 fall 5
    server cp03.ocp2.k8s.work 192.168.210.23:6443 check inter 2000 rise 2 fall 5

backend ocp1-machine-config-server
    balance source
    mode tcp
    option ssl-hello-chk
    option httpclose
    server bootstrap01.ocp1.k8s.work 192.168.200.11:22623 check inter 2000 rise 2 fall 5
    server cp01.ocp1.k8s.work 192.168.200.21:22623 check inter 2000 rise 2 fall 5
    server cp02.ocp1.k8s.work 192.168.200.22:22623 check inter 2000 rise 2 fall 5
    server cp03.ocp1.k8s.work 192.168.200.23:22623 check inter 2000 rise 2 fall 5

backend ocp2-machine-config-server
    balance source
    mode tcp
    option ssl-hello-chk
    option httpclose
    server bootstrap01.ocp2.k8s.work 192.168.210.11:22623 check inter 2000 rise 2 fall 5
    server cp01.ocp2.k8s.work 192.168.210.21:22623 check inter 2000 rise 2 fall 5
    server cp02.ocp2.k8s.work 192.168.210.22:22623 check inter 2000 rise 2 fall 5
    server cp03.ocp2.k8s.work 192.168.210.23:22623 check inter 2000 rise 2 fall 5

backend ocp1-ingress-http
    balance source
    mode tcp
    server cp01.ocp1.k8s.work 192.168.200.21:80 check inter 2000 rise 2 fall 5
    server cp02.ocp1.k8s.work 192.168.200.22:80 check inter 2000 rise 2 fall 5
    server cp03.ocp1.k8s.work 192.168.200.23:80 check inter 2000 rise 2 fall 5

backend ocp2-ingress-http
    balance source
    mode tcp
    server cp01.ocp2.k8s.work 192.168.210.21:80 check inter 2000 rise 2 fall 5
    server cp02.ocp2.k8s.work 192.168.210.22:80 check inter 2000 rise 2 fall 5
    server cp03.ocp2.k8s.work 192.168.210.23:80 check inter 2000 rise 2 fall 5

backend ocp1-ingress-https
    balance source
    mode tcp
    option ssl-hello-chk
    option httpclose
    server cp01.ocp1.k8s.work 192.168.200.21:443 check inter 2000 rise 2 fall 5
    server cp02.ocp1.k8s.work 192.168.200.22:443 check inter 2000 rise 2 fall 5
    server cp03.ocp1.k8s.work 192.168.200.23:443 check inter 2000 rise 2 fall 5

backend ocp2-ingress-https
    balance source
    mode tcp
    option ssl-hello-chk
    option httpclose
    server cp01.ocp2.k8s.work 192.168.210.21:443 check inter 2000 rise 2 fall 5
    server cp02.ocp2.k8s.work 192.168.210.22:443 check inter 2000 rise 2 fall 5
    server cp03.ocp2.k8s.work 192.168.210.23:443 check inter 2000 rise 2 fall 5

#---------------------------------------------------------------------

```

Vagrantfile:
```
# Vagrantfile API/syntax version. Don't touch unless you know what you're doing!
VAGRANTFILE_API_VERSION = "2"

# Require a recent version of vagrant otherwise some have reported errors 
# setting host names on boxes
Vagrant.require_version ">= 1.7.2"

CONFIG = {
  "domain" =>               ENV['CLUSTER_DOMAIN'] || "ocp1.k8s.work",
  "network_name" =>         ENV['CLUSTER_NETWORK'] || "ocp1",
  "network_cidr" =>         ENV['CLUSTER_CIDR'] || "192.168.200.0/24",
  "domain_mac_seed" =>      ENV['DOMAIN_MAC_SEED'] || "52:54:00:00",
  "num_bootstrap_nodes" =>  ENV['NUM_BOOTSTRAP_NODES'] || "1",
  "bootstrap_cores" =>      ENV['BOOTSTRAP_CORES'] || "2",
  "bootstrap_memory" =>     ENV['BOOTSTRAP_MEMORY'] || "4096",
  "bootstrap_disk" =>       ENV['BOOTSTRAP_DISK'] || "40G", 
  "bootstrap_mac" =>        ENV['BOOTSTRAP_MAC'] || ":01",
  "num_cp_nodes" =>         ENV['NUM_CP'] || "3",
  "cp_cores" =>             ENV['CP_CORES'] || "8",
  "cp_memory"  =>           ENV['CP_MEMORY'] || "16384",
  "cp_disk" =>              ENV['CP_DISK'] || "100G", 
  "cp_mac" =>               ENV['CP_MAC'] || ":02",
  "num_worker_nodes" =>     ENV['NUM_WORKER'] || "0",
  "worker_cores" =>         ENV['WORKER_CORES'] || "4",
  "worker_memory"  =>       ENV['WORKER_MEMORY'] || "8192",
  "worker_disk" =>          ENV['WORKER_DISK'] || "100G", 
  "worker_mac" =>           ENV['WORKER_MAC'] || ":03",
}


Vagrant.configure(VAGRANTFILE_API_VERSION) do |config|
  config.ssh.insert_key = false

  BS = (CONFIG['num_bootstrap_nodes']).to_i
  (1..BS).each do |i|
    vm_name = "bootstrap0#{i}"
    vm_cpu = CONFIG['bootstrap_cores']
    vm_memory = CONFIG['bootstrap_memory']
    vm_disk = CONFIG['bootstrap_disk']
    vm_mac = "#{CONFIG['domain_mac_seed']}#{CONFIG['bootstrap_mac']}:0#{i}"
    config.vm.define vm_name do |node|
      node.vm.hostname = "#{vm_name}.#{CONFIG['domain']}"
      node.ssh.username = "core"
      node.vm.network :private_network,
        :type => "dhcp",
        :libvirt__network_name => "#{CONFIG['network_name']}",
        :libvirt__always_destroy => false,
        :autostart => false,
        :mac => vm_mac
      node.vm.provider :libvirt do |domain|
        domain.uri = 'qemu+unix:///system'
        domain.driver = 'kvm'
        domain.memory = "#{vm_memory}".to_i
        domain.cpus = "#{vm_cpu}".to_i
        domain.storage :file, :size => "#{vm_disk}", :type => 'qcow2'
        boot_network = {'network' => CONFIG['network_name']}
        domain.boot 'hd'
        domain.boot boot_network
        domain.mgmt_attach = false
      end
    end
  end

  CP = (CONFIG['num_cp_nodes']).to_i
  (1..CP).each do |i|
    vm_name = "cp0#{i}"
    vm_cpu = CONFIG['cp_cores']
    vm_memory = CONFIG['cp_memory']
    vm_disk = CONFIG['cp_disk']
    vm_mac = "#{CONFIG['domain_mac_seed']}#{CONFIG['cp_mac']}:0#{i}"
    config.vm.define vm_name do |node|
      node.vm.hostname = "#{vm_name}.#{CONFIG['domain']}"
      node.ssh.username = "core"
      node.vm.network :private_network,
        :type => "dhcp",
        :libvirt__network_name => "#{CONFIG['network_name']}",
        :libvirt__always_destroy => false,
        :autostart => false,
        :mac => vm_mac
      node.vm.provider :libvirt do |domain|
        domain.uri = 'qemu+unix:///system'
        domain.driver = 'kvm'
        domain.memory = "#{vm_memory}".to_i
        domain.cpus = "#{vm_cpu}".to_i
        domain.storage :file, :size => "#{vm_disk}", :type => 'qcow2'
        boot_network = {'network' => CONFIG['network_name']}
        domain.boot 'hd'
        domain.boot boot_network
        domain.mgmt_attach = false
      end
    end
  end


  WORKER = (CONFIG['num_worker_nodes']).to_i
  (1..WORKER).each do |i|
    vm_name = "worker0#{i}"
    vm_cpu = CONFIG['worker_cores']
    vm_memory = CONFIG['worker_memory']
    vm_disk = CONFIG['worker_disk']
    vm_mac = "#{CONFIG['domain_mac_seed']}#{CONFIG['worker_mac']}:0#{i}"
    config.vm.define vm_name, autostart: false  do |node|
      node.vm.hostname = "#{vm_name}.#{CONFIG['domain']}"
      node.ssh.username = "core"
      node.vm.network :private_network,
        :type => "dhcp",
        :libvirt__network_name => "#{CONFIG['network_name']}",
        :libvirt__always_destroy => false,
        :autostart => false,
        :mac => vm_mac
      node.vm.provider :libvirt do |domain|
        domain.management_network_mac = vm_mac
        domain.management_network_name = CONFIG['network_name']
        domain.management_network_address = CONFIG['network_cidr']
        domain.uri = 'qemu+unix:///system'
        domain.driver = 'kvm'
        domain.memory = "#{vm_memory}".to_i
        domain.cpus = "#{vm_cpu}".to_i
        domain.storage :file, :size => "#{vm_disk}", :type => 'qcow2'
        boot_network = {'network' => CONFIG['network_name']}
        domain.boot 'hd'
        domain.boot boot_network
        domain.mgmt_attach = false
      end
    end
  end
end

```

Join [Cilium Slack](http://slack.cilium.io) or follow [@ciliumproject](https://twitter.com/ciliumproject) for notifications!