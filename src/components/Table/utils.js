/**
 * Helper function that creates the header object to pass into a react-table.
 * @param {string} header is the table header name
 * @param {function(*)} accessor  is the key to use to look up the prop in the passed in data
 * @param {object} options is any additional options to include in the header object
 * @returns {{Header: string, accessor: string}}
 */
// eslint-disable-next-line import/prefer-default-export
export const createHeader = (header, accessor, options = {}) => {
  return {
    Header: header,
    accessor, // accessor is the "key" in the data
    ...options,
  };
};

const OFFICE_TABLE_FILTER_CACHE_ID = 'office_table_queue_filters';

/**
 * Helper method to store key:value pair to session storage.
 * This used to store queue filters for a given session.
 *
 * @param {string} key is session storage key
 * @param {json} value  is session storage value to cache
 */
export const setTableQueueFilterSessionStorageValue = (key, value) => {
  const cache = window.sessionStorage.getItem(OFFICE_TABLE_FILTER_CACHE_ID);
  let json = null;
  if (cache) {
    json = JSON.parse(cache);
    json[key] = value;
  } else {
    json = {};
    json[key] = value;
  }
  window.sessionStorage.setItem(OFFICE_TABLE_FILTER_CACHE_ID, JSON.stringify(json));
};

/**
 * Helper method to retrieve value of key:pair stored in session storage.
 * This used to retrieve queue filters for a given session.
 *
 * @param {string} key is session storage key
 * @param {json} value  is session storage value to cache
 * @returns {[]}
 */
export const getTableQueueFilterSessionStorageValue = (key) => {
  let cache = window.sessionStorage.getItem(OFFICE_TABLE_FILTER_CACHE_ID);
  if (!cache) {
    cache = '{}';
  }
  const json = JSON.parse(cache);
  if (key in json) {
    return json[key];
  }
  json[key] = [];
  window.sessionStorage.setItem(OFFICE_TABLE_FILTER_CACHE_ID, JSON.stringify(json));
  return json[key];
};
