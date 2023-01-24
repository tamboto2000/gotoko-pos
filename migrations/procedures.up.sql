-- procedure to get list of products
DELIMITER //
CREATE PROCEDURE IF NOT EXISTS getProductList(in_category_id INT, in_qs TEXT, in_limit INT, in_skip INT)
BEGIN
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
	
	ORDER BY (
		CASE WHEN in_qs IS NOT NULL AND in_qs != ''
        THEN
			MATCH(p.name) AGAINST(in_qs WITH QUERY EXPANSION)
		ELSE
			p.id
		END
    ) DESC
	LIMIT in_limit OFFSET in_skip;
END //
DELIMITER ;

-- procedure to get products count based on category and search keyword
DELIMITER //
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
END //
DELIMITER ;

-- procedure to update a product
DELIMITER //
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
END //
DELIMITER ;

-- procedure to update a category
DELIMITER //
CREATE PROCEDURE IF NOT EXISTS updateCategory(
	in_id INT,	
    in_name VARCHAR(100)    
)
BEGIN
	UPDATE categories SET
		name = CASE WHEN in_name IS NOT NULL AND in_name != '' THEN in_name ELSE name END
	WHERE id = in_id;
END //
DELIMITER ;

-- procedure to update a payment method
DELIMITER //
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
END //
DELIMITER ;

-- procedure to add order
DELIMITER //
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
    SET v_total_paid = in_order->'$.totalPaid';
    SET v_total_return = 0;
    SET v_products = orderSubTotal(in_order->'$.products');
    SET v_products = v_products->'$.products';
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
            in_order->'$.paymentId',
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
                v_product->'$.productId',
                v_product->'$.qty',
                CASE WHEN v_product->'$.discount' IS NOT NULL THEN v_product->'$.discount.discountId' ELSE NULL END,
                v_product->'$.totalFinalPrice',
                v_product->'$.totalNormalPrice'
            );
            
			SET loop_i = loop_i + 1;
		END WHILE;        
    COMMIT;
    
    SET v_order = JSON_OBJECT(
		'order', JSON_OBJECT(
			'orderId', v_last_id,
			'cashiersId', in_cashier_id,
            'paymentTypesId', in_order->'$.paymentId',
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
END //
DELIMITER ;