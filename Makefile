.PHONY = halt clean install start 

halt:
	sudo service pagemail stop

clean:
	sudo rom -f /home/ec2-user/server
	sudo rm -rf /var/www/pagemail/* /home/ec2-user/dist

install:
	sudo cp -r /home/ec2-user/dist/* /var/www/pagemail/
	sudo chmod a+x /home/ec2-user/server

start:
	sudo service pagemail start
	sudo service nginx restart
