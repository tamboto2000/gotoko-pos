DROP TABLE IF EXISTS cashiers;
CREATE TABLE IF NOT EXISTS cashiers (
	id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    passcode TEXT NOT NULL,    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

DROP TABLE IF EXISTS cashier_sessions;
CREATE TABLE IF NOT EXISTS cashier_sessions (
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    cashier_id INT NOT NULL,
    issued_at BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

DROP TABLE IF EXISTS products;
CREATE TABLE IF NOT EXISTS products (
	id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    stock INT UNSIGNED NOT NULL,
    price INT UNSIGNED NOT NULL,
    image_url TEXT NOT NULL,
    category_id INT NOT NULL,    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,    
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FULLTEXT(name),
    INDEX(category_id)
);

DROP TABLE IF EXISTS product_skus;
CREATE TABLE IF NOT EXISTS product_skus (
	product_id INT NOT NULL,
    sku VARCHAR(5) NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,    
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

DROP TABLE IF EXISTS categories;
CREATE TABLE IF NOT EXISTS categories (
	id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,    
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

DROP TABLE IF EXISTS discounts;
CREATE TABLE IF NOT EXISTS discounts (
	id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    product_id INT NOT NULL,
    min_qty INT UNSIGNED NOT NULL,
    type ENUM('BUY_N' ,'PERCENT') NOT NULL,
    result INT UNSIGNED NOT NULL,
    expired_at BIGINT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,    
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

DROP TABLE IF EXISTS payment_methods;
CREATE TABLE IF NOT EXISTS payment_methods (
	id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    type ENUM('CASH', 'E-WALLET', 'EDC') NOT NULL,
    logo_url TEXT,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,    
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

DROP TABLE IF EXISTS order_items;
CREATE TABLE IF NOT EXISTS order_items (
	id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    order_id INT NOT NULL,
    product_id INT NOT NULL,
    qty INT NOT NULL,
    discount_id INT,
    total_final_price INT NOT NULL,
    total_normal_price INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

DROP TABLE IF EXISTS orders;
CREATE TABLE IF NOT EXISTS orders (
	id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    cashier_id INT NOT NULL,
    payment_id INT NOT NULL,    
    total_price INT UNSIGNED NOT NULL,
    total_paid INT UNSIGNED NOT NULL,
    total_return INT UNSIGNED NOT NULL,
    receipt_id VARCHAR(5) NOT NULL,
    is_pdf_downloaded BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- -- function to calculate subtotal per product
-- CREATE FUNCTION IF NOT EXISTS orderItemSubTotal(in_prod_id INT, in_qty INT)
-- RETURNS JSON
-- BEGIN
-- 	DECLARE v_prod_id INT;
--     DECLARE v_prod_name VARCHAR(100);    
-- 	DECLARE v_price INT;    
--     DECLARE v_stock INT;
--     DECLARE v_image_url TEXT;
--     DECLARE v_disc_id INT;
--     DECLARE v_disc_min_qty INT;
--     DECLARE v_disc_type ENUM('BUY_N' ,'PERCENT');
--     DECLARE v_disc_res INT;
--     DECLARE v_disc_expire BIGINT;
--     DECLARE v_total_price INT;
--     DECLARE V_final_price INT;
--     DECLARE v_disc JSON;
-- 	SELECT 
-- 		p.id,
--         p.name,
--         p.image_url,
-- 		p.price,
--         p.stock,
--         d.id,        
--         d.min_qty,
--         d.type,
--         d.result,
--         d.expired_at
-- 	INTO
-- 		v_prod_id,
--         v_prod_name,
--         v_image_url,
-- 		v_price,
--         v_stock,
--         v_disc_id,
--         v_disc_min_qty,
--         v_disc_type,
--         v_disc_res,
--         v_disc_expire
-- 	FROM products p
--     LEFT JOIN discounts d ON d.product_id = p.id
--     WHERE p.id = in_prod_id;
    
--     -- set initial total price
--     SET v_total_price = v_price * in_qty;
--     SET v_final_price = v_price * in_qty;
    
--     -- check if eligible for discount
--     -- if yes, use discount
-- 	IF v_disc_min_qty IS NOT NULL AND in_qty >= v_disc_min_qty AND v_disc_expire > UNIX_TIMESTAMP() THEN
-- 		SET v_disc = JSON_OBJECT(
-- 			'discountId', v_disc_id,
--             'qty', v_disc_min_qty,
--             'type', v_disc_type,
--             'result', v_disc_res,
--             'expiredAt', FROM_UNIXTIME(v_disc_expire, '%Y-%m-%dT%H:%i:%s.%fZ'),
--             'expiredAtFormat', FROM_UNIXTIME(v_disc_expire, '%d %b %Y')            
--         );
        
-- 		IF v_disc_type = 'BUY_N' THEN
-- 			IF in_qty != v_disc_min_qty THEN
-- 				SET v_disc = NULL;
-- 			ELSE					
-- 			 -- SET v_final_price = v_total_price - v_disc_res;
--              SET v_final_price = v_disc_res;
--              SET v_disc = JSON_SET(v_disc, '$.stringFormat', CONCAT('Buy ',v_disc_min_qty,' only Rp. ',FORMAT(v_disc_res, 0, 'de_DE')));
-- 			END IF;
-- 		ELSEIF v_disc_type = 'PERCENT' THEN
-- 			SET v_final_price = v_total_price - ((v_total_price / 100) * v_disc_res);
--             SET v_disc = JSON_SET(v_disc, '$.stringFormat', CONCAT('Discount ', v_disc_res, '%', ' Rp. ', FORMAT(v_final_price, 0, 'de_DE')));
-- 		END IF;
-- 	END IF;
		    
--     RETURN JSON_OBJECT(
-- 		'productId', v_prod_id,
--         'name', v_prod_name,
--         'stock', v_stock,
--         'price', v_price,
--         'image', v_image_url,
--         'qty', in_qty,
-- 		'totalNormalPrice', ROUND(v_total_price), 
-- 		'totalFinalPrice', ROUND(v_final_price),
--         'discount', v_disc
--     );
-- END;

-- -- function to calculate order subtotal
-- CREATE FUNCTION IF NOT EXISTS orderSubTotal(in_items JSON)
-- RETURNS JSON
-- BEGIN
-- 	DECLARE v_items_count INT;    
--     DECLARE loop_i INT;
--     DECLARE v_order_items JSON;
--     DECLARE v_subtotal INT;
--     DECLARE v_temp_order_item JSON;
    
--     SET v_subtotal = 0;
--     SET v_items_count = JSON_LENGTH(in_items);
--     SET v_order_items = JSON_ARRAY();
    
--     SET loop_i = 0;
--     WHILE loop_i < v_items_count DO
-- 		SET v_temp_order_item = orderItemSubTotal(
-- 			JSON_EXTRACT(in_items, CONCAT('$[',loop_i,'].productId')),
--             JSON_EXTRACT(in_items, CONCAT('$[',loop_i,'].qty'))
--         );
        
--         SET v_subtotal = v_subtotal + v_temp_order_item->'$.totalFinalPrice';
-- 		SET v_order_items = JSON_ARRAY_APPEND(v_order_items, '$', v_temp_order_item);
    
-- 		SET loop_i = loop_i+1;
--     END WHILE;	
    
--     RETURN JSON_OBJECT(
-- 		'subtotal', v_subtotal,
--         'products', v_order_items
--     );
-- END;

-- procedure to get list of products
DROP PROCEDURE IF EXISTS getProductList;
CREATE PROCEDURE IF NOT EXISTS getProductList(in_category_id INT, in_qs TEXT, in_limit INT, in_skip INT)
BEGIN
	IF in_limit <= 0 THEN
		SET in_limit = 2147483647;
	END IF;
	
	SELECT 
		DISTINCT p.id,
		p.name, 
		pk.sku,
		p.stock, 
		p.price,
		p.image_url,
		p.category_id,
		c.name AS category_name,
		d.min_qty AS discount_min_qty,
		d.type AS discount_type,
		d.result AS discount_result,
		d.expired_at AS discount_expired_at    
	FROM products p
	INNER JOIN product_skus pk ON pk.product_id = p.id
	LEFT JOIN categories c ON c.id = p.category_id
	LEFT JOIN discounts d ON d.product_id = p.id
	WHERE 
		(
			CASE WHEN in_category_id IS NOT NULL AND in_category_id > 0
			THEN 
				p.category_id = in_category_id 
			ELSE 
				TRUE 
			END
		) AND (
			CASE WHEN in_qs IS NOT NULL AND in_qs != ''
			THEN 
				MATCH(p.name) AGAINST(in_qs WITH QUERY EXPANSION) 
			ELSE 
				TRUE 
			END
		)
	
-- 	ORDER BY (
-- 		CASE WHEN in_qs IS NOT NULL AND in_qs != ''
--         THEN
-- 			MATCH(p.name) AGAINST(in_qs WITH QUERY EXPANSION)
-- 		ELSE
-- 			p.id
-- 		END
--     ) DESC
	LIMIT in_limit OFFSET in_skip; 
END;

-- procedure to get products count based on category and search keyword
CREATE PROCEDURE IF NOT EXISTS getProductCount(in_category_id INT, in_qs TEXT)
BEGIN
	SELECT COUNT(id) FROM products 
    WHERE (
		CASE WHEN in_qs IS NOT NULL AND in_qs != ''
		THEN 
			MATCH(name) AGAINST(in_qs WITH QUERY EXPANSION)
		ELSE
			TRUE
		END
    ) AND (
		CASE WHEN in_category_id IS NOT NULL AND in_category_id > 0
		THEN 
			category_id = in_category_id
		ELSE
			TRUE
		END
    );
END;

-- procedure to update a product
DROP PROCEDURE IF EXISTS updateProduct;
CREATE PROCEDURE IF NOT EXISTS updateProduct(
	in_id INT,
	in_category_id INT, 
    in_name VARCHAR(100), 
    in_image_url TEXT,
    in_price INT UNSIGNED,
    in_stock INT UNSIGNED
)
BEGIN
	UPDATE products SET
		category_id = CASE WHEN in_category_id IS NOT NULL AND in_category_id > 0 THEN in_category_id ELSE category_id END,
        name = CASE WHEN in_name IS NOT NULL AND in_name != '' THEN in_name ELSE name END,
        image_url = CASE WHEN in_image_url IS NOT NULL AND in_image_url != '' THEN in_image_url ELSE image_url END,
        price = CASE WHEN in_price IS NOT NULL AND in_price > 0 THEN in_price ELSE price END,
        stock = CASE WHEN in_stock IS NOT NULL AND in_stock >= 0 THEN in_stock ELSE stock END
	WHERE id = in_id;
END;

-- procedure to update a category
CREATE PROCEDURE IF NOT EXISTS updateCategory(
	in_id INT,	
    in_name VARCHAR(100)    
)
BEGIN
	UPDATE categories SET
		name = CASE WHEN in_name IS NOT NULL AND in_name != '' THEN in_name ELSE name END
	WHERE id = in_id;
END;

-- procedure to update a payment method
DROP PROCEDURE IF EXISTS updatePaymentMethod;
CREATE PROCEDURE IF NOT EXISTS updatePaymentMethod(
	in_id INT,	
    in_name VARCHAR(50),
    in_type ENUM('CASH', 'E-WALLET', 'EDC'),
    in_logo TEXT    
)
BEGIN
	UPDATE payment_methods SET
		name = CASE WHEN in_name IS NOT NULL AND in_name != '' THEN in_name ELSE name END,
        type = CASE WHEN in_type IS NOT NULL AND in_type != '' THEN in_type ELSE type END,
        logo_url = CASE WHEN in_logo IS NOT NULL AND in_logo != '' THEN in_logo ELSE logo_url END
	WHERE id = in_id;
END;


-- procedure to add order
DROP PROCEDURE IF EXISTS addOrder;
CREATE PROCEDURE IF NOT EXISTS addOrder(in_cashier_id INT, in_receipt_id VARCHAR(5), in_order JSON)
BEGIN
	DECLARE v_total_price INT;
    DECLARE v_total_paid INT;
    DECLARE v_total_return INT;
    DECLARE v_products JSON;
    DECLARE v_products_len INT;
    DECLARE v_order JSON;
    DECLARE v_last_id INT;
    DECLARE v_created_at TIMESTAMP;
    DECLARE v_updated_at TIMESTAMP;
    DECLARE v_product JSON;
    
    DECLARE loop_i INT;
	DECLARE track_no INT DEFAULT 0;
    DECLARE EXIT HANDLER FOR SQLEXCEPTION
    BEGIN        
        GET DIAGNOSTICS CONDITION 1 @`errno` = MYSQL_ERRNO, @`sqlstate` = RETURNED_SQLSTATE, @`text` = MESSAGE_TEXT;
        SET @full_error = CONCAT('ERROR ', @`errno`, ' (', @`sqlstate`, '): ', @`text`);
        SELECT track_no, @full_error, @`text`;    
        ROLLBACK;
    END;
    
    SET v_total_price = 0;
    -- SET v_total_paid = in_order->'$.totalPaid';
    SET v_total_paid = JSON_EXTRACT(in_order, '$.totalPaid');
    SET v_total_return = 0;
    CALL orderSubTotal(JSON_EXTRACT(in_order, '$.products'), v_products);
    SET v_products = JSON_EXTRACT(v_products, '$.products');
    SET v_products_len = JSON_LENGTH(v_products);
    
    SET loop_i = 0;    
    WHILE loop_i < v_products_len DO
		SET v_total_price = v_total_price + JSON_EXTRACT(v_products, CONCAT('$[',loop_i,'].totalFinalPrice'));    
		SET loop_i = loop_i + 1;
    END WHILE;
    
    SET v_total_return = v_total_paid - v_total_price;
    
    START TRANSACTION;
		INSERT INTO orders (
			cashier_id,
			payment_id,    
			total_price,
			total_paid,
			total_return,
			receipt_id
        ) VALUES (
			in_cashier_id,
            JSON_EXTRACT(in_order, '$.paymentId'),
            v_total_price,
            v_total_paid,
            v_total_return,
            in_receipt_id
        );
                
        SET v_last_id := LAST_INSERT_ID();
        
        SELECT 
			created_at, 
			updated_at 
        INTO
			v_created_at,
            v_updated_at
        FROM orders
        WHERE id = v_last_id;
        
        SET loop_i = 0;    
		WHILE loop_i < v_products_len DO
			SET v_product = JSON_EXTRACT(v_products, CONCAT('$[',loop_i,']'));
			INSERT INTO order_items (
				order_id,
				product_id,
				qty,
				discount_id,
				total_final_price,
				total_normal_price
            ) VALUES (
				v_last_id,
                JSON_EXTRACT(v_product, '$.productId'),
                JSON_EXTRACT(v_product, '$.qty'),
                CASE WHEN JSON_EXTRACT(v_product, '$.discount') IS NOT NULL THEN JSON_EXTRACT(v_product, '$.discount.discountId') ELSE NULL END,
                JSON_EXTRACT(v_product, '$.totalFinalPrice'),
                JSON_EXTRACT(v_product, '$.totalNormalPrice')
            );
            
			SET loop_i = loop_i + 1;
		END WHILE;        
    COMMIT;
    
    SET v_order = JSON_OBJECT(
		'order', JSON_OBJECT(
			'orderId', v_last_id,
			'cashiersId', in_cashier_id,
            'paymentTypesId', JSON_EXTRACT(in_order, '$.paymentId'),
            'totalPrice', v_total_price,
            'totalPaid', v_total_paid,
            'totalReturn', v_total_return,
            'receiptId', in_receipt_id,
            'updatedAt', DATE_FORMAT(v_updated_at, '%Y-%m-%dT%H:%i:%s.%fZ'),
            'createdAt', DATE_FORMAT(v_created_at, '%Y-%m-%dT%H:%i:%s.%fZ')            
        ),
        'products', v_products
    );
    
    SELECT v_order;
END;