CREATE DATABASE tenant_dev;
CREATE DATABASE link_dev;

GRANT ALL PRIVILEGES ON tenant_dev.* TO 'user'@'%';
GRANT ALL PRIVILEGES ON link_dev.* TO 'user'@'%';

FLUSH PRIVILEGES;