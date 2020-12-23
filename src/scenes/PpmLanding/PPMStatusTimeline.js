import React from 'react';
import { get, includes } from 'lodash';
import moment from 'moment';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';

import { StatusTimeline } from './StatusTimeline';

import {
  getSignedCertification,
  selectPaymentRequestCertificationForMove,
} from 'shared/Entities/modules/signed_certifications';
import { selectCurrentMove } from 'store/entities/selectors';

const PpmStatuses = {
  Submitted: 'SUBMITTED',
  Approved: 'APPROVED',
  PaymentRequested: 'PAYMENT_REQUESTED',
  Completed: 'COMPLETED',
};

const PpmStatusTimelineCodes = {
  Submitted: 'SUBMITTED',
  PpmApproved: 'PPM_APPROVED',
  InProgress: 'IN_PROGRESS',
  PaymentRequested: 'PAYMENT_REQUESTED',
  PaymentReviewed: 'PAYMENT_REVIEWED',
};

export class PPMStatusTimeline extends React.Component {
  componentDidMount() {
    const { moveId } = this.props;
    this.props.getSignedCertification(moveId);
  }

  static determineActualMoveDate(ppm) {
    const approveDate = get(ppm, 'approve_date');
    const originalMoveDate = get(ppm, 'original_move_date');
    const actualMoveDate = get(ppm, 'actual_move_date');
    // if there's no approve date, then the PPM hasn't been approved yet
    // and the in progress date should not be shown
    if (!approveDate) {
      return;
    }
    // if there's an actual move date that is known and passed, show it
    // else show original move date if it has passed
    if (actualMoveDate && moment(actualMoveDate, 'YYYY-MM-DD').isSameOrBefore()) {
      return actualMoveDate;
    }
    if (moment(originalMoveDate, 'YYYY-MM-DD').isSameOrBefore()) {
      return originalMoveDate;
    }
  }

  isCompleted(statusCode) {
    const { ppm } = this.props;
    const moveIsApproved = includes(
      [PpmStatuses.Approved, PpmStatuses.PaymentRequested, PpmStatuses.Completed],
      ppm.status,
    );
    const moveInProgress = moment(ppm.original_move_date, 'YYYY-MM-DD').isSameOrBefore();
    const moveIsComplete = includes([PpmStatuses.PaymentRequested, PpmStatuses.Completed], ppm.status);

    switch (statusCode) {
      case PpmStatusTimelineCodes.Submitted:
        return true;
      case PpmStatusTimelineCodes.Approved:
        return moveIsApproved;
      case PpmStatusTimelineCodes.InProgress:
        return (moveInProgress && ppm.status === PpmStatuses.Approved) || moveIsComplete;
      case PpmStatusTimelineCodes.PaymentRequested:
        return moveIsComplete;
      case PpmStatusTimelineCodes.PaymentReviewed:
        return ppm.status === PpmStatuses.Completed;
      default:
        console.log('Unknown status');
    }
  }

  getStatuses() {
    const { ppm, signedCertification } = this.props;
    const actualMoveDate = PPMStatusTimeline.determineActualMoveDate(ppm);
    const approveDate = get(ppm, 'approve_date');
    const submitDate = get(ppm, 'submit_date');
    const paymentRequestedDate = signedCertification && signedCertification.date ? signedCertification.date : null;
    return [
      {
        name: 'Submitted',
        code: PpmStatusTimelineCodes.Submitted,
        dates: [submitDate],
        completed: this.isCompleted(PpmStatusTimelineCodes.Submitted),
      },
      {
        name: 'Approved',
        code: PpmStatusTimelineCodes.PpmApproved,
        dates: [approveDate],
        completed: this.isCompleted(PpmStatusTimelineCodes.Approved),
      },
      {
        name: 'In progress',
        code: PpmStatusTimelineCodes.InProgress,
        dates: [actualMoveDate],
        completed: this.isCompleted(PpmStatusTimelineCodes.InProgress),
      },
      {
        name: 'Payment requested',
        code: PpmStatusTimelineCodes.PaymentRequested,
        dates: [paymentRequestedDate],
        completed: this.isCompleted(PpmStatusTimelineCodes.PaymentRequested),
      },
      {
        name: 'Payment reviewed',
        code: PpmStatusTimelineCodes.PaymentReviewed,
        completed: this.isCompleted(PpmStatusTimelineCodes.PaymentReviewed),
      },
    ];
  }

  render() {
    const statuses = this.getStatuses();
    return <StatusTimeline statuses={statuses} showEstimated={false} />;
  }
}

PPMStatusTimeline.propTypes = {
  ppm: PropTypes.object.isRequired,
};

function mapStateToProps(state) {
  const move = selectCurrentMove(state);
  const moveId = move?.id;

  return {
    signedCertification: selectPaymentRequestCertificationForMove(state, moveId),
    moveId,
  };
}

const mapDispatchToProps = {
  getSignedCertification,
};

export default connect(mapStateToProps, mapDispatchToProps)(PPMStatusTimeline);
