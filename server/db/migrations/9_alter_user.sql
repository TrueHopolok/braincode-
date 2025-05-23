-- Should work with previous fill data
ALTER TABLE User
ADD CONSTRAINT CHK_name CHECK(LENGTH(name) > 3);