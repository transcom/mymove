import React, { Component } from 'react';
import { withRouter } from 'react-router';
import ReactTable from 'react-table-6';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { get } from 'lodash';
import Alert from 'shared/Alert';
import { formatTimeAgo } from 'shared/formatters';
import { logOut as logOutAction } from 'store/auth/actions';
import { SHIPMENT_OPTIONS } from 'shared/constants';
import { defaultColumns } from './queueTableColumns';

import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import 'react-table-6/react-table.css';

class QueueTable extends Component {
  constructor() {
    super();
    this.state = {
      data: [],
      origDutyStationData: [],
      destDutyStationData: [],
      pages: null,
      loading: true,
      refreshing: false, // only true when the user clicks the refresh button
      lastLoadedAt: new Date(),
      lastLoadedAtText: formatTimeAgo(new Date()),
      interval: setInterval(() => {
        this.setState((prevState) => ({
          lastLoadedAtText: formatTimeAgo(prevState.lastLoadedAt),
        }));
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

  componentWillUnmount() {
    clearInterval(this.state.interval);
  }

  openMove(rowInfo) {
    this.props.history.push(`/queues/new/moves/${rowInfo.original.id}`, {
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
      // grab all destination duty station and remove duplicates
      // this will build on top of the current duty stations list we see from the data
      let origDutyStationDataSet = new Set(this.getOriginDutyStations());
      let destDutyStationDataSet = new Set(this.getDestinationDutyStations());
      body.forEach((value) => {
        if (value.origin_duty_location_name !== undefined && value.origin_duty_location_name !== '') {
          origDutyStationDataSet.add(value.origin_duty_location_name);
        }
        if (value.destination_duty_station_name !== undefined && value.destination_duty_station_name !== '') {
          destDutyStationDataSet.add(value.destination_duty_station_name);
        }
      });

      // Only update the queue list if the request that is returning
      // is for the same queue as the most recent request.
      if (this.state.loadingQueue === loadingQueueType) {
        this.setState({
          data: body,
          origDutyStationData: [...origDutyStationDataSet].sort(),
          destDutyStationData: [...destDutyStationDataSet].sort(),
          pages: 1,
          loading: false,
          refreshing: false,
          lastLoadedAt: new Date(),
        });
      }
    } catch (e) {
      this.setState({
        data: [],
        origDutyStationData: [],
        destDutyStationData: [],
        pages: 1,
        loading: false,
        refreshing: false,
        lastLoadedAt: new Date(),
      });
      // redirect to home page if unauthorized
      if (e.status === 401) {
        this.props.logOut();
      }
    }
  }

  refresh() {
    clearInterval(this.state.interval);

    this.setState({
      refreshing: true,
      lastLoadedAt: new Date(),
      interval: setInterval(() => {
        this.setState((prevState) => ({
          lastLoadedAtText: formatTimeAgo(prevState.lastLoadedAt),
        }));
      }, 5000),
    });

    this.fetchData();
  }

  getDestinationDutyStations = () => {
    return this.state.destDutyStationData;
  };

  getOriginDutyStations = () => {
    return this.state.origDutyStationData;
  };

  render() {
    const titles = {
      new: 'New moves',
      troubleshooting: 'Troubleshooting',
      ppm_payment_requested: 'Payment requested',
      all: 'All moves',
      ppm_completed: 'Completed moves',
      ppm_approved: 'Approved moves',
    };

    const showColumns = defaultColumns(this);

    const defaultSort = (queueType) => {
      if (['all'].includes(queueType)) {
        return [{ id: 'locator', asc: true }];
      }
      return [{ id: 'move_date', asc: true }];
    };

    this.state.data.forEach((row) => {
      row.shipments = SHIPMENT_OPTIONS.PPM;

      if (row.ppm_status !== null) {
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
        <h1 className="queue-heading">{titles[this.props.queueType]}</h1>
        <div className="queue-table">
          <span className="staleness-indicator" data-testid="staleness-indicator">
            Last updated {formatTimeAgo(this.state.lastLoadedAt)}
          </span>
          <span className={'refresh' + (this.state.refreshing ? ' focused' : '')} title="Refresh" aria-label="Refresh">
            <FontAwesomeIcon
              data-testid="refreshQueue"
              className="link-blue"
              icon="sync-alt"
              onClick={this.refresh.bind(this)}
              color="blue"
              size="lg"
              spin={!this.state.refreshing && this.state.loading}
            />
          </span>
          <ReactTable
            columns={showColumns}
            data={this.state.data}
            loading={this.state.loading} // Display the loading overlay when we need it
            defaultSorted={defaultSort(this.props.queueType)}
            pageSize={this.state.data.length}
            className="-striped -highlight"
            showPagination={false}
            getTrProps={(state, rowInfo) => ({
              'data-testid': 'queueTableRow',
              onDoubleClick: () => this.openMove(rowInfo),
              onClick: () => this.openMove(rowInfo),
            })}
            getTheadFilterThProps={() => {
              return {
                style: {
                  display: 'flex',
                  flexDirection: 'column',
                  justifyContent: 'center',
                  position: 'inherit',
                  overflow: 'inherit',
                },
              };
            }}
            getTableProps={() => {
              return {
                style: { overflow: 'inherit' },
              };
            }}
          />
        </div>
      </div>
    );
  }
}

const mapStateToProps = (state) => {
  return {
    showFlashMessage: get(state, 'flashMessages.display', false),
    flashMessageLines: get(state, 'flashMessages.messageLines', false),
  };
};

function mapDispatchToProps(dispatch) {
  return bindActionCreators({ logOut: logOutAction }, dispatch);
}

export default withRouter(connect(mapStateToProps, mapDispatchToProps)(QueueTable));
