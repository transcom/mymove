-- Add new weight-related service item param keys.
INSERT INTO service_item_param_keys
	(id, key, description, type, origin, created_at, updated_at)
VALUES ('d87d82da-3ac2-44e8-bce0-cb4de40f9a72', 'WeightAdjusted', 'Adjusted weight (by TIO)', 'INTEGER', 'SYSTEM',
		now(), now()),
	   ('1e6257e9-757d-4d59-8846-727dd8a055e7', 'WeightReweigh', 'Reweigh weight', 'INTEGER', 'PRIME', now(), now());

-- Add new keys to all service items with weight-based pricing.
INSERT INTO service_params
	(id, service_id, service_item_param_key_id, created_at, updated_at)
VALUES ('4a359f1a-7e93-46b8-b868-0e9c956556f8', (SELECT id FROM re_services WHERE code = 'DDASIT'), (SELECT id FROM service_item_param_keys WHERE key = 'WeightAdjusted'), now(), now()),
	   ('7f6faeef-756f-4373-af1f-116324cc3fe7', (SELECT id FROM re_services WHERE code = 'DDASIT'), (SELECT id FROM service_item_param_keys WHERE key = 'WeightReweigh'), now(), now()),
	   ('c0bd2d5f-3ed2-454b-bbdc-0a7cb6154cce', (SELECT id FROM re_services WHERE code = 'DDDSIT'), (SELECT id FROM service_item_param_keys WHERE key = 'WeightAdjusted'), now(), now()),
	   ('704e75b1-e2a2-45d8-93b0-29756b34a2c6', (SELECT id FROM re_services WHERE code = 'DDDSIT'), (SELECT id FROM service_item_param_keys WHERE key = 'WeightReweigh'), now(), now()),
	   ('6f7c8cf7-b369-4b9b-ad40-f38b5c2adcf0', (SELECT id FROM re_services WHERE code = 'DDFSIT'), (SELECT id FROM service_item_param_keys WHERE key = 'WeightAdjusted'), now(), now()),
	   ('27d3a8f3-ae49-41ff-a825-e8d296f16f3e', (SELECT id FROM re_services WHERE code = 'DDFSIT'), (SELECT id FROM service_item_param_keys WHERE key = 'WeightReweigh'), now(), now()),
	   ('05ec8b38-7879-4835-bdbf-e6eb8424c149', (SELECT id FROM re_services WHERE code = 'DDP'), (SELECT id FROM service_item_param_keys WHERE key = 'WeightAdjusted'), now(), now()),
	   ('c4a736ba-466e-48ea-9724-c7ccc51936d2', (SELECT id FROM re_services WHERE code = 'DDP'), (SELECT id FROM service_item_param_keys WHERE key = 'WeightReweigh'), now(), now()),
	   ('ca5d3f35-5d4f-4ecb-a132-f309ba289b91', (SELECT id FROM re_services WHERE code = 'DDSHUT'), (SELECT id FROM service_item_param_keys WHERE key = 'WeightAdjusted'), now(), now()),
	   ('771b78fe-b6d8-4edd-9642-dd52a97bb98d', (SELECT id FROM re_services WHERE code = 'DDSHUT'), (SELECT id FROM service_item_param_keys WHERE key = 'WeightReweigh'), now(), now()),
	   ('10654340-1000-4a5f-9ba3-d2ac4e6a3c2e', (SELECT id FROM re_services WHERE code = 'DLH'), (SELECT id FROM service_item_param_keys WHERE key = 'WeightAdjusted'), now(), now()),
	   ('1ed12c9d-0b3f-4ebd-86f2-80fc666efb31', (SELECT id FROM re_services WHERE code = 'DLH'), (SELECT id FROM service_item_param_keys WHERE key = 'WeightReweigh'), now(), now()),
	   ('b1d8863f-c906-474a-9bc0-a95251167df5', (SELECT id FROM re_services WHERE code = 'DOASIT'), (SELECT id FROM service_item_param_keys WHERE key = 'WeightAdjusted'), now(), now()),
	   ('8f868715-01ff-43ef-a3fb-96d57a9728fa', (SELECT id FROM re_services WHERE code = 'DOASIT'), (SELECT id FROM service_item_param_keys WHERE key = 'WeightReweigh'), now(), now()),
	   ('970d4afb-db8b-4b68-9e33-4f5f9e23aa05', (SELECT id FROM re_services WHERE code = 'DOFSIT'), (SELECT id FROM service_item_param_keys WHERE key = 'WeightAdjusted'), now(), now()),
	   ('eb761dfa-3291-4166-a460-3196212bb50e', (SELECT id FROM re_services WHERE code = 'DOFSIT'), (SELECT id FROM service_item_param_keys WHERE key = 'WeightReweigh'), now(), now()),
	   ('4e174246-bcce-4c1e-abc8-d83360180b11', (SELECT id FROM re_services WHERE code = 'DOP'), (SELECT id FROM service_item_param_keys WHERE key = 'WeightAdjusted'), now(), now()),
	   ('1c48d3c4-895a-4994-b6b2-85197d1e7d8f', (SELECT id FROM re_services WHERE code = 'DOP'), (SELECT id FROM service_item_param_keys WHERE key = 'WeightReweigh'), now(), now()),
	   ('cabf461c-5b96-4d50-baba-a19ef4f229af', (SELECT id FROM re_services WHERE code = 'DOPSIT'), (SELECT id FROM service_item_param_keys WHERE key = 'WeightAdjusted'), now(), now()),
	   ('ed767489-63fe-4c96-943e-6002f6215c02', (SELECT id FROM re_services WHERE code = 'DOPSIT'), (SELECT id FROM service_item_param_keys WHERE key = 'WeightReweigh'), now(), now()),
	   ('b1fb32e0-8c7a-467e-9b03-33fc26b5a2b8', (SELECT id FROM re_services WHERE code = 'DOSHUT'), (SELECT id FROM service_item_param_keys WHERE key = 'WeightAdjusted'), now(), now()),
	   ('1656fb06-b7f7-40b5-bc7d-636847cae0cd', (SELECT id FROM re_services WHERE code = 'DOSHUT'), (SELECT id FROM service_item_param_keys WHERE key = 'WeightReweigh'), now(), now()),
	   ('b1419707-90f3-4085-8ef9-534adda08a39', (SELECT id FROM re_services WHERE code = 'DPK'), (SELECT id FROM service_item_param_keys WHERE key = 'WeightAdjusted'), now(), now()),
	   ('09984861-13af-4a1d-b3e0-d7446316f50e', (SELECT id FROM re_services WHERE code = 'DPK'), (SELECT id FROM service_item_param_keys WHERE key = 'WeightReweigh'), now(), now()),
	   ('3fcab2c9-8219-4f29-ac93-b0db02b79bf3', (SELECT id FROM re_services WHERE code = 'DSH'), (SELECT id FROM service_item_param_keys WHERE key = 'WeightAdjusted'), now(), now()),
	   ('10814916-b3b8-48b3-89ef-1bbe4efcadbe', (SELECT id FROM re_services WHERE code = 'DSH'), (SELECT id FROM service_item_param_keys WHERE key = 'WeightReweigh'), now(), now()),
	   ('f34625e4-f820-4736-9921-92de51a12578', (SELECT id FROM re_services WHERE code = 'DUPK'), (SELECT id FROM service_item_param_keys WHERE key = 'WeightAdjusted'), now(), now()),
	   ('0b6f1257-ae67-4ee9-868d-7940ca9013cf', (SELECT id FROM re_services WHERE code = 'DUPK'), (SELECT id FROM service_item_param_keys WHERE key = 'WeightReweigh'), now(), now()),
	   ('1d7c6d2d-03ba-4d39-bbb6-0543475457b1', (SELECT id FROM re_services WHERE code = 'FSC'), (SELECT id FROM service_item_param_keys WHERE key = 'WeightAdjusted'), now(), now()),
	   ('6df0c2c1-9763-4c8b-9055-0912721cd6e6', (SELECT id FROM re_services WHERE code = 'FSC'), (SELECT id FROM service_item_param_keys WHERE key = 'WeightReweigh'), now(), now()),
	   ('92264493-ab39-4c6b-9e54-e09f6739f196', (SELECT id FROM re_services WHERE code = 'ICOLH'), (SELECT id FROM service_item_param_keys WHERE key = 'WeightAdjusted'), now(), now()),
	   ('991d8749-ba77-4533-b504-320adf7b8a03', (SELECT id FROM re_services WHERE code = 'ICOLH'), (SELECT id FROM service_item_param_keys WHERE key = 'WeightReweigh'), now(), now()),
	   ('063a5fe2-3c7e-4c9b-8abd-875f39797a29', (SELECT id FROM re_services WHERE code = 'ICOUB'), (SELECT id FROM service_item_param_keys WHERE key = 'WeightAdjusted'), now(), now()),
	   ('9ee81df9-9ca2-4b59-b752-623737de30bb', (SELECT id FROM re_services WHERE code = 'ICOUB'), (SELECT id FROM service_item_param_keys WHERE key = 'WeightReweigh'), now(), now()),
	   ('259237b9-3403-4303-9daa-4877126a6086', (SELECT id FROM re_services WHERE code = 'IDASIT'), (SELECT id FROM service_item_param_keys WHERE key = 'WeightAdjusted'), now(), now()),
	   ('4a61e9b4-b30c-4b86-a2e1-5616c7cdb8fc', (SELECT id FROM re_services WHERE code = 'IDASIT'), (SELECT id FROM service_item_param_keys WHERE key = 'WeightReweigh'), now(), now()),
	   ('a5ad9d30-d4db-4c22-9d69-0e95461c6c27', (SELECT id FROM re_services WHERE code = 'IDDSIT'), (SELECT id FROM service_item_param_keys WHERE key = 'WeightAdjusted'), now(), now()),
	   ('03a5f550-d9b5-4a02-99fd-9e07fc1827b7', (SELECT id FROM re_services WHERE code = 'IDDSIT'), (SELECT id FROM service_item_param_keys WHERE key = 'WeightReweigh'), now(), now()),
	   ('a8a2b922-b36c-4da0-b227-602634c174a6', (SELECT id FROM re_services WHERE code = 'IDFSIT'), (SELECT id FROM service_item_param_keys WHERE key = 'WeightAdjusted'), now(), now()),
	   ('a23989b1-96a6-49c2-95d3-70b1c3253d44', (SELECT id FROM re_services WHERE code = 'IDFSIT'), (SELECT id FROM service_item_param_keys WHERE key = 'WeightReweigh'), now(), now()),
	   ('743439f2-0737-4a23-b83e-8321f9106a65', (SELECT id FROM re_services WHERE code = 'IDSHUT'), (SELECT id FROM service_item_param_keys WHERE key = 'WeightAdjusted'), now(), now()),
	   ('32dd5d0e-8ca2-4b42-b2ef-73acd6b6c664', (SELECT id FROM re_services WHERE code = 'IDSHUT'), (SELECT id FROM service_item_param_keys WHERE key = 'WeightReweigh'), now(), now()),
	   ('8a196b23-32be-47ff-be59-230ab82431e8', (SELECT id FROM re_services WHERE code = 'IHPK'), (SELECT id FROM service_item_param_keys WHERE key = 'WeightAdjusted'), now(), now()),
	   ('58b94891-4dfe-4248-a276-6d3caee2a22e', (SELECT id FROM re_services WHERE code = 'IHPK'), (SELECT id FROM service_item_param_keys WHERE key = 'WeightReweigh'), now(), now()),
	   ('5aa63524-5c34-4d40-8ed4-c684b06ffbb2', (SELECT id FROM re_services WHERE code = 'IHUPK'), (SELECT id FROM service_item_param_keys WHERE key = 'WeightAdjusted'), now(), now()),
	   ('8c4d6820-ce36-41ed-9c68-d88a23013ed8', (SELECT id FROM re_services WHERE code = 'IHUPK'), (SELECT id FROM service_item_param_keys WHERE key = 'WeightReweigh'), now(), now()),
	   ('0a6aedcb-188a-4b7a-8e52-2ca48261605b', (SELECT id FROM re_services WHERE code = 'IOASIT'), (SELECT id FROM service_item_param_keys WHERE key = 'WeightAdjusted'), now(), now()),
	   ('b4a63b8f-6245-4aac-800e-726ff3428434', (SELECT id FROM re_services WHERE code = 'IOASIT'), (SELECT id FROM service_item_param_keys WHERE key = 'WeightReweigh'), now(), now()),
	   ('4824ff90-76e0-4563-bc67-fde7149221b2', (SELECT id FROM re_services WHERE code = 'IOCLH'), (SELECT id FROM service_item_param_keys WHERE key = 'WeightAdjusted'), now(), now()),
	   ('d934f347-7a4b-4bac-aeba-3a4b9e688b93', (SELECT id FROM re_services WHERE code = 'IOCLH'), (SELECT id FROM service_item_param_keys WHERE key = 'WeightReweigh'), now(), now()),
	   ('63734cd2-32e2-4e9c-aaee-3478c7ae1cf8', (SELECT id FROM re_services WHERE code = 'IOCUB'), (SELECT id FROM service_item_param_keys WHERE key = 'WeightAdjusted'), now(), now()),
	   ('1d5decdf-4d94-4eea-9ca8-166fe9f829dc', (SELECT id FROM re_services WHERE code = 'IOCUB'), (SELECT id FROM service_item_param_keys WHERE key = 'WeightReweigh'), now(), now()),
	   ('3bd79570-0def-454b-bd03-28f3089692ec', (SELECT id FROM re_services WHERE code = 'IOFSIT'), (SELECT id FROM service_item_param_keys WHERE key = 'WeightAdjusted'), now(), now()),
	   ('f6638b7c-7ce2-420e-ba72-4592fa065015', (SELECT id FROM re_services WHERE code = 'IOFSIT'), (SELECT id FROM service_item_param_keys WHERE key = 'WeightReweigh'), now(), now()),
	   ('79f00586-80f2-4175-b478-d8accc8dca5c', (SELECT id FROM re_services WHERE code = 'IOOLH'), (SELECT id FROM service_item_param_keys WHERE key = 'WeightAdjusted'), now(), now()),
	   ('89d4271f-30ea-4864-b3c0-2f3fe0a76625', (SELECT id FROM re_services WHERE code = 'IOOLH'), (SELECT id FROM service_item_param_keys WHERE key = 'WeightReweigh'), now(), now()),
	   ('cbde2cf0-8df9-4505-9d28-0c7d9c6da71c', (SELECT id FROM re_services WHERE code = 'IOOUB'), (SELECT id FROM service_item_param_keys WHERE key = 'WeightAdjusted'), now(), now()),
	   ('2ed2ac16-c510-4dac-a807-92630bd655da', (SELECT id FROM re_services WHERE code = 'IOOUB'), (SELECT id FROM service_item_param_keys WHERE key = 'WeightReweigh'), now(), now()),
	   ('454a3460-bc4a-4b3d-96ad-1bcf64d4a79a', (SELECT id FROM re_services WHERE code = 'IOPSIT'), (SELECT id FROM service_item_param_keys WHERE key = 'WeightAdjusted'), now(), now()),
	   ('ac00232a-c553-4183-875f-0ebf091cdbe6', (SELECT id FROM re_services WHERE code = 'IOPSIT'), (SELECT id FROM service_item_param_keys WHERE key = 'WeightReweigh'), now(), now()),
	   ('9140fd9e-1d8b-4a3e-98c5-7ba7d3a22bea', (SELECT id FROM re_services WHERE code = 'IOSHUT'), (SELECT id FROM service_item_param_keys WHERE key = 'WeightAdjusted'), now(), now()),
	   ('9b9212d4-4b3f-4ee5-9256-a2c0754a89d8', (SELECT id FROM re_services WHERE code = 'IOSHUT'), (SELECT id FROM service_item_param_keys WHERE key = 'WeightReweigh'), now(), now()),
	   ('1dd2498a-d2be-48d6-9946-d8e5900065a6', (SELECT id FROM re_services WHERE code = 'IUBPK'), (SELECT id FROM service_item_param_keys WHERE key = 'WeightAdjusted'), now(), now()),
	   ('c26de36b-f073-4db3-88a6-6739b9ca20e4', (SELECT id FROM re_services WHERE code = 'IUBPK'), (SELECT id FROM service_item_param_keys WHERE key = 'WeightReweigh'), now(), now()),
	   ('eb249bb6-0a14-4adc-8a05-dc4445fe3769', (SELECT id FROM re_services WHERE code = 'IUBUPK'), (SELECT id FROM service_item_param_keys WHERE key = 'WeightAdjusted'), now(), now()),
	   ('bb58ef24-4b0c-4253-9d78-8802ffd53057', (SELECT id FROM re_services WHERE code = 'IUBUPK'), (SELECT id FROM service_item_param_keys WHERE key = 'WeightReweigh'), now(), now()),
	   ('a3df5df2-84db-409a-a13c-c8fc61e13ace', (SELECT id FROM re_services WHERE code = 'NSTH'), (SELECT id FROM service_item_param_keys WHERE key = 'WeightAdjusted'), now(), now()),
	   ('2fd21c6b-c348-49f1-8848-a6cf03d4290e', (SELECT id FROM re_services WHERE code = 'NSTH'), (SELECT id FROM service_item_param_keys WHERE key = 'WeightReweigh'), now(), now()),
	   ('7282674f-7f4e-473e-bcd9-125191e443da', (SELECT id FROM re_services WHERE code = 'NSTUB'), (SELECT id FROM service_item_param_keys WHERE key = 'WeightAdjusted'), now(), now()),
	   ('16f4c508-8edb-48b8-aa83-49941693fd70', (SELECT id FROM re_services WHERE code = 'NSTUB'), (SELECT id FROM service_item_param_keys WHERE key = 'WeightReweigh'), now(), now());

-- Rename WeightActual to WeightOriginal.
UPDATE service_item_param_keys
SET key         = 'WeightOriginal',
	description = 'Original weight'
WHERE key = 'WeightActual';

-- Rename WeightBilledActual to WeightBilled.
UPDATE service_item_param_keys
SET key         = 'WeightBilled',
	description = 'Billed weight'
WHERE key = 'WeightBilledActual';

-- Add is_optional to our service_params metadata to if a parameter is optional for a given service item.
-- Default everything to false (required) initially as that's the typical case, then we will make certain
-- params optional.
ALTER TABLE service_params
	ADD COLUMN is_optional BOOLEAN NOT NULL DEFAULT false;

COMMENT ON COLUMN service_params.is_optional IS 'True if this parameter is optional for this service item.';

-- Note the optional weight-based params (applies to all weight-based service items currently).
UPDATE service_params
SET is_optional = true
WHERE service_item_param_key_id IN (
	SELECT id
	FROM service_item_param_keys
	WHERE key IN ('WeightAdjusted', 'WeightEstimated', 'WeightReweigh')
);
