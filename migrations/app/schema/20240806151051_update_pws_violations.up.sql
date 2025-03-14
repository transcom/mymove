alter table report_violations
Drop CONSTRAINT report_violations_violation_id_fkey;

truncate pws_violations;

INSERT INTO
    public.pws_violations (
        id,
        display_order,
        paragraph_number,
        title,
        category,
        sub_category,
        requirement_summary,
        requirement_statement,
        is_kpi,
        additional_data_elem
    )
VALUES
    (
        '9cdc8dc3-6cf4-46fb-b272-1468ef40796f',
        1,
        '1.2.3',
        'Customer Support',
        'Pre-Move Services',
        'Customer Support',
        'Provide 24/7 customer support capability',
        'The contractor shall maintain a 24-hour, 7-day per week customer support capability for all issues pertaining to a customer''s move. The contractor shall staff the customer support capability such that all customer calls, to include Interactive Voice Response (IVR) time, will be answered within four (4) minutes. Wait time is cumulative and does not restart if a call is handed off or escalated to Level 2/3.',
        false,
        ''
    ),
    (
        'c359ebc3-a506-4f41-8f91-409d59c97b22',
        2,
        '1.2.5.1',
        'Point of Contact (POC)',
        'Pre-Move Services',
        'Customer Support',
        'Provide a single point of contact (POC)',
        'The contractor shall assign, during initial communication with each customer, a single POC responsible for coordination and communication throughout all phases of the move. The POC''s contact information shall be maintained throughout the entire shipment process and until all associated actions are final.',
        false,
        ''
    ),
    (
        'eec7dd42-97e2-4c6c-b9e4-3d1f88bafe41',
        3,
        '1.2.5.2',
        'Pre-Move Services',
        'Pre-Move Services',
        'Counseling',
        'Establish contact with the customer within one (1) government business day (GBD) of receiving an order',
        'When ordered, the contractor shall establish contact with the customer within one (1) GBD of receiving an order and customer contact information and shall then provide ordered services.',
        false,
        ''
    ),
    (
        'cdcda061-442e-45dc-b6bc-4a122b91f2ea',
        4,
        '1.2.5.2',
        'Customer Counseling',
        'Pre-Move Services',
        'Counseling',
        'Provide accurate entitlement counseling and forms',
        'The contractor shall provide accurate entitlement and shipment counseling in accordance with all applicable US Government, DoD, Joint, Military Service-specific, and International regulations and instructions to all customers regardless of location and accessibility, to include customers in remote locations, or customers who lack the ability to conduct counseling via face-to-face communication or the Internet. Upon counseling completion, the contractor shall provide an electronic or hard copy form as a record of the customer application for shipment or storage of personal property. The inventory shall include all details listed in Appendix H. The Government may issue task orders for site-specific in-person counseling services IAW Appendix D, para D.5. Customers shall be offered either an in-person or virtual pre-move survey as part of the counseling process.',
        false,
        ''
    ),
    (
        '75a40e3b-4c8f-4750-b8da-925bf0200b28',
        5,
        '1.2.5.2.1',
        'Personally Procured Moves (PPM)',
        'Pre-Move Services',
        'Counseling',
        'Provide accurate PPM counseling',
        'When ordered, the contractor shall provide accurate entitlement on PPMs during counseling in accordance with all applicable US Government, DoD, Joint, Military Service-specific, and International regulations and instructions. The contractor shall calculate the estimate of what it would cost the contractor to perform the relocation and provide the customer the incentive estimate for the PPM.  The estimate shall be provided to the customer at the time of counseling. The contractor shall provide the customer with the updated cost based on actual shipment weight(s) once the customer completes their PPM. Both the estimate and any updates shall be calculated based on the total weight of all shipments executed against the remaining weight entitlement and adjusted accordingly with the form including all details listed in Appendix H. The contractor shall provide the customer with the government designated form(s) for PPMs, and advise the customer of all documentation needed for reimbursement from the military services (e.g. full and empty weight tickets).',
        false,
        ''
    ),
    (
        '39680e15-81eb-40b9-966e-2cbbf9235724',
        6,
        '1.2.5.3',
        'Scheduling',
        'Pre-Move Services',
        'Counseling',
        'Schedule relocation during counseling',
        'The contractor shall schedule shipment relocation services during customer counseling, which must be conducted IAW scheduling timelines referenced in 1.2.5.3.1 below. Based on customer request, the contractor shall provide a pickup date spread for required packing, pickup, and delivery during counseling and firm dates after counseling as summarized below.',
        true,
        'observedPickupSpreadDates'
    ),
    (
        'e1ee1719-a6d5-49b0-ad3b-c4dac0a3f16f',
        7,
        '1.2.5.3.1',
        'Scheduling Requirements',
        'Pre-Move Services',
        'Counseling',
        'Schedule relocation using pickup spread rules',
        'During Customer Counseling, or Scheduling if Customer Counseling is not ordered, the contractor shall provide a pickup date spread in accordance with the timelines in the table below, agreeable to the customer, which shall not to exceed seven (7) consecutive calendar days from the members requested pickup date contained within that spread. The contractor shall document the start and end of the spread, and the customer''s acceptance. The contractor shall provide a firm schedule for all applicable relocation services in accordance with the timelines in the table below (See PWS Paragraph 1.2.5.3.1.). The contractor shall ensure all firm dates are within the previously agreed upon spread. An â€œapproved orderâ€ is an order sent to the contractor after the Ordering Officer (OO) validates the requirement.',
        true,
        'observedPickupSpreadDates'
    ),
    (
        '661f2950-3e21-489d-be0b-2a60922e3af2',
        8,
        '1.2.5.4.1',
        'Weight Estimates',
        'Pre-Move Services',
        'Weight Estimate',
        'Provide accurate weight estimate',
        'The contractor shall provide the government and customer weight estimates on all shipments no later than 10 days prior to the first scheduled pack or pickup date. For shipments ordered less than ten (10) days prior to first scheduled pack or pickup date, weight estimates must be provided no later than three (3) days prior to first scheduled pack or pickup date. For shipments ordered less than three (3) days prior to the first scheduled pack or pickup date, weight estimates must be provided no later than one (1) day prior to first scheduled pack or pickup date. The government will only pay costs associated with shipments up to 110% of the estimated weight.',
        false,
        ''
    ),
    (
        '61245a25-6684-434d-aa11-13eb0725c5fb',
        9,
        '1.2.6.4',
        'Items Requiring Government Pre-Approval',
        'Physical Move Services',
        'Additional Services',
        'Obtain pre-approvals for additional services',
        'The following services referenced in 1.2.6.4.1 and 1.2.6.4.2. must be approved by the government prior to performance. Requests for approval shall be sent to the OO at the responsible origin or destination of the shipment.',
        false,
        ''
    ),
    (
        '9a0adb53-8f54-45ef-8e87-c10d83b4f70b',
        10,
        '1.2.6.4.1',
        'Crating',
        'Physical Move Services',
        'Additional Services',
        'Provide appropriate crating services',
        'Upon approval, the contractor shall perform crating services for items such as mirrors, paintings, glass or marble tabletops and similar fragile articles, and taxidermy when crates are not provided by the customer or when the customer provided crates are not serviceable. This does not include cases, footlockers, passenger bags, cartons, boxes, tri-wall containers, liftvans, and barrels that may be placed in a cargo transporter (commercial sea vans; container express cargo transporters and other transoceanic cargo transporters) for ocean or air transport. The customer retains ownership of all crates. (Per Attachment 2, Pricing Rate Table, the price for crating services is for the construction of new crates only).',
        false,
        ''
    );

INSERT INTO
    public.pws_violations (
        id,
        display_order,
        paragraph_number,
        title,
        category,
        sub_category,
        requirement_summary,
        requirement_statement,
        is_kpi,
        additional_data_elem
    )
VALUES
    (
        'eb749a9a-eb3b-429d-b935-220e13cab5ca',
        11,
        '1.2.6.4.2',
        'Shuttles',
        'Physical Move Services',
        'Additional Services',
        'Provide appropriate shuttle services',
        'Upon approval, the contractor shall perform shuttle services to pick up or deliver shipments when the origin or destination delivery location is inaccessible due to building design, nonexistent or inaccessible roadway, inadequate or unsafe public or private road, overhead obstruction, deterioration of roadway due to rain, flood, or snow, construction, or other obstacles preventing the linehaul truck from accessing the pickup or delivery location. A shuttle is defined as a truck-to-truck transfer between a larger and smaller vehicle (or vice versa) that allows for safe pickup or delivery from the nearest safely accessible point to the pickup or final delivery, not a truck-to-warehouse or warehouse-to-truck transfer.',
        false,
        ''
    ),
    (
        '40b757d4-f0a1-44f8-998c-b8d306a787c0',
        12,
        '1.2.5.4',
        'Documentation',
        'Physical Move Services',
        'Inventory & Documentation',
        'Prepare accurate and legible documentation',
        'The contractor shall prepare and retain accurate and legible documentation (written or electronic) which reflects the true condition of all household goods. Documentation shall include, but is not limited to, weight estimates, inventory sheets, warehouse receipt, warehouse exception sheets, pickup and delivery confirmations, certified weight tickets, entitlement and any changes to such, customer notifications, record of loss and damage, claims, and record of all correspondence between contractor and customer.',
        false,
        ''
    ),
    (
        'ea602740-08aa-4392-b42c-65e674a61a92',
        13,
        '1.2.6.1',
        'Inventory',
        'Physical Move Services',
        'Inventory & Documentation',
        'Properly prepare inventories electronically and provide hard copies when requested',
        'The contractor shall prepare all shipment inventories in accordance with International Organization for Standardization (ISO) Standard 17451-1. The contractor shall separately weigh or cube and annotate Professional Books, Papers & Equipment (PBP&E); Organizational Clothing and Individual Equipment (OCIE); and required medical equipment in accordance with government regulations. The contractor''s IT system shall be utilized 100% of the time to electronically collect inventory and condition information for each customer''s shipment. Serial numbers from electronics, major appliances, firearms, and other items shall be scanned for accurate collection and documentation. If requested, a hard paper copy inventory shall be provided to the customer.',
        false,
        ''
    ),
    (
        '82672078-23a8-45fd-8fd5-1c880a2f158a',
        14,
        '1.2.6.7.2',
        'Transfer of Custody',
        'Physical Move Services',
        'Inventory & Documentation',
        'Conduct transfer of custody and retain documents',
        'When custody of a shipment is transferred to or from the contractor to another contractor, the contractor transferring custody shall furnish the contractor receiving custody with two (2) legible duplicate copies of the shipment inventory. A joint inspection shall be performed at any point liability for shipment transfers to or from the contractor and another service provider or the customer at no cost to the government. In the event a difference of opinion arises between the contractor and the receiving party regarding shortage, overage, or the condition of any element of the inventory, the contractor shall annotate such discrepancies accordingly. If no new damage or loss is discovered, the inspection documents shall state ''no differences noted.'' The absence of any annotation beside an inventory item denotes that the container and items were received in good condition. The contractor shall sign and date the completed inspection documents, obtain a signature from and provide a completed copy to the receiving party, and retain a copy for the customer''s file.',
        false,
        ''
    ),
    (
        '369e36c7-5b79-4213-9560-c6e0f8098de5',
        15,
        '1.2.6.11',
        'Weight Tickets',
        'Physical Move Services',
        'Inventory & Documentation',
        'Obtain proper certified weight tickets',
        'The contractor shall obtain certified, legible, and unaltered weight tickets for each shipment or piece of a shipment if transported separately by weighing on a certified weight scale as defined in the CFR Title 49, Part 375.103. Weighing shall be conducted as defined in the CFR Title 49, Part 375.509 and shall comply with all applicable local, state, federal, and foreign country laws.  The contractor shall retain all weight tickets, and make the information contained therein available to the customer and the government.  All weight tickets shall be certified by the weigh master, and shall contain name and location of scale, date, all weight entries (tare, gross and net weights), task order number, and bill of lading number. All invoices presented to collect any shipment charges dependent on the weight transported shall be accompanied by true copies of all weight tickets obtained in the determination of the shipment weight. For partial NTS shipment release, the contractor shall provide certified weight tickets to the NTS service provider and the government. When an NTS shipment is released from storage, all invoices shall be based on the lowest weight of all weight tickets for that NTS shipment.  This includes handling, delivery, and reweigh tickets.',
        false,
        ''
    ),
    (
        '0a4222d0-e85b-403a-b3f9-24a17d222aaa',
        16,
        '1.2.6.14',
        'Safeguarding PII for International Shipments',
        'Physical Move Services',
        'Inventory & Documentation',
        'Safeguard customer PII for international shipments',
        'IAW Homeland Security Customs and Border Protection guidance for safeguarding Personally Identifiable Information (PII); the contractor shall ensure its associated port agents, overseas general agents, and other responsible parties do not include the customer''s Social Security Number (SSN); the customer''s rank or grade, the words â€œDOD Personal Property, DOD Shipment or Military Shipment.â€',
        false,
        ''
    ),
    (
        '2ff1116a-17a9-4624-a462-72c3cbbaefc5',
        17,
        '1.2.6.1.2',
        'Inventory Barcoding',
        'Physical Move Services',
        'Inventory & Documentation',
        'Items must be barcoded. Barcodes must be scanned during inventory process.',
        'Each item, crate, and carton shall be affixed with a unique barcoded sticker label and tag number to enable in-app scanning and enhance load and delivery inventory confirmation. No items shall be loaded on the truck unless the item has been barcoded and inventoried. At destination, each item or carton shall be scanned as it is unloaded from the truck to ensure all pieces are delivered.',
        false,
        ''
    ),
    (
        '0220ef66-03f2-419b-a4e4-e6d87c51191c',
        18,
        '1.2.6.2',
        'Organizational Clothing & Individual Equipment',
        'Physical Move Services',
        'Inventory & Documentation',
        'Identify and separate Organizational Clothing & Individual Equipment (OCIE)',
        'OCIE is clothing and equipment issued to customers for use in the performance of duty. It is common for customers to personally purchase items for use in their duties that appear to be OCIE items but are not. These items are commonly referred to as â€œpersonal kitâ€. The contractor shall request that the customer identify personal kit items. The contractor shall separate personal kit items from OCIE for inventory and claims purposes. The contractor shall conduct an inventory of OCIE at pack-out and delivery. The contractor shall identify OCIE as â€œM-PROâ€ on the inventory.',
        false,
        ''
    ),
    (
        '58c9e425-9cfc-4741-8865-f3f8f7efa2df',
        19,
        '1.2.6.12',
        'Automatic Reweigh',
        'Physical Move Services',
        'Inventory & Documentation',
        'Adhere to automatic reweigh requirement',
        'The contractor shall reweigh any shipment or combination of shipments where the customer has been identified as exceeding or being within 10% or closer to their total weight entitlement.',
        false,
        ''
    ),
    (
        '572af70e-3829-4444-af53-f7bd185f3bf8',
        20,
        '1.2.6.12',
        'Reweigh Invoicing',
        'Physical Move Services',
        'Inventory & Documentation',
        'Invoice on the lesser of the weights when a reweigh is performed',
        'When a reweigh is performed, the contractor shall invoice on the lesser of the two weights.  In the event the contractor fails to perform a reweigh, the contractor shall be limited to invoicing at the customer’s remaining total weight entitlement for all shipments or the weight documented on a certified weight ticket(s), whichever is less.',
        false,
        ''
    );

INSERT INTO
    public.pws_violations (
        id,
        display_order,
        paragraph_number,
        title,
        category,
        sub_category,
        requirement_summary,
        requirement_statement,
        is_kpi,
        additional_data_elem
    )
VALUES
    (
        '0adf15e2-220d-4239-9086-8bcf89a711d7',
        21,
        '1.2.6.12',
        'Reweigh Accomodation',
        'Physical Move Services',
        'Inventory & Documentation',
        'Customer and COR Reweigh Accomodation',
        'The contractor shall accommodate the customer or the COR when either party makes a request to witness a reweigh, by providing the location and the date and time in order to give a reasonable opportunity for the interested parties to be present. ',
        false,
        ''
    ),
    (
        '7b7c0fa5-daa8-494c-afcc-30891fa0b777',
        36,
        '1.2.6.6.4',
        'Packing Upholstered Furniture (Shipment Preparation for NTS)',
        'Physical Move Services',
        'Packing/Unpacking',
        'Packing Requirements (Upholstered Furniture)',
        'Upholstered furniture, to include wicker and wood frame with cushions, shall be placed right side up on all legs in suitable containers covered by plastic or paper and secured with tape, shrink wrap or equivalent materials so that nothing touches or presses against the upholstery. Removable cushions shall be packed with the master pieces.',
        false,
        ''
    ),
    (
        'f9438e53-667a-47e4-b87d-8841be4e1eee',
        22,
        '1.2.6.13',
        'Customs Clearance',
        'Physical Move Services',
        'Inventory & Documentation',
        'International shipment clearance, inspections, and certifications requirements',
        'The contractor shall perform all customs clearance, agricultural inspections and certifications, and other related services that pertain to and influence the movement of personal property (gun control, quarantine, pest infestation, etc.) in accordance with all applicable local, state, federal, and foreign country laws and DoD regulations. DoD consignment requirements are in the Personal Property Consignment Instruction Guide (PPCIG). Shipments entering the United States must comply with Title 19, Section 148 of the Code of Federal Regulations.',
        false,
        ''
    ),
    (
        'eab29339-a609-4df1-9f4a-b78ba1bcf7b7',
        23,
        '1.2.6.3',
        'Packing/Loading',
        'Physical Move Services',
        'Packing/Unpacking',
        'Protect against real and personal property damage',
        'The contractor shall prepare, pack, unpack, load, and unload all personal property to protect all real and personal property against loss or damage.',
        false,
        ''
    ),
    (
        'c57b9fd0-7821-4157-b9fc-7fe13f71e7d8',
        24,
        '1.2.6.3.1',
        'Packing Materials',
        'Physical Move Services',
        'Packing/Unpacking',
        'Use appropriate packing materials',
        'The contractor shall provide packing materials that are new or in sound condition, except in the case when the customer has provided original or specially designed packaging that the contractor has inspected and accepted as being as good or in sound condition. When allowed, and if material is not new, all marks pertaining to any previous shipment must be obliterated. The contractor shall use furniture pads or other appropriate materials to wrap or protect all other items not packed in boxes, containers, or cartons. The use of any type of protective material does not reduce the level of contractor liability for any lost or damaged items. New packing material shall be used for mattresses, box springs, linens, bedding, and clothing.',
        false,
        ''
    ),
    (
        '8eee0093-75c6-435e-899b-5ba98cacb41a',
        25,
        '1.2.6.6',
        'Shipment Preparation for Non-Temporary Storage (NTS)',
        'Physical Move Services',
        'Packing/Unpacking',
        'Containerize NTS shipments',
        'The contractor shall prepare and load property going into NTS in containers at residence for shipment to NTS. The contractor shall seal all containers, using tamper-proof seals, at the residence. Power-driven equipment, motorcycles, boats, trailers, over size items, and overstuffed furniture may be shipped uncrated.',
        false,
        ''
    ),
    (
        'e8d36775-bdef-4feb-9324-6c54089b8bc4',
        26,
        '1.2.6.15.1',
        'Unpacking and Re-assembly',
        'Physical Move Services',
        'Packing/Unpacking',
        'Properly unpack and reassemble personal property',
        'Unloading and unpacking at destination includes the one-time laying of rugs and the one-time placement of furniture and like items in a room or dwelling designated by the customer or their representative. All articles disassembled by the contractor or originating from NTS shall be reassembled. If hardware is missing, the contractor shall obtain appropriate hardware to reassemble. On a one-time basis, all barrels, boxes, cartons, and crates shall be unpacked (upon request) and the contents placed in a room designated by the customer. This includes the placement of articles in closets, cabinets, cupboards, or on shelving in the kitchen when convenient and consistent with safety of the article(s) and proximity of the area desired by the customer, but does not include arranging the articles in a manner desired by the customer.',
        false,
        ''
    ),
    (
        '68f072aa-269c-4951-811e-a7eb22857deb',
        27,
        '1.2.6.15.3',
        'Debris removal',
        'Physical Move Services',
        'Packing/Unpacking',
        'Remove debris from residence',
        'All debris incident to the packing, unpacking, loading, or unloading of the delivered shipment shall be removed on the date(s) of delivery, unless otherwise waived by the customer.',
        false,
        ''
    ),
    (
        '4dcfc166-ce57-4c6d-b361-5dc9dcde2698',
        28,
        '1.2.6.5',
        'Restricted Items',
        'Physical Move Services',
        'Packing/Unpacking',
        'Improper acceptance of restricted items',
        'The contractor shall not knowingly provide service for any item defined as restricted by law, policy or agency of the U.S. Government or any foreign entity in an international point-to-point move.',
        false,
        ''
    ),
    (
        '8efbb25a-7182-4123-b739-61af57c33168',
        29,
        '1.2.6.3.3',
        'Packing/Unpacking Unaccompanied Baggage',
        'Physical Move Services',
        'Packing/Unpacking',
        'Unaccompanied Baggage (UB) packing / unpacking done in accordance with JTR. UB rates apply.',
        'Unaccompanied baggage packing and unpacking rates shall be used when the task order includes unaccompanied baggage, and for all items being transported under an unaccompanied baggage rate. Unaccompanied baggage shall be packed and unpacked in accordance with the Joint Travel Regulation.',
        false,
        ''
    );

INSERT INTO
    public.pws_violations (
        id,
        display_order,
        paragraph_number,
        title,
        category,
        sub_category,
        requirement_summary,
        requirement_statement,
        is_kpi,
        additional_data_elem
    )
VALUES
    (
        '7f934ba0-2322-4642-99d7-359ba413882a',
        30,
        '1.2.6.6.1',
        'Cartons & Packing Material (Shipment Preparation for NTS)',
        'Physical Move Services',
        'Packing/Unpacking',
        'Cartons and other packing material must meet contract specifications. (Shipment Preparation for NTS)',
        'All cartons and wrapping material shall be in new or sound condition and adequate for the use employed. New packing material shall be used for mattresses, box springs, linens, bedding, and clothing. After packing, cartons shall be closed and sealed by taping lengthwise at all joints. Cartons shall have a minimum average bursting strength of 200 pounds per square inch and dish packs shall have a minimum average bursting strength of 350 pounds per square inch. Cartons shall be stacked in an upright position to minimize crushing, with the exception of mattress cartons. Plastic containers (tote or similar) and similar types of containers shall not be used. However, if items are packed by the customer in plastic or similar type containers, the contractor may pack these containers in an approved carton if a carton is available that will accommodate the container. If the plastic container cannot be packed in an approved carton, the contractor shall empty and pack the contents into an appropriate approved carton.',
        false,
        ''
    ),
    (
        'b41a9f7b-07ed-481b-8c59-0e26f442aca5',
        31,
        '1.2.6.6.3',
        'Mattresses & Box Springs (Shipment Preparation for NTS)',
        'Physical Move Services',
        'Packing/Unpacking',
        'Mattresses and box springs must be placed in cartons and sealed. (Shipment Preparation for NTS)',
        'All mattresses and box springs, except those in hide-a-beds or sofa beds, shall be placed in cartons and completely sealed.',
        false,
        ''
    ),
    (
        '30a4e886-a4f6-4cd8-85b7-bf72d52b3b2e',
        32,
        '1.2.6.6.5',
        'Rugs, Rug Pads, Carpet (Shipment Preparation for NTS)',
        'Physical Move Services',
        'Packing/Unpacking',
        'Rugs, rug pads, carpet must be rolled, covered and taped. (Shipment Preparation for NTS)',
        'All rugs, rug pads and carpets shall be properly rolled (not folded) and covered by paper and secured with tape or equivalent materials.',
        false,
        ''
    ),
    (
        '0e596cb5-fe04-4214-9644-8432c0aeaf72',
        33,
        '1.2.6.6.8',
        'Appliance packing rules (Shipment Preparation for NTS)',
        'Physical Move Services',
        'Packing/Unpacking',
        'Nothing shall be packed in appliances excepting integral parts. (Shipment Preparation for NTS)',
        'Nothing shall be packed in washers, dryers, refrigerators, freezers, stoves, or other major appliances except such items as electrical cords, connecting hoses and similar items that are required as an integral part of the appliance in its normal operation.',
        false,
        ''
    ),
    (
        '6a556fac-8026-4b60-b2f9-182a84fa1afc',
        34,
        '1.2.6.3',
        'Disassembly / Reassembly',
        'Physical Move Services',
        'Packing/Unpacking',
        'Properly disassemble and reassemble of original pieces',
        'The contractor shall disassemble items only to the extent necessary for shipment and the contractor shall be responsible for subsequent reassembly of all original pieces.',
        false,
        ''
    ),
    (
        '7e8d80b2-0ff9-4156-8cc4-3dd639eb87e2',
        35,
        '1.2.6.6.2',
        'Packing Linens, Clothing, Bedding, etc. (Shipment Preparation for NTS)',
        'Physical Move Services',
        'Packing/Unpacking',
        'Packing Requirements (Linens, Clothing, etc.)',
        'Linens, towels, bedding, draperies, and other items of this type shall be packed into wardrobe type cartons and shall be completely sealed. Clothing shall not be stored in closet bags. Hangers shall be removed from clothing packed in flat wardrobes.  ',
        false,
        ''
    ),
    (
        'ae531db9-2c00-4843-81c1-b868f8e025ab',
        37,
        '1.2.6.6.7',
        'Removal of Items from Drawers, Hampers, Bureaus (Shipment Preparation for NTS)',
        'Physical Move Services',
        'Packing/Unpacking',
        'All articles must be removed from chests of drawers, hampers, bureaus',
        'All articles shall be removed from chests of drawers, bureaus, clothes hampers, and other similar items.',
        false,
        ''
    ),
    (
        'cb85bfc6-21c7-4ba5-9664-37bb5da829cf',
        38,
        '1.2.6.6.9',
        'Power-Driven Equipment',
        'Physical Move Services',
        'Packing/Unpacking',
        'Proper handling of power-driven equipment',
        'The contractor shall verify that power-driven equipment, boats and motorcycles have been drained of all gasoline, the cables disconnected from the battery terminals, and the cable ends secured and protected with electrical tape.  Batteries may be shipped with the power-driven equipment. The contractor shall verify boat drain plugs have been removed and if not permanently attached to the boat, placed in a cloth bag and tied to the boat. Motorcycle keys shall remain in the customer’s file to facilitate handling and movement.',
        false,
        ''
    ),
    (
        '4fda8337-c645-4cec-a4fe-d8a2844ab35c',
        39,
        '1.2.5.3.2',
        'Changes to Schedule',
        'Physical Move Services',
        'Shipment Schedule',
        'Accommodate changes to schedule',
        'The contractor shall accommodate all requests for a change of schedule that are received prior to delivery. Examples [not all inclusive] for changes to schedule may include: termination of shipment, rescheduling of pickup or delivery dates, diversion of shipment to a different destination, more than one pickup location for a shipment, more than one delivery location for a shipment.',
        false,
        ''
    ),
    (
        '8d8847db-9cb5-4ef9-bbd2-e226df840224',
        40,
        '1.2.5.3.3',
        'Cancellations',
        'Physical Move Services',
        'Shipment Schedule',
        'Accommodate shipment cancellation',
        'The contractor shall accommodate shipment cancellation up to the day of scheduled packing or pickup without cost or obligation to the government, provided packing has not begun.',
        false,
        ''
    );

INSERT INTO
    public.pws_violations (
        id,
        display_order,
        paragraph_number,
        title,
        category,
        sub_category,
        requirement_summary,
        requirement_statement,
        is_kpi,
        additional_data_elem
    )
VALUES
    (
        '8510cf69-4632-413a-bdbf-94499887de3e',
        41,
        '1.2.5.3.4',
        'Diversions (D)/Terminations (T)/Reshipment (R)',
        'Physical Move Services',
        'Shipment Schedule',
        'Properly process diversions',
        'The contractor shall process all D/T/R based on the location of the shipment when notified, and will invoice IAW PWS, Appendix F, for services completed. Subsequent movement of the shipment(s) will be made in the most cost-effective manner based on the date and location of the shipment when the order modification was received.',
        false,
        ''
    ),
    (
        'ef7282b4-f3a3-42e4-9ef6-68e86b88f045',
        42,
        '1.2.2.3',
        'Shipment ITV',
        'Physical Move Services',
        'Shipment Schedule',
        'Provide In-Transit Visibility',
        'The contractor''s IT system shall provide geofencing tracking services within its mobile application that provides customers with real-time location tracking of the shipment once the supplier is within at least ten (10) miles of the customer''s residence. Outside of the at least ten (10) mile radius in which geofencing tracking services will be provided, the contractor''s IT system shall provide detailed point-to-point status updates including arrival and departure time and estimated arrival at destination no later than one (1) Government Business Day (GBD) from a shipment location or status change.',
        false,
        ''
    ),
    (
        '5842a049-7c4d-43ee-a37d-9ab80e1dce4e',
        43,
        '1.2.6.7',
        'Pickup',
        'Physical Move Services',
        'Shipment Schedule',
        'Pick up shipment on time',
        'The contractor shall pickup all pieces of a shipment on the scheduled pickup date. The shipment is not considered an on-time pickup if the contractor changes the date at any time without approval of the customer.',
        true,
        'observedPickupDate'
    ),
    (
        '0976f6f9-fda2-4e66-b69c-7cb413228a2a',
        44,
        '1.2.6.8',
        'Hours of Operation',
        'Physical Move Services',
        'Shipment Schedule',
        'Adhere to Hours of Operation',
        'The contractor shall not begin pickup or delivery at the customer’s residence before 0800 hours or after 1700 hours without prior approval of the customer or the government. The contractor shall provide information to the customer and the government on the afternoon preceding the scheduled pickup or delivery as to whether the service will be performed in the morning (0800 to 1200) or in the afternoon (1200 to 1700) of the following day.  The contractor shall not begin any service that will not allow completion by 2100 hours without prior approval of the government.  Shipments shall not be scheduled for pickup or delivery on Non-Government Business Days, U.S. Federal holidays, or foreign national holidays unless there is a mutual agreement between the government and the contractor.  Unless otherwise stated, all references to “days” are government business days (GBD).  IAW the DTR, a GBD is defined as a business day (i.e. Monday through Friday) that is not a federal holiday.',
        false,
        ''
    ),
    (
        '74c62377-4b76-489b-a2f5-5e2c40f343dd',
        45,
        '1.2.6.9',
        'Transport',
        'Physical Move Services',
        'Shipment Schedule',
        'Deliver by Required Delivery Date',
        'The contractor shall transport shipments, including non-standard shipments (Appendix D); from origin to destination so as to ensure delivery by the RDD as determined by the domestic and international (to include household goods shipments and unaccompanied baggage shipments) transit times located in Appendix C, attached hereto.',
        true,
        'observedDeliveryDate'
    ),
    (
        '4296c046-24cb-44d8-8404-e26517b2abf0',
        46,
        '1.2.6.15',
        'Delivery',
        'Physical Move Services',
        'Shipment Schedule',
        'Deliver by Required Delivery Date',
        'The contractor shall deliver and unload all pieces of a shipment as scheduled by the RDD.',
        true,
        'observedDeliveryDate'
    ),
    (
        'd10315b9-c63f-47d5-8176-fe1a73da61c9',
        47,
        '1.2.8.2',
        'Scheduling Notifications',
        'Physical Move Services',
        'Shipment Schedule',
        'Notify customer of scheduled dates',
        'The contractor shall notify the customer of all scheduled dates as soon as known for counseling, packing, unpacking, pickup, delivery, and all other dates for which interaction with the contractor by the customer is required.',
        false,
        ''
    ),
    (
        '73a5b32b-131f-4fe2-85bf-6b0636128619',
        48,
        '1.2.8.5',
        'Inbound Shipment Notification',
        'Physical Move Services',
        'Shipment Schedule',
        'Notify customer at least 24 hours in advance of delivery',
        'The contractor shall notify and confirm with the customer no later than twenty-four (24) hours in advance of shipment delivery. The contractor shall not deliver a customer’s personal property to SIT without customer approval unless the contractor has documented two (2) unsuccessful attempts to contact the customer. Each attempt must document a proposed First Available Delivery Date (FADD). The attempts must be made at least eight (8) hours apart, and no later than twenty-four (24) hours in advance of the proposed FADD.',
        false,
        ''
    ),
    (
        '7eaa7756-6c8f-46df-a211-01e9f6c0ebfe',
        49,
        '1.2.8.6',
        'Quality Assurance Schedule Notification',
        'Physical Move Services',
        'Shipment Schedule',
        'Provide QA Forecast',
        'On a daily basis, NLT 0800 local installation time, the contractor shall provide the government a rolling 30-day Shipment Schedule containing the schedule for all shipments being packed, picked-up, or delivered for the purposes of scheduling government QAE. The report shall contain all dates, shall be filterable by installation, city, county, state, country, and shall contain the address location of the origin or destination activities to be observed. Any direct deliveries scheduled for same day, the contractor shall make notification of delivery within one (1) hour to the destination activity.',
        false,
        ''
    ),
    (
        'bea64c6e-e6cc-4965-be17-1c67412d26a0',
        50,
        '1.2.5.3.6.',
        'Shipment Suitability',
        'Physical Move Services',
        'Shipment Schedule',
        'Determine Shipment Suitability',
        ' If, prior to pick up, the shipment is determined to be in a condition that makes it likely to permeate, contaminate, or otherwise cause damage to other HHGs or equipment, the contractor shall coordinate with the GSR as soon as the condition is identified.',
        false,
        ''
    );

INSERT INTO
    public.pws_violations (
        id,
        display_order,
        paragraph_number,
        title,
        category,
        sub_category,
        requirement_summary,
        requirement_statement,
        is_kpi,
        additional_data_elem
    )
VALUES
    (
        '6cf22021-e914-4509-9610-7df292e5caaa',
        51,
        '1.2.8.5.',
        'Missed RDD Notification',
        'Physical Move Services',
        'Shipment Schedule',
        'Notify customer of missed RDD and revised RDD',
        'If an inbound shipment is projected to fail to meet the firm RDD agreed to during counseling or scheduling, the contractor shall notify the customer at the earliest practicable time or no later than one (1) day and provide a revised RDD. ',
        false,
        ''
    ),
    (
        '782748b0-beae-43b3-a7b3-11141d870b58',
        52,
        '1.2.6.12',
        'Reweighs',
        'Physical Move Services',
        'Shipment Weights',
        'Conduct reweigh when requested',
        'When requested by the customer or the COR, the contractor shall conduct a reweigh before the actual commencement of unloading for delivery.',
        false,
        ''
    ),
    (
        '6342f7ce-76f8-45fb-9fe6-0a374094d0ac',
        53,
        '1.2.8.3',
        'Weight Notifications',
        'Physical Move Services',
        'Shipment Weights',
        'Notify customer and government of shipment weight in a timely fashion',
        'The contractor shall notify the customer and the government of the actual weight of each shipment within one (1) GBD of shipment pickup or prior to delivery or placement into SIT, whichever is earlier.',
        false,
        ''
    ),
    (
        '842a8923-6602-479f-9cca-c77adeb1365a',
        54,
        '1.2.8.4',
        'Excess Cost Notifications',
        'Physical Move Services',
        'Shipment Weights',
        'Notify customer of excess cost',
        'If a customer is at risk for excess costs based on any shipment or combination of shipments exceeding or being within 10% or closer to their total weight entitlement or any other entitlement, the contractor shall notify the customer within one (1) day of discovery. Notification shall include that the customer is responsible for any excess costs that may be incurred, provide an estimated excess cost amount, and obtain written acknowledgment from the customer',
        false,
        ''
    ),
    (
        'a158499f-1dc9-488f-a0d5-e46698546d94',
        55,
        '1.2.6.15.4',
        'Storage-in-Transit (SIT)',
        'Physical Move Services',
        'Storage',
        'Request ordering and payment of storage in accordance with SIT eligibility window',
        'The contractor''s period of SIT eligibility begins on the First Available Delivery Date (FADD) and ends by the 5th day after the requested delivery date from storage or the actual delivery date, whichever is earlier',
        false,
        ''
    ),
    (
        '1b9630a4-9400-424c-a05a-caa44be49a98',
        56,
        '1.2.6.16',
        'Storage',
        'Physical Move Services',
        'Storage',
        'Provide adequate storage facilities',
        'The contractor shall provide warehouse storage facilities to accommodate SIT as required in accordance with all local, state, federal, and country fire, safety and construction codes, standards and ordinances, ensuring that all stored shipments are adequately protected. For SIT facilities residing in a multi-occupancy structure, the SIT provider''s storage area will be separated from other occupants of the building by a firewall or partition having a fire resistance rating sufficient to protect the warehouse from the fire exposure of the other occupant. The minimum separation shall be a solid wall or partition, without windows, doors or other openings, having a fire resistance rating of not less than one hour. The construction, upkeep, purchase, lease or rental of any commercial structure, land, or equipment for the storage facility shall be the responsibility of the contractor. All SIT facilities shall maintain at least an operational Class 3 supervised detection and reporting system. All facilities shall meet all requirements for insurance rate credit by the Insurance Services Office (ISO) or other cognizant fire insurance rating organization for an other than wood frame or pole building and shall provide a fire wall separation resistance rating sufficient to protect the warehouse from the fire exposure of another occupant. If host country standards, practices, or customs conflict with SIT standards, exceptions may be granted by the Government Representative. All storage facilities shall be located above the 100-year flood plain for the area.',
        false,
        ''
    ),
    (
        'baa74ce7-e4e2-4fd8-936e-ca34a2c07fdd',
        57,
        '1.2.6.16.1',
        'Shipment Hostage',
        'Physical Move Services',
        'Storage',
        'No contractor entity shall hold a shipment hostage',
        'The contractor and all subcontractors performing services under this contract acknowledge that holding shipments hostage is a violation of USC Title 37, Section 453, at subparagraph (c)(5) which provides, ''No carrier, port agent, warehouseman, freight forwarder, or other person involved in the transportation of property may have a lien on, or hold, impound, or otherwise interfere with the movement of baggage and household goods being transported under this section.â€',
        false,
        ''
    ),
    (
        '8455abfd-b33d-4c81-afb0-0510964afeae',
        58,
        '1.2.8.7.1',
        'Advance Notice of SIT Expirations & Extensions',
        'Physical Move Services',
        'Storage',
        'Provide notice of SIT expiration',
        'Thirty (30) days prior to expiration of any SIT entitlement, the contractor shall provide the customer written notification via traceable means of the upcoming expiration and seek a desired disposition from the customer. The notification shall include, at a minimum, the exact date responsibility for storage charges and fees transfers to the customer, all costs and fees the customer can expect to incur, and changes in insurance coverage.',
        false,
        ''
    ),
    (
        '1f00c087-8d12-4252-812e-9e484f13ef7a',
        59,
        '1.2.8.7.2',
        'SIT Extension',
        'Physical Move Services',
        'Storage',
        'Request SIT extension for customer',
        'If the customer requests delivery after the SIT expiration date or requests an extension of storage, the contractor shall prepare and submit a written request to extend the storage period at government expense. Upon receipt of an approved request to extend a customer’s storage period, the contractor shall update the shipment to reflect the storage period extension and notify the customer of the new expiration date.',
        false,
        ''
    ),
    (
        'ab9531cd-d56e-44fe-b68e-98801b1cfe45',
        60,
        '1.2.8.7.3',
        'Conversion to Customer''s Expense',
        'Physical Move Services',
        'Storage',
        'Properly convert storage lot to customer expense',
        'Upon approval by the Government and expiration of the SIT entitlement, the contractor shall consider the DoD customer''s property converted to customer''s expense. Once converted, the contractor shall provide the customer written notification by traceable means, within five (5) days from the date their account converted to customer''s expense.',
        false,
        ''
    );

INSERT INTO
    public.pws_violations (
        id,
        display_order,
        paragraph_number,
        title,
        category,
        sub_category,
        requirement_summary,
        requirement_statement,
        is_kpi,
        additional_data_elem
    )
VALUES
    (
        'bcdbc5b7-d9a7-427a-a3a5-641d3d9c2283',
        61,
        '1.2.8.7.4',
        'Disposition of Converted Shipments',
        'Physical Move Services',
        'Storage',
        'Properly dispose of storage lot',
        'The contractor shall seek authorization from the customer by way of a notarized authorization to dispose of the property. If authorization is not obtained, the contractor shall follow all applicable local, state and federal laws when disposing of lots converted to customer''s expense.',
        false,
        ''
    ),
    (
        '2c49c890-3ff1-4fc6-8f13-04bce33283f2',
        62,
        '1.2.6.16',
        'Prevent Exposure of SIT Shipments',
        'Physical Move Services',
        'Storage',
        'Prevent exposure of SIT shipments to harm',
        'The contractor shall prevent exposure of all shipments to vermin, dust, mold, mildew, moisture, hazardous chemicals, as well as prevent exposure to extreme heat, cold, humidity, and direct sunlight.',
        false,
        ''
    ),
    (
        '3cb0182b-97dd-4644-8528-3f0f856bfeb2',
        63,
        '1.2.6.15.2',
        'NTS Pickup Coordination',
        'Physical Move Services',
        'Storage',
        'Coordinate NTS pickup with NTS provider',
        'The contractor shall coordinate the pickup of NTS shipments from a storage facility with the NTS provider in order to meet the RDD.',
        false,
        ''
    ),
    (
        '1589f130-3cdb-48b2-b927-ceeec32ba9bb',
        64,
        '1.2.6.6.6',
        'NTS Firearm Requirements',
        'Physical Move Services',
        'Storage',
        'Identify Firearms to NTS provider',
        'All firearms shall be identified to the NTS provider upon delivery to the storage facility.',
        false,
        ''
    ),
    (
        'ad27a537-e201-48ca-97ec-e07790858a6f',
        65,
        '1.2.6.7.1',
        'NTS Shipment Pickup/Delivery Requirements',
        'Physical Move Services',
        'Storage',
        'Properly pack, pickup, deliver NTS Shipments',
        'The contractor shall be responsible for packing, pickup, and delivery of NTS shipments.',
        false,
        ''
    ),
    (
        'ac2b2a43-704d-4393-8dfb-66f062eadfb6',
        66,
        '1.2.6.7.1',
        'NTS Warehouse Coordination',
        'Physical Move Services',
        'Storage',
        'NTS Warehouse Location Coordination',
        'The contractor shall coordinate with the government to determine the warehouse location for each shipment going into NTS.',
        false,
        ''
    ),
    (
        '3118597d-32ad-4ef3-8589-377d35ff0857',
        67,
        '1.2.8.7.1.',
        'Joint Inspection - Storage Expiration',
        'Physical Move Services',
        'Storage',
        'Provide Documentation of Joint Inspection, Notifications, and Correspondence Re: SIT Expiration',
        'The contractor shall provide a copy of the Joint Inspection, upon request, and shall retain a copy of all notifications and correspondence in the customer’s file.',
        false,
        ''
    ),
    (
        '6143d59c-368c-4053-9cd6-ba2a8bdf7b99',
        68,
        '1.2.8.7.1.',
        'SIT Expiration Notification (15 Days)',
        'Physical Move Services',
        'Storage',
        'Provide customer and government notification of SIT entitlement expiration (15 days)',
        'If desired disposition is not obtained, at fifteen (15) days prior to expiration of SIT entitlement the contractor shall repeat the above notification to the customer and OO.',
        false,
        ''
    ),
    (
        '2cf57923-2fb6-46d0-9807-0487d3f0120f',
        69,
        '1.2.5.3.5',
        'Installation Scheduling',
        'Physical Move Services',
        'Workforce/Sub-Contractor Management',
        'Adhere to military facility / base access guidelines',
        'The contractor shall schedule all pickups or deliveries in accordance with specific installation or facility requirements. Any delay due to personnel disqualification from specific installation access or failure to follow published access guidelines shall be considered an unacceptable delay.',
        false,
        ''
    ),
    (
        '80bd56ce-4f43-4ac5-bbbe-36092cc5ea77',
        70,
        '1.2.1.1',
        'Background Checks/Records',
        'Physical Move Services',
        'Workforce/Sub-Contractor Management',
        'Conduct background checks',
        'Prior to engaging in any services identified herein, the contractor shall ensure a background check is conducted (at contractor expense); IAW industry standard, for all persons performing under this contract whose role involves interacting with a customer or handling or transporting shipments. The contractor shall provide employment records to Government upon request, to the extent allowed by law.',
        false,
        ''
    );

INSERT INTO
    public.pws_violations (
        id,
        display_order,
        paragraph_number,
        title,
        category,
        sub_category,
        requirement_summary,
        requirement_statement,
        is_kpi,
        additional_data_elem
    )
VALUES
    (
        'ed6b10f1-2e3e-4621-83a0-87c79cb63d34',
        71,
        '1.2.1.1',
        'Background Checks/Records',
        'Physical Move Services',
        'Workforce/Sub-Contractor Management',
        'Government disqualification of employee',
        'The government has the right to prevent certain employees from performing under the contract due to an unfavorable background check.',
        false,
        ''
    ),
    (
        'dfaa8ffe-fd32-41c7-9654-fd73fd35d6f1',
        72,
        '1.2.1.2',
        'Workforce Requirements',
        'Physical Move Services',
        'Workforce/Sub-Contractor Management',
        'Personnel are trained in assigned duties',
        'The contractor shall ensure all employees remain trained and qualified in their assigned duties.',
        false,
        ''
    ),
    (
        'ff36f06d-a8d9-48ef-a176-f35ebe75dced',
        73,
        '1.2.1.2',
        'Workforce Requirements',
        'Physical Move Services',
        'Workforce/Sub-Contractor Management',
        'Prohibition of smoking',
        'Smoking is prohibited in the customer''s residence or within 50 feet of personal property during all phases of shipment and storage.',
        false,
        ''
    ),
    (
        '15ea4a6c-fb74-4d20-a100-f83f93855a31',
        74,
        '1.2.1.2',
        'Workforce Requirements',
        'Physical Move Services',
        'Workforce/Sub-Contractor Management',
        'Possession or Under Influence of Drugs or Alcohol',
        'The contractor shall ensure all employees and sub-contractors are free from possession of and not under the influence of drugs or alcohol while in a customer''s residence or handling a customer''s personal property.',
        false,
        ''
    ),
    (
        'e5ee6bb9-6ff2-43e0-a33e-7fc4fc4f1b56',
        75,
        '1.2.1.2.1',
        'Defense Personal Property Program (DP3) Performance History',
        'Physical Move Services',
        'Workforce/Sub-Contractor Management',
        'Use of disqualified entity (3 years)',
        'Contractor shall ensure no entity that has been disqualified or revoked from DP3 within three (3) years of move execution date will perform work under this contract.',
        false,
        ''
    ),
    (
        '6f99065d-95d1-4297-8e58-b0f2fb9b2672',
        76,
        '1.2.1.3',
        'Customer Interaction',
        'Physical Move Services',
        'Workforce/Sub-Contractor Management',
        'At least one crewmember fluent in English language',
        'At least one crewmember or warehouse employee, where applicable, shall be fluent in English for the purposes of customer interaction.',
        false,
        ''
    ),
    (
        '73120c54-1b89-4567-a48f-7fafa7920acc',
        77,
        '1.2.1.3',
        'Customer Interaction',
        'Physical Move Services',
        'Workforce/Sub-Contractor Management',
        'Appearance of contractor personnel',
        'All personnel shall be clean and neat and be easily identifiable as company personnel.',
        false,
        ''
    ),
    (
        'f6e80377-191f-4e20-8c8f-b26d6e2195a0',
        78,
        '1.2.1.3',
        'Customer Interaction',
        'Physical Move Services',
        'Workforce/Sub-Contractor Management',
        'Replacement of employees exhibiting unprofessional behavior',
        'The contractor shall replace any individuals exhibiting unprofessional behavior, when requested by the customer or a government representative.',
        false,
        ''
    ),
    (
        'b363cbee-9c4b-4096-9c9b-c85f7b4707b5',
        79,
        '1.2.1.4',
        'Driver Identification/Qualification Requirements',
        'Physical Move Services',
        'Workforce/Sub-Contractor Management',
        'Use of qualified drivers',
        'The contractor shall ensure all drivers who perform under this contract are qualified and licensed in accordance with local, state, federal, and foreign country or international laws.',
        false,
        ''
    ),
    (
        '97518b1e-df0e-4f8d-b576-c09a5764d55e',
        80,
        '1.2.2.3',
        'Crew Photographs',
        'Physical Move Services',
        'Workforce/Sub-Contractor Management',
        'Provide Crew Photographs and TSP Contacts',
        'The contractor’s IT system shall also provide customers with current photographs of the crew assigned to each move prior to crew arrival and shall provide customers with the capability to contact the service provider directly.',
        false,
        ''
    );

INSERT INTO
    public.pws_violations (
        id,
        display_order,
        paragraph_number,
        title,
        category,
        sub_category,
        requirement_summary,
        requirement_statement,
        is_kpi,
        additional_data_elem
    )
VALUES
    (
        '0385aa57-15b8-4270-99e2-6619a75ed221',
        81,
        '1.2.1.2.',
        'Base Access',
        'Physical Move Services',
        'Workforce/Sub-Contractor Management',
        'Ensure Base Access',
        'The contractor shall ensure all persons interacting with customers under this contract on and off base meet the specific requirements for local installation access as listed in DoD Manual 5200.08.',
        false,
        ''
    ),
    (
        'fdcde2ff-592f-4c00-ab3e-c6c4147abc2b',
        82,
        '1.2.1.2.',
        'English Language',
        'Physical Move Services',
        'Workforce/Sub-Contractor Management',
        'Ensure Use of English',
        'English shall be the only language used with regard to this contract for written correspondence, discussions and other business transactions.',
        false,
        ''
    ),
    (
        '51d5b062-44a4-4d17-9827-45aa8933a2df',
        83,
        '1.2.1.3.',
        'Appearance of Contractor Personnel',
        'Physical Move Services',
        'Workforce/Sub-Contractor Management',
        'Neat and clean personnel',
        'All personnel shall be clean and neat and be easily identifiable as company personnel.',
        false,
        ''
    ),
    (
        '1261c17d-5229-4004-a17c-ed7765c7d491',
        84,
        '1.2.7.2.2',
        'Claims Settlement',
        'Liability',
        'Loss & Damage',
        'Respond to claims in a timely fashion',
        'The contractor shall pay, deny, or make an offer on all claims valued at $1000 or less within 30 calendar days of receipt of the claim and of all other claims within 60 calendar days of receipt IAW the Claims and Liability Rules (Appendix E).',
        true,
        'observedClaimsResponseDate'
    ),
    (
        '9c95459e-009f-434f-a4ae-9a2079370625',
        85,
        '1.2.7.2.1',
        'Scope of Liability',
        'Liability',
        'Loss & Damage',
        'Liability for Full Replacement Value',
        'The contractor shall be liable for all loss or damage up to Full Replacement Value (FRV) for all shipments from the point of origin to the point custody transfers to a customer as defined in the Claims and Liability Rules (Appendix E). For the contractor to claim any exemptions, contractor must prove it was free from negligence. The contractor accepts full responsibility for performance of all of its employees, subcontractors, and agents. In the event of any damage to public or private property from acts or omissions of persons performing under this contract, the contractor shall immediately repair and correct damages at contractor''s expense.',
        false,
        ''
    ),
    (
        'f506da8c-a6d0-4629-a6a6-6a79b4ab2588',
        86,
        '1.2.6.16.13',
        'Mold Remediation',
        'Liability',
        'Loss & Damage',
        'Payment of mold remediation services',
        'Services for mold remediation will normally be at the expense of the contractor, however, service payments may be authorized when the Government determines the mitigating contractor is not liable for the damage. Contractor shall request the service authorization from the local Ordering OO.',
        false,
        ''
    ),
    (
        '96c9b54d-ac2e-4424-998e-a0a56fae5cc9',
        87,
        '1.2.6.16.11',
        'Customer Elects to Inspect Remediated Items',
        'Liability',
        'Loss & Damage',
        'Facilitate customer inspection and acceptance of remediated items',
        'If the customer does not accept the remediation on any item during the inspection, that item shall be separated from the accepted items. If the contractor agrees with the customer that those items are unacceptable, the contractor shall deliver the accepted items and process claims on the unacceptable items for compensation at FRV.',
        false,
        ''
    ),
    (
        'b78c6a76-c803-45e9-ab1d-65d46c96f6e6',
        88,
        '1.2.6.16.10',
        'Delivery of Remediated Items',
        'Liability',
        'Loss & Damage',
        'Provide notification and delivery of remediated items',
        'Before delivery, contractor shall notify the customer and destination QAE or COR that the items have been remediated, are ready for delivery, and provide a reasonable opportunity to inspect the remediated items before delivery begins.',
        false,
        ''
    ),
    (
        'b61528eb-52a5-4ad2-bbce-3b81e7d4cf90',
        89,
        '1.2.6.16.8',
        'Shipment Inspection',
        'Liability',
        'Loss & Damage',
        'Inspect and remove items of sentimental value',
        'The contractor shall offer the customer an opportunity to inspect the shipment and remove items of sentimental or special value at the owner''s discretion in coordination with the responsible QAE or COR.',
        false,
        ''
    ),
    (
        '0cd3c46f-6b4f-4af4-9af7-df5386a81088',
        90,
        '1.2.6.16.4',
        'Possible Contamination',
        'Liability',
        'Loss & Damage',
        'Notify COR of contaminated containers',
        'The contractor shall contact the responsible contracting officer representative (COR) when containers show signs of possible contamination, for example water saturation or mold growth on the exterior. The contractor shall be responsible for arranging for all testing and mitigation. If testing determines mold is present, the contractor shall contact the servicing MCO and the responsible OO for guidance. If mold is suspected, the contractor shall notify the customer, the servicing Military Claims Officer (MCO), and the responsible Ordering Officer (OO) who will authorize the appropriate testing. ',
        false,
        ''
    );

INSERT INTO
    public.pws_violations (
        id,
        display_order,
        paragraph_number,
        title,
        category,
        sub_category,
        requirement_summary,
        requirement_statement,
        is_kpi,
        additional_data_elem
    )
VALUES
    (
        'dae0eb63-9d59-47df-9175-7e5bb39d7e76',
        91,
        '1.2.6.16.3',
        'Damage Mitigation',
        'Liability',
        'Loss & Damage',
        'Takes reasonable steps to reduce damage (Liability exclusions)',
        'In the event a shipment is damaged as a result of any one of the excluded causes listed in Appendix E, para E.3., Exclusions from Liability, the contractor shall take reasonable steps to mitigate the extent of the damage.',
        false,
        ''
    ),
    (
        '495b0508-596d-43b0-b1da-e0a6421222d3',
        92,
        '1.2.6.16.12.',
        'Refusal of Remediated Items',
        'Liability',
        'Loss & Damage',
        'Handling Ccstomer refusal of remediated items and placement in storage',
        'If customers refuse delivery of remediated items after delivery of those items begins, the contractor shall transport those items to a storage facility at the contractor’s discretion.',
        false,
        ''
    ),
    (
        'f5f07320-aebd-4b25-be22-6c8c2f736ed2',
        93,
        '1.2.6.16.3',
        'Damage Mitigation - CO Direction',
        'Liability',
        'Loss & Damage',
        'Undertake mitigation steps',
        'The contractor shall undertake specific mitigation steps as directed by CO.',
        false,
        ''
    ),
    (
        'bc76ad0b-3b23-4c85-a441-15b8a4d0ef12',
        94,
        '1.2.6.16.5.',
        'Mold Remediation Estimate',
        'Liability',
        'Loss & Damage',
        'Proper mold remediation estimate',
        'Prior to undertaking any remediation work, the contractor shall procure the services of a qualified mold remediation firm and obtain a written estimate, unless otherwise directed by the COR.  The contractor shall provide a copy of the estimate to the QAE, COR, MCO, and customer.',
        false,
        ''
    ),
    (
        '7894320f-0904-4e8b-be70-4835041f1798',
        95,
        '1.2.6.16.6.',
        'Uncontaminated Items Delivery',
        'Liability',
        'Loss & Damage',
        'Deliver uncontaminated items',
        'The contractor shall deliver any uncontaminated items to the destination.',
        false,
        ''
    ),
    (
        '35aa6a32-385b-4b4c-9d29-bf133415d109',
        96,
        '1.2.6.16.7.',
        'Pictures and Inventory (Remediation)',
        'Liability',
        'Loss & Damage',
        'Provide pictures and Inventory of Salvageable & Non-salvageable Items',
        'The contractor shall provide pictures and an inventory of each category, salvageable & non-salvageable, if requested by the Government.',
        false,
        ''
    ),
    (
        '1829a978-b53f-4caf-aa5a-07daa1537839',
        97,
        '1.2.6.16.9.',
        'Disposal of un-remediated contaminated items',
        'Liability',
        'Loss & Damage',
        'Appropriately dispose of Un-remediated Contaminated Items',
        'The contractor shall be responsible for appropriately disposing of the un-remediated portion of the contaminated items.',
        false,
        ''
    ),
    (
        '4a873236-df31-442b-a676-0d72e1f10002',
        98,
        '1.2.7.2.4',
        'Hardship Expenses',
        'Liability',
        'Inconvenience & Hardship Claims',
        'Pay for customer''s hardship expenses due to service failure',
        'In the event the contractor fails to perform IAW the agreed to schedule, the contractor shall reimburse the customer for any out of pocket expenses incurred which are determined unavoidable and unrecoverable under any other means by the COR. These amounts shall be in addition to amounts paid in relation to an inconvenience claim.',
        false,
        ''
    ),
    (
        'ae2c5021-ad8b-4521-8391-03a70d123dd6',
        99,
        '1.2.7.2.3',
        'Inconvenience Claim',
        'Liability',
        'Inconvenience & Hardship Claims',
        'Payment for inconvenience claims',
        'The contractor shall pay the customer a daily amount equal to the applicable pickup or delivery location government per diem (to exclude lodging) for all individuals on the relocation order according to the JTR for all days past the missed pickup or delivery date. The contractor shall, in addition, pay the customer the applicable daily amount for each day the customer is awaiting delivery out of SIT if not completed on customer''s first requested date and scheduled delivery date is not within five (5) GBDs (within ten (10) GBDs for shipments with a requested delivery date between June 15 through August 15).',
        false,
        ''
    ),
    (
        '07e2e788-3937-44ee-a5b8-87e5b085f766',
        100,
        '1.2.6.16.16',
        'Inconvenience Claim Liability',
        'Liability',
        'Inconvenience & Hardship Claims',
        'Inconvenience Claim Liability',
        'Contractor may be liable for an inconvenience claim until the items are available for delivery.',
        false,
        ''
    );

alter table report_violations add CONSTRAINT report_violations_violation_id_fkey FOREIGN KEY (violation_id) REFERENCES public.pws_violations (id);