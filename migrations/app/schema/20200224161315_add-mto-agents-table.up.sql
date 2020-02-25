CREATE TYPE mto_agents_type AS ENUM (
    'RELEASING_AGENT',
    'RECEIVING_AGENT'
    );

create table mto_agents
(
    id uuid PRIMARY KEY NOT NULL,
    move_task_order_id uuid REFERENCES move_task_orders,
    agent_type mto_agents_type,
    first_name text,
    last_name text,
    email text,
    phone text,
    created_at timestamp not null,
    updated_at timestamp not null
);