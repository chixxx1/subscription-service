CREATE TABLE IF NOT EXISTS subscription (
	id					 SERIAL			  PRIMARY KEY,
	service_name VARCHAR(255)	NOT NULL,
	price				 INTEGER		  NOT NULL			CHECK (price >= 0),
	user_id			 UUID         NOT NULL,
	start_date 	 DATE				  NOT NULL,
	end_date		 DATE,
	created_at	 TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
	updated_at   TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_subscr_user_id ON subscription(user_id);
CREATE INDEX idx_subscr_service_name ON subscription(service_name);
CREATE INDEX idx_subscr_dates ON subscription(start_date, end_date);
