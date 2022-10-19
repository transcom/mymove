import reducer, { getDocsByStatusAndType, getIncentiveActionType, getTabularExpenses } from './ducks';
describe('office ppm reducer', () => {
  describe('GET_PPM_INCENTIVE', () => {
    it('handles SUCCESS', () => {
      const newState = reducer(null, {
        type: getIncentiveActionType.success,
        payload: { gcc: 123400, incentive_percentage: 12400 },
      });

      expect(newState).toEqual({
        isLoading: false,
        hasErrored: false,
        hasSucceeded: true,
        calculation: { gcc: 123400, incentive_percentage: 12400 },
      });
    });
    it('handles START', () => {
      const newState = reducer(null, {
        type: getIncentiveActionType.start,
      });
      expect(newState).toEqual({
        isLoading: true,
        hasErrored: false,
        hasSucceeded: false,
      });
    });
    it('handles FAILURE', () => {
      const newState = reducer(null, {
        type: getIncentiveActionType.failure,
      });
      expect(newState).toEqual({
        isLoading: false,
        hasErrored: true,
        hasSucceeded: false,
      });
    });
  });
  describe('CLEAR_PPM_INCENTIVE', () => {
    it('handles SUCCESS', () => {
      const newState = reducer(null, {
        type: 'CLEAR_PPM_INCENTIVE',
      });

      expect(newState).toEqual({
        calculation: null,
      });
    });
  });
});
describe('getTabularExpenses', () => {
  const schema = {
    type: 'string',
    title: 'Moving Expense Type',
    enum: [
      'CONTRACTED_EXPENSE',
      'RENTAL_EQUIPMENT',
      'PACKING_MATERIALS',
      'WEIGHING_FEE',
      'GAS',
      'TOLLS',
      'OIL',
      'OTHER',
    ],
    'x-display-value': {
      CONTRACTED_EXPENSE: 'Contracted Expense',
      RENTAL_EQUIPMENT: 'Rental Equipment',
      PACKING_MATERIALS: 'Packing Materials',
      WEIGHING_FEE: 'Weighing Fees',
      GAS: 'Gas',
      TOLLS: 'Tolls',
      OIL: 'Oil',
      OTHER: 'Other',
    },
  };
  describe('when there is no expense data', () => {
    it('return and empty array', () => {
      expect(getTabularExpenses(null, null)).toEqual([]);
    });
  });
  describe('when there are a few categories', () => {
    const expenseData = {
      categories: [
        {
          category: 'CONTRACTED_EXPENSE',
          payment_methods: {
            GTCC: 600,
          },
          total: 600,
        },
        {
          category: 'RENTAL_EQUIPMENT',
          payment_methods: {
            OTHER: 500,
          },
          total: 500,
        },
        {
          category: 'TOLLS',
          payment_methods: {
            OTHER: 500,
          },
          total: 500,
        },
      ],
      grand_total: {
        payment_method_totals: {
          GTCC: 600,
          OTHER: 1000,
        },
        total: 1600,
      },
    };
    const result = getTabularExpenses(expenseData, schema);
    it('should fill in all categories', () => {
      expect(result.map((r) => r.type)).toEqual([
        'Contracted Expense',
        'Rental Equipment',
        'Packing Materials',
        'Weighing Fees',
        'Gas',
        'Tolls',
        'Oil',
        'Other',
        'Total',
      ]);
    });
    it('should extract GTCC', () => {
      expect(result[0].GTCC).toEqual(600);
    });
    it('should extract other', () => {
      expect(result[1].other).toEqual(500);
    });

    it('should include total as last object in array', () => {
      expect(result[result.length - 1]).toEqual({
        GTCC: 600,
        other: 1000,
        total: 1600,
        type: 'Total',
      });
    });

    it('should reshape by category', () => {
      expect(result).toEqual([
        { GTCC: 600, other: null, total: 600, type: 'Contracted Expense' },
        {
          GTCC: null,
          other: 500,
          total: 500,
          type: 'Rental Equipment',
        },
        {
          GTCC: null,
          other: null,
          total: null,
          type: 'Packing Materials',
        },
        {
          GTCC: null,
          other: null,
          total: null,
          type: 'Weighing Fees',
        },
        { GTCC: null, other: null, total: null, type: 'Gas' },
        { GTCC: null, other: 500, total: 500, type: 'Tolls' },
        { GTCC: null, other: null, total: null, type: 'Oil' },
        { GTCC: null, other: null, total: null, type: 'Other' },
        { GTCC: 600, other: 1000, total: 1600, type: 'Total' },
      ]);
    });
  });
  describe('getDocsByStatusAndType', () => {
    it('should filter documents by status and type to exclude', () => {
      const documents = [
        {
          move_document_type: 'EXPENSE',
          status: 'AWAITING_REVIEW',
        },
        {
          move_document_type: 'STORAGE',
          status: 'HAS_ISSUE',
        },
        {
          move_document_type: 'EXPENSE',
          status: 'OK',
        },
        {
          move_document_type: 'STORAGE',
          status: 'OK',
        },
      ];
      const filteredDocs = getDocsByStatusAndType(documents, 'OK', 'STORAGE');
      expect(filteredDocs).toEqual([
        {
          move_document_type: 'EXPENSE',
          status: 'AWAITING_REVIEW',
        },
      ]);
    });

    it('should filter documents by status to exclude when a type is missing', () => {
      const documents = [
        {
          move_document_type: 'EXPENSE',
          status: 'AWAITING_REVIEW',
        },
        {
          move_document_type: 'STORAGE',
          status: 'HAS_ISSUE',
        },
        {
          move_document_type: 'EXPENSE',
          status: 'OK',
        },
        {
          move_document_type: 'STORAGE',
          status: 'OK',
        },
      ];
      const filteredDocs = getDocsByStatusAndType(documents, 'OK');
      expect(filteredDocs).toEqual([
        {
          move_document_type: 'EXPENSE',
          status: 'AWAITING_REVIEW',
        },
        {
          move_document_type: 'STORAGE',
          status: 'HAS_ISSUE',
        },
      ]);
    });

    it('should filter documents by type to exclude when a status is missing', () => {
      const documents = [
        {
          move_document_type: 'EXPENSE',
          status: 'AWAITING_REVIEW',
        },
        {
          move_document_type: 'STORAGE',
          status: 'HAS_ISSUE',
        },
        {
          move_document_type: 'EXPENSE',
          status: 'OK',
        },
        {
          move_document_type: 'STORAGE',
          status: 'OK',
        },
      ];
      const filteredDocs = getDocsByStatusAndType(documents, null, 'STORAGE');
      expect(filteredDocs).toEqual([
        {
          move_document_type: 'EXPENSE',
          status: 'AWAITING_REVIEW',
        },
        {
          move_document_type: 'EXPENSE',
          status: 'OK',
        },
      ]);
    });
  });
});
