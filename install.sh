#!/bin/bash

export base=$1
if [ -z $1 ]; then
	echo "error: must give base directory as first arg"
	exit 1
fi
if [ -z $env ]; then 
	if [ -z $2 ]; then
		echo "error: no environment set"
		exit 1
	else
		export env=$2
	fi
fi

# setup env
if [ $env = "prd" ]; then
	echo "installing prd binary"
	export svc=pagemail
elif [ $env = "stg" ]; then
	echo "installing stg binary"
	export svc=pagemail.staging
else
	echo "unknown environment $env"
	exit 1
fi

# install binary into $env
service $svc stop
cp $base/services/$svc.service /etc/systemd/system/$svc.service
cp $base/build/main /home/ec2/$env/pagemail/main
chmod a+x /home/ec2/$env/pagemail/main
systemctl daemon-reload

# install test sites
rm /var/www/testsites/*
cp $base/test_pages/* /var/www/testsites

# install nginx configs
for f in pm staging_pm test v2 www_pm
do
	rm /etc/nginx/conf.d/$f.conf
	cp $base/nginx/pm.conf /etc/nginx/conf.d/pm.conf
done

# get ssl certificates
certbot certonly --nginx -d "v2.pagemail.io,www.pagemail.io,pagemail.io,staging.pagemail.io" --expand --non-interactive
service nginx restart
