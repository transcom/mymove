import React from 'react';
import { useParams, useLocation } from 'react-router-dom';
import { Button, Grid, GridContainer } from '@trussworks/react-uswds';

import styles from '../TXOMoveInfo/TXOTab.module.scss';

import evaluationReportsStyles from './EvaluationReports.module.scss';

import { useEvaluationReportsQueries } from 'hooks/queries';
import ShipmentEvaluationReports from 'components/Office/EvaluationReportTable/ShipmentEvaluationReports';
import EvaluationReportTable from 'components/Office/EvaluationReportTable/EvaluationReportTable';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import Alert from 'shared/Alert';

const EvaluationReports = () => {
  const { moveCode } = useParams();
  const location = useLocation();

  const { shipmentEvaluationReports, counselingEvaluationReports, shipments, isLoading, isError } =
    useEvaluationReportsQueries(moveCode);

  if (isLoading) {
    return <LoadingPlaceholder />;
  }
  if (isError) {
    return <SomethingWentWrong />;
  }

  return (
    <div className={styles.tabContent}>
      <GridContainer>
        {location.state?.showDeleteSuccess && (
          <div className={evaluationReportsStyles.alert}>
            <Alert type="success">Your report has been canceled</Alert>
          </div>
        )}
        {location.state?.showSaveDraftSuccess && (
          <div className={evaluationReportsStyles.alert}>
            <Alert type="success">Your draft report has been saved</Alert>
          </div>
        )}
        <Grid row>
          <h1>Quality assurance reports</h1>
        </Grid>
        <GridContainer className={evaluationReportsStyles.evaluationReportSection}>
          <Grid row className={evaluationReportsStyles.counselingHeadingContainer}>
            <h2>Counseling QAE reports ({counselingEvaluationReports.length})</h2>
            <Button>Create report</Button>
          </Grid>
          <Grid row>
            <EvaluationReportTable
              reports={counselingEvaluationReports}
              emptyText="No QAE reports have been submitted for counseling."
            />
          </Grid>
        </GridContainer>
        <GridContainer className={evaluationReportsStyles.evaluationReportSection}>
          <Grid row>
            <ShipmentEvaluationReports
              reports={shipmentEvaluationReports}
              shipments={shipments}
              emptyText="No QAE reports have been submitted for this shipment"
            />
          </Grid>
        </GridContainer>
      </GridContainer>
    </div>
  );
};

export default EvaluationReports;
