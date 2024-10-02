--INC1859709 - Update postal code to GBLOC
-- BGAC glboc was transistioned to BGNC
--B-21500
update postal_code_to_gblocs
set
    gbloc = 'BGNC'
where
    postal_code like '233%'