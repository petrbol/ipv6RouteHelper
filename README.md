## ISC KEA6 IPv6 route helper

This route helper receive information from ISC KEA IPv6 DHCP server using High Availability configuration and insert IPv6 prefixes to the routing table.
Persistent after reboot, default routes are restored with default timer settings. Compiled app included in repo.
Developed for https://github.com/petrbol/RouterConfigurationTool/

#### How to use
```
./rctKea6RouteHelper -h
Usage of ./rctKea6RouteHelper:
  -address string
        default listen ip address (default "127.0.0.1")
  -expiration int
        automatic remove not renewed route after timeout (default 86400)
  -leasesFile string
        where to store backup leases file (default "/etc/rct/savedLeaseRouteHelper.json")
  -port int
        listen TCP port for Kea messages (default 8082)
  -priority int
        default DHCP router helper priority (default 4096)
```
#### Systemd example
```
[Unit]
Description=rctKea6RouteHelper

[Service]
ExecStart=/usr/sbin/rctKea6RouteHelper -leasesFile /etc/savedLeaseRouteHelper.json
Restart=on-abort

[Install]
WantedBy=multi-user.target
```
#### Kea6 HA configuration to send information to the rctKea6RouteHelper
```
    "hooks-libraries": [
      {
        "library": "/usr/lib/x86_64-linux-gnu/kea/hooks/libdhcp_lease_cmds.so",
        "parameters": { }
      },
      {
        "library": "/usr/lib/x86_64-linux-gnu/kea/hooks/libdhcp_ha.so",
        "parameters": {
          "high-availability": [ {
            "this-server-name": "server1",
            "mode": "passive-backup",
            "wait-backup-ack": false,
            "peers": [{
              "name": "server1",
              "url": "http://127.0.0.1:8000/",
              "role": "primary"
            }, {
              "name": "server2",
              "url": "http://127.0.0.1:8082/",
              "role": "backup"
            }]
          } ]
        }
      }
    ],
```
