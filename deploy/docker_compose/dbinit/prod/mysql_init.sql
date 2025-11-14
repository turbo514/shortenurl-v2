CREATE DATABASE tenant_prod;
CREATE DATABASE link_prod;

GRANT ALL PRIVILEGES ON tenant_prod.* TO 'user'@'%';
GRANT ALL PRIVILEGES ON link_prod.* TO 'user'@'%';

FLUSH PRIVILEGES;