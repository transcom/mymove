-- Jira: MB-16327

UPDATE transportation_offices
    SET name = 'PPPO Tobyhanna Army Depot - USA',
        gbloc = 'AGFM'
    WHERE id = '46898e12-8657-4ece-bb89-9a9e94815db9';

UPDATE addresses
    SET street_address_1 = '11 Hap Arnold Blvd',
        street_address_2 = NULL,
        city = 'Coolbaugh Township',
        state = 'PA',
        postal_code = '18466'
    WHERE id = (SELECT address_id FROM transportation_offices where id = '46898e12-8657-4ece-bb89-9a9e94815db9');
