/* eslint-disable import/prefer-default-export */
// utility functions related to API interactions

// eslint-disable-next-line security/detect-object-injection
export const mapObjectToArray = (obj) => Object.keys(obj).map((i) => obj[i]);

export const getQueriesStatus = (queries) => {
  // Queries should be an array of statuses returned by useQuery (https://react-query.tanstack.com/docs/api#usequery)
  return {
    isLoading: queries.some((q) => q.status === 'loading'),
    isError: queries.some((q) => q.status === 'error'),
    isSuccess: queries.every((q) => q.status === 'success'),
    errors: queries.reduce((errors, q) => (q.error ? [...errors, q.error] : errors), []),
  };
};
