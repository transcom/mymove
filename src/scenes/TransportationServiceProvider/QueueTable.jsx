import React, { Component } from 'react';
import { withRouter } from 'react-router';
import ReactTable from 'react-table';
import { connect } from 'react-redux';
import { get } from 'lodash';
import 'react-table/react-table.css';
import { RetrieveShipmentsForTSP } from './api.js';

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
    RetrieveShipmentsForTSP(this.props.queueType).then(
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
      new: 'New Shipments',
      all: 'All Shipments',
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
                Header: 'GBL',
                accessor: 'source_gbloc',
              },
              {
                Header: 'Requested Pickup Date',
                accessor: 'requested_pickup_date',
              },
              {
                Header: 'Pickup Date',
                accessor: 'pickup_date',
              },
              {
                Header: 'Delivery Date',
                accessor: 'delivery_date',
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

const mapStateToProps = state => ({});

export default withRouter(connect(mapStateToProps)(QueueTable));
