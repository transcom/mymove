import React from 'react';
import * as PropTypes from 'prop-types';
import { GridContainer, Grid } from '@trussworks/react-uswds';
import 'styles/office.scss';

import evaluationReportStyles from './EvaluationReportMoveInfo.module.scss';

import DataTable from 'components/DataTable';
import { CustomerShape } from 'types';
import { ORDERS_BRANCH_OPTIONS, ORDERS_RANK_OPTIONS } from 'constants/orders';

const EvaluationReportMoveInfo = ({ customerInfo, grade }) => {
  const customerInfoTableBody = (
    <>
      {customerInfo.last_name}, {customerInfo.first_name}
      <br />
      {customerInfo.phone}
      <br />
      {ORDERS_RANK_OPTIONS[grade]}
      <br />
      {ORDERS_BRANCH_OPTIONS[customerInfo.agency] ? ORDERS_BRANCH_OPTIONS[customerInfo.agency] : customerInfo.agency}
    </>
  );

  const officeUserInfoTableBody = (
    <>
      {customerInfo.last_name}, {customerInfo.first_name}
      <br />
      {customerInfo.phone}
      <br />
      {customerInfo.email}
    </>
  );

  return (
    <GridContainer className={evaluationReportStyles.cardContainer}>
      <Grid row>
        <Grid col desktop={{ col: 8 }}>
          <h2>Move information</h2>
        </Grid>
        <Grid className={evaluationReportStyles.qaeAndCustomerInfo} col desktop={{ col: 2 }}>
          <DataTable columnHeaders={['Customer information']} dataRow={[customerInfoTableBody]} />
          <DataTable columnHeaders={['QAE']} dataRow={[officeUserInfoTableBody]} />
        </Grid>
      </Grid>
    </GridContainer>
  );
};

EvaluationReportMoveInfo.propTypes = {
  customerInfo: CustomerShape.isRequired,
  grade: PropTypes.string.isRequired,
};

export default EvaluationReportMoveInfo;
