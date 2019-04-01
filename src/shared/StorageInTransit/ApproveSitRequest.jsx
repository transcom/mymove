import React, { Component } from 'react';
import FontAwesomeIcon from '@fortawesome/react-fontawesome';
import faCheck from '@fortawesome/fontawesome-free-solid/faCheck';
import './StorageInTransit.css';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import StorageInTransitOfficeApprovalForm from './StorageInTransitOfficeApprovalForm.jsx';
import './StorageInTransit.css';

export class ApproveSitRequest extends Component {
  state = {
    showForm: false,
  };

  openForm = () => {
    this.setState({ showForm: true });
  };

  closeForm = () => {
    this.setState({ showForm: false });
  };

  render() {
    if (this.state.showForm)
      return (
        <div className="storage-in-transit-panel-modal">
          <div className="title">Approve SIT Request</div>
          <StorageInTransitOfficeApprovalForm />
          <div className="usa-grid-full align-center-vertical">
            <div className="usa-width-one-half">
              <p className="cancel-link">
                <a className="usa-button-secondary" onClick={this.closeForm}>
                  Cancel
                </a>
              </p>
            </div>
            <div className="usa-width-one-half align-right">
              <button
                className="button usa-button-primary storage-in-transit-request-form-send-request-button"
                disabled={!this.props.formEnabled}
              >
                Approve
              </button>
            </div>
          </div>
        </div>
      );
    return (
      <span className="approve-sit">
        <a className="approve-sit-link" onClick={this.openForm}>
          <FontAwesomeIcon className="icon" icon={faCheck} />
          Approve
        </a>
      </span>
    );
  }
}

ApproveSitRequest.propTypes = {};

function mapStateToProps(state) {
  return {};
}
function mapDispatchToProps(dispatch) {
  // Bind an action, which submit the form by its name
  return bindActionCreators({}, dispatch);
}

export default connect(mapStateToProps, mapDispatchToProps)(ApproveSitRequest);
