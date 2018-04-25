import React from 'react';
import { render } from 'react-dom';
import _ from 'lodash';
import { makeData } from './Utils';
import { RetrieveDutyStations } from './api.js';

// Import React Table
import ReactTable from 'react-table';
import 'react-table/react-table.css';

const rawData = RetrieveDutyStations();
console.log(rawData);

// Example based on the one here: https://react-table.js.org/#/story/server-side-data
const requestData = (pageSize, page, sorted, filtered) => {
  return new Promise((resolve, reject) => {
    // You can retrieve your data however you want, in this case, we will just use some local data.
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
    rawData.then(makeRows => {
      console.log(makeRows.slice(pageSize * page, pageSize * page + pageSize));

      const res = {
        rows: makeRows.slice(pageSize * page, pageSize * page + pageSize),
        pages: Math.ceil(makeRows.length / pageSize),
      };

      // Here we'll simulate a server response with 500ms of delay.
      setTimeout(() => resolve(res), 500);
    });
  });
};

class Admin extends React.Component {
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
    const { data, pages, loading } = this.state;
    return (
      <div>
        <ReactTable
          columns={[
            {
              Header: 'Flair',
              Cell: row => (
                <div>
                  <img src="https://upload.wikimedia.org/wikipedia/commons/thumb/f/f1/Heart_coraz%C3%B3n.svg/130px-Heart_coraz%C3%B3n.svg.png" />
                </div>
              ),
            },
            {
              Header: 'ID',
              accessor: 'id',
            },
            {
              Header: 'Name',
              accessor: 'name',
            },
            {
              Header: 'Branch',
              accessor: 'branch',
            },
            {
              Header: 'Created At',
              accessor: 'created_at',
            },
            {
              Header: 'Updated At',
              accessor: 'updated_at',
            },
          ]}
          manual // Forces table not to paginate or sort automatically, so we can handle it server-side
          data={data}
          pages={pages} // Display the total number of pages
          loading={loading} // Display the loading overlay when we need it
          onFetchData={this.fetchData} // Request new data when things change
          // filterable
          defaultPageSize={50}
          className="-striped -highlight"
        />
        <br />
      </div>
    );
  }
}

export default Admin;
