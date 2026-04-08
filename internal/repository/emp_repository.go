package repository

// 1. สร้าง service สำหรับเพิ่ม/แก้ไข employee
// 2. สร้าง service สำหรับค้นหาพนักงานด้วยชื่อนามสกุลหรือแผนก
// 3. สร้าง service สำหรับหาจำนวนพนักงานในแต่ละแผนก โดย filter แผนกได้

// CREATE TABLE departments (
//     id SERIAL PRIMARY KEY,
//     name VARCHAR(100) NOT NULL,
//     work_floor INT NOT NULL
// );

// CREATE TABLE employees (
//     id SERIAL PRIMARY KEY,
//     firstname VARCHAR(100) NOT NULL,
//     lastname VARCHAR(100) NOT NULL,
// 	created_at TIMESTAMPTZ DEFAULT NOW(),
//     updated_at TIMESTAMPTZ DEFAULT NOW(),
//     department_id INT,
//     CONSTRAINT fk_departments
//         FOREIGN KEY(department_id)
//         REFERENCES departments(id)
//         ON DELETE SET NULL
// );

// INSERT INTO departments (name, work_floor) VALUES
// ('Developer', 1),
// ('CO', 2),
// ('Admin', 3),
// ('Tester', 4),
// ('SA', 5),
// ('UX/UI', 2);
