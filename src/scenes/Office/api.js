import _ from 'lodash';

import { getClient, checkResponse } from 'shared/api';

export const requestData = (queueType, pageSize, page, sorted, filtered) => {
  return new Promise((resolve, reject) => {
    const rawData = RetrieveMovesForOffice(queueType);

    let filteredData = [rawData]; // Ok, this is an object, not an array
    console.log(typeof filteredData);

    // You can use the filters in your request, but you are responsible for applying them.
    if (filtered.length) {
      filteredData = filtered.reduce((filteredSoFar, nextFilter) => {
        return filteredSoFar.filter(row => {
          return (row[nextFilter.id] + '').includes(nextFilter.value);
        });
      }, filteredData);
    }
    // You can also use the sorting in your request, but again, you are responsible for applying it.
    const sortedData = _.orderBy(
      filteredData,
      sorted.map(sort => {
        return row => {
          if (row[sort.id] === null || row[sort.id] === undefined) {
            return -Infinity;
          }
          return typeof row[sort.id] === 'string'
            ? row[sort.id].toLowerCase()
            : row[sort.id];
        };
      }),
      sorted.map(d => (d.desc ? 'desc' : 'asc')),
    );

    // You must return an object containing the rows of the current page, and optionally the total pages number.
    const res = {
      rows: filteredData.slice(pageSize * page, pageSize * page + pageSize),
      pages: Math.ceil(filteredData.length / pageSize),
    };
  });
};

export async function GetSpec() {
  const client = await getClient();
  return client.spec;
}

export async function RetrieveMovesForOffice(queueType) {
  const client = await getClient();
  const response = await client.apis.queues.showQueue({
    queueType,
  });
  checkResponse(response, 'failed to retrieve moves due to server error');
  return response.body;
}
