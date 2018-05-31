-- method_of_receipt should be OTHER_DD not OTHER
UPDATE reimbursements SET method_of_receipt = 'OTHER_DD' where method_of_receipt = 'OTHER';