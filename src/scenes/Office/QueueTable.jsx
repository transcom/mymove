import React, { Component } from 'react';
import { withRouter } from 'react-router';
import ReactTable from 'react-table';
import { connect } from 'react-redux';
import { get } from 'lodash';
import 'react-table/react-table.css';
import { RetrieveMovesForOffice } from './api.js';
import Alert from 'shared/Alert';
import { selectServiceMemberForMove } from 'shared/Entities/modules/serviceMembers';

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

  async fetchData() {
    const loadingQueueType = this.props.queueType;

    this.setState({
      data: [],
      pages: null,
      loading: true,
      loadingQueue: loadingQueueType,
    });

    // Catch any errors here and render an empty queue
    try {
      const body = await RetrieveMovesForOffice(this.props.queueType);

      // Only update the queue list if the request that is returning
      // is for the same queue as the most recent request.
      if (this.state.loadingQueue === loadingQueueType) {
        this.setState({
          data: body,
          pages: 1,
          loading: false,
        });
      }
    } catch (e) {
      this.setState({
        data: [],
        pages: 1,
        loading: false,
      });
    }
  }

  render() {
    const titles = {
      new: 'New Moves',
      troubleshooting: 'Troubleshooting',
      ppm: 'PPMs',
      hhg_accepted: 'Accepted HHGs',
      hhg_delivered: 'Delivered HHGs',
      hhg_completed: 'Completed HHGs',
      all: 'All Moves',
    };

    this.state.data.forEach(row => {
      if (row.ppm_status === 'PAYMENT_REQUESTED') {
        row.synthetic_status = row.ppm_status;
      } else {
        row.synthetic_status = row.status;
      }
    });

    return (
      <div>
        {this.props.flashMessage ? (
          <Alert type="success" heading="Success">
            Move #{this.props.moveLocator} for {this.props.lastName}, {this.props.firstName} has been canceled <br />
            An email confirmation has been sent to the customer.
            <br />
          </Alert>
        ) : null}
        <h1>Queue: {titles[this.props.queueType]}</h1>
        <div className="queue-table">
          <ReactTable
            columns={[
              {
                Header: 'Status',
                accessor: 'synthetic_status',
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
              onDoubleClick: e => this.props.history.push(`new/moves/${rowInfo.original.id}`),
            })}
          />
        </div>
      </div>
    );
  }
}

const mapStateToProps = state => {
  const moveId = get(state, 'office.officeMove.id');
  const serviceMember = selectServiceMemberForMove(state, moveId);
  return {
    flashMessage: get(state, 'office.flashMessage', false),
    moveLocator: get(state, 'office.officeMove.locator', 'Unloaded'),
    firstName: serviceMember.first_name,
    lastName: serviceMember.last_name,
  };
};

export default withRouter(connect(mapStateToProps)(QueueTable));
