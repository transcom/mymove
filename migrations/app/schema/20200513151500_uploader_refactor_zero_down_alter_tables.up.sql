alter table proof_of_service_docs
    drop column upload_id;

alter table uploads
    drop column document_id,
    drop column uploader_id;

alter table invoices
    drop column upload_id;

drop table contractor;
