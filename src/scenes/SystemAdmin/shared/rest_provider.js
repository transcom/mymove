import { stringify } from 'query-string';
import { diff } from 'deep-object-diff';
import { snakeCase } from 'lodash';
import {
  CREATE,
  DELETE,
  DELETE_MANY,
  fetchUtils,
  GET_LIST,
  GET_MANY,
  GET_MANY_REFERENCE,
  GET_ONE,
  UPDATE,
  UPDATE_MANY,
} from 'react-admin';

/**
 * Maps react-admin queries to a simple REST API
 *
 * The REST dialect is similar to the one of FakeRest
 * @see https://github.com/marmelab/FakeRest
 * @example
 * GET_LIST     => GET http://my.api.url/posts?sort=['title','ASC']&range=[0, 24]
 * GET_ONE      => GET http://my.api.url/posts/123
 * GET_MANY     => GET http://my.api.url/posts?filter={ids:[123,456,789]}
 * UPDATE       => PATCH http://my.api.url/posts/123
 * CREATE       => POST http://my.api.url/posts
 * DELETE       => DELETE http://my.api.url/posts/123
 */
const restProvider = (apiUrl, httpClient = fetchUtils.fetchJson) => {
  /**
   * @param {String} type One of the constants appearing at the top if this file, e.g. 'UPDATE'
   * @param {String} resource Name of the resource to fetch, e.g. 'posts'
   * @param {Object} params The data request params, depending on the type
   * @returns {Object} { url, options } The HTTP request parameters
   */
  const convertDataRequestToHTTP = (type, resource, params) => {
    let url = '';
    const options = {};
    options.headers = new Headers({ Accept: 'application/json' }); // required
    switch (type) {
      case GET_LIST: {
        const { page, perPage } = params.pagination;
        const { field, order } = params.sort;
        const query = {
          sort: snakeCase(field),
          order: order === 'ASC' ? true : false,
          page: page,
          perPage: perPage,
          filter: JSON.stringify(params.filter),
        };
        url = `${apiUrl}/${resource}?${stringify(query)}`;
        break;
      }
      case GET_ONE:
        url = `${apiUrl}/${resource}/${params.id}`;
        break;
      case GET_MANY: {
        if (resource === 'offices') {
          const query = {
            page: 1,
            perPage: 10000,
          };
          url = `${apiUrl}/${resource}?${stringify(query)}`;
          break;
        }
        if (params.ids !== undefined) {
          params.id = params.ids[0];
        }
        const query = {
          filter: JSON.stringify({
            id: params.id,
          }),
        };
        url = `${apiUrl}/${resource}?${stringify(query)}`;
        break;
      }
      case GET_MANY_REFERENCE: {
        const { page, perPage } = params.pagination;
        const { field, order } = params.sort;
        const query = {
          sort: snakeCase(field),
          order: order === 'ASC' ? true : false,
          page: page,
          perPage: perPage,
          filter: JSON.stringify({
            ...params.filter,
            [params.target]: params.id,
          }),
        };
        url = `${apiUrl}/${resource}?${stringify(query)}`;
        break;
      }
      case UPDATE:
        url = `${apiUrl}/${resource}/${params.id}`;
        options.method = 'PATCH';
        options.headers.set('If-Match', params.data?.eTag); // for optimistic locking / concurrency control
        const paramsDiff = diff(params.previousData, params.data);
        if (paramsDiff.roles) {
          paramsDiff.roles = params.data.roles;
        }
        if (params.data.active != null) {
          paramsDiff.active = params.data.active;
        }
        options.body = JSON.stringify(paramsDiff);
        break;
      case CREATE:
        url = `${apiUrl}/${resource}`;
        options.method = 'POST';
        options.body = JSON.stringify(params.data);
        break;
      case DELETE:
        url = `${apiUrl}/${resource}/${params.id}`;
        options.method = 'DELETE';
        break;
      default:
        throw new Error(`Unsupported fetch action type ${type}`);
    }
    return {
      url,
      options,
    };
  };

  /**
   * @param {Object} response HTTP response from fetch()
   * @param {String} type One of the constants appearing at the top if this file, e.g. 'UPDATE'
   * @param {String} resource Name of the resource to fetch, e.g. 'posts'
   * @param {Object} params The data request params, depending on the type
   * @returns {Object} Data response
   */
  const convertHTTPResponse = (response, type, resource, params) => {
    const { headers, json } = response;
    switch (type) {
      case GET_LIST:
      case GET_MANY_REFERENCE:
        if (!headers.has('content-range')) {
          throw new Error(
            'The Content-Range header is missing in the HTTP Response. The simple REST data provider expects responses for lists of resources to contain this header with the total number of results to build the pagination. If you are using CORS, did you declare Content-Range in the Access-Control-Expose-Headers header?',
          );
        }
        return {
          data: json,
          total: parseInt(headers.get('content-range').split('/').pop(), 10),
        };
      case CREATE:
        return {
          data: {
            ...params.data,
            id: json.id,
          },
        };
      case DELETE_MANY: {
        return {
          data: json || [],
        };
      }
      default:
        return {
          data: json,
        };
    }
  };

  /**
   * @param {string} type Request type, e.g GET_LIST
   * @param {string} resource Resource name, e.g. "posts"
   * @param {Object} payload Request parameters. Depends on the request type
   * @returns {Promise} the Promise for a data response
   */
  return (type, resource, params) => {
    // simple-rest doesn't handle filters on UPDATE route, so we fallback to calling UPDATE n times instead
    if (type === UPDATE_MANY) {
      return Promise.all(
        params.ids.map((id) =>
          httpClient(`${apiUrl}/${resource}/${id}`, {
            method: 'PATCH',
            body: JSON.stringify(params.data),
          }),
        ),
      ).then((responses) => ({
        data: responses.map((response) => response.json),
      }));
    }
    // simple-rest doesn't handle filters on DELETE route, so we fallback to calling DELETE n times instead
    if (type === DELETE_MANY) {
      return Promise.all(
        params.ids.map((id) =>
          httpClient(`${apiUrl}/${resource}/${id}`, {
            method: 'DELETE',
          }),
        ),
      ).then((responses) => ({
        data: responses.map((response) => response.json),
      }));
    }

    const { url, options } = convertDataRequestToHTTP(type, resource, params);
    return httpClient(url, options).then((response) => convertHTTPResponse(response, type, resource, params));
  };
};

export default restProvider;
