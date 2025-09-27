-- Add reset password permission
INSERT INTO permissions (name_ar, name_en, code, created_at, updated_at)
VALUES
    ('إعادة تعيين كلمة المرور', 'Reset User Password', 'users:reset-password', NOW(), NOW());

-- Assign reset password permission to super_admin role
INSERT INTO permission_role (permission_id, role_id)
SELECT p.id, r.id
FROM permissions p, roles r
WHERE p.code = 'users:reset-password' AND r.code = 'super_admin';
