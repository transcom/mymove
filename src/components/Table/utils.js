/**
 * Helper function that creates the header object to pass into a react-table.
 * @param {string} header is the table header name
 * @param {string} accessor  is the key to use to look up the prop in the passed in data
 * @param {object} options is any additional options to include in the header object
 * @returns {{Header: string, accessor: string}}
 */
const createHeader = (header, accessor, options = {}) => {
  return {
    Header: header,
    accessor, // accessor is the "key" in the data
    ...options,
  };
};

/**
 * Creates a string header to pass into a react-table.
 * @param {string} header is the table header name
 * @param {string} accessor  is the key to use to look up the prop in the passed in data
 * @returns {{Header: string, accessor: string}}
 */
// eslint-disable-next-line import/prefer-default-export
export const createStringHeader = (header, accessor) => {
  return createHeader(header, accessor);
};
