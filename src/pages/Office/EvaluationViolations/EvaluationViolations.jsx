import React, { useState } from 'react';
import 'styles/office.scss';
import { GridContainer, Grid, Button } from '@trussworks/react-uswds';
import classnames from 'classnames';
import { useMutation } from 'react-query';
import { useParams, useHistory } from 'react-router';

import styles from '../TXOMoveInfo/TXOTab.module.scss';

import evaluationViolationsStyles from './EvaluationViolations.module.scss';

import ConnectedDeleteEvaluationReportConfirmationModal from 'components/ConfirmationModals/DeleteEvaluationReportConfirmationModal';
import { useShipmentEvaluationReportQueries } from 'hooks/queries';
import { deleteEvaluationReport } from 'services/ghcApi';

const EvaluationViolations = () => {
  const { moveCode, reportId } = useParams();
  const history = useHistory();

  const { evaluationReport } = useShipmentEvaluationReportQueries(reportId);
  const mtoRefId = evaluationReport.moveReferenceID;

  const handleBackToEvalForm = () => {
    // TODO: Save as draft before rerouting
    history.push(`/moves/${moveCode}/evaluation-reports/${reportId}`);
  };

  const [isDeleteModelOpen, setIsDeleteModelOpen] = useState(false);

  const [deleteEvaluationReportMutation] = useMutation(deleteEvaluationReport);

  const toggleCancelModel = () => {
    setIsDeleteModelOpen(!isDeleteModelOpen);
  };

  const cancelReport = async () => {
    // Close the modal
    setIsDeleteModelOpen(!isDeleteModelOpen);

    // Mark as deleted in database
    await deleteEvaluationReportMutation(reportId);

    // Reroute back to eval report page, include flag to know to show alert
    history.push(`/moves/${moveCode}/evaluation-reports`, { showDeleteSuccess: true });
  };

  return (
    <>
      <ConnectedDeleteEvaluationReportConfirmationModal
        isOpen={isDeleteModelOpen}
        closeModal={toggleCancelModel}
        submitModal={cancelReport}
      />
      <GridContainer className={classnames(styles.tabContent, evaluationViolationsStyles.tabContent)}>
        <GridContainer>
          <div className={styles.pageHeader}>
            <h1>{evaluationReport.type} report</h1>
            <div className={styles.pageHeaderDetails}>
              <h6>REPORT ID #{reportId}</h6>
              <h6>MOVE CODE {moveCode}</h6>
              <h6>MTO REFERENCE ID #{mtoRefId}</h6>
            </div>
          </div>
        </GridContainer>
        <GridContainer className={evaluationViolationsStyles.cardContainer}>
          <Grid row>
            <Grid col desktop={{ col: 8 }}>
              <h2>Select violations</h2>
            </Grid>
          </Grid>
        </GridContainer>
        <GridContainer className={evaluationViolationsStyles.buttonContainer}>
          <Grid row>
            <Grid col>
              <div className={evaluationViolationsStyles.buttonRow}>
                <Button
                  className={classnames(evaluationViolationsStyles.backToEvalButton, 'usa-button--unstyled')}
                  type="button"
                  onClick={handleBackToEvalForm}
                >
                  {'< Back to Evaluation form'}
                </Button>
                <div className={evaluationViolationsStyles.grow} />
                <Button className="usa-button--unstyled" type="button" onClick={toggleCancelModel}>
                  Cancel
                </Button>
                <Button data-testid="saveDraft" type="submit" className="usa-button--secondary">
                  Save draft
                </Button>
                <Button disabled>Review and submit</Button>
              </div>
            </Grid>
          </Grid>
        </GridContainer>
      </GridContainer>
    </>
  );
};

export default EvaluationViolations;
