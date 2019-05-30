import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import { isValid, isSubmitting, submit, hasSubmitSucceeded } from 'redux-form';
import { bindActionCreators } from 'redux';

import ReleaseFromSitForm from './ReleaseFromSitForm.jsx';
import { updateSitReleaseFromSit } from 'shared/Entities/modules/storageInTransits';
import { formName as ReleaseFromSitFormName } from 'shared/StorageInTransit/ReleaseFromSitForm.jsx';
import './StorageInTransit.css';

export class ReleaseFromSit extends Component {
  constructor(props) {
    super(props);
    this.state = { storageInTransit: {} };
  }

  closeForm = () => {
    this.props.onClose();
  };

  onSubmit = values => {
    this.props.updateSitReleaseFromSit(this.props.sit.shipment_id, this.props.sit.id, values);
  };

  submitReleaseFromSit = () => {
    this.props.submitForm();
  };

  componentDidMount() {
    const { actual_start_date } = this.props.sit;
    this.setState({
      storageInTransit: Object.assign({}, this.props.sit, { actual_start_date: actual_start_date }),
    });
  }

  componentDidUpdate(prevProps) {
    if (this.props.hasSubmitSucceeded && !prevProps.hasSubmitSucceeded) {
      this.props.onClose();
    }
  }

  render() {
    return (
      <div className="storage-in-transit-panel-modal">
        <div className="title">Release shipment from SIT</div>
        <ReleaseFromSitForm minDate={this.state.storageInTransit.actual_start_date} onSubmit={this.onSubmit} />

        <div className="usa-grid-full align-center-vertical">
          <div className="usa-width-one-half">
            <p className="cancel-link">
              <a className="usa-button-secondary" data-cy="release-from-sit-cancel" onClick={this.closeForm}>
                Cancel
              </a>
            </p>
          </div>
          <div className="usa-width-one-half align-right">
            <button
              className="button usa-button-primary"
              data-cy="release-from-sit-button"
              disabled={!this.props.formEnabled}
              onClick={this.submitReleaseFromSit}
            >
              Done
            </button>
          </div>
        </div>
      </div>
    );
  }
}

ReleaseFromSit.propTypes = {
  sit: PropTypes.object.isRequired,
  formEnabled: PropTypes.bool.isRequired,
  updateSitReleaseFromSit: PropTypes.func.isRequired,
  onClose: PropTypes.func.isRequired,
};

function mapStateToProps(state) {
  return {
    formEnabled: isValid(ReleaseFromSitFormName)(state) && !isSubmitting(ReleaseFromSitFormName)(state),
    hasSubmitSucceeded: hasSubmitSucceeded(ReleaseFromSitFormName)(state),
  };
}
function mapDispatchToProps(dispatch) {
  // Bind an action, which submit the form by its name
  return bindActionCreators({ submitForm: () => submit(ReleaseFromSitFormName), updateSitReleaseFromSit }, dispatch);
}

export default connect(mapStateToProps, mapDispatchToProps)(ReleaseFromSit);
