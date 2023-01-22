-- procedure to get list of products
DROP PROCEDURE IF EXISTS getProductList;
DELIMITER //
CREATE PROCEDURE getProductList(in_category_id INT, in_qs TEXT, in_limit INT, in_skip INT)
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
DROP PROCEDURE IF EXISTS getProductCount;
DELIMITER //
CREATE PROCEDURE getProductCount(in_category_id INT, in_qs TEXT)
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
DROP PROCEDURE IF EXISTS updateProduct;
DELIMITER //
CREATE PROCEDURE updateProduct(
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