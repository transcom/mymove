import React from 'react';
import 'styles/office.scss';
import { GridContainer, Grid } from '@trussworks/react-uswds';
import classnames from 'classnames';
import { useParams } from 'react-router';

import styles from '../TXOMoveInfo/TXOTab.module.scss';

import shipmentEvaluationReportStyles from './ShipmentEvaluationReport.module.scss';

import ShipmentEvaluationForm from 'components/Office/ShipmentEvaluationForm/ShipmentEvaluationForm';
import { useShipmentEvaluationReportQueries } from 'hooks/queries';
import { formatQAReportID } from 'utils/formatters';
import DataTable from 'components/DataTable';
import { CustomerShape } from 'types';
import { OrdersShape } from 'types/customerShapes';
import { ORDERS_BRANCH_OPTIONS, ORDERS_RANK_OPTIONS } from 'constants/orders';
import { shipmentTypeLabels } from 'content/shipments';
import EvaluationReportShipmentDisplay from 'components/Office/EvaluationReportShipmentDisplay/EvaluationReportShipmentDisplay';

const ShipmentEvaluationReport = ({ customerInfo, orders }) => {
  const { moveCode, reportId } = useParams();

  const { evaluationReport, mtoShipment } = useShipmentEvaluationReportQueries(reportId);
  const mtoRefId = evaluationReport.moveReferenceID;

  const shipmentDisplayInfo = (shipment) => {
    return {
      ...shipment,
      heading: shipmentTypeLabels[shipment.shipmentType],
      isDiversion: shipment.diversion,
      shipmentStatus: shipment.status,
      destinationAddress: shipment.destinationAddress,
    };
  };

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

  const officeUserInfoTableBody = evaluationReport.officeUser ? (
    <>
      {evaluationReport.officeUser.lastName}, {evaluationReport.officeUser.firstName}
      <br />
      {evaluationReport.officeUser.phone}
      <br />
      {evaluationReport.officeUser.email}
    </>
  ) : (
    ''
  );

  return (
    <div className={classnames(styles.tabContent, shipmentEvaluationReportStyles.tabContent)}>
      <GridContainer>
        <div className={styles.pageHeader}>
          <h1>Shipment report</h1>
          <div className={styles.pageHeaderDetails}>
            <h6>REPORT ID {formatQAReportID(reportId)}</h6>
            <h6>MOVE CODE #{moveCode}</h6>
            <h6>MTO REFERENCE ID #{mtoRefId}</h6>
          </div>
        </div>

        <GridContainer className={shipmentEvaluationReportStyles.cardContainer}>
          <Grid row>
            <Grid col desktop={{ col: 8 }}>
              <h2>Shipment information</h2>
              {mtoShipment.id && (
                <EvaluationReportShipmentDisplay
                  isSubmitted
                  shipmentId={mtoShipment.id}
                  displayInfo={shipmentDisplayInfo(mtoShipment)}
                  shipmentType={mtoShipment.shipmentType}
                />
              )}
            </Grid>
            <Grid className={shipmentEvaluationReportStyles.qaeAndCustomerInfo} col desktop={{ col: 2 }}>
              <DataTable columnHeaders={['Customer information']} dataRow={[customerInfoTableBody]} />
              <DataTable columnHeaders={['QAE']} dataRow={[officeUserInfoTableBody]} />
            </Grid>
          </Grid>
        </GridContainer>

        <ShipmentEvaluationForm evaluationReport={evaluationReport} />
      </GridContainer>
    </div>
  );
};

ShipmentEvaluationReport.propTypes = {
  customerInfo: CustomerShape.isRequired,
  orders: OrdersShape.isRequired,
};

export default ShipmentEvaluationReport;
