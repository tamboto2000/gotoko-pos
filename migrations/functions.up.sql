-- procedure to calculate subtotal per product
DROP PROCEDURE IF EXISTS orderItemSubTotal;
DELIMITER //
CREATE PROCEDURE orderItemSubTotal(IN in_prod_id INT, IN in_qty INT, OUT out_order_item JSON)
BEGIN
	DECLARE v_prod_id INT;
    DECLARE v_prod_name VARCHAR(100);    
	DECLARE v_price INT;    
    DECLARE v_stock INT;
    DECLARE v_image_url TEXT;
    DECLARE v_disc_id INT;
    DECLARE v_disc_min_qty INT;
    DECLARE v_disc_type ENUM('BUY_N' ,'PERCENT');
    DECLARE v_disc_res INT;
    DECLARE v_disc_expire BIGINT;
    DECLARE v_total_price INT;
    DECLARE V_final_price INT;
    DECLARE v_disc JSON;
	SELECT 
		p.id,
        p.name,
        p.image_url,
		p.price,
        p.stock,
        d.id,        
        d.min_qty,
        d.type,
        d.result,
        d.expired_at
	INTO
		v_prod_id,
        v_prod_name,
        v_image_url,
		v_price,
        v_stock,
        v_disc_id,
        v_disc_min_qty,
        v_disc_type,
        v_disc_res,
        v_disc_expire
	FROM products p
    LEFT JOIN discounts d ON d.product_id = p.id
    WHERE p.id = in_prod_id;
    
    -- set initial total price
    SET v_total_price = v_price * in_qty;
    SET v_final_price = v_price * in_qty;
    
    -- check if eligible for discount
    -- if yes, use discount
	IF v_disc_min_qty IS NOT NULL AND in_qty >= v_disc_min_qty AND v_disc_expire > UNIX_TIMESTAMP() THEN
		SET v_disc = JSON_OBJECT(
			'discountId', v_disc_id,
            'qty', v_disc_min_qty,
            'type', v_disc_type,
            'result', v_disc_res,
            'expiredAt', FROM_UNIXTIME(v_disc_expire, '%Y-%m-%dT%H:%i:%s.%fZ'),
            'expiredAtFormat', FROM_UNIXTIME(v_disc_expire, '%d %b %Y')            
        );
        
		IF v_disc_type = 'BUY_N' THEN
			IF in_qty != v_disc_min_qty THEN
				SET v_disc = NULL;
			ELSE					
			 -- SET v_final_price = v_total_price - v_disc_res;
             SET v_final_price = v_disc_res;
             SET v_disc = JSON_SET(v_disc, '$.stringFormat', CONCAT('Buy ',v_disc_min_qty,' only Rp. ',FORMAT(v_disc_res, 0, 'de_DE')));
			END IF;
		ELSEIF v_disc_type = 'PERCENT' THEN
			SET v_final_price = v_total_price - ((v_total_price / 100) * v_disc_res);
            SET v_disc = JSON_SET(v_disc, '$.stringFormat', CONCAT('Discount ', v_disc_res, '%', ' Rp. ', FORMAT(v_final_price, 0, 'de_DE')));
		END IF;
	END IF;
		    
    SELECT JSON_OBJECT(
		'productId', v_prod_id,
        'name', v_prod_name,
        'stock', v_stock,
        'price', v_price,
        'image', v_image_url,
        'qty', in_qty,
		'totalNormalPrice', ROUND(v_total_price), 
		'totalFinalPrice', ROUND(v_final_price),
        'discount', v_disc
    ) INTO out_order_item;
END //
DELIMITER ;

-- function to calculate order subtotal
DELIMITER //
CREATE FUNCTION IF NOT EXISTS orderSubTotal(in_items JSON)
RETURNS JSON
BEGIN
	DECLARE v_items_count INT;    
    DECLARE loop_i INT;
    DECLARE v_order_items JSON;
    DECLARE v_subtotal INT;
    DECLARE v_temp_order_item JSON;
    
    SET v_subtotal = 0;
    SET v_items_count = JSON_LENGTH(in_items);
    SET v_order_items = JSON_ARRAY();
    
    SET loop_i = 0;
    WHILE loop_i < v_items_count DO
		SET v_temp_order_item = orderItemSubTotal(
			JSON_EXTRACT(in_items, CONCAT('$[',loop_i,'].productId')),
            JSON_EXTRACT(in_items, CONCAT('$[',loop_i,'].qty'))
        );
        
        SET v_subtotal = v_subtotal + v_temp_order_item->'$.totalFinalPrice';
		SET v_order_items = JSON_ARRAY_APPEND(v_order_items, '$', v_temp_order_item);
    
		SET loop_i = loop_i+1;
    END WHILE;	
    
    RETURN JSON_OBJECT(
		'subtotal', v_subtotal,
        'products', v_order_items
    );
END //
DELIMITER ;