import React from 'react';
import { useParams, useLocation } from 'react-router-dom';
import { Button, Grid, GridContainer } from '@trussworks/react-uswds';

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
    <GridContainer>
      {location.state?.showDeleteSuccess && (
        <div className={evaluationReportsStyles.alert}>
          <Alert type="success">Your report has been canceled</Alert>
        </div>
      )}
      <Grid row>
        <h1>Quality assurance reports</h1>
      </Grid>
      <Grid row>
        <h2>Counseling QAE reports ({counselingEvaluationReports.length})</h2>
        <Button>Create report</Button>
        <EvaluationReportTable reports={counselingEvaluationReports} />
      </Grid>
      <Grid row>
        <ShipmentEvaluationReports reports={shipmentEvaluationReports} shipments={shipments} />
      </Grid>
    </GridContainer>
  );
};

export default EvaluationReports;
