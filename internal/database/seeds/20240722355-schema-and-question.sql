-- Shop Schema
INSERT INTO dp_schemas (schema_id, description, initial_sql)
VALUES (
    'shop',
    'The schema that is for a shop',
    'CREATE TABLE products (
        product_id INT PRIMARY KEY,
        product_name VARCHAR(100),
        price DECIMAL(10, 2),
        stock INT
    );

    CREATE TABLE customers (
        customer_id INT PRIMARY KEY,
        customer_name VARCHAR(100),
        email VARCHAR(100)
    );

    CREATE TABLE orders (
        order_id INT PRIMARY KEY,
        customer_id INT,
        order_date DATE,
        total DECIMAL(10, 2),
        FOREIGN KEY (customer_id) REFERENCES customers(customer_id)
    );

    INSERT INTO products (product_id, product_name, price, stock) VALUES
    (1, ''Laptop'', 999.99, 10),
    (2, ''Mouse'', 19.99, 100);

    INSERT INTO customers (customer_id, customer_name, email) VALUES
    (1, ''Alice'', ''alice@example.com''),
    (2, ''Bob'', ''bob@example.com'');

    INSERT INTO orders (order_id, customer_id, order_date, total) VALUES
    (1, 1, ''2024-01-01'', 1019.98);'
);
-- Library Schema
INSERT INTO dp_schemas (schema_id, description, initial_sql)
VALUES (
    'library',
    'The schema that is for a library',
    'CREATE TABLE books (
        book_id INT PRIMARY KEY,
        title VARCHAR(200),
        author VARCHAR(100),
        genre VARCHAR(50),
        published_date DATE
    );

    CREATE TABLE members (
        member_id INT PRIMARY KEY,
        member_name VARCHAR(100),
        membership_date DATE
    );

    CREATE TABLE borrowings (
        borrowing_id INT PRIMARY KEY,
        book_id INT,
        member_id INT,
        borrowing_date DATE,
        return_date DATE,
        FOREIGN KEY (book_id) REFERENCES books(book_id),
        FOREIGN KEY (member_id) REFERENCES members(member_id)
    );

    INSERT INTO books (book_id, title, author, genre, published_date) VALUES
    (1, ''1984'', ''George Orwell'', ''Dystopian'', ''1949-06-08''),
    (2, ''To Kill a Mockingbird'', ''Harper Lee'', ''Fiction'', ''1960-07-11'');

    INSERT INTO members (member_id, member_name, membership_date) VALUES
    (1, ''John Doe'', ''2023-01-15''),
    (2, ''Jane Smith'', ''2022-12-20'');

    INSERT INTO borrowings (borrowing_id, book_id, member_id, borrowing_date, return_date) VALUES
    (1, 1, 1, ''2024-06-01'', ''2024-06-15'');'
);
-- School Schema
INSERT INTO dp_schemas (schema_id, description, initial_sql)
VALUES (
    'school',
    'The schema that is for a school',
    'CREATE TABLE students (
        student_id INT PRIMARY KEY,
        student_name VARCHAR(100),
        date_of_birth DATE,
        grade INT
    );

    CREATE TABLE teachers (
        teacher_id INT PRIMARY KEY,
        teacher_name VARCHAR(100),
        subject VARCHAR(50)
    );

    CREATE TABLE classes (
        class_id INT PRIMARY KEY,
        class_name VARCHAR(50),
        teacher_id INT,
        schedule_time TIME,
        FOREIGN KEY (teacher_id) REFERENCES teachers(teacher_id)
    );

    CREATE TABLE enrollments (
        enrollment_id INT PRIMARY KEY,
        student_id INT,
        class_id INT,
        enrollment_date DATE,
        FOREIGN KEY (student_id) REFERENCES students(student_id),
        FOREIGN KEY (class_id) REFERENCES classes(class_id)
    );

    INSERT INTO students (student_id, student_name, date_of_birth, grade) VALUES
    (1, ''Emily'', ''2010-05-12'', 5),
    (2, ''Daniel'', ''2011-07-23'', 4);

    INSERT INTO teachers (teacher_id, teacher_name, subject) VALUES
    (1, ''Mr. Brown'', ''Mathematics''),
    (2, ''Ms. Green'', ''Science'');

    INSERT INTO classes (class_id, class_name, teacher_id, schedule_time) VALUES
    (1, ''Math 101'', 1, ''09:00:00''),
    (2, ''Science 101'', 2, ''10:00:00'');

    INSERT INTO enrollments (enrollment_id, student_id, class_id, enrollment_date) VALUES
    (1, 1, 1, ''2024-07-01''),
    (2, 2, 2, ''2024-07-01'');'
);
-- Insert into dp_questions
INSERT INTO dp_questions (
    schema_id,
    type,
    difficulty,
    title,
    description,
    answer
)
VALUES (
    'shop',
    '條件查詢',
    'easy',
    'Find a product in the shop',
    'Write a SQL query to find the ''Laptop'' product in the shop schema.',
    'SELECT * FROM products WHERE product_name = ''Laptop'';'
),
(
    'shop',
    '條件查詢',
    'easy',
    'List all customers',
    'Write a SQL query to list all customers in the shop schema.',
    'SELECT * FROM customers;'
),
(
    'shop',
    '條件查詢',
    'medium',
    'Find orders by a specific customer',
    'Write a SQL query to find all orders placed by the customer with ID 1.',
    'SELECT * FROM orders WHERE customer_id = 1;'
),
(
    'shop',
    '群組應用',
    'medium',
    'Find products below a certain stock level',
    'Write a SQL query to find all products with a stock level less than 20.',
    'SELECT * FROM products WHERE stock < 20;'
),
(
    'library',
    '群組應用',
    'easy',
    'List all books',
    'Write a SQL query to list all books in the library schema.',
    'SELECT * FROM books;'
),
(
    'library',
    '群組應用',
    'medium',
    'Find a book by title',
    'Write a SQL query to find the book ''1984'' by George Orwell.',
    'SELECT * FROM books WHERE title = ''1984'';'
),
(
    'library',
    '子查詢應用',
    'easy',
    'List all members',
    'Write a SQL query to list all members in the library schema.',
    'SELECT * FROM members;'
),
(
    'library',
    '子查詢應用',
    'medium',
    'Find borrowings by a specific member',
    'Write a SQL query to find all borrowings by the member with ID 1.',
    'SELECT * FROM borrowings WHERE member_id = 1;'
),
(
    'library',
    '子查詢應用',
    'hard',
    'Find overdue borrowings',
    'Write a SQL query to find all borrowings \
    where the return date is past due (today''s date is ''2024-07-16'').',
    'SELECT * FROM borrowings WHERE return_date < ''2024-07-16'';'
),
(
    'school',
    '子查詢+群組綜合應用',
    'easy',
    'List all students',
    'Write a SQL query to list all students in the school schema.',
    'SELECT * FROM students;'
),
(
    'school',
    '子查詢+群組綜合應用',
    'medium',
    'Find students in a specific grade',
    'Write a SQL query to find all students in grade 5.',
    'SELECT * FROM students WHERE grade = 5;'
),
(
    'school',
    '子查詢+群組綜合應用',
    'easy',
    'List all teachers',
    'Write a SQL query to list all teachers in the school schema.',
    'SELECT * FROM teachers;'
),
(
    'school',
    '進階外部查詢',
    'medium',
    'Find classes taught by a specific teacher',
    'Write a SQL query to find all classes taught by the teacher with ID 1.',
    'SELECT * FROM classes WHERE teacher_id = 1;'
),
(
    'school',
    '進階外部查詢',
    'easy',
    'List all classes',
    'Write a SQL query to list all classes in the school schema.',
    'SELECT * FROM classes;'
),
(
    'school',
    '進階外部查詢',
    'medium',
    'Find enrollments by a specific student',
    'Write a SQL query to find all enrollments for the student with ID 1.',
    'SELECT * FROM enrollments WHERE student_id = 1;'
),
(
    'school',
    '聯集應用',
    'hard',
    'Find students enrolled in a specific class',
    'Write a SQL query to find all students enrolled in the class \
    with ID 1.',
    'SELECT s.student_id, s.student_name FROM students s \
    JOIN enrollments e ON s.student_id = e.student_id WHERE e.class_id = 1;'
),
(
    'shop',
    '聯集應用',
    'hard',
    'Calculate total sales',
    'Write a SQL query to calculate the total sales in the shop.',
    'SELECT SUM(total) AS total_sales FROM orders;'
),
(
    'library',
    '聯集應用',
    'medium',
    'Find books by genre',
    'Write a SQL query to find all books in the genre ''Fiction''.',
    'SELECT * FROM books WHERE genre = ''Fiction'';'
),
(
    'school',
    '進階Exists指令應用',
    'hard',
    'Find students older than a specific age',
    'Write a SQL query to find all students older than 10 \
    years old (assuming today''s date is ''2024-07-16'').',
    'SELECT * FROM students WHERE date_of_birth < ''2014-07-16'';'
);
