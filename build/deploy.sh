cp build/env.pi tmp/.env
rsync -rav -e ssh --exclude="*.go" tmp/ danielstuessy@filmscanner.local:/app/
