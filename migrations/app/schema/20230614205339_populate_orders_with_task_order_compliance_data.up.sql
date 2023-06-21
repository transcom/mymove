UPDATE orders
SET supply_and_services_cost_estimate = 'Prices for services under this task order will be in accordance with rates provided in GHC Attachment 2 - Pricing Rate Table. It is the responsibility of the contractor to provide the estimated weight quantity to apply to services on this task order, when applicable (See Attachment 1 - Performance Work Statement).',
	packing_and_shipping_instructions = (SELECT CONCAT_WS(' ', 'Packaging, packing, and shipping instructions as identified in the Conformed Copy of', contract_number, 'Attachment 1 Performance Work Statement') FROM contractors c WHERE TYPE = 'Prime'),
	method_of_payment = 'Payment will be made using the Third-Party Payment System (TPPS) Automated Payment System',
	naics = '488510 - FREIGHT TRANSPORTATION ARRANGEMENT';

ALTER TABLE orders
ALTER COLUMN supply_and_services_cost_estimate SET NOT NULL,
ALTER COLUMN packing_and_shipping_instructions SET NOT NULL,
ALTER COLUMN method_of_payment SET NOT NULL,
ALTER COLUMN naics SET NOT NULL;
