CREATE TYPE payment_status_enum AS ENUM ('Pending', 'Processed', 'Failed');
CREATE TYPE payment_methods_enum AS ENUM ('Bank_transfer', 'E-payment');
CREATE TYPE master_status_enum AS ENUM ('Enabled', 'Disabled');
CREATE TYPE question_type_enum AS ENUM('Radio', 'Checkbox');

CREATE TABLE roles(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    modified_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE users(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email TEXT NOT NULL UNIQUE,
    role_id UNIQUE REFERENCES roles(id),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    modified_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);



CREATE TBALE stipends(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    student_id UUID REFERENCES students(id),
    amount INT NOT NULL, 
    payment_date TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    payment payment_status_enum,
    payment_method payment_methods_enum,
    journal_number INT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    modified_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE banks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    modified_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE student_bank_details(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    student_id UUID REFERENCES students(id), 
    bank_id UUID REFERENCES banks(ID),
    account_number INT NOT NULL UNIQUE,
    account_holder_name TEXT NOT NULL, s
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    modified_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

# student services

CREATE TABLE colleges (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL UNIQUE,
    modified_by UUID REFERENCES users(id),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    modified_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE programs(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL UNIQUE,
    college_id UUID REFERENCES colleges(id),
    modified_by UUID REFERENCES users(id),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    modified_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE students(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL, 
    rub_id_card_number INT NOT NULL UNIQUE,
    email INT NOT NULL UNIQUE,
    phone_number INT NOT NULL,
    date_of_bith DATA,
    college_id UUID REFERENCES colleges(id),
    modified_by UUID REFERENCES users(id),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    modified_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);


