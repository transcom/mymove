LOCK TABLE uploads, contractor, proof_of_service_docs, invoices IN SHARE MODE;

create type upload_type as enum (
    'PRIME',
    'USER'
    );

alter table uploads
    alter column uploader_id drop not null,
    add column upload_type upload_type,
    drop constraint uploads_document_id_fkey,
    drop constraint uploads_uploader_id_fkey;


update uploads
    set upload_type = 'USER';

alter table uploads
    alter column upload_type set not null;

alter table proof_of_service_docs
    alter column upload_id drop not null;

create table user_uploads
(
    id uuid not null
        primary key,
    document_id uuid
        constraint user_uploads_document_id_fkey
            references documents,
    uploader_id uuid not null constraint user_uploads_uploader_id_fkey references users,
    upload_id uuid not null constraint user_uploads_uploads_id_fkey references uploads on delete restrict,
    created_at timestamp not null,
    updated_at timestamp not null,
        deleted_at timestamp with time zone

);

create index if not exists user_uploads_uploader_id_idx
    on user_uploads (uploader_id);

create index if not exists user_uploads_document_id_idx
    on user_uploads (document_id);

create index user_uploads_deleted_at_idx
    on user_uploads (deleted_at);

create table contractors
(
    id uuid not null
        constraint contractors_pkey
            primary key,
    created_at timestamp with time zone default now() not null,
    updated_at timestamp with time zone default now() not null,
    name varchar(80) not null,
    contract_number varchar(80) not null
        constraint contractors_contract_number_key
            unique,
    type varchar(80) not null
);

insert into contractors (id, created_at, updated_at, name,  contract_number, type)
select id, created_at, updated_at, name,  contract_number, type from contractor;

create table prime_uploads
(
    id uuid not null
        primary key,
    proof_of_service_docs_id uuid
        constraint prime_uploads_proof_of_service_docs_id_fkey
            references proof_of_service_docs not null,
    contractor_id uuid not null constraint prime_uploads_contractor_id_fkey references contractors,
    upload_id uuid not null constraint prime_uploads_uploads_id_fkey references uploads on delete restrict,
    created_at timestamp not null,
    updated_at timestamp not null,
    deleted_at timestamp with time zone

);

create index if not exists prime_uploads_proof_of_service_docs_id_idx
    on prime_uploads (proof_of_service_docs_id);

create index if not exists prime_uploads_contractor_id_idx
    on prime_uploads (contractor_id);

create index prime_uploads_deleted_at_idx
    on prime_uploads (deleted_at);

insert into user_uploads (id, document_id, uploader_id, upload_id, created_at, updated_at, deleted_at)
select uuid_generate_v4(), uploads.document_id, uploads.uploader_id, uploads.id, created_at, updated_at, deleted_at from uploads;

update uploads
    set uploader_id = null;

insert into prime_uploads (id, proof_of_service_docs_id, contractor_id, upload_id, created_at, updated_at)
select uuid_generate_v4(), proof_of_service_docs.id, '5db13bb4-6d29-4bdb-bc81-262f4513ecf6' ,proof_of_service_docs.upload_id, created_at, updated_at from proof_of_service_docs;

alter table invoices
    drop constraint invoices_uploads_id_fk,
    add column user_uploads_id uuid null constraint invoices_user_uploads_id_fkey references user_uploads on delete restrict;

update invoices set user_uploads_id = (select id from user_uploads where user_uploads.upload_id = invoices.upload_id);

alter table proof_of_service_docs
    drop constraint proof_of_service_docs_upload_id_fkey;

