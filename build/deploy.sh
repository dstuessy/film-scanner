ssh systemd-manager@filmscanner.local "systemctl stop filmscanner.service" && \
rsync -razP -e ssh --exclude="*.go" tmp/pi/ danielstuessy@filmscanner.local:/app/ && \
ssh systemd-manager@filmscanner.local "systemctl start filmscanner.service" && \
ssh systemd-manager@filmscanner.local "systemctl status filmscanner.service"
