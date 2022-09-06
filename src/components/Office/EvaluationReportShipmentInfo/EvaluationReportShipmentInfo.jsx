import React from 'react';
import { PropTypes } from 'prop-types';
import { GridContainer } from '@trussworks/react-uswds';
import classnames from 'classnames';

import evaluationReportStyles from './EvaluationReportShipmentInfo.module.scss';

import styles from 'components/Office/EvaluationReportPreview/EvaluationReportPreview.module.scss';
import 'styles/office.scss';
import DataTable from 'components/DataTable';
import { ORDERS_BRANCH_OPTIONS, ORDERS_RANK_OPTIONS } from 'constants/orders';
import { shipmentTypeLabels } from 'content/shipments';
import EvaluationReportShipmentDisplay from 'components/Office/EvaluationReportShipmentDisplay/EvaluationReportShipmentDisplay';

const EvaluationReportShipmentInfo = ({ shipments, report, customerInfo, grade }) => {
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
      {ORDERS_RANK_OPTIONS[grade]}
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
      <div className={evaluationReportStyles.sidebySideContainer}>
        <div>
          <h2>{report.type === 'SHIPMENT' ? 'Shipment' : 'Move'} information</h2>
          {shipments.length > 0 &&
            shipments.map((shipment) => (
              <div
                key={shipment.id}
                className={classnames(styles.shipmentDisplayContainer, evaluationReportStyles.shipmentCardColumn)}
              >
                <EvaluationReportShipmentDisplay
                  isSubmitted
                  key={shipment.id}
                  shipmentId={shipment.id}
                  displayInfo={shipmentDisplayInfo(shipment)}
                  shipmentType={shipment.shipmentType}
                />
              </div>
            ))}
        </div>
        <div className={evaluationReportStyles.qaeAndCustomerInfo} data-testid="qaeAndCustomerInfo">
          <DataTable columnHeaders={['QAE']} dataRow={[officeUserInfoTableBody]} />
          <DataTable columnHeaders={['Customer information']} dataRow={[customerInfoTableBody]} />
        </div>
      </div>
    </GridContainer>
  );
};

EvaluationReportShipmentInfo.propTypes = {
  report: PropTypes.object.isRequired,
  shipments: PropTypes.arrayOf(PropTypes.object).isRequired,
  customerInfo: PropTypes.object.isRequired,
  grade: PropTypes.string.isRequired,
};

export default EvaluationReportShipmentInfo;
