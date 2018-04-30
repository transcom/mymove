import React, { Component } from 'react';
import { render } from 'react-dom';
import _ from 'lodash';
import { FunctionNameGoesHere } from './api.js';

import ReactTable from 'react-table';
import 'react-table/react-table.css';

const rawData = FunctionNameGoesHere();

const requestData = (pageSize, page, sorted, filtered) => {
  return new Promise((resolve, reject) => {
    let filteredData = rawData;

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

export default class QueueTable extends Component {
  constructor() {
    super();
    this.state = {
      data: [],
      pages: null,
      loading: true,
    };
    this.fetchData = this.fetchData.bind(this);
  }

  fetchData(state, instance) {
    // Whenever the table model changes, or the user sorts or changes pages, this method gets called and passed the current table model.
    // You can set the `loading` prop of the table to true to use the built-in one or show you're own loading bar if you want.
    this.setState({ loading: true });
    // Request the data however you want.  Here, we'll use our mocked service we created earlier
    requestData(state.pageSize, state.page, state.sorted, state.filtered).then(
      res => {
        // Now just get the rows of data to your React Table (and update anything else like total pages or loading)
        this.setState({
          data: res.rows,
          pages: res.pages,
          loading: false,
        });
      },
    );
  }

  render() {
    return (
      <div>
        <h1>Header that changes based on props</h1>
        <div>
          <ReactTable
            columns={[
              {
                Header: 'Status',
                accessor: 'status',
              },
              {
                Header: 'Locator #',
                accessor: 'locator',
              },
              {
                Header: 'Customer name',
                accessor: 'customer_name',
              },
              {
                Header: 'DOD ID',
                accessor: 'edipi',
              },
              {
                Header: 'Rank',
                accessor: 'rank',
              },
              {
                Header: 'Move type',
                accessor: 'orders_type',
              },
              {
                Header: 'Move date',
                accessor: 'move_date',
              },
              {
                Header: 'Created',
                accessor: 'created_at',
              },
              // {
              //   Header: 'Last modified',
              //   accessor: d =>
              //     <div dangerouslySetInnerHTML={{
              //       __html: d.last_modified_date + " by " + d.last_modified_name
              //     }}
              // },
            ]}
            // manual // Forces table not to paginate or sort automatically, so we can handle it server-side
            // data={data}
            // pages={pages} // Display the total number of pages
            // loading={loading} // Display the loading overlay when we need it
            onFetchData={this.fetchData} // Request new data when things change
            // filterable
            defaultPageSize={50}
            className="-striped -highlight"
          />
        </div>
      </div>
    );
  }
}

// export default QueueTable;
