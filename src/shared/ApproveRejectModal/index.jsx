import React, { Component } from 'react';
import PropTypes from 'prop-types';

export class ApproveRejectModal extends Component {
  state = {
    showRejectionInput: false,
    rejectBtnIsDisabled: false,
    rejectionReason: null,
  };

  handleApproveClick = () => {
    this.props.approveBtnOnClick();
  };

  handleRejectionChange = rejectionReason => {
    this.setState({ rejectionReason });

    if (!rejectionReason) {
      this.setState({ rejectBtnIsDisabled: true });
    } else {
      this.setState({ rejectBtnIsDisabled: false });
    }
  };

  handleRejectionClick = () => {
    if (!this.state.showRejectionInput) {
      this.setState({ showRejectionInput: true, rejectBtnIsDisabled: true });
    } else if (!this.state.rejectionReason) {
      console.error('Rejection reason empty. Please fill out rejection reason.');
    } else {
      this.props.rejectBtnOnClick(this.state.rejectionReason);
    }
  };

  render() {
    const { hideModal } = this.props;
    return (
      !hideModal && (
        <>
          <div>
            <button onClick={this.handleApproveClick}>Approve</button>
            <button onClick={this.handleRejectionClick} disabled={this.state.rejectBtnIsDisabled}>
              Reject
            </button>
          </div>
          <div>
            {this.state.showRejectionInput && (
              <label>
                Rejection reason
                <input
                  name="rejectionReason"
                  onChange={event => this.handleRejectionChange(event.target.value)}
                ></input>
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
  hideModal: PropTypes.bool,
};

ApproveRejectModal.defaultProps = {
  hideModal: false,
};
