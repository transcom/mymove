import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { milmoveLog, MILMOVE_LOG_LEVEL } from 'utils/milmoveLog';

export class ApproveRejectModal extends Component {
  state = {
    showApproveBtn: true,
    showRejectionToggleBtn: true,
    showRejectionInput: false,
    rejectBtnIsDisabled: true,
    rejectionReason: null,
  };

  handleApproveClick = () => {
    this.props.approveBtnOnClick();
  };

  handleRejectionChange = (rejectionReason) => {
    this.setState({ rejectionReason });

    if (!rejectionReason) {
      this.setState({ rejectBtnIsDisabled: true });
    } else {
      this.setState({ rejectBtnIsDisabled: false });
    }
  };

  handleRejectionClick = () => {
    if (!this.state.rejectionReason) {
      milmoveLog(MILMOVE_LOG_LEVEL.ERROR, 'Rejection reason empty. Please fill out rejection reason.');
    } else {
      this.props.rejectBtnOnClick(this.state.rejectionReason);
    }
  };

  handleRejectionCancelClick = () => {
    this.setState({
      showApproveBtn: true,
      showRejectionToggleBtn: true,
      showRejectionInput: false,
      rejectBtnIsDisabled: true,
      rejectionReason: null,
    });
  };

  handleRejectionToggleClick = () => {
    this.setState({ showRejectionInput: true, showRejectionToggleBtn: false, showApproveBtn: false });
  };

  render() {
    const { showModal } = this.props;
    return (
      showModal && (
        <>
          <div>
            {this.state.showApproveBtn && <button onClick={this.handleApproveClick}>Approve</button>}
            {this.state.showRejectionToggleBtn && (
              <button data-testid="rejectionToggle" onClick={this.handleRejectionToggleClick}>
                Reject
              </button>
            )}
          </div>
          <div>
            {this.state.showRejectionInput && (
              <label htmlFor="rejectionReason">
                Rejection reason
                <input
                  name="rejectionReason"
                  id="rejectionReason"
                  onChange={(event) => this.handleRejectionChange(event.target.value)}
                />
                <button
                  data-testid="rejectionButton"
                  onClick={this.handleRejectionClick}
                  disabled={this.state.rejectBtnIsDisabled}
                >
                  Reject
                </button>
                <button onClick={this.handleRejectionCancelClick}>Cancel</button>
              </label>
            )}
          </div>
        </>
      )
    );
  }
}

ApproveRejectModal.propTypes = {
  /** REQUIRED. Function that is passed in to the onClick prop of the approve button */
  approveBtnOnClick: PropTypes.func.isRequired,
  /** REQUIRED. Function that is passed in to the onClick prop of the reject button */
  rejectBtnOnClick: PropTypes.func.isRequired,
  /** OPTIONAL. Boolean to hide the modal or not */
  showModal: PropTypes.bool,
};

ApproveRejectModal.defaultProps = {
  showModal: true,
};
