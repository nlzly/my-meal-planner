-- User Queries

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: CreateUser :exec
INSERT INTO users (id, email, name, picture) 
VALUES ($1, $2, $3, $4);

-- name: UpdateUser :exec
UPDATE users 
SET name = $2, picture = $3, updated_at = now()
WHERE id = $1;


-- Meal Plan Queries

-- name: CreateMealPlan :exec
INSERT INTO meal_plans (id, name, description, created_by)
VALUES ($1, $2, $3, $4);

-- name: GetMealPlanByID :one
SELECT * FROM meal_plans WHERE id = $1;

-- name: GetMealPlansByUser :many
SELECT mp.*
FROM meal_plans mp
JOIN meal_plan_access mpa ON mpa.meal_plan_id = mp.id
WHERE mpa.user_id = $1;

-- name: UpdateMealPlan :exec
UPDATE meal_plans
SET name = $2, description = $3, updated_at = now()
WHERE id = $1;

-- name: DeleteMealPlan :exec
DELETE FROM meal_plans WHERE id = $1;


-- Meal Queries

-- name: CreateMeal :exec
INSERT INTO meals (id, meal_plan_id, name, description, day, meal_type)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: GetMealsByPlanID :many
SELECT * FROM meals WHERE meal_plan_id = $1;

-- name: GetMealByID :one
SELECT * FROM meals WHERE id = $1;

-- name: UpdateMeal :exec
UPDATE meals
SET name = $2, description = $3, day = $4, meal_type = $5, updated_at = now()
WHERE id = $1;

-- name: DeleteMeal :exec
DELETE FROM meals WHERE id = $1;


-- MealPlanAccess Queries

-- name: GrantMealPlanAccess :exec
INSERT INTO meal_plan_access (id, user_id, meal_plan_id, role)
VALUES ($1, $2, $3, $4);

-- name: GetMealPlanAccess :many
SELECT * FROM meal_plan_access WHERE meal_plan_id = $1;

-- name: GetUserMealPlanAccess :one
SELECT * FROM meal_plan_access 
WHERE user_id = $1 AND meal_plan_id = $2;


-- Share Link Queries

-- name: CreateShareLink :exec
INSERT INTO share_links (id, meal_plan_id, created_by, role, expires_at)
VALUES ($1, $2, $3, $4, $5);

-- name: GetShareLinkByID :one
SELECT * FROM share_links WHERE id = $1;

-- name: DeleteExpiredShareLinks :exec
DELETE FROM share_links WHERE expires_at < now();
