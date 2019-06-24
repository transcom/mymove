import React, { Component } from 'react';

import { ProgressTimeline, ProgressTimelineStep } from 'shared/ProgressTimeline';

import WizardHeader from '../WizardHeader';
import Link from 'react-router-dom/es/Link';

class Review extends Component {
  render() {
    const moveId = this.props.match.params.moveId;
    const weightTicketsPage = `/moves/${moveId}/ppm-weight-ticket`;
    const expensePage = `/moves/${moveId}/ppm-expenses`;
    return (
      <>
        <WizardHeader
          title="Review"
          right={
            <ProgressTimeline>
              <ProgressTimelineStep name="Weight" completed />
              <ProgressTimelineStep name="Expenses" completed />
              <ProgressTimelineStep name="Review" current />
            </ProgressTimeline>
          }
        />
        <div className="usa-grid ">
          <ul>
            <li>
              <Link to={weightTicketsPage} data-cy="weight-ticket-link">
                Weight Ticket
              </Link>
            </li>
            <li>
              <Link to={expensePage} data-cy="expense-link">
                Expenses
              </Link>
            </li>
          </ul>
        </div>
      </>
    );
  }
}

export default Review;
