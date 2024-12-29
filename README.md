# switchbot-contact-sensor_exporter

Exporter designed to collect and build metrics for `SwitchBot Contact Sensor`

To start collecting metrics you need to know MAC address of your sensor device.
You can use scanner application to scan for available sensor devices:

```
go build .
sudo ./switchbot-contact-sensor_exporter scanner
```
```
Scanning for 5s...
Done
Found SwitchBot devices:
MAC: f7:xx:xx:xx:xx:xx, ServiceData: 64006481ffff000588, isOpen: 0
```

Having found MAC addreses you can configure exporter service. Example config is presented in `config_example.json` file.

Start expoter:

```
sudo ./switchbot-contact-sensor_exporter exporter
```

Options:

- `-c`
  - Path to config file.
  - default `/etc/switchbot-contact-sensor_exporter/config.json`
- `--httpListenAddress`
  - Address to bind to.
  - default: `0.0.0.0:9353`
- `--btScanDuration`
  - Duration in seconds for which exporter listens to sensor data
  - default: 2s
- `--btScanInterval`
  - How often should exporter run sensor data listener
  - default :500ms

Exporter sample

```
# HELP battery Battery level
# TYPE battery gauge
battery{location="kitchen",mac="f7:76:a8:xx:xx:xx"} 100
# HELP leave_open Leave open
# TYPE leave_open gauge
leave_open{location="kitchen",mac="f7:76:a8:xx:xx:xx"} 0
# HELP open Open
# TYPE open gauge
open{location="kitchen",mac="f7:76:a8:xx:xx:xx"} 0
# HELP time Time
# TYPE time gauge
time{location="kitchen",mac="f7:76:a8:xx:xx:xx"} 11
```