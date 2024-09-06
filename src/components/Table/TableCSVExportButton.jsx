import React, { useState, useRef, useContext } from 'react';
import { CSVLink } from 'react-csv';
import { Button } from '@trussworks/react-uswds';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import moment from 'moment';
import PropTypes from 'prop-types';

import SelectedGblocContext from 'components/Office/GblocSwitcher/SelectedGblocContext';

const TableCSVExportButton = ({
  labelText,
  filePrefix,
  totalCount,
  tableColumns,
  hiddenColumns,
  queueFetcher,
  queueFetcherKey,
  paramSort,
  paramFilters,
  className,
  isHeadquartersUser,
}) => {
  const [isLoading, setIsLoading] = useState(false);
  const [csvRows, setCsvRows] = useState([]);
  const csvLinkRef = useRef(null);
  const { id: sortColumn, desc: sortOrder } = paramSort.length ? paramSort[0] : {};

  const gblocContext = useContext(SelectedGblocContext);
  const { selectedGbloc } = isHeadquartersUser && gblocContext ? gblocContext : { selectedGbloc: undefined };

  const formatDataForExport = (data, columns = tableColumns) => {
    const formattedData = [];
    data.forEach((row) => {
      const formattedRow = {};
      columns
        .filter((column) => !hiddenColumns.includes(column.id))
        .forEach((column) => {
          if (column.exportValue) {
            formattedRow[column.Header] = column.exportValue(row);
          } else if (typeof column.accessor === 'function') {
            formattedRow[column.Header] = column.accessor(row);
          } else {
            formattedRow[column.Header] = row[column.accessor];
          }
        });
      formattedData.push(formattedRow);
    });

    return formattedData;
  };

  const handleCsvExport = async () => {
    setIsLoading(true);
    const response = await queueFetcher(queueFetcherKey, {
      sort: sortColumn,
      order: sortOrder ? 'desc' : 'asc',
      filters: paramFilters,
      currentPageSize: totalCount,
      viewAsGBLOC: selectedGbloc,
    });

    const formattedData = formatDataForExport(response[queueFetcherKey]);
    setCsvRows(formattedData);

    csvLinkRef.current?.click();
    setIsLoading(false);
  };

  return (
    <p>
      <Button
        className={className}
        onClick={handleCsvExport}
        data-test-id="csv-export-btn-visible"
        disabled={!totalCount}
        aria-disabled={!totalCount}
        tabIndex={0}
      >
        <span data-test-id="csv-export-btn-text">{labelText}</span>{' '}
        <FontAwesomeIcon icon={isLoading ? 'spinner' : 'download'} spin={isLoading} />
      </Button>
      <CSVLink
        data-test-id="csv-export-btn-hidden"
        filename={`${filePrefix}-${moment().toISOString()}.csv`}
        data={csvRows}
        className="hidden"
        tabIndex={-1}
      >
        <span ref={csvLinkRef} />
      </CSVLink>
    </p>
  );
};

TableCSVExportButton.propTypes = {
  // labelText is the label displayed on this export to CSV button
  labelText: PropTypes.string,
  // filePrefix is a string used in the exported file's name before a timestamp
  filePrefix: PropTypes.string,
  // totalCount is the number of items in the queue, used to send an accurate page size in the request
  totalCount: PropTypes.number.isRequired,
  // tableColumns is the columns object used by the table and contains column header text, an accessor key or function
  tableColumns: PropTypes.arrayOf(PropTypes.object).isRequired,
  // hiddenColumns is an array of column ids to not include in the export
  hiddenColumns: PropTypes.arrayOf(PropTypes.string),
  // queueFetcher is the function to handle refetching non-paginated queue data
  queueFetcher: PropTypes.func.isRequired,
  // queueFetcherKey is the key the queue data is stored under in the retrun value of queueFetchers
  queueFetcherKey: PropTypes.string.isRequired,
  // paramSort is the sort column and order currently applied to the queue
  paramSort: PropTypes.array,
  // paramSort is the filter columns and values currently applied to the queue
  paramFilters: PropTypes.array,
  // isHeadquartersUser identifies if the active role is a headquarters user to allow switching GBLOCs
  isHeadquartersUser: PropTypes.bool,
};

TableCSVExportButton.defaultProps = {
  labelText: 'Export to CSV',
  filePrefix: 'Moves',
  hiddenColumns: [],
  paramSort: [],
  paramFilters: [],
  isHeadquartersUser: false,
};

export default TableCSVExportButton;
