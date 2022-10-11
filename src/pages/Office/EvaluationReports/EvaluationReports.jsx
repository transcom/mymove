import React, { useState } from 'react';
import { useParams, useHistory, useLocation } from 'react-router-dom';
import { Button, Grid, GridContainer } from '@trussworks/react-uswds';
import PropTypes from 'prop-types';
import { useMutation, queryCache } from 'react-query';

import styles from '../TXOMoveInfo/TXOTab.module.scss';

import evaluationReportsStyles from './EvaluationReports.module.scss';

import { useEvaluationReportsQueries } from 'hooks/queries';
import ShipmentEvaluationReports from 'components/Office/EvaluationReportTable/ShipmentEvaluationReports';
import EvaluationReportTable from 'components/Office/EvaluationReportTable/EvaluationReportTable';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import Alert from 'shared/Alert';
import { CustomerShape } from 'types';
import { createCounselingEvaluationReport, deleteEvaluationReport } from 'services/ghcApi';
import { COUNSELING_EVALUATION_REPORTS, SHIPMENT_EVALUATION_REPORTS } from 'constants/queryKeys';
import { milmoveLog, MILMOVE_LOG_LEVEL } from 'utils/milmoveLog';
import Restricted from 'components/Restricted/Restricted';
import { permissionTypes } from 'constants/permissions';

const EvaluationReports = ({ customerInfo, grade, destinationDutyLocationPostalCode }) => {
  const { moveCode } = useParams();
  const location = useLocation();
  const history = useHistory();
  const [reportToDelete, setReportToDelete] = useState(undefined);
  const [isDeleteModalOpen, setIsDeleteModalOpen] = useState(false);

  const { shipmentEvaluationReports, counselingEvaluationReports, shipments, isLoading, isError } =
    useEvaluationReportsQueries(moveCode);

  const [deleteEvaluationReportMutation] = useMutation(deleteEvaluationReport);

  const deleteReport = async () => {
    // Close the modal
    setIsDeleteModalOpen(!isDeleteModalOpen);

    const reportID = reportToDelete.id;

    // Mark as deleted in database
    await deleteEvaluationReportMutation(reportID, {
      onError: (error) => {
        const errorMsg = error?.response?.body;
        milmoveLog(MILMOVE_LOG_LEVEL.LOG, errorMsg);
      },
      onSuccess: () => {
        // Reroute back to eval report page, include flag to show success alert
        history.push(`/moves/${moveCode}/evaluation-reports`, { showDeleteSuccess: true });
        queryCache
          .refetchQueries([COUNSELING_EVALUATION_REPORTS])
          .then(queryCache.refetchQueries(SHIPMENT_EVALUATION_REPORTS).then());
      },
    });
  };

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

    history.push(`/moves/${moveCode}/evaluation-reports/${reportId}`);
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
            <Alert type="success">Your report has been deleted</Alert>
          </div>
        )}
        {location.state?.showCanceledSuccess && (
          <div className={evaluationReportsStyles.alert}>
            <Alert type="success">Your report has been canceled</Alert>
          </div>
        )}
        {location.state?.showSaveDraftSuccess && (
          <div className={evaluationReportsStyles.alert}>
            <Alert type="success">Your draft report has been saved</Alert>
          </div>
        )}
        {location.state?.showSubmitSuccess && (
          <div className={evaluationReportsStyles.alert}>
            <Alert type="success">Your report has been successfully submitted</Alert>
          </div>
        )}
        <Grid row>
          <h1>Quality assurance reports</h1>
        </Grid>
        <GridContainer className={evaluationReportsStyles.evaluationReportSection}>
          <Grid row className={evaluationReportsStyles.counselingHeadingContainer}>
            <h2>Counseling QAE reports ({counselingEvaluationReports.length})</h2>
            <Restricted to={permissionTypes.createEvaluationReport}>
              <Button onClick={() => handleCounselingCreateClick()}>Create report</Button>
            </Restricted>
          </Grid>
          <Grid row>
            <EvaluationReportTable
              reports={counselingEvaluationReports}
              moveCode={moveCode}
              customerInfo={customerInfo}
              grade={grade}
              destinationDutyLocationPostalCode={destinationDutyLocationPostalCode}
              shipments={shipments}
              emptyText="No QAE reports have been submitted for counseling."
              setReportToDelete={setReportToDelete}
              setIsDeleteModalOpen={setIsDeleteModalOpen}
              isDeleteModalOpen={isDeleteModalOpen}
              deleteReport={deleteReport}
            />
          </Grid>
        </GridContainer>
        <GridContainer className={evaluationReportsStyles.evaluationReportSection}>
          <Grid row>
            <ShipmentEvaluationReports
              reports={shipmentEvaluationReports}
              shipments={shipments}
              moveCode={moveCode}
              customerInfo={customerInfo}
              grade={grade}
              destinationDutyLocationPostalCode={destinationDutyLocationPostalCode}
              emptyText="No QAE reports have been submitted for this shipment"
              setReportToDelete={setReportToDelete}
              setIsDeleteModalOpen={setIsDeleteModalOpen}
              isDeleteModalOpen={isDeleteModalOpen}
              deleteReport={deleteReport}
            />
          </Grid>
        </GridContainer>
      </GridContainer>
    </div>
  );
};

EvaluationReports.propTypes = {
  customerInfo: CustomerShape.isRequired,
  grade: PropTypes.string.isRequired,
  destinationDutyLocationPostalCode: PropTypes.string.isRequired,
};

export default EvaluationReports;
