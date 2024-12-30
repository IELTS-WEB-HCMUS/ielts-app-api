CREATE TABLE otps (
    id SERIAL PRIMARY KEY,                          -- ID tự tăng, là khóa chính
    target VARCHAR(255) NOT NULL,                  -- Email hoặc số điện thoại
    type VARCHAR(50) NOT NULL,                     -- Loại OTP (reset_password, verify_email, v.v.)
    otp_code VARCHAR(6) NOT NULL,                  -- Mã OTP
    expired_at TIMESTAMP NOT NULL,                 -- Thời gian hết hạn
    is_verified BOOLEAN DEFAULT FALSE,            -- Trạng thái đã xác thực
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Thời gian tạo
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Thời gian cập nhật
    CONSTRAINT unique_target_type UNIQUE (target, type) -- Ràng buộc duy nhất cho target và type
);

ALTER TABLE otps DROP CONSTRAINT unique_target_type;

CREATE TABLE otp_attempts (
    id SERIAL PRIMARY KEY,
    otp_id INT NOT NULL REFERENCES otps(id) ON DELETE CASCADE,
    value VARCHAR(6) NOT NULL,
    is_success BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);