

-- INSERT INTO orders (id, title, description, deadline) values('asdasd','asd','asdasfasdasd','2 jan 2022');
-- INSERT INTO requirements (request, expectedoutcome, orderid, status) values('asdasd','asdasd','asdasd','0');

-- INSERT INTO users (id, username, email, pswd) values('1','elloy','elloy@gmail.com',sha256('100300'));

-- CREATE TABLE orders(id varchar(37) PRIMARY KEY, title varchar(50),description varchar(255),deadline timestamp );

-- CREATE TABLE requirements(id SERIAL PRIMARY KEY,request varchar(50),expectedoutcome varchar(50),orderid varchar(37),userid varchar(37),status bool,FOREIGN KEY(orderid) REFERENCES orders(id),FOREIGN KEY (userid) references users(id));

-- CREATE TABLE users(id varchar(37) PRIMARY KEY, username varchar(50),email varchar(50),pswd varchar (100));
drop table if exists image_submissions;
drop table if exists submissions;
drop table if exists tasks;


drop table if exists requirements ;
DROP table if exists orders;
DROP table if exists users;

CREATE TABLE orders(
    id varchar(37) PRIMARY KEY,
    title varchar(50),
    description varchar(255),
    deadline timestamp
);

CREATE TABLE users(
	id varchar(37) PRIMARY KEY,
    username varchar(50),
    pswd varchar (256),
    email varchar(30),
    user_role varchar(7)
);
CREATE TABLE requirements(
    id SERIAL PRIMARY KEY,
    request varchar(50),
    expected_outcome varchar(50),
    order_id varchar(37),
    status smallint,
    FOREIGN KEY(order_id) REFERENCES orders(id)
);

CREATE TABLE tasks(
	ID varchar(37) PRIMARY KEY,
	user_id varchar(37),
    assigner_id varchar(37),
	requirement_id int,
    fulfillment_status smallint,
    note varchar(200),
    allowed bool,
    num_of_prerequisite int,
    deadline timestamp,
    total_reviewer smallint,
    FOREIGN KEY (requirement_id) REFERENCES requirements(id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (assigner_id) REFERENCES users(id)
);

CREATE TABLE forwarded_review(
	reviewer_id varchar(37),
    task_id varchar(37),
    FOREIGN KEY (reviewer_id) REFERENCES users(id),
    FOREIGN KEY (task_id) REFERENCES tasks(id)
);


CREATE TABLE prerequisite(
	task_id varchar(37),
    prerequisite varchar(37)
);

CREATE TABLE submissions(
	id varchar(37) PRIMARY KEY,
    submit_time timestamp,
    message varchar(255),
    task_id varchar(37),
    FOREIGN KEY (task_id) REFERENCES tasks(id)
);

CREATE TABLE image_submissions(
	id int,
	image bytea,
    submission_id varchar(37),
    FOREIGN KEY (submission_id) REFERENCES submissions(id)
);

CREATE TABLE submission_messages(
    task_id varchar(37),
    user_id varchar(37),
    msg varchar(1024),
    FOREIGN KEY (task_id) REFERENCES tasks(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);



INSERT INTO users 
(id, username, pswd, email, user_role)
VALUES ('cd75bf2e-0876-46b4-a7a2-355ba2e8e034', 'elloy', sha256('100300'), 'elloy@elloy.com', 'Admin');

INSERT INTO users 
(id, username, pswd, email, user_role)
VALUES ('10b16316-ec54-4fdf-9a30-8deded11f633', 'jorich', sha256('100300'), 'jorich@elloy.com', 'User');


