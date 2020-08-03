CREATE TABLE categories
   ( category_id SERIAL CONSTRAINT category_id_pk PRIMARY KEY,
category_name character varying(25) CONSTRAINT category_name_nn NOT NULL,
category_active boolean DEFAULT false
   );


CREATE TABLE products
( product_id SERIAL CONSTRAINT product_id_pk PRIMARY KEY,
product_name character varying(25) CONSTRAINT product_name_nn NOT NULL,
supply_product_price numeric(8,2),
sale_product_price numeric(8,2),
category_id INTEGER CONSTRAINT category_id_fk REFERENCES categories (category_id),
product_quantity INTEGER DEFAULT 0,
product_active boolean DEFAULT false
   );

ALTER TABLE products
ADD CONSTRAINT supply_product_price_check CHECK (supply_product_price >= 0.0);

ALTER TABLE products
ADD CONSTRAINT sale_product_price_check CHECK (sale_product_price >= 0.0);