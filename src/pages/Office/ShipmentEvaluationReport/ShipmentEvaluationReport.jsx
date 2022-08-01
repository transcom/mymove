import React from 'react';
import 'styles/office.scss';
import { GridContainer, Grid } from '@trussworks/react-uswds';
import classnames from 'classnames';
import { useParams } from 'react-router';

import styles from '../TXOMoveInfo/TXOTab.module.scss';
import ShipmentEvaluationForm from '../../../components/Office/ShipmentEvaluationForm/ShipmentEvaluationForm';

import shipmentEvaluationReportStyles from './ShipmentEvaluationReport.module.scss';

import { useShipmentEvaluationReportQueries } from 'hooks/queries';
import ShipmentDisplay from 'components/Office/ShipmentDisplay/ShipmentDisplay';
import DataTable from 'components/DataTable';
import { CustomerShape } from 'types';
import { OrdersShape } from 'types/customerShapes';
import { ORDERS_BRANCH_OPTIONS, ORDERS_RANK_OPTIONS } from 'constants/orders';
import { shipmentTypeLabels } from 'content/shipments';

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
      destinationAddress: shipment.destinationAddress || '-',
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
    <div className={classnames(styles.tabContent, shipmentEvaluationReportStyles.tabContent)}>
      <GridContainer>
        <div className={styles.pageHeader}>
          <h1>Shipment report</h1>
          <div className={styles.pageHeaderDetails}>
            <h6>REPORT ID #{reportId}</h6>
            <h6>MOVE CODE {moveCode}</h6>
            <h6>MTO REFERENCE ID #{mtoRefId}</h6>
          </div>
        </div>
      </GridContainer>
      <GridContainer className={shipmentEvaluationReportStyles.cardContainer}>
        <Grid row>
          <Grid col desktop={{ col: 8 }}>
            <h2>Shipment information</h2>
            {mtoShipment.id && (
              <ShipmentDisplay
                isSubmitted
                shipmentId={mtoShipment.id}
                displayInfo={shipmentDisplayInfo(mtoShipment)}
                shipmentType={mtoShipment.type}
              />
            )}
          </Grid>
          <Grid className={shipmentEvaluationReportStyles.qaeAndCustomerInfo} col desktop={{ col: 2 }}>
            <DataTable columnHeaders={['Customer information']} dataRow={[customerInfoTableBody]} />
            <DataTable columnHeaders={['QAE']} dataRow={[officeUserInfoTableBody]} />
          </Grid>
        </Grid>
      </GridContainer>

      <ShipmentEvaluationForm />
    </div>
  );
};

ShipmentEvaluationReport.propTypes = {
  customerInfo: CustomerShape.isRequired,
  orders: OrdersShape.isRequired,
};

export default ShipmentEvaluationReport;
