import React from 'react';
import 'styles/office.scss';
import { GridContainer, Grid, Button } from '@trussworks/react-uswds';
import classnames from 'classnames';
import { useParams, useHistory } from 'react-router';

import styles from '../TXOMoveInfo/TXOTab.module.scss';

import evaluationViolationsStyles from './EvaluationViolations.module.scss';

import { useShipmentEvaluationReportQueries, usePWSViolationsQueries } from 'hooks/queries';
import QaeReportHeader from 'components/Office/QaeReportHeader/QaeReportHeader';
import ViolationAccordion from 'components/Office/ViolationAccordion/ViolationAccordion';

const EvaluationViolations = () => {
  const { moveCode, reportId } = useParams();
  const history = useHistory();

  const { evaluationReport } = useShipmentEvaluationReportQueries(reportId);
  const { violations } = usePWSViolationsQueries();

  const handleBackToEvalForm = () => {
    // TODO: Save as draft before rerouting
    history.push(`/moves/${moveCode}/evaluation-reports/${reportId}`);
  };

  const cancelForViolations = () => {
    history.push(`/moves/${moveCode}/evaluation-reports`);
  };

  // Get distinct categories
  const categories = [...new Set(violations.map((item) => item.category))];

  return (
    <div className={classnames(styles.tabContent, evaluationViolationsStyles.tabContent)}>
      <GridContainer>
        <QaeReportHeader report={evaluationReport} />
        <GridContainer className={evaluationViolationsStyles.cardContainer}>
          <Grid row>
            <Grid col>
              <h2>Select violations</h2>
              <p>Select the paragraph from the Performance Work Statement (PWS) that the GHC Prime has violated.</p>
            </Grid>
          </Grid>
          {categories.map((category) => (
            <ViolationAccordion
              category={category}
              violations={violations.filter((violation) => violation.category === category)}
              key={`${category}-category`}
            />
          ))}
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
    </div>
  );
};

export default EvaluationViolations;
