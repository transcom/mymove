--Set auto approved to false for IOSFSC and IDSFSC
update re_service_items
set is_auto_approved = false
where service_id in ('81e29d0c-02a6-4a7a-be02-554deb3ee49e', '690a5fc1-0ea5-4554-8294-a367b5daefa9');