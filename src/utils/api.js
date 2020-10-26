/*  import/prefer-default-export */
// utility functions related to API interactions

export const getQueriesStatus = (queries) => {
  // Queries should be an array of statuses returned by useQuery (https://react-query.tanstack.com/docs/api#usequery)
  return {
    isLoading: queries.some((q) => q.isLoading),
    isError: queries.some((q) => q.isError),
    isSuccess: queries.every((q) => q.isSuccess),
    errors: queries.reduce((errors, q) => (q.error ? [...errors, q.error] : errors), []),
  };
};
