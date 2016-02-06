cp dashboard.nginx /etc/nginx/sites-available/
cp dashboard.init /etc/init.d/
ln -s /etc/nginx/sites-available/dashboard.nginx /etc/nginx/sites-enabled/dashboard.nginx
service nginx reload && service nginx restart
cd ../
go build
service dashboard.init start
