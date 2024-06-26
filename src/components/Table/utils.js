import { BRANCH_OPTIONS, MOVE_STATUS_OPTIONS } from 'constants/queues';
import { PAYMENT_REQUEST_STATUS_LABELS } from 'constants/paymentRequestStatus';

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

export const OFFICE_TABLE_QUEUE_SESSION_STORAGE_ID = 'office_table_queue_filters';
export const TEMPLATE_OFFICE_TABLE_QUEUE_FILTER_CACHE_OBJECT = {
  filters: [],
  sortParam: null,
  page: null,
  pageSize: null,
};

/**
 * Method to store table queue filters for a given session.
 * @param {string} key is session storage key
 * @param {[object]} value  is session storage value to cache
 */
export const setTableQueueFilterSessionStorageValue = (key, value) => {
  const cache = window.sessionStorage.getItem(OFFICE_TABLE_QUEUE_SESSION_STORAGE_ID);
  let json = null;
  if (cache) {
    json = JSON.parse(cache);
    json[key].filters = value;
  } else {
    json = {};
    json[key] = TEMPLATE_OFFICE_TABLE_QUEUE_FILTER_CACHE_OBJECT;
    json[key].filters = value;
  }
  window.sessionStorage.setItem(OFFICE_TABLE_QUEUE_SESSION_STORAGE_ID, JSON.stringify(json));
};

/**
 * Method to retrieve table queue filters for a given session.
 * @param {string} key is session storage key
 * @param {json} value  is session storage value to cache
 * @returns {[object]}
 */
export const getTableQueueFilterSessionStorageValue = (key) => {
  let cache = window.sessionStorage.getItem(OFFICE_TABLE_QUEUE_SESSION_STORAGE_ID);
  if (!cache) {
    const newJson = {};
    newJson[key] = TEMPLATE_OFFICE_TABLE_QUEUE_FILTER_CACHE_OBJECT;
    cache = JSON.stringify(newJson);
  }
  const json = JSON.parse(cache);
  if (key in json) {
    return json[key].filters;
  }
  json[key] = TEMPLATE_OFFICE_TABLE_QUEUE_FILTER_CACHE_OBJECT;
  window.sessionStorage.setItem(OFFICE_TABLE_QUEUE_SESSION_STORAGE_ID, JSON.stringify(json));
  return json[key].filters;
};

/**
 * Method to cache page size of table queue by key.
 * @param {string} key is session storage key
 * @param {int} pageSize is page size to cache
 */
export const setTableQueuePageSizeSessionStorageValue = (key, pageSize) => {
  const cache = window.sessionStorage.getItem(OFFICE_TABLE_QUEUE_SESSION_STORAGE_ID);
  let json = null;
  if (cache) {
    json = JSON.parse(cache);
    json[key].pageSize = pageSize;
  } else {
    json = {};
    json[key] = TEMPLATE_OFFICE_TABLE_QUEUE_FILTER_CACHE_OBJECT;
    json[key].pageSize = pageSize;
  }
  window.sessionStorage.setItem(OFFICE_TABLE_QUEUE_SESSION_STORAGE_ID, JSON.stringify(json));
};

/**
 * Method to retrieve cached page size of table queue by key.
 * @param {string} key is session storage key
 * @returns {int} null by default
 */
export const getTableQueuePageSizeSessionStorageValue = (key) => {
  let cache = window.sessionStorage.getItem(OFFICE_TABLE_QUEUE_SESSION_STORAGE_ID);
  if (!cache) {
    const newJson = {};
    newJson[key] = TEMPLATE_OFFICE_TABLE_QUEUE_FILTER_CACHE_OBJECT;
    cache = JSON.stringify(newJson);
  }
  const json = JSON.parse(cache);
  if (key in json) {
    return json[key].pageSize;
  }
  json[key] = TEMPLATE_OFFICE_TABLE_QUEUE_FILTER_CACHE_OBJECT;
  window.sessionStorage.setItem(OFFICE_TABLE_QUEUE_SESSION_STORAGE_ID, JSON.stringify(json));
  return json[key].pageSize;
};

/**
 * Method to cache current page of table queue by key.
 * @param {string} key is session storage key
 * @param {int} page is page to cache
 */
export const setTableQueuePageSessionStorageValue = (key, page) => {
  const cache = window.sessionStorage.getItem(OFFICE_TABLE_QUEUE_SESSION_STORAGE_ID);
  let json = null;
  if (cache) {
    json = JSON.parse(cache);
    json[key].page = page;
  } else {
    json = {};
    json[key] = TEMPLATE_OFFICE_TABLE_QUEUE_FILTER_CACHE_OBJECT;
    json[key].page = page;
  }
  window.sessionStorage.setItem(OFFICE_TABLE_QUEUE_SESSION_STORAGE_ID, JSON.stringify(json));
};

/**
 * Method to retrieve cached page of table queue by key.
 * @param {string} key is session storage key
 * @returns {int} null by default
 */
export const getTableQueuePageSessionStorageValue = (key) => {
  let cache = window.sessionStorage.getItem(OFFICE_TABLE_QUEUE_SESSION_STORAGE_ID);
  if (!cache) {
    const newJson = {};
    newJson[key] = TEMPLATE_OFFICE_TABLE_QUEUE_FILTER_CACHE_OBJECT;
    cache = JSON.stringify(newJson);
  }
  const json = JSON.parse(cache);
  if (key in json) {
    return json[key].page;
  }
  json[key] = TEMPLATE_OFFICE_TABLE_QUEUE_FILTER_CACHE_OBJECT;
  window.sessionStorage.setItem(OFFICE_TABLE_QUEUE_SESSION_STORAGE_ID, JSON.stringify(json));
  return json[key].page;
};

/**
 * Method to cache sort parameters of table queue by key.
 * @param {string} key is session storage key
 * @param {[object]} sortParam is sortParam to cache
 */
export const setTableQueueSortParamSessionStorageValue = (key, sortParam) => {
  const cache = window.sessionStorage.getItem(OFFICE_TABLE_QUEUE_SESSION_STORAGE_ID);
  let json = null;
  if (cache) {
    json = JSON.parse(cache);
    json[key].sortParam = sortParam;
  } else {
    json = {};
    json[key] = TEMPLATE_OFFICE_TABLE_QUEUE_FILTER_CACHE_OBJECT;
    json[key].sortParam = sortParam;
  }
  window.sessionStorage.setItem(OFFICE_TABLE_QUEUE_SESSION_STORAGE_ID, JSON.stringify(json));
};

/**
 * Method to retrieve cached sort parameters of table queue by key.
 * @param {string} key is session storage key
 * @returns {[object]} null by default
 */
export const getTableQueueSortParamSessionStorageValue = (key) => {
  let cache = window.sessionStorage.getItem(OFFICE_TABLE_QUEUE_SESSION_STORAGE_ID);
  if (!cache) {
    const newJson = {};
    newJson[key] = TEMPLATE_OFFICE_TABLE_QUEUE_FILTER_CACHE_OBJECT;
    cache = JSON.stringify(newJson);
  }
  const json = JSON.parse(cache);
  if (key in json && json[key].sortParam !== null) {
    return json[key].sortParam;
  }
  json[key] = TEMPLATE_OFFICE_TABLE_QUEUE_FILTER_CACHE_OBJECT;
  window.sessionStorage.setItem(OFFICE_TABLE_QUEUE_SESSION_STORAGE_ID, JSON.stringify(json));
  return json[key].sortParam;
};

/**
 * Get display label for given id.
 * @param {string} value is option id
 * @returns {string} display label
 */
export const getSelectionOptionLabel = (value) => {
  // Loop though known options and attempt to retrieve display.
  let label = BRANCH_OPTIONS.filter((option) => value === option.value).map((option) => option.label);
  if (label.length > 0) {
    return label;
  }
  label = MOVE_STATUS_OPTIONS.filter((option) => value === option.value).map((option) => option.label);
  if (label.length > 0) {
    return label;
  }

  if (value in PAYMENT_REQUEST_STATUS_LABELS) {
    return PAYMENT_REQUEST_STATUS_LABELS[value];
  }

  // Nothing was found for value. Determine if missing known options.
  return 'N/A';
};
