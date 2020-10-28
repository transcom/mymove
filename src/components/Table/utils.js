/**
 * Helper function that creates the header object to pass into a react-table.
 * @param {string} header is the table header name
 * @param {string} accessor  is the key to use to look up the prop in the passed in data
 * @param {object} options is any additional options to include in the header object
 * @returns {{Header: string, accessor: string}}
 */
//  import/prefer-default-export
export const createHeader = (header, accessor, options = {}) => {
  return {
    Header: header,
    accessor, // accessor is the "key" in the data
    ...options,
  };
};

/**
 * Text filter function that searches with "startsWith".
 * @param rows The rows of data.
 * @param id The column ID name.
 * @param filterValue The filter value.
 * @returns {*} True, value matches.
 */
export const textFilter = (rows, id, filterValue) => {
  return rows.filter((row) => {
    const rowValue = row.values[`${id}`];
    return rowValue !== undefined ? String(rowValue).toLowerCase().startsWith(String(filterValue).toLowerCase()) : true;
  });
};
