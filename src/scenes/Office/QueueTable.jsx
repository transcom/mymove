import React, { Component } from 'react';
import ReactTable from 'react-table';
import 'react-table/react-table.css';
import { RetrieveMovesForOffice } from './api.js';

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
    RetrieveMovesForOffice(this.props.queueType).then(
      response => {
        this.setState({
          data: response,
          pages: 1,
          loading: false,
        });
      },
      response => {
        // TODO: add error handling
      },
    );
  }

  render() {
    return (
      <div>
        <h1 style={{ textTransform: 'capitalize' }}>
          {this.props.queueType} Moves Queue
        </h1>
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
              {
                Header: 'Last modified',
                accessor: 'last_modified_date',
              },
            ]}
            data={this.state.data}
            loading={this.state.loading} // Display the loading overlay when we need it
            onFetchData={this.fetchData} // Request new data when things change
            pageSize={this.state.data.length}
            className="-striped -highlight"
            getTrProps={(state, rowInfo, column, instance) => ({
              onDoubleClick: e => alert('A row was clicked!'),
            })}
          />
        </div>
      </div>
    );
  }
}
