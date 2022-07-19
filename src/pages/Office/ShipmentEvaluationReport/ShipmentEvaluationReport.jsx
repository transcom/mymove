import React from 'react';
import 'styles/office.scss';
import { GridContainer, Grid, Button } from '@trussworks/react-uswds';
import classnames from 'classnames';
import { useParams } from 'react-router';

import styles from '../TXOMoveInfo/TXOTab.module.scss';

import shipmentEvaluationReportStyles from './ShipmentEvaluationReport.module.scss';

// import Alert from 'shared/Alert';

const mtoRefId = 'TODO'; // move?.referenceId

const ShipmentEvaluationReport = () => {
  const { moveCode, reportId } = useParams();

  return (
    <div className={classnames(styles.tabContent, shipmentEvaluationReportStyles.tabContent)}>
      <GridContainer>
        <Grid row>
          <Grid col desktop={{ col: 8, offset: 2 }}>
            {/* {showDeletionSuccess && <Alert type="success">Your remark has been deleted.</Alert>} */}
            <h1>Shipment report</h1>
            <div className={styles.pageHeaderDetails}>
              <h6>REPORT ID #{reportId}</h6>
              <h6>MOVE CODE {moveCode}</h6>
              <h6>MTO REFERENCE ID {mtoRefId}</h6>
            </div>
          </Grid>
        </Grid>
      </GridContainer>
      <GridContainer className={shipmentEvaluationReportStyles.cardContainer}>
        <Grid row>
          <Grid col desktop={{ col: 8, offset: 2 }}>
            <h2>Shipment information</h2>
          </Grid>
        </Grid>
      </GridContainer>

      <GridContainer className={shipmentEvaluationReportStyles.cardContainer}>
        <Grid row>
          <Grid col desktop={{ col: 8, offset: 2 }}>
            <h2>Evaluation form</h2>
          </Grid>
        </Grid>
      </GridContainer>
      <div style={{ display: 'flex', float: 'right' }}>
        <Button className="usa-button--unstyled">Cancel</Button>
        <Button className="usa-button--secondary">Save draft</Button>
        <Button type="submit">Submit</Button>
      </div>
    </div>
  );
};

export default ShipmentEvaluationReport;
