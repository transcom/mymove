import React, { Component } from 'react';
import { withRouter } from 'react-router';
import ReactTable from 'react-table';
import { connect } from 'react-redux';
import { get } from 'lodash';
import 'react-table/react-table.css';
import { RetrieveMovesForOffice } from './api.js';
import Alert from 'shared/Alert';

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
    const titles = {
      new: 'New Moves',
      troubleshooting: 'Troubleshooting',
      ppm: 'PPMs',
      all: 'All Moves',
    };

    return (
      <div>
        {this.props.flashMessage ? (
          <Alert type="success" heading="Success">
            Move #{this.props.moveLocator} for {this.props.lastName},{' '}
            {this.props.firstName} has been cancelled <br />
            An email confirmation has been sent to the customer.<br />
          </Alert>
        ) : null}
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
                this.props.history.push(`new/moves/${rowInfo.original.id}`),
            })}
          />
        </div>
      </div>
    );
  }
}

const mapStateToProps = state => ({
  flashMessage: get(state, 'office.flashMessage', false),
  moveLocator: get(state, 'office.officeMove.locator', 'Unloaded'),
  firstName: get(state, 'office.officeServiceMember.first_name', 'Unloaded'),
  lastName: get(state, 'office.officeServiceMember.last_name', 'Unloaded'),
});

export default withRouter(connect(mapStateToProps)(QueueTable));
