import React, { Component } from 'react';
import { withRouter } from 'react-router';
import ReactTable from 'react-table';
import { connect } from 'react-redux';
import { capitalize } from 'lodash';
import 'react-table/react-table.css';
import { RetrieveShipmentsForTSP } from './api.js';
import { formatDate, formatDateTimeWithTZ, formatTimeAgo } from 'shared/formatters';
import FontAwesomeIcon from '@fortawesome/react-fontawesome';
import faSyncAlt from '@fortawesome/fontawesome-free-solid/faSyncAlt';

class QueueTable extends Component {
  constructor() {
    super();
    this.state = {
      data: [],
      pages: null,
      loading: true,
      refreshing: false, // only true when the user clicks the refresh button
      lastLoadedAt: new Date(),
      lastLoadedAtText: formatTimeAgo(new Date()),
      interval: setInterval(() => {
        this.setState({
          lastLoadedAtText: formatTimeAgo(this.state.lastLoadedAt),
        });
      }, 5000),
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

  openShipment(rowInfo) {
    this.props.history.push(`/shipments/${rowInfo.original.id}`, {
      referrerPathname: this.props.history.location.pathname,
    });
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
          refreshing: false,
          lastLoadedAt: new Date(),
        });
      }
    } catch (e) {
      this.setState({
        data: [],
        pages: 1,
        loading: false,
        refreshing: false,
        lastLoadedAt: new Date(),
      });
    }
  }

  refresh() {
    clearInterval(this.state.interval);

    this.setState({
      refreshing: true,
      lastLoadedAt: new Date(),
      interval: setInterval(() => {
        this.setState({
          lastLoadedAtText: formatTimeAgo(this.state.lastLoadedAt),
        });
      }, 5000),
    });

    this.fetchData();
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
        <h1 className="queue-heading">Queue: {titles[this.props.queueType]}</h1>
        <div className="queue-table">
          <span className="staleness-indicator" data-cy="staleness-indicator">
            Last updated {formatTimeAgo(this.state.lastLoadedAt)}
          </span>
          <span className={'refresh' + (this.state.refreshing ? ' focused' : '')} title="Refresh" aria-label="Refresh">
            <FontAwesomeIcon
              data-cy="refreshQueue"
              className="link-blue"
              icon={faSyncAlt}
              onClick={this.refresh.bind(this)}
              color="blue"
              size="lg"
              spin={!this.state.refreshing && this.state.loading}
            />
          </span>

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
                Cell: row => <span data-cy="locator">{row.value}</span>,
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
              'data-cy': 'queueTableRow',
              onDoubleClick: () => this.openShipment(rowInfo),
              onClick: () => this.openShipment(rowInfo),
            })}
          />
        </div>
      </div>
    );
  }
}

const mapStateToProps = state => ({});

export default withRouter(connect(mapStateToProps)(QueueTable));
