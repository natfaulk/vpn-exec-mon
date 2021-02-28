# Go Program Manager

I needed a program to run behind a VPN on a linux machine. However, the VPN would intermittently go down and the program stop. It is a pain to constantly check if the VPN is still up and then rerun the program.

This fairly simple program restarts the VPN when necessary, as well as starting the executable whenever it stops.

## Checking the VPN connectivity
If the VPN is connected, local network addresses should be unreachable.  
In my particular circumstances, I know that 1.1.1.1 is unreachable without a VPN - the powers that be have decided to block external DNS quries for some reason, therefore it is a good test to check for general internet connectivity.

There are probably many edge cases that are not covered, but it works pretty well for my purposes. It could probably have easily been implemented using a simple bash / python script, but Go is pretty pleasant to write.

## Usage
Only tested on Linux, possibly works on OSX, would need significant rewrites to work on Windows.  

In the scripts directory, two scripts are needed: `vpnstart.sh` and `vpnend.sh`. These start and stop the VPN.

Add a .env file with the following parameters:  
`LOCAL_ADDR=xxx.xxx.xxx.xxx` This is a local network address that should only be accessible when VPN disconnected.

`REMOTE_ADDR=xxx.xxx.xxx.xxx` This is a remote address that should only be accessible when the VPN is connected. Either an address that is blocked by the network normally, or an address on the remote network. i.e wireguard often uses `10.7.0.0/24` for the virtual network, with `10.7.0.1` being the remote host.

Finally run `./vpn-exec-mon myprogram arg1 arg2 etc...` where the arguements for vpn-exec-mon are the program name followed by the arguements. One *may* need to run this as root (i.e. sudo ./vpn-exec-mon) depending on whether the vpn scripts require root access - which is not ideal but should be possible to find ways around this.

