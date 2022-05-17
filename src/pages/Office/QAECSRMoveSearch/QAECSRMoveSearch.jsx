import React, { useMemo, useState } from 'react';
import { Formik, Field } from 'formik';
import { withRouter } from 'react-router-dom';
import { Button } from '@trussworks/react-uswds';
import * as Yup from 'yup';

import styles from './QAECSRMoveSearch.module.scss';

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
  createHeader('DoD ID', 'customer.dodID', {
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
  // createHeader('Origin duty location', 'originDutyLocation.name', {
  //   id: 'originDutyLocation',
  //   isFilterable: true,
  // }),
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
    then: Yup.string().length(6),
    // TODO need to figure out how to do this only when dodID is selected
    otherwise: Yup.string().length(10),
  }),
});

// const QAECSRMoveSearch = ({ history }) => {
const QAECSRMoveSearch = () => {
  const [search, setSearch] = useState({ moveCode: '9TR9JG', dodID: null });

  const handleClick = () => {
    // history.push(`/moves/${values.locator}/details`);
    // console.log('click');
  };
  const onSubmit = (values) => {
    // console.log('onSubmit', values);
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
  };

  const { searchResult, isLoading, isError } = useQAECSRMoveSearchQueries({
    moveCode: search.moveCode,
  });
  // const { totalCount = 0, data = [], page = 1, perPage = 20 } = searchResult;
  const { data = [] } = searchResult;
  const tableColumns = useMemo(() => columns(true), []);
  const tableData = useMemo(() => data, [data]);
  return (
    <>
      <Formik
        initialValues={{ searchType: 'moveCode', searchText: '' }}
        onSubmit={onSubmit}
        validationSchema={validationSchema}
        validateOnMount
        validateOnChange
      >
        {(formik) => {
          return (
            <Form className={formStyles.form} onSubmit={formik.handleSubmit}>
              {formik.isValid && <p>isValid</p>}
              {formik.isSubmitting && <p>isSubmitting</p>}
              <p>{formik.values.searchType}</p>
              <p>{formik.values.searchText}</p>
              <p>{JSON.stringify(formik.errors)}</p>
              <div role="group" aria-labelledby="my-radio-group" className="usa-radio">
                <label htmlFor="radio-picked-one">
                  <Field id="radio-picked-one" type="radio" name="searchType" value="moveCode" />
                  Move Code
                </label>
                <label htmlFor="radio-picked-two">
                  <Field id="radio-picked-two" type="radio" name="searchType" value="dodID" />
                  DOD ID
                </label>
              </div>
              <TextField id="foobar" label="Search" name="searchText" />
              <Button type="submit">Search</Button>
            </Form>
          );
        }}
      </Formik>
      {isLoading && <LoadingPlaceholder />}
      {isError && <SomethingWentWrong />}
      <div className={styles.MoveQueue}>
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
          title="All moves"
          handleClick={handleClick}
          useQueries={useQAECSRMoveSearchQueries}
          searchKey="moveCode"
          searchValue={search.moveCode}
          data={tableData}
        />
      </div>
    </>
  );
};

// QAECSRMoveSearch.propTypes = {
//   history: HistoryShape.isRequired,
// };

export default withRouter(QAECSRMoveSearch);
