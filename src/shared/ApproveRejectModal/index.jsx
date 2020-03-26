import React, { Component } from 'react';
import PropTypes from 'prop-types';

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
      console.error('Rejection reason empty. Please fill out rejection reason.');
    } else {
      this.props.rejectBtnOnClick(this.state.rejectionReason);
    }
  };

  handleRejectionCancelClick = () => {
    this.setState({
      showApproveBtn: true,
      showRejectionToggleBtn: true,
      rejectBtnIsDisabled: true,
      rejectionReason: null,
      showRejectionInput: false,
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
            {this.state.showRejectionToggleBtn && <button onClick={this.handleRejectionToggleClick}>Reject</button>}
          </div>
          <div>
            {this.state.showRejectionInput && (
              <label>
                Rejection reason
                <input name="rejectionReason" onChange={(event) => this.handleRejectionChange(event.target.value)} />
                <button onClick={this.handleRejectionClick} disabled={this.state.rejectBtnIsDisabled}>
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
