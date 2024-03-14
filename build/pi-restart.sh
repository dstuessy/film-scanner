ssh systemd-manager@filmscanner.local "systemctl restart filmscanner.service" && \
ssh systemd-manager@filmscanner.local "systemctl status filmscanner.service"
