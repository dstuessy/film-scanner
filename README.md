# Open Scanner

## Deployment on Raspberry Pi

Deploy the build artifact once the setup below has been completed.

``` sh
$ ./build/deploy.sh
```

## Raspberry Pi setup

SSH into raspberry pi.

Add the following to `/lib/udev/rules.d/99-systemd.rules`.

```
KERNEL=="video0", SYMLINK="video0", TAG+="systemd"
```

Reload the udev rules without restarting the system

``` sh
$ sudo udevadm control --reload-rules && udevadm trigger
```

Enable the deployed service file on the raspberry pi

``` sh
$ sudo systemctl enable /app/filmscanner.service
```

Follow instructions to start and stop the filmscanner remotely

[https://sleeplessbeastie.eu/2021/03/03/how-to-manage-systemd-services-remotely/](https://sleeplessbeastie.eu/2021/03/03/how-to-manage-systemd-services-remotely/)

**Note:** Make sure the authorized\_keys file does not restric the command.
