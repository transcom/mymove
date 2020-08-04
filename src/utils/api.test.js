import { mapObjectToArray, getQueriesStatus } from './api';

describe('mapObjectToArray', () => {
  it('returns an array of the values from a given object', () => {
    const testObject = {
      test1: 'my value',
      test2: 'second value',
      test3: { name: 'object value' },
      test4: false,
      test5: [1, 4, 5],
    };

    const testArray = ['my value', 'second value', { name: 'object value' }, false, [1, 4, 5]];

    expect(mapObjectToArray(testObject)).toEqual(testArray);
  });

  it('doesn’t crash if the object is empty', () => {
    const testObject = {};
    const testArray = [];

    expect(mapObjectToArray(testObject)).toEqual(testArray);
  });
});

describe('getQueriesStatus', () => {
  it('returns isLoading true if any queries are loading', () => {
    const testQueries = [{ status: 'idle' }, { status: 'loading' }, { status: 'idle' }];

    const result = {
      isLoading: true,
      isError: false,
      isSuccess: false,
      errors: [],
    };

    expect(getQueriesStatus(testQueries)).toEqual(result);
  });

  it('returns isLoading false if no queries are loading', () => {
    const testQueries = [{ status: 'idle' }, { status: 'idle' }, { status: 'idle' }];

    const result = {
      isLoading: false,
      isError: false,
      isSuccess: false,
      errors: [],
    };

    expect(getQueriesStatus(testQueries)).toEqual(result);
  });

  it('returns isError true if any queries are errored', () => {
    const testQueries = [{ status: 'success' }, { status: 'idle' }, { status: 'error' }];

    const result = {
      isLoading: false,
      isError: true,
      isSuccess: false,
      errors: [],
    };

    expect(getQueriesStatus(testQueries)).toEqual(result);
  });

  it('returns isError false if no queries are errored', () => {
    const testQueries = [{ status: 'idle' }, { status: 'idle' }, { status: 'idle' }];

    const result = {
      isLoading: false,
      isError: false,
      isSuccess: false,
      errors: [],
    };

    expect(getQueriesStatus(testQueries)).toEqual(result);
  });

  it('returns isError true if any queries are errored', () => {
    const testQueries = [{ status: 'success' }, { status: 'idle' }, { status: 'error' }];

    const result = {
      isLoading: false,
      isError: true,
      isSuccess: false,
      errors: [],
    };

    expect(getQueriesStatus(testQueries)).toEqual(result);
  });

  it('returns isSuccess false if not all queries are success', () => {
    const testQueries = [{ status: 'success' }, { status: 'success' }, { status: 'idle' }];

    const result = {
      isLoading: false,
      isError: false,
      isSuccess: false,
      errors: [],
    };

    expect(getQueriesStatus(testQueries)).toEqual(result);
  });

  it('returns isSuccess true if all queries are success', () => {
    const testQueries = [{ status: 'success' }, { status: 'success' }, { status: 'success' }];

    const result = {
      isLoading: false,
      isError: false,
      isSuccess: true,
      errors: [],
    };

    expect(getQueriesStatus(testQueries)).toEqual(result);
  });

  it('reduces all errors into a single array', () => {
    const testQueries = [
      { status: 'error', error: new Error('Test API Error 1') },
      { status: 'success' },
      { status: 'error', error: new Error('Test API Error 2') },
    ];

    const result = {
      isLoading: false,
      isError: true,
      isSuccess: false,
      errors: [testQueries[0].error, testQueries[2].error],
    };

    expect(getQueriesStatus(testQueries)).toEqual(result);
  });
});
