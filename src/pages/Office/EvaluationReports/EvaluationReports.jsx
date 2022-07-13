import React from 'react';
import { useParams } from 'react-router-dom';
import { Grid, GridContainer } from '@trussworks/react-uswds';

import { useEvaluationReportsQueries } from 'hooks/queries';
import ShipmentEvaluationReports from 'components/Office/EvaluationReportTable/ShipmentEvaluationReports';
import EvaluationReportTable from 'components/Office/EvaluationReportTable/EvaluationReportTable';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';

const EvaluationReports = () => {
  const { moveCode } = useParams();
  const { evaluationReports, shipments, isLoading, isError } = useEvaluationReportsQueries(moveCode);

  if (isLoading) {
    return <LoadingPlaceholder />;
  }
  if (isError) {
    return <SomethingWentWrong />;
  }

  return (
    <GridContainer>
      <Grid row>
        <h1>Quality assurance reports</h1>
      </Grid>
      <Grid row>
        <h2>Counseling QAE reports (?)</h2>
        <EvaluationReportTable reports={evaluationReports.filter((r) => r.type === 'COUNSELING')} />
      </Grid>
      <Grid row>
        <ShipmentEvaluationReports reports={evaluationReports} shipments={shipments} />
      </Grid>
    </GridContainer>
  );
};

export default EvaluationReports;
