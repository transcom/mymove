import React, { useState } from 'react';
import 'styles/office.scss';
import { GridContainer, Grid, Button } from '@trussworks/react-uswds';
import classnames from 'classnames';
import { useParams, useHistory } from 'react-router';
import { queryCache, useMutation } from 'react-query';

import styles from '../TXOMoveInfo/TXOTab.module.scss';

import shipmentEvaluationReportStyles from './ShipmentEvaluationReport.module.scss';

import ConnectedDeleteEvaluationReportConfirmationModal from 'components/ConfirmationModals/DeleteEvaluationReportConfirmationModal';
import { deleteEvaluationReport } from 'services/ghcApi';
import { SHIPMENT_EVALUATION_REPORTS } from 'constants/queryKeys';

const mtoRefId = 'TODO'; // move?.referenceId

const ShipmentEvaluationReport = () => {
  const { moveCode, reportId } = useParams();
  const history = useHistory();

  const [deleteEvaluationReportMutation] = useMutation(deleteEvaluationReport, {
    onSuccess: async () => {
      await queryCache.invalidateQueries([SHIPMENT_EVALUATION_REPORTS, moveCode]);
    },
  });

  const [isDeleteModelOpen, setIsDeleteModelOpen] = useState(false);

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
          <Button className="usa-button--unstyled" onClick={toggleCancelModel}>
            Cancel
          </Button>
          <Button className="usa-button--secondary">Save draft</Button>
          <Button type="submit">Submit</Button>
        </div>
      </div>
    </>
  );
};

export default ShipmentEvaluationReport;
