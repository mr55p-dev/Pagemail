#!/bin/bash
set -e

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

# install service
systemctl stop $svc
cp $base/deploy/services/$svc.service /etc/systemd/system/$svc.service
chmod a+x /home/ec2-user/$env/pagemail/main
systemctl daemon-reload

# install test sites
rm /var/www/testsites/*
cp $base/deploy/test_pages/* /var/www/testsites

# install nginx configs
for f in $(ls $base/deploy/nginx)
do
	cp $base/deploy/nginx/$f /etc/nginx/conf.d/$f
done

# get ssl certificates
certbot --nginx -d "v2.pagemail.io,www.pagemail.io,pagemail.io,staging.pagemail.io" --expand --non-interactive
systemctl restart nginx
systemctl start $svc
