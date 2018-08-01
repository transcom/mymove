import React, { Component } from 'react';
import { withRouter } from 'react-router';
import ReactTable from 'react-table';
import { connect } from 'react-redux';
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

  async fetchData() {
    this.setState({
      data: [],
      pages: null,
      loading: true,
    });

    const body = await RetrieveShipmentsForTSP(this.props.queueType);

    this.setState({
      data: body,
      pages: 1,
      loading: false,
    });
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
                Header: 'Locator',
                accessor: 'move.locator',
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
