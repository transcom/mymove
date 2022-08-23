import React from 'react';
import { PropTypes } from 'prop-types';
import { GridContainer, Grid } from '@trussworks/react-uswds';
import 'styles/office.scss';

import styles from '../EvaluationReportTable/EvaluationReportContainer.module.scss';

import evaluationReportStyles from './EvaluationReportShipmentInfo.module.scss';

import DataTable from 'components/DataTable';
import { ORDERS_BRANCH_OPTIONS, ORDERS_RANK_OPTIONS } from 'constants/orders';
import { shipmentTypeLabels } from 'content/shipments';
import EvaluationReportShipmentDisplay from 'components/Office/EvaluationReportShipmentDisplay/EvaluationReportShipmentDisplay';

const EvaluationReportShipmentInfo = ({ shipments, report, customerInfo, orders }) => {
  const shipmentDisplayInfo = (shipment) => ({
    ...shipment,
    heading: shipmentTypeLabels[shipment.shipmentType],
    isDiversion: shipment.diversion,
    shipmentStatus: shipment.status,
    destinationAddress: shipment.destinationAddress,
  });

  const customerInfoTableBody = (
    <>
      {customerInfo.last_name}, {customerInfo.first_name}
      <br />
      {customerInfo.phone}
      <br />
      {ORDERS_RANK_OPTIONS[orders.grade]}
      <br />
      {ORDERS_BRANCH_OPTIONS[customerInfo.agency] ? ORDERS_BRANCH_OPTIONS[customerInfo.agency] : customerInfo.agency}
    </>
  );

  const officeUserInfoTableBody = report.officeUser ? (
    <>
      {report.officeUser.lastName}, {report.officeUser.firstName}
      <br />
      {report.officeUser.phone}
      <br />
      {report.officeUser.email}
    </>
  ) : (
    ''
  );

  return (
    <GridContainer className={evaluationReportStyles.cardContainer}>
      <Grid row>
        <Grid col desktop={{ col: 8 }}>
          <h2>{report.type === 'SHIPMENT' ? 'Shipment' : 'Move'} information</h2>
          {shipments.length > 0 &&
            shipments.map((shipment) => (
              <div key={shipment.id} className={styles.shipmentDisplayContainer}>
                <EvaluationReportShipmentDisplay
                  isSubmitted
                  key={shipment.id}
                  shipmentId={shipment.id}
                  displayInfo={shipmentDisplayInfo(shipment)}
                  shipmentType={shipment.shipmentType}
                />
              </div>
            ))}
        </Grid>
        <Grid
          className={evaluationReportStyles.qaeAndCustomerInfo}
          col
          desktop={{ col: 2 }}
          data-testid="qaeAndCustomerInfo"
        >
          <DataTable columnHeaders={['Customer information']} dataRow={[customerInfoTableBody]} />
          <DataTable columnHeaders={['QAE']} dataRow={[officeUserInfoTableBody]} />
        </Grid>
      </Grid>
    </GridContainer>
  );
};

EvaluationReportShipmentInfo.propTypes = {
  report: PropTypes.object.isRequired,
  shipments: PropTypes.arrayOf(PropTypes.object).isRequired,
  customerInfo: PropTypes.object.isRequired,
  orders: PropTypes.object.isRequired,
};

export default EvaluationReportShipmentInfo;
