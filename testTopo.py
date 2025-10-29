#!/usr/bin/python

"""
Simple Mininet script to create a topology with:
- 1 switch (s1)
- 3 hosts (h1, h2, h3)
- Opens xterm for each host

Usage:
    sudo python mininet_single3.py
"""

from mininet.net import Mininet
from mininet.topo import SingleSwitchTopo
from mininet.node import OVSBridge
from mininet.cli import CLI
from mininet.log import setLogLevel, info
from mininet.term import makeTerm

import os
import shutil

NODES = 3

def singleSwitchTopo():
    """Create a single switch topology with 3 hosts using built-in topology"""
    
    # Create topology with single switch and 3 hosts
    topo = SingleSwitchTopo(k=NODES)
    
    # Create network with OVSBridge switch (no controller needed)
    net = Mininet(topo=topo, switch=OVSBridge, controller=None)
    
    info('*** Starting network\n')
    net.start()
    
    info('*** Opening xterm for each host\n')
    for i in range(NODES):
        copy_and_replace("./torrent",f"./h{i+1}/torrent")
    for i,host in enumerate(net.hosts):
        dir_name = host.name
        if not os.path.exists(dir_name):
            os.makedirs(dir_name)
            info(f'Created directory: {dir_name}\n')
        
        # Open xterm in the host's directory
        makeTerm(host, title=f'Host {host.name}', cmd=f'cd {dir_name}; bash')
    
    info('*** Running CLI\n')
    CLI(net)
    
    info('*** Stopping network\n')
    net.stop()

def copy_and_replace(source_path, destination_path):
    if os.path.exists(destination_path):
        os.remove(destination_path)
    shutil.copy2(source_path, destination_path)


if __name__ == '__main__':
    
    setLogLevel('info')
    singleSwitchTopo()

