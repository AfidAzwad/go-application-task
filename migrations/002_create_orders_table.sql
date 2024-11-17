-- Create the order_status_enum type if it does not exist
DO $$
BEGIN
    -- Check if 'order_status_enum' type does not exist, and create it if necessary
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'order_status_enum') THEN
CREATE TYPE order_status_enum AS ENUM ('pending', 'completed', 'cancelled');
END IF;
END;
$$ LANGUAGE plpgsql;

-- Create the orders table if it does not exist
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'orders') THEN
CREATE TABLE orders (
                        id SERIAL PRIMARY KEY,
                        consignment_id VARCHAR(255) UNIQUE NOT NULL,
                        store_id INT NOT NULL,
                        merchant_order_id VARCHAR(255) DEFAULT '',
                        recipient_name VARCHAR(255) NOT NULL,
                        recipient_phone VARCHAR(255) NOT NULL,
                        recipient_address VARCHAR(255) NOT NULL,
                        recipient_city INT NOT NULL,
                        recipient_zone INT NOT NULL,
                        recipient_area INT NOT NULL,
                        delivery_type INT NOT NULL,
                        item_type INT NOT NULL,
                        special_instruction VARCHAR(255) DEFAULT '',
                        item_quantity INT NOT NULL,
                        item_weight FLOAT NOT NULL,
                        amount_to_collect FLOAT NOT NULL,
                        item_description VARCHAR(255) DEFAULT '',
                        order_status order_status_enum NOT NULL DEFAULT 'pending',
                        delivery_fee FLOAT NOT NULL,
                        cod_fee FLOAT NOT NULL,
                        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                        user_id INT REFERENCES users(id) ON DELETE CASCADE
);

-- Add unique index for consignment_id
CREATE UNIQUE INDEX IF NOT EXISTS idx_consignment_id ON orders (consignment_id);

END IF;
END;
$$ LANGUAGE plpgsql;

