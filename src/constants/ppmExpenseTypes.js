export const expenseTypes = {
  CONTRACTED_EXPENSE: 'CONTRACTED_EXPENSE',
  OIL: 'OIL',
  PACKING_MATERIALS: 'PACKING_MATERIALS',
  RENTAL_EQUIPMENT: 'RENTAL_EQUIPMENT',
  STORAGE: 'STORAGE',
  TOLLS: 'TOLLS',
  WEIGHING_FEE: 'WEIGHING_FEE',
  OTHER: 'OTHER',
};

export const expenseTypesArr = [
  'CONTRACTED_EXPENSE',
  'OIL',
  'PACKING_MATERIALS',
  'RENTAL_EQUIPMENT',
  'STORAGE',
  'TOLLS',
  'WEIGHING_FEE',
  'OTHER',
];

export const expenseTypeLabels = {
  CONTRACTED_EXPENSE: 'Contracted expense',
  OIL: 'Oil',
  PACKING_MATERIALS: 'Packing materials',
  RENTAL_EQUIPMENT: 'Rental equipment',
  STORAGE: 'Storage',
  TOLLS: 'Tolls',
  WEIGHING_FEE: 'Weighing fee',
  OTHER: 'Other',
};

export const getExpenseTypeValue = (key) => expenseTypeLabels[key];

export const llvmExpenseTypes = {
  'Contracted expense': 'CONTRACTED_EXPENSE',
  Oil: 'OIL',
  'Packing materials': 'PACKING_MATERIALS',
  'Rental equipment': 'RENTAL_EQUIPMENT',
  Storage: 'STORAGE',
  Tolls: 'TOLLS',
  'Weighing fee': 'WEIGHING_FEE',
  Other: 'OTHER',
};

export const ppmExpenseTypes = [
  { value: 'Contracted expense', key: 'CONTRACTED_EXPENSE' },
  { value: 'Oil', key: 'OIL' },
  { value: 'Packing materials', key: 'PACKING_MATERIALS' },
  { value: 'Rental equipment', key: 'RENTAL_EQUIPMENT' },
  { value: 'Storage', key: 'STORAGE' },
  { value: 'Tolls', key: 'TOLLS' },
  { value: 'Weighing fee', key: 'WEIGHING_FEE' },
  { value: 'Other', key: 'OTHER' },
];
