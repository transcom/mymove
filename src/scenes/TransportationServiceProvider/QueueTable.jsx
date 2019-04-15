import React, { Component } from 'react';
import { withRouter } from 'react-router';
import ReactTable from 'react-table';
import { connect } from 'react-redux';
import { capitalize } from 'lodash';
import 'react-table/react-table.css';
import { RetrieveShipmentsForTSP } from './api.js';
import { formatDate, formatDateTimeWithTZ } from 'shared/formatters';

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
    const loadingQueueType = this.props.queueType;

    this.setState({
      data: [],
      pages: null,
      loading: true,
      loadingQueue: loadingQueueType,
    });

    // Catch any errors here and render an empty queue
    try {
      const body = await RetrieveShipmentsForTSP(this.props.queueType);

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
      new: 'New Shipments',
      accepted: 'Accepted Shipments',
      completed: 'Completed Shipments',
      approved: 'Approved Shipments',
      in_transit: 'In Transit Shipments',
      delivered: 'Delivered Shipments',
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
                Cell: row => <span className="status">{capitalize(row.value.replace('_', ' '))}</span>,
              },
              {
                Header: 'GBL',
                accessor: 'gbl_number',
              },
              {
                Header: 'Customer name',
                accessor: 'service_member',
                Cell: row => (
                  <span className="customer_name">
                    {row.value.last_name}, {row.value.first_name}
                  </span>
                ),
              },
              {
                Header: 'Locator #',
                accessor: 'move.locator',
              },
              {
                Header: 'Channel',
                accessor: 'traffic_distribution_list',
                Cell: row => (
                  <span className="channel">
                    {row.value.source_rate_area} to Region {row.value.destination_region}
                  </span>
                ),
              },
              {
                Header: 'Pickup Date',
                id: 'pickup_date',
                accessor: d =>
                  d.actual_pickup_date ||
                  d.pm_survey_planned_pickup_date ||
                  d.requested_pickup_date ||
                  d.original_pickup_date,
                Cell: row => <span className="pickup_date">{formatDate(row.value)}</span>,
              },
              {
                Header: 'Delivered on',
                id: 'delivered_on',
                accessor: 'actual_delivery_date',
                Cell: row => <span className="delivered_on">{formatDate(row.value)}</span>,
              },
              {
                Header: 'Last modified',
                accessor: 'updated_at',
                Cell: row => <span className="updated_at">{formatDateTimeWithTZ(row.value)}</span>,
              },
            ]}
            data={this.state.data}
            defaultSorted={[{ id: 'pickup_date', asc: true }]}
            loading={this.state.loading} // Display the loading overlay when we need it
            pageSize={this.state.data.length}
            className="-striped -highlight"
            showPagination={false}
            getTrProps={(state, rowInfo) => ({
              onDoubleClick: e => this.props.history.push(`/shipments/${rowInfo.original.id}`),
            })}
          />
        </div>
      </div>
    );
  }
}

const mapStateToProps = state => ({});

export default withRouter(connect(mapStateToProps)(QueueTable));
