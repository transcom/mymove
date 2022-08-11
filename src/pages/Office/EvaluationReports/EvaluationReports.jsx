import React from 'react';
import { useParams, useHistory, useLocation } from 'react-router-dom';
import { Button, Grid, GridContainer } from '@trussworks/react-uswds';
import { useMutation, queryCache } from 'react-query';

import styles from '../TXOMoveInfo/TXOTab.module.scss';

import evaluationReportsStyles from './EvaluationReports.module.scss';

import { useEvaluationReportsQueries } from 'hooks/queries';
import ShipmentEvaluationReports from 'components/Office/EvaluationReportTable/ShipmentEvaluationReports';
import EvaluationReportTable from 'components/Office/EvaluationReportTable/EvaluationReportTable';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import Alert from 'shared/Alert';
import { createCounselingEvaluationReport } from 'services/ghcApi';
import { COUNSELING_EVALUATION_REPORTS } from 'constants/queryKeys';
import { milmoveLog, MILMOVE_LOG_LEVEL } from 'utils/milmoveLog';

const EvaluationReports = () => {
  const { moveCode } = useParams();
  const location = useLocation();
  const history = useHistory();

  const { shipmentEvaluationReports, counselingEvaluationReports, shipments, isLoading, isError } =
    useEvaluationReportsQueries(moveCode);

  const [createCounselingEvaluationReportMutation] = useMutation(createCounselingEvaluationReport, {
    onSuccess: () => {
      queryCache.invalidateQueries([COUNSELING_EVALUATION_REPORTS, moveCode]);
    },
    onError: (error) => {
      const errorMsg = error?.response?.body;
      milmoveLog(MILMOVE_LOG_LEVEL.LOG, errorMsg);
    },
  });

  const handleCounselingCreateClick = async () => {
    const report = await createCounselingEvaluationReportMutation({ moveCode });
    const reportId = report?.id;

    history.push(`/moves/${moveCode}/counseling-evaluation-reports/${reportId}`);
  };

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
            <Button onClick={() => handleCounselingCreateClick()}>Create report</Button>
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
