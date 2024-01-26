-- Create a trigger function to update the updated column
CREATE OR REPLACE FUNCTION update_updated()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;


-- Create a macro to apply the trigger to a table
CREATE OR REPLACE FUNCTION apply_update_trigger(table_name TEXT)
RETURNS VOID AS $$
BEGIN
 IF NOT EXISTS (
    SELECT 1
    FROM information_schema.triggers
    WHERE trigger_schema = 'public'
      AND trigger_name = format('trigger_update_updated_%I', table_name)
  ) THEN
    EXECUTE format('
        CREATE TRIGGER trigger_update_updated_%I
        BEFORE UPDATE ON %I
        FOR EACH ROW
        EXECUTE FUNCTION update_updated()
    ', table_name, table_name);
  END IF;
END;
$$ LANGUAGE plpgsql;



CREATE TABLE IF NOT EXISTS users(
  id TEXT NOT NULL UNIQUE,
  username VARCHAR(15) NOT NULL UNIQUE,
  email VARCHAR(30) NOT NULL UNIQUE,
  password TEXT NOT NULL,
  created TIMESTAMP NOT NULL DEFAULT NOW(),
  updated TIMESTAMP NOT NULL DEFAULT NOW(),
  PRIMARY KEY(id)
);

SELECT apply_update_trigger('users');

CREATE TABLE IF NOT EXISTS roles(
  id SERIAL NOT NULL UNIQUE,
  role VARCHAR(15) NOT NULL UNIQUE,
  PRIMARY KEY(id)
);

INSERT INTO roles (role) VALUES ('customer') ON CONFLICT (role)
DO NOTHING;
INSERT INTO roles (role) VALUES ('moderator') ON CONFLICT (role)
DO NOTHING;
INSERT INTO roles (role) VALUES ('admin') ON CONFLICT (role)
DO NOTHING;

CREATE TABLE IF NOT EXISTS  users_roles(
  userid TEXT NOT NULL,
  roleid INT NOT NULL,
  CONSTRAINT fk_ur
  FOREIGN KEY (userid)
  REFERENCES users(id),
  CONSTRAINT fk_ru
  FOREIGN KEY (roleid)
  REFERENCES roles(id),
  PRIMARY KEY(userid, roleid)
);

CREATE TABLE IF NOT EXISTS  categories(
  id TEXT NOT NULL UNIQUE,
  name VARCHAR(30) NOT NULL,
  PRIMARY KEY(id)
);

CREATE TABLE IF NOT EXISTS  products(
  id TEXT NOT NULL UNIQUE,
  name VARCHAR(30) NOT NULL,
  description TEXT NOT NULL,
  price INT NOT NULL,
  image TEXT NOT NULL,
  featured BOOLEAN NOT NULL DEFAULT false,
  published BOOLEAN NOT NULL DEFAULT true,
  category TEXT NOT NULL,
  weighed BOOLEAN NOT NULL DEFAULT true,
  created TIMESTAMP NOT NULL DEFAULT NOW(),
  updated TIMESTAMP NOT NULL DEFAULT NOW(),
  CONSTRAINT fk_cp
  FOREIGN KEY (category)
  REFERENCES categories(id),
  PRIMARY KEY(id)
);

SELECT apply_update_trigger('products');

CREATE TABLE IF NOT EXISTS customers(
  id TEXT NOT NULL UNIQUE,
  fullname VARCHAR(30) NOT NULL,
  email VARCHAR(30) NOT NULL,
  address TEXT NOT NULL,
  phone VARCHAR(15) NOT NULL,
  created TIMESTAMP NOT NULL DEFAULT NOW(),
  updated TIMESTAMP NOT NULL DEFAULT NOW(),
  PRIMARY KEY(id)
);

SELECT apply_update_trigger('customers');



DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'payment') THEN
        CREATE TYPE PAYMENT AS ENUM ('cash', 'stripe', 'paypal');
    END IF;
END$$;


CREATE TABLE IF NOT EXISTS orders(
  id TEXT NOT NULL UNIQUE,
  customer TEXT NOT NULL,
  pickuptime TIMESTAMP NOT NULL,
  fulfilled BOOLEAN NOT NULL DEFAULT false,
  method PAYMENT NOT NULL,
  created TIMESTAMP NOT NULL DEFAULT NOW(),
  updated TIMESTAMP NOT NULL DEFAULT NOW(),
  CONSTRAINT fk_csm
  FOREIGN KEY (customer)
  REFERENCES customers(id),
  PRIMARY KEY(id)
);

SELECT apply_update_trigger('orders');

CREATE TABLE IF NOT EXISTS purchases(
  id TEXT NOT NULL UNIQUE,
  productid TEXT NOT NULL,
  quantity INT NOT NULL,
  created TIMESTAMP NOT NULL DEFAULT NOW(),
  updated TIMESTAMP NOT NULL DEFAULT NOW(),
  orderid TEXT NOT NULL,
  CONSTRAINT fk_op
  FOREIGN KEY (productid)
  REFERENCES products(id),
  CONSTRAINT fk_os
  FOREIGN KEY (orderid)
  REFERENCES orders(id),
  PRIMARY KEY(id)
);

SELECT apply_update_trigger('purchases');


CREATE TABLE IF NOT EXISTS visits(
  id TEXT NOT NULL UNIQUE,
  ip TEXT NOT NULL,
  views INT NOT NULL,
  duration INT NOT NULL,
  sauce TEXT NOT NULL,
  agent TEXT NOT NULL,
  date TIMESTAMP NOT NULL DEFAULT NOW(),
  PRIMARY KEY(id)
);

