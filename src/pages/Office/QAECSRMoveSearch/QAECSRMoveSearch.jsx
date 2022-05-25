import React, { useMemo, useState } from 'react';
import { Formik, Field } from 'formik';
import { withRouter } from 'react-router-dom';
import { Button, Radio } from '@trussworks/react-uswds';
import * as Yup from 'yup';
import classnames from 'classnames';

import styles from './QAECSRMoveSearch.module.scss';

import { HistoryShape } from 'types/router';
import formStyles from 'styles/form.module.scss';
import { Form } from 'components/form/Form';
import TextField from 'components/form/fields/TextField/TextField';
import { createHeader } from 'components/Table/utils';
import { useQAECSRMoveSearchQueries } from 'hooks/queries';
import { serviceMemberAgencyLabel } from 'utils/formatters';
import MultiSelectCheckBoxFilter from 'components/Table/Filters/MultiSelectCheckBoxFilter';
import SelectFilter from 'components/Table/Filters/SelectFilter';
import { BRANCH_OPTIONS, MOVE_STATUS_OPTIONS, MOVE_STATUS_LABELS } from 'constants/queues';
import SearchResultsTable from 'components/Table/SearchResultsTable';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';

const columns = (showBranchFilter = true) => [
  createHeader('Move code', 'locator', {
    id: 'locator',
    isFilterable: true,
  }),
  createHeader('DOD ID', 'customer.dodID', {
    id: 'dodID',
    isFilterable: true,
  }),
  createHeader(
    'Customer name',
    (row) => {
      return `${row.customer.last_name}, ${row.customer.first_name}`;
    },
    {
      id: 'lastName',
      isFilterable: true,
    },
  ),
  createHeader(
    'Status',
    (row) => {
      return MOVE_STATUS_LABELS[`${row.status}`];
    },
    {
      id: 'status',
      isFilterable: true,
      // eslint-disable-next-line react/jsx-props-no-spreading
      Filter: (props) => <MultiSelectCheckBoxFilter options={MOVE_STATUS_OPTIONS} {...props} />,
    },
  ),
  createHeader(
    'Origin ZIP',
    (row) => {
      return `${row.originDutyLocation.address.postalCode}`;
    },
    {
      id: 'originZIP',
      isFilterable: true,
    },
  ),
  createHeader(
    'Destination ZIP',
    (row) => {
      return `${row.destinationDutyLocation.address.postalCode}`;
    },
    {
      id: 'destinationZIP',
      isFilterable: true,
    },
  ),
  createHeader(
    'Branch',
    (row) => {
      return serviceMemberAgencyLabel(row.customer.agency);
    },
    {
      id: 'branch',
      isFilterable: showBranchFilter,
      Filter: (props) => (
        // eslint-disable-next-line react/jsx-props-no-spreading
        <SelectFilter options={BRANCH_OPTIONS} {...props} />
      ),
    },
  ),
  createHeader('# of shipments', 'shipmentsCount', { disableSortBy: true }),
];

const validationSchema = Yup.object().shape({
  searchType: Yup.string().required('searchtype error'),
  searchText: Yup.string().when('searchType', {
    is: 'moveCode',
    then: Yup.string().length(6, 'Move Code must be exactly 6 characters'),
    otherwise: Yup.string().length(10, 'DOD ID must be exactly 10 characters'),
  }),
});

const QAECSRMoveSearch = ({ history }) => {
  const [search, setSearch] = useState({ moveCode: null, dodID: null });
  const [searchHappened, setSearchHappened] = useState(false);

  const handleClick = (values) => {
    history.push(`/moves/${values.locator}/details`);
  };
  const onSubmit = (values) => {
    const payload = {
      moveCode: null,
      dodID: null,
    };
    if (values.searchType === 'moveCode') {
      payload.moveCode = values.searchText;
    } else if (values.searchType === 'dodID') {
      payload.dodID = values.searchText;
    }
    setSearch(payload);
    setSearchHappened(true);
  };

  const { searchResult, isLoading, isError } = useQAECSRMoveSearchQueries({
    moveCode: search.moveCode,
    dodID: search.dodID,
  });

  const { data = [] } = searchResult;
  const tableColumns = useMemo(() => columns(true), []);
  const tableData = useMemo(() => data, [data]);
  return (
    <div className={styles.QAECSRMoveSearchPage}>
      <h1>Search for a move</h1>
      <Formik
        initialValues={{ searchType: 'moveCode', searchText: '' }}
        onSubmit={onSubmit}
        validationSchema={validationSchema}
      >
        {(formik) => {
          return (
            <Form
              className={classnames(formStyles.form, styles.QAECSRMoveSearch)}
              onSubmit={formik.handleSubmit}
              role="search"
            >
              <p>What do you want to search for?</p>
              <div role="group" className={formStyles.radioGroup}>
                <Field
                  as={Radio}
                  id="radio-picked-movecode"
                  type="radio"
                  name="searchType"
                  value="moveCode"
                  title="Move Code"
                  label="Move Code"
                />
                <Field
                  as={Radio}
                  id="radio-picked-dodid"
                  type="radio"
                  name="searchType"
                  value="dodID"
                  title="DOD ID"
                  label="DOD ID"
                />
              </div>
              <div className={classnames(styles.searchBar)}>
                <TextField
                  id="searchText"
                  className="usa-search__input"
                  label="Search"
                  name="searchText"
                  type="search"
                />
                <Button className={classnames(styles.searchButton)} type="submit" disabled={!formik.isValid}>
                  Search
                </Button>
              </div>
            </Form>
          );
        }}
      </Formik>
      {isLoading && <LoadingPlaceholder />}
      {isError && <SomethingWentWrong />}
      {searchHappened && (
        <SearchResultsTable
          showFilters
          showPagination
          defaultCanSort
          defaultSortedColumns={[{ id: 'status', desc: false }]}
          disableMultiSort
          manualSortBy={false}
          manualFilters={false}
          disableSortBy={false}
          columns={tableColumns}
          title="Results"
          handleClick={handleClick}
          useQueries={useQAECSRMoveSearchQueries}
          data={tableData}
        />
      )}
    </div>
  );
};

QAECSRMoveSearch.propTypes = {
  history: HistoryShape.isRequired,
};

export default withRouter(QAECSRMoveSearch);
