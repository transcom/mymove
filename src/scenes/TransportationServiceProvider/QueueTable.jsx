import React, { Component } from 'react';
import { withRouter } from 'react-router';
import ReactTable from 'react-table';
import { connect } from 'react-redux';
import { get } from 'lodash';
import 'react-table/react-table.css';
import { RetrieveMovesForTSP } from './api.js';

class QueueTable extends Component {
  constructor() {
    super();
    this.state = {
      data: [],
      pages: null,
      loading: true,
    };
    this.fetchData = this.fetchData.bind(this);
  }

  componentDidMount() {
    this.fetchData();
  }

  componentDidUpdate(prevProps) {
    if (this.props.queueType !== prevProps.queueType) {
      this.fetchData();
    }
  }

  static defaultProps = {
    moveLocator: '',
    firstName: '',
    lastName: '',
  };

  fetchData() {
    RetrieveMovesForTSP(this.props.queueType).then(
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
    const titles = {
      all: 'All Moves',
      other: 'OTHER',
    };

    return (
      <div>
        <h1>Queue: {titles[this.props.queueType]}</h1>
        <div className="queue-table">
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
            pageSize={this.state.data.length}
            className="-striped -highlight"
            getTrProps={(state, rowInfo) => ({
              onDoubleClick: e =>
                this.props.history.push(
                  `${this.props.queueType}/moves/${rowInfo.original.id}`,
                ),
            })}
          />
        </div>
      </div>
    );
  }
}

const mapStateToProps = state => ({
  moveLocator: get(state, 'office.officeMove.locator', 'Unloaded'),
  firstName: get(state, 'office.officeServiceMember.first_name', 'Unloaded'),
  lastName: get(state, 'office.officeServiceMember.last_name', 'Unloaded'),
});

export default withRouter(connect(mapStateToProps)(QueueTable));
