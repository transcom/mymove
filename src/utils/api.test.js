import { getQueriesStatus } from './api';

describe('getQueriesStatus', () => {
  it('returns isLoading true if any queries are loading', () => {
    const testQueries = [
      { status: 'idle', isInitialLoading: false, isError: false, isSuccess: false },
      { status: 'loading', isInitialLoading: true, isError: false, isSuccess: false },
      { status: 'idle', isInitialLoading: false, isError: false, isSuccess: false },
    ];

    const result = {
      isLoading: true,
      isError: false,
      isSuccess: false,
      errors: [],
    };

    expect(getQueriesStatus(testQueries)).toEqual(result);
  });

  it('returns isLoading false if no queries are loading', () => {
    const testQueries = [
      { status: 'idle', isInitialLoading: false, isError: false, isSuccess: false },
      { status: 'idle', isInitialLoading: false, isError: false, isSuccess: false },
      { status: 'idle', isInitialLoading: false, isError: false, isSuccess: false },
    ];

    const result = {
      isLoading: false,
      isError: false,
      isSuccess: false,
      errors: [],
    };

    expect(getQueriesStatus(testQueries)).toEqual(result);
  });

  it('returns isError true if any queries are errored', () => {
    const testQueries = [
      { status: 'success', isInitialLoading: false, isError: false, isSuccess: true },
      { status: 'idle', isInitialLoading: false, isError: false, isSuccess: false },
      { status: 'error', isInitialLoading: false, isError: true, isSuccess: false, error: 'Test error' },
    ];

    const result = {
      isLoading: false,
      isError: true,
      isSuccess: false,
      errors: ['Test error'],
    };

    expect(getQueriesStatus(testQueries)).toEqual(result);
  });

  it('returns isError false if no queries are errored', () => {
    const testQueries = [
      { status: 'idle', isInitialLoading: false, isError: false, isSuccess: false },
      { status: 'idle', isInitialLoading: false, isError: false, isSuccess: false },
      { status: 'idle', isInitialLoading: false, isError: false, isSuccess: false },
    ];

    const result = {
      isLoading: false,
      isError: false,
      isSuccess: false,
      errors: [],
    };

    expect(getQueriesStatus(testQueries)).toEqual(result);
  });

  it('returns isError true if any queries are errored', () => {
    const testQueries = [
      { status: 'success', isInitialLoading: false, isError: false, isSuccess: true },
      { status: 'idle', isInitialLoading: false, isError: false, isSuccess: false },
      { status: 'error', isInitialLoading: false, isError: true, isSuccess: false, error: 'Test error' },
    ];

    const result = {
      isLoading: false,
      isError: true,
      isSuccess: false,
      errors: ['Test error'],
    };

    expect(getQueriesStatus(testQueries)).toEqual(result);
  });

  it('returns isSuccess false if not all queries are success', () => {
    const testQueries = [
      { status: 'success', isInitialLoading: false, isError: false, isSuccess: true },
      { status: 'success', isInitialLoading: false, isError: false, isSuccess: true },
      { status: 'idle', isInitialLoading: false, isError: false, isSuccess: false },
    ];

    const result = {
      isLoading: false,
      isError: false,
      isSuccess: false,
      errors: [],
    };

    expect(getQueriesStatus(testQueries)).toEqual(result);
  });

  it('returns isSuccess true if all queries are success', () => {
    const testQueries = [
      { status: 'success', isInitialLoading: false, isError: false, isSuccess: true },
      { status: 'success', isInitialLoading: false, isError: false, isSuccess: true },
      { status: 'success', isInitialLoading: false, isError: false, isSuccess: true },
    ];

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
      {
        status: 'error',
        isInitialLoading: false,
        isError: true,
        isSuccess: false,
        error: new Error('Test API Error 1'),
      },
      { status: 'success', isInitialLoading: false, isError: false, isSuccess: true },
      {
        status: 'error',
        isInitialLoading: false,
        isError: true,
        isSuccess: false,
        error: new Error('Test API Error 2'),
      },
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
