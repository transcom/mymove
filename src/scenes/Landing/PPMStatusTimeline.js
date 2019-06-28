import React from 'react';
import { get, includes } from 'lodash';
import { StatusTimeline } from './StatusTimeline';
import moment from 'moment';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import {
  getSignedCertification,
  selectPaymentRequestCertificationForMove,
} from 'shared/Entities/modules/signed_certifications';

const statuses = {
  Submitted: 'SUBMITTED',
  Approved: 'APPROVED',
  PaymentRequested: 'PAYMENT_REQUESTED',
  Completed: 'COMPLETED',
};

const codes = {
  Submitted: 'SUBMITTED',
  PpmApproved: 'PPM_APPROVED',
  InProgress: 'IN_PROGRESS',
  PaymentRequested: 'PAYMENT_REQUESTED',
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

  isCompleted(code) {
    const { ppm } = this.props;
    const moveIsApproved = includes([statuses.Approved, statuses.PaymentRequested, statuses.Completed], ppm.status);
    const moveInProgress = moment(ppm.original_move_date, 'YYYY-MM-DD').isSameOrBefore();
    const moveIsComplete = includes(['PAYMENT_REQUESTED', 'COMPLETED'], ppm.status);

    switch (code) {
      case codes.Submitted:
        return true;
      case codes.Approved:
        return moveIsApproved;
      case codes.InProgress:
        return (moveInProgress && ppm.status === statuses.Approved) || moveIsComplete;
      case codes.PaymentRequested:
        return moveIsComplete;
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
        code: codes.Submitted,
        dates: [submitDate],
        completed: this.isCompleted(codes.Submitted),
      },
      {
        name: 'Approved',
        code: codes.PpmApproved,
        dates: [approveDate],
        completed: this.isCompleted(codes.Approved),
      },
      {
        name: 'In progress',
        code: codes.InProgress,
        dates: [actualMoveDate],
        completed: this.isCompleted(codes.InProgress),
      },
      {
        name: 'Payment requested',
        code: codes.PaymentRequested,
        dates: [paymentRequestedDate],
        completed: this.isCompleted(codes.PaymentRequested),
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
  const move = state.moves.currentMove || state.moves.latestMove || {};
  const moveId = move.id || null;
  return {
    signedCertification: selectPaymentRequestCertificationForMove(state, move.id),
    moveId,
  };
}

const mapDispatchToProps = {
  getSignedCertification,
};

export default connect(mapStateToProps, mapDispatchToProps)(PPMStatusTimeline);
