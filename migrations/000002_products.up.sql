
CREATE TABLE products (
    id SERIAL PRIMARY KEY,                  
    name TEXT NOT NULL,             
    description TEXT,                        
    category VARCHAR(100),                   
    image_url VARCHAR(255),                  
    average_rating DECIMAL(3, 2), 
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, 
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE reviews (
    id SERIAL PRIMARY KEY,                     
    product_id INT NOT NULL,                 
    rating INT,
    helpful_count INT,             
    comment TEXT,                               
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, 
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, 
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE 
);
