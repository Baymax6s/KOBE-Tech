UPDATE user_profiles 
SET object_key = '', is_uploaded = FALSE 
WHERE user_id IN (1, 2, 3, 4);