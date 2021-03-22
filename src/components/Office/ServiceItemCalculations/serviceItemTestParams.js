const ContractCode = {
  eTag: 'MjAyMS0wMy0xOFQwMTozMTo1MS4yNDYwMDRa',
  id: 'f2a3e73f-6450-43d6-a783-181501cfab22',
  key: 'ContractCode',
  origin: 'SYSTEM',
  paymentServiceItemID: '28039a62-387d-479f-b50f-e0041b7e6e22',
  type: 'STRING',
  value: '1',
};

const ContractYearName = {
  eTag: 'MjAyMS0wMy0xOFQwMTozMTo1MS4yNTYzOVo=',
  id: '0ccef02a-59da-44d7-8258-f0e24c6c9b97',
  key: 'ContractYearName',
  origin: 'PRICER',
  paymentServiceItemID: '28039a62-387d-479f-b50f-e0041b7e6e22',
  type: 'STRING',
  value: 'Contract Year Name',
};

const EscalationCompounded = {
  eTag: 'MjAyMS0wMy0xOFQwMTozMTo1MS4yNzc2MTda',
  id: '6c7f1673-1ada-44fe-aa9b-e921d6e15f0e',
  key: 'EscalationCompounded',
  origin: 'PRICER',
  paymentServiceItemID: '28039a62-387d-479f-b50f-e0041b7e6e22',
  type: 'DECIMAL',
  value: '1.033',
};

const IsPeak = {
  eTag: 'MjAyMS0wMy0xOFQwMTozMTo1MS4yODc3NDla',
  id: '83f24c0d-25ab-465a-b60b-d27bfb77b41a',
  key: 'IsPeak',
  origin: 'PRICER',
  paymentServiceItemID: '28039a62-387d-479f-b50f-e0041b7e6e22',
  type: 'BOOLEAN',
  value: 'FALSE',
};

const PriceRateOrFactor = {
  eTag: 'MjAyMS0wMy0xOFQwMTozMTo1MS4yOTc2ODJa',
  id: 'b3ca0c12-fea3-4dd1-b228-30c1cc007452',
  key: 'PriceRateOrFactor',
  origin: 'PRICER',
  paymentServiceItemID: '28039a62-387d-479f-b50f-e0041b7e6e22',
  type: 'DECIMAL',
  value: '1.033',
};

const RequestedPickupDate = {
  eTag: 'MjAyMS0wMy0xOFQwMTozMTo1MS4zMDY2Nzha',
  id: '0e908b35-e61b-47c5-b4bc-f1649aa1cdc2',
  key: 'RequestedPickupDate',
  origin: 'PRIME',
  paymentServiceItemID: '28039a62-387d-479f-b50f-e0041b7e6e22',
  type: 'DATE',
  value: '2020-03-11',
};

const ServiceAreaOrigin = {
  eTag: 'MjAyMS0wMy0xOFQwMTozMTo1MS4zMTY5NDha',
  id: '87e77d29-d8c9-4b74-b45f-6842cd3ef970',
  key: 'ServiceAreaOrigin',
  origin: 'SYSTEM',
  paymentServiceItemID: '28039a62-387d-479f-b50f-e0041b7e6e22',
  type: 'STRING',
  value: '176',
};

const WeightActual = {
  eTag: 'MjAyMS0wMy0xOFQwMTozMTo1MS4zMjY2NDVa',
  id: '70abd9bc-afaa-4e4d-ad15-d3e55b57d2fb',
  key: 'WeightActual',
  origin: 'PRIME',
  paymentServiceItemID: '28039a62-387d-479f-b50f-e0041b7e6e22',
  type: 'INTEGER',
  value: '8500',
};

const WeightBilledActual = {
  eTag: 'MjAyMS0wMy0xOFQwMTozMTo1MS4zMzU1Njda',
  id: '5a993802-1504-4415-9b18-fdb1fdfd201c',
  key: 'WeightBilledActual',
  origin: 'SYSTEM',
  paymentServiceItemID: '28039a62-387d-479f-b50f-e0041b7e6e22',
  type: 'INTEGER',
  value: '8500',
};

const WeightEstimated = {
  eTag: 'MjAyMS0wMy0xOFQwMTozMTo1MS4zNDQxMTda',
  id: '02438e39-de6c-4c64-b817-9932ee319a4c',
  key: 'WeightEstimated',
  origin: 'PRIME',
  paymentServiceItemID: '28039a62-387d-479f-b50f-e0041b7e6e22',
  type: 'INTEGER',
  value: '8000',
};

const ZipDestAddress = {
  eTag: 'MjAyMS0wMy0xOFQwMTozMTo1MS4zNTI1MDZa',
  id: 'b26fcc8f-2c06-4b00-8b51-4715a2eb0f33',
  key: 'ZipDestAddress',
  origin: 'PRIME',
  paymentServiceItemID: '28039a62-387d-479f-b50f-e0041b7e6e22',
  type: 'STRING',
  value: '91910',
};

const ZipPickupAddress = {
  eTag: 'MjAyMS0wMy0xOFQwMTozMTo1MS4zNjA5MTha',
  id: 'dcfa55b2-3106-4e1b-af4a-f19d82b5f446',
  key: 'ZipPickupAddress',
  origin: 'PRIME',
  paymentServiceItemID: '28039a62-387d-479f-b50f-e0041b7e6e22',
  type: 'STRING',
  value: '32210',
};

const zip3 = {
  eTag: 'MjAyMS0wMy0xOFQwMTozMTo1MS4yNjY4M1o=',
  id: 'b4ba804d-f661-4df1-a488-11da9668647b',
  key: 'DistanceZip3',
  origin: 'SYSTEM',
  paymentServiceItemID: '28039a62-387d-479f-b50f-e0041b7e6e22',
  type: 'INTEGER',
  value: '210',
};

const testParams = {
  DomesticLongHaul: [
    ContractCode,
    ContractYearName,
    EscalationCompounded,
    IsPeak,
    PriceRateOrFactor,
    RequestedPickupDate,
    ServiceAreaOrigin,
    WeightActual,
    WeightBilledActual,
    WeightEstimated,
    zip3,
    ZipDestAddress,
    ZipPickupAddress,
  ],
  DomesticShortHaul: [],
  DomesticOriginPrice: [],
  DomesticDestinationPrice: [],
  DomesticOrigin1stSIT: [],
  DomesticDestination1stSIT: [],
  DomesticOriginaAdditionalSIT: [],
  DomesticDestinationaAdditionalSIT: [],
  DomesticOriginaSITDelivery: [],
  DomesticDestinationaSITDelivery: [],
  DomesticPacking: [],
  DomesticUnpacking: [],
  DomesticCrating: [],
  DomesticCratingStandalone: [],
  DomesticUncrating: [],
  DomesticOriginShuttleService: [],
  DomesticDestinationShuttleService: [],
  NonStandardHHG: [],
  NonStandardUB: [],
  FuelSurchage: [],
  DomesticMobileHomeFactor: [],
  DomesticTowAwayBoatFactor: [],
  DomesticNTSPackingFactor: [],
};

export default testParams;
