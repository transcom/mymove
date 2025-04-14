export const expenseTypes = {
  CONTRACTED_EXPENSE: 'CONTRACTED_EXPENSE',
  OIL: 'OIL',
  PACKING_MATERIALS: 'PACKING_MATERIALS',
  RENTAL_EQUIPMENT: 'RENTAL_EQUIPMENT',
  SMALL_PACKAGE: 'SMALL_PACKAGE',
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
  'SMALL_PACKAGE',
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
  SMALL_PACKAGE: 'Small package reimbursement',
  STORAGE: 'Storage',
  TOLLS: 'Tolls',
  WEIGHING_FEE: 'Weighing fee',
  OTHER: 'Other',
};

export const getExpenseTypeValue = (key) => expenseTypeLabels[key];

export const llvmExpenseTypes = Object.fromEntries(
  Object.entries(expenseTypeLabels).map(([key, value]) => [value, key]),
);

export const ppmExpenseTypes = [
  { value: 'Contracted expense', key: 'CONTRACTED_EXPENSE' },
  { value: 'Oil', key: 'OIL' },
  { value: 'Packing materials', key: 'PACKING_MATERIALS' },
  { value: 'Rental equipment', key: 'RENTAL_EQUIPMENT' },
  { value: 'Small package reimbursement', key: 'SMALL_PACKAGE' },
  { value: 'Storage', key: 'STORAGE' },
  { value: 'Tolls', key: 'TOLLS' },
  { value: 'Weighing fee', key: 'WEIGHING_FEE' },
  { value: 'Other', key: 'OTHER' },
];
