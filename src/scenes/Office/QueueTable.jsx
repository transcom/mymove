import React, { Component } from 'react';
import { withRouter } from 'react-router';
import ReactTable from 'react-table';
import { connect } from 'react-redux';
import { get, capitalize } from 'lodash';
import 'react-table/react-table.css';
import Alert from 'shared/Alert';
import { formatDate, formatDateTimeWithTZ, formatTimeAgo } from 'shared/formatters';
import FontAwesomeIcon from '@fortawesome/react-fontawesome';
import faClock from '@fortawesome/fontawesome-free-solid/faClock';
import './office.scss';
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

  openMove(rowInfo) {
    this.props.history.push(`new/moves/${rowInfo.original.id}`, {
      referrerPathname: this.props.history.location.pathname,
    });
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
      const body = await this.props.retrieveMoves(this.props.queueType);

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
      new: 'New Moves',
      troubleshooting: 'Troubleshooting',
      ppm: 'PPMs',
      hhg_accepted: 'Accepted HHGs',
      hhg_delivered: 'Delivered HHGs',
      all: 'All Moves',
    };

    this.state.data.forEach(row => {
      if (this.props.queueType === 'new' && row.ppm_status && row.hhg_status) {
        row.shipments = 'HHG, PPM';
      } else if (row.ppm_status && !row.hhg_status) {
        row.shipments = 'PPM';
      } else {
        row.shipments = 'HHG';
      }

      if (this.props.queueType === 'ppm' && row.ppm_status !== null) {
        row.synthetic_status = row.ppm_status;
      } else {
        row.synthetic_status = row.status;
      }
    });

    return (
      <div>
        {this.props.showFlashMessage ? (
          <Alert type="success" heading="Success">
            {this.props.flashMessageLines.join('\n')}
            <br />
          </Alert>
        ) : null}
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
                Header: <FontAwesomeIcon icon={faClock} />,
                id: 'clockIcon',
                accessor: row => row.synthetic_status,
                Cell: row =>
                  row.value === 'PAYMENT_REQUESTED' || row.value === 'SUBMITTED' ? (
                    <span data-cy="ppm-queue-icon">
                      <FontAwesomeIcon icon={faClock} style={{ color: 'orange' }} />
                    </span>
                  ) : (
                    ''
                  ),
                width: 50,
                show: this.props.queueType === 'ppm',
              },
              {
                Header: 'Status',
                accessor: 'synthetic_status',
                Cell: row => (
                  <span className="status" data-cy="status">
                    {capitalize(row.value && row.value.replace('_', ' '))}
                  </span>
                ),
              },
              {
                Header: 'Customer name',
                accessor: 'customer_name',
              },
              {
                Header: 'DoD ID',
                accessor: 'edipi',
              },
              {
                Header: 'Rank',
                accessor: 'rank',
                Cell: row => <span className="rank">{row.value && row.value.replace('_', '-')}</span>,
              },
              {
                Header: 'Shipments',
                accessor: 'shipments',
                show: this.props.queueType === 'new',
              },
              {
                Header: 'Locator #',
                accessor: 'locator',
                Cell: row => <span data-cy="locator">{row.value}</span>,
              },
              {
                Header: 'GBL',
                accessor: 'gbl_number',
                show: this.props.queueType !== 'ppm',
              },
              {
                Header: 'Move date',
                accessor: 'move_date',
                Cell: row => <span className="move_date">{formatDate(row.value)}</span>,
              },
              {
                Header: 'Last modified',
                accessor: 'last_modified_date',
                Cell: row => <span className="updated_at">{formatDateTimeWithTZ(row.value)}</span>,
                show: this.props.queueType !== 'new',
              },
              {
                Header: 'Submitted',
                accessor: 'submitted_date',
                Cell: row => <span className="submitted_date">{formatDateTimeWithTZ(row.value)}</span>,
                show: this.props.queueType === 'new',
              },
            ]}
            data={this.state.data}
            loading={this.state.loading} // Display the loading overlay when we need it
            defaultSorted={[{ id: 'move_date', asc: true }]}
            pageSize={this.state.data.length}
            className="-striped -highlight"
            showPagination={false}
            getTrProps={(state, rowInfo) => ({
              'data-cy': 'queueTableRow',
              onDoubleClick: () => this.openMove(rowInfo),
              onClick: () => this.openMove(rowInfo),
            })}
          />
        </div>
      </div>
    );
  }
}

const mapStateToProps = state => {
  return {
    showFlashMessage: get(state, 'flashMessages.display', false),
    flashMessageLines: get(state, 'flashMessages.messageLines', false),
  };
};

export default withRouter(connect(mapStateToProps)(QueueTable));
