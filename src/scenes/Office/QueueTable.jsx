/* eslint-disable security/detect-object-injection */
import React, { useEffect, useState } from 'react';
import { withRouter } from 'react-router';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { get } from 'lodash';
import Alert from 'shared/Alert';
import { formatTimeAgo } from 'shared/formatters';
import { setUserIsLoggedIn } from 'shared/Data/users';
import { SHIPMENT_OPTIONS } from 'shared/constants';
import { defaultColumns } from './queueTableColumns';
import TableQueue from 'components/Table/TableQueue';

import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import 'react-table-6/react-table.css';

const titles = {
  new: 'New moves',
  troubleshooting: 'Troubleshooting',
  ppm_payment_requested: 'Payment requested',
  all: 'All moves',
  ppm_completed: 'Completed moves',
  ppm_approved: 'Approved moves',
};

const QueueTable = ({ flashMessageLines, history, showFlashMessage, queueType, retrieveMoves, setUserIsLoggedIn }) => {
  const [data, setData] = useState([]);
  const [origDutyStationData, setOrigDutyStationData] = useState([]);
  const [destDutyStationData, setDestDutyStationData] = useState([]);
  const [isError, setIsError] = useState(false);
  const [isSuccess, setIsSuccess] = useState(true);
  const [loading, setLoading] = useState(true);
  const [refreshing, setRefreshing] = useState(true);
  const [lastLoadedAt, setLastLoadedAt] = useState(new Date());
  const [lastLoadedAtText, setLastLoadedAtText] = useState(formatTimeAgo(new Date()));
  const [loadingQueue, setLoadingQueue] = useState(queueType);

  useEffect(() => {
    const id = setInterval(() => {
      setLastLoadedAtText(formatTimeAgo(lastLoadedAt));
    }, 5000);

    return () => {
      clearInterval(id);
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  useEffect(() => {
    setLastLoadedAtText(formatTimeAgo(lastLoadedAt));
  }, [lastLoadedAt]);

  useEffect(() => {
    const fetchData = async () => {
      setData([]);
      setLoading(true);
      setLoadingQueue(queueType);

      try {
        const body = await retrieveMoves(queueType);
        // grab all destination duty station and remove duplicates
        // this will build on top of the current duty stations list we see from the data
        const sortOrigDutyStationData = (origDutyStationData) => {
          const origDutyStationDataSet = new Set(origDutyStationData);
          body.forEach((value) => {
            if (value.origin_duty_station_name !== undefined && value.origin_duty_station_name !== '') {
              origDutyStationDataSet.add(value.origin_duty_station_name);
            }
          });
          return [...origDutyStationDataSet].sort();
        };

        const sortDestDutyStationData = (destDutyStationData) => {
          const destDutyStationDataSet = new Set(destDutyStationData);
          body.forEach((value) => {
            if (value.destination_duty_station_name !== undefined && value.destination_duty_station_name !== '') {
              destDutyStationDataSet.add(value.destination_duty_station_name);
            }
          });
          return [...destDutyStationDataSet].sort();
        };

        // Only update the queue list if the request that is returning
        // is for the same queue as the most recent request.
        if (loadingQueue === queueType) {
          setData(body);
          setOrigDutyStationData(sortOrigDutyStationData);
          setDestDutyStationData(sortDestDutyStationData);
          setLoading(false);
          setRefreshing(false);
          setLastLoadedAt(new Date());
          setIsError(false);
          setIsSuccess(true);
        }
      } catch (e) {
        setData([]);
        setOrigDutyStationData([]);
        setDestDutyStationData([]);
        setLoading(false);
        setRefreshing(false);
        setLastLoadedAt(new Date());
        setIsError(true);
        setIsSuccess(false);

        // redirect to home page if unauthorized
        if (e.status === 401) {
          setUserIsLoggedIn(false);
        }
      }
    };
    if (refreshing || loadingQueue !== queueType) {
      fetchData();
    }
  }, [loadingQueue, queueType, retrieveMoves, setUserIsLoggedIn, refreshing]);

  const refresh = () => {
    setRefreshing(true);
    setLastLoadedAt(new Date());
  };

  const openMove = (rowInfo) => {
    history.push(`/queues/new/moves/${rowInfo.id}`, {
      referrerPathname: history.location.pathname,
    });
  };

  const useQuery = ({ sort, order, filters = [], currentPage = 1, currentPageSize = 20 }) => {
    return {
      queueResult: { data: data, totalCount: data.length, page: 1, perPage: 100 },
      isLoading: loading,
      isError: isError,
      isSuccess: isSuccess,
    };
  };

  const showColumns = defaultColumns(origDutyStationData, destDutyStationData);

  const defaultSort = () => {
    if (['all'].includes(queueType)) {
      return [{ id: 'locator', asc: true }];
    }
    return [{ id: 'move_date', asc: true }];
  };

  data.forEach((row) => {
    row.shipments = SHIPMENT_OPTIONS.PPM;

    if (row.ppm_status !== null) {
      row.synthetic_status = row.ppm_status;
    } else {
      row.synthetic_status = row.status;
    }
  });

  const defaultSortColumn = defaultSort();

  return (
    <div>
      {showFlashMessage ? (
        <Alert type="success" heading="Success">
          {flashMessageLines.join('\n')}
          <br />
        </Alert>
      ) : null}
      <div className="queue-table">
        <span className="staleness-indicator" data-testid="staleness-indicator">
          Last updated {lastLoadedAtText}
        </span>
        <span className={'refresh' + (refreshing ? ' focused' : '')} title="Refresh" aria-label="Refresh">
          <FontAwesomeIcon
            data-testid="refreshQueue"
            className="link-blue"
            icon="sync-alt"
            onClick={refresh}
            color="blue"
            size="lg"
            spin={!refreshing && loading}
          />
        </span>
        <TableQueue
          showFilters
          showPagination={false}
          manualSortBy
          defaultCanSort
          defaultSortedColumns={defaultSortColumn}
          disableMultiSort
          disableSortBy={false}
          columns={showColumns}
          title={titles[queueType]}
          handleClick={openMove}
          useQueries={useQuery}
        />
      </div>
    </div>
  );
};

QueueTable.defaultProps = {
  moveLocator: '',
  firstName: '',
  lastName: '',
};

const mapStateToProps = (state) => {
  return {
    showFlashMessage: get(state, 'flashMessages.display', false),
    flashMessageLines: get(state, 'flashMessages.messageLines', false),
  };
};

function mapDispatchToProps(dispatch) {
  return bindActionCreators({ setUserIsLoggedIn }, dispatch);
}

export default withRouter(connect(mapStateToProps, mapDispatchToProps)(QueueTable));
