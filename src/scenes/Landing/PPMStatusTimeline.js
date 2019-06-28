import React from 'react';
import { includes } from 'lodash';
import { getDates, StatusTimeline } from './StatusTimeline';
import moment from 'moment';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import {
  getSignedCertification,
  selectPaymentRequestCertificationForMove,
} from '../../shared/Entities/modules/signed_certifications';

export class PPMStatusTimeline extends React.Component {
  componentDidMount() {
    const { moveId } = this.props;
    this.props.getSignedCertification(moveId);
  }

  getStatuses() {
    return [
      { name: 'Submitted', code: 'SUBMITTED', date_type: 'submit_date' },
      { name: 'Approved', code: 'PPM_APPROVED', date_type: 'approve_date' },
      { name: 'In progress', code: 'IN_PROGRESS', date_type: 'actual_move_date' },
      { name: 'Payment requested', code: 'PAYMENT_REQUESTED' },
    ];
  }

  getCompletedStatus(status) {
    const { ppm } = this.props;

    if (status === 'SUBMITTED') {
      return true;
    }

    if (status === 'PPM_APPROVED') {
      return includes(['APPROVED', 'PAYMENT_REQUESTED', 'COMPLETED'], ppm.status);
    }

    if (status === 'IN_PROGRESS') {
      const moveInProgress = moment(ppm.original_move_date, 'YYYY-MM-DD').isSameOrBefore();
      return (moveInProgress && ppm.status === 'APPROVED') || includes(['PAYMENT_REQUESTED', 'COMPLETED'], ppm.status);
    }

    if (status === 'PAYMENT_REQUESTED') {
      return includes(['PAYMENT_REQUESTED', 'COMPLETED'], ppm.status);
    }
  }

  addDates(statuses) {
    const { signedCertification } = this.props;
    const paymentRequestedDate = signedCertification && signedCertification.date ? signedCertification.date : null;
    return statuses.map(status => {
      if (status.code === 'PAYMENT_REQUESTED') {
        return {
          ...status,
          dates: [paymentRequestedDate],
        };
      }
      return {
        ...status,
        dates: [getDates(this.props.ppm, status.date_type)],
      };
    });
  }

  addCompleted(statuses) {
    return statuses.map(status => {
      return {
        ...status,
        completed: this.getCompletedStatus(status.code),
      };
    });
  }

  render() {
    const statuses = this.addDates(this.addCompleted(this.getStatuses()));
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
