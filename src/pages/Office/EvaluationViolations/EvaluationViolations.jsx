import React from 'react';
import 'styles/office.scss';
import { GridContainer, Grid, Button } from '@trussworks/react-uswds';
import classnames from 'classnames';
import { useParams, useHistory } from 'react-router';

import styles from '../TXOMoveInfo/TXOTab.module.scss';

import evaluationViolationsStyles from './EvaluationViolations.module.scss';

import { useEvaluationReportQueries } from 'hooks/queries';
import QaeReportHeader from 'components/Office/QaeReportHeader/QaeReportHeader';

const EvaluationViolations = () => {
  const { moveCode, reportId } = useParams();
  const history = useHistory();

  const { evaluationReport } = useEvaluationReportQueries(reportId);

  const handleBackToEvalForm = () => {
    // TODO: Save as draft before rerouting
    history.push(`/moves/${moveCode}/evaluation-reports/${reportId}`);
  };

  const cancelForViolations = () => {
    history.push(`/moves/${moveCode}/evaluation-reports`);
  };

  return (
    <GridContainer className={classnames(styles.tabContent, evaluationViolationsStyles.tabContent)}>
      <QaeReportHeader report={evaluationReport} />
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
              <Button className="usa-button--unstyled" type="button" onClick={cancelForViolations}>
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
  );
};

export default EvaluationViolations;
