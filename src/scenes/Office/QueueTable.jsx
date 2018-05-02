import React, { Component } from 'react';
import { render } from 'react-dom';
import ReactTable from 'react-table';
import 'react-table/react-table.css';
import { requestData, RetrieveMovesForOffice } from './api.js';

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
    RetrieveMovesForOffice(this.props.queueType)
      // requestData(
      //   this.props.queueType,
      //   state.pageSize,
      //   state.page,
      //   state.sorted,
      //   state.filtered,
      // )
      .then(
        response => {
          // Now just get the rows of data to your React Table (and update anything else like total pages or loading)
          console.log('response, ', response);
          this.setState({
            data: response.rows,
            pages: 1,
            loading: false,
          });
        },
        response => {
          console.log('hit the second res function!');
        },
      );
  }

  render() {
    return (
      <div>
        <h1>{this.props.queueType} Moves Queue</h1>
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
                accessor: 'first_name',
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
              {
                Header: 'Last modified',
                accessor: 'updated_at',
              },
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
